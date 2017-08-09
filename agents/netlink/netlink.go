//
// Copyright 2017 Nippon Telegraph and Telephone Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// netlink is a Netlink Agent for Lagopus2
package netlink

import (
	"fmt"
	"github.com/lagopus/vsw/utils/notifier"
	"github.com/lagopus/vsw/vswitch"
	"github.com/vishvananda/netlink"
	"io/ioutil"
	"net"
	"os"
	"sync"
	"syscall"
	"unsafe"
)

var log = vswitch.Logger
var netlinkInstance *NetlinkAgent

var stateToNudState = map[int]vswitch.NudState{
	netlink.NUD_NONE:       vswitch.NudStateNone,
	netlink.NUD_INCOMPLETE: vswitch.NudStateIncomplete,
	netlink.NUD_REACHABLE:  vswitch.NudStateReachable,
	netlink.NUD_STALE:      vswitch.NudStateStale,
	netlink.NUD_DELAY:      vswitch.NudStateDelay,
	netlink.NUD_PROBE:      vswitch.NudStateProbe,
	netlink.NUD_FAILED:     vswitch.NudStateFailed,
	netlink.NUD_NOARP:      vswitch.NudStateNoArp,
	netlink.NUD_PERMANENT:  vswitch.NudStatePermanent,
}

type NetlinkAgent struct {
	vswch      chan notifier.Notification
	ruch       chan netlink.RouteUpdate
	ndch       chan NeighUpdate
	nldone     chan struct{}
	vrfs       map[string]*netlink.Vrf
	taps       map[string]*netlink.Tuntap
	links      map[int]*vswitch.VifInfo
	tables     map[int]*vswitch.VrfInfo
	files      map[vswitch.VifIndex]*os.File
	filesMutex sync.RWMutex
}

const rulePref = 1000

func enableIpForwarding(link string) error {
	path := "/proc/sys/net/ipv4/conf/" + link + "/forwarding"
	err := ioutil.WriteFile(path, []byte("1\n"), 0644)
	if err != nil {
		log.Fatalf("Netlink Agent: Can't enable IP forwarding on %s: %v", link, err)
	}
	return err
}

func (n *NetlinkAgent) addVRF(vrf *vswitch.VrfInfo) {
	log.Printf("Netlink Agent: Adding VRF for %v", vrf.Name())

	var tableID int
	fmt.Sscanf(vrf.Name(), "vrf%d", &tableID)

	// Create VRF
	nv := &netlink.Vrf{
		LinkAttrs: netlink.LinkAttrs{
			Name:   vrf.Name(),
			TxQLen: 1000,
		},
		Table: uint32(tableID),
	}
	if err := netlink.LinkAdd(nv); err != nil {
		log.Fatalf("Netlink Agent: Adding VRF %s failed: %v", vrf.Name(), err)
		return
	}

	log.Printf("Netlink Agent: %s has LinkIndex %d", vrf.Name(), nv.Index)

	// Create Rules
	rule := netlink.NewRule()
	rule.IifName = vrf.Name()
	rule.Priority = rulePref
	rule.Table = tableID

	if err := netlink.RuleAdd(rule); err != nil {
		netlink.LinkDel(nv)
		log.Fatalf("Netlink Agent: Adding iif Rule to VRF %s failed: %v", vrf.Name(), err)
		return
	}

	rule.OifName = rule.IifName
	rule.IifName = ""

	if err := netlink.RuleAdd(rule); err != nil {
		netlink.LinkDel(nv)
		log.Fatalf("Netlink Agent: Adding oif Rule to VRF %s failed: %v", vrf.Name(), err)
		return
	}

	// Bring the Link Up
	if err := netlink.LinkSetUp(nv); err != nil {
		netlink.LinkDel(nv)
		log.Fatalf("Netlink Agent: Bringing up VRF %s failed: %v", vrf.Name(), err)
		return
	}

	// Enable forwarding
	if err := enableIpForwarding(vrf.Name()); err != nil {
		netlink.LinkDel(nv)
		return
	}

	n.vrfs[nv.LinkAttrs.Name] = nv
	n.tables[tableID] = vrf
}

func (n *NetlinkAgent) deleteRule(vrfName string) {
	// delete related rule
	rule := netlink.NewRule()
	rule.IifName = vrfName
	if err := netlink.RuleDel(rule); err != nil {
		log.Fatalf("Netlink Agent: Deleting Rule iif= %s failed: %v", vrfName, err)
	}
	rule.OifName = rule.IifName
	rule.IifName = ""
	if err := netlink.RuleDel(rule); err != nil {
		log.Fatalf("Netlink Agent: Deleting Rule oif= %s failed: %v", vrfName, err)
	}
}

func (n *NetlinkAgent) deleteVRF(vrf *vswitch.VrfInfo) {
	if link, ok := n.vrfs[vrf.String()]; ok {
		netlink.LinkDel(link)
		n.deleteRule(link.LinkAttrs.Name)
	}
}

func (n *NetlinkAgent) deleteAllVRF() {
	for _, link := range n.vrfs {
		netlink.LinkDel(link)
		n.deleteRule(link.LinkAttrs.Name)
	}
}

func (n *NetlinkAgent) handleVRFNoti(t notifier.Type, vrf *vswitch.VrfInfo) {
	switch t {
	case notifier.Add:
		n.addVRF(vrf)
	case notifier.Delete:
		n.deleteVRF(vrf)
	default:
		// nop
	}
}

// For TunTap
const (
	sizeOfIfReq = 40
	IFNAMSIZ    = 16
)

type ifReq struct {
	Name  [IFNAMSIZ]byte
	Flags uint16
	pad   [sizeOfIfReq - IFNAMSIZ - 2]byte
}

func (n *NetlinkAgent) addTap(vrf *vswitch.VrfInfo, vif *vswitch.VifInfo) {
	log.Printf("Netlink Agent: Adding Tap for %v", vif)

	// Create a Tap
	nt := &netlink.Tuntap{
		LinkAttrs: netlink.LinkAttrs{
			Name:   vif.String(),
			TxQLen: 1000,
		},
		Mode:  netlink.TUNTAP_MODE_TAP,
		Flags: netlink.TUNTAP_ONE_QUEUE | netlink.TUNTAP_NO_PI,
	}

	if err := netlink.LinkAdd(nt); err != nil {
		log.Fatalf("Netlink Agent: Adding TAP %v failed: %v", vif, err)
		return
	}

	// Set Master
	master := &netlink.Bridge{LinkAttrs: netlink.LinkAttrs{Name: vrf.Name()}}
	if err := netlink.LinkSetMaster(nt, master); err != nil {
		log.Fatalf("Netlink Agent: Setting TAP %v's master to %s failed: %v", vif, vrf.Name(), err)
		netlink.LinkDel(nt)
		return
	}

	// Enable ARP
	if err := netlink.LinkSetARPOn(nt); err != nil {
		log.Fatalf("Netlink Agent: Enable ARP on TAP %v failed: %v", vif, err)
		netlink.LinkDel(nt)
		return
	}

	// Enable forwarding
	if err := enableIpForwarding(vif.String()); err != nil {
		netlink.LinkDel(nt)
		return
	}

	// Open created Tap
	file, err := os.OpenFile("/dev/net/tun", os.O_RDWR, 0)
	if err != nil {
		log.Printf("Netlink Agent: Can't open TAP for %v", vif)
		netlink.LinkDel(nt)
		return
	}

	var req ifReq
	req.Flags = uint16(nt.Flags) | uint16(nt.Mode)
	copy(req.Name[:15], nt.Name)

	_, _, errno := syscall.Syscall(syscall.SYS_IOCTL, file.Fd(), uintptr(syscall.TUNSETIFF), uintptr(unsafe.Pointer(&req)))
	if errno != 0 {
		log.Printf("Netlink Agent: Getting TAP %v failed: %v", vif.String(), errno)
		file.Close()
		//	netlink.LinkDel(nt)
		return
	}

	n.links[nt.Index] = vif
	n.taps[nt.LinkAttrs.Name] = nt

	n.filesMutex.Lock()
	n.files[vif.VifIndex()] = file
	n.filesMutex.Unlock()
}

func (n *NetlinkAgent) deleteTap(vif *vswitch.VifInfo) {
	log.Printf("Netlink Agent: Deleting %v", vif)

	if tap, ok := n.taps[vif.String()]; ok {
		netlink.LinkDel(tap)
		delete(n.taps, vif.String())
		delete(n.links, tap.Index)
	}

	if file, ok := n.files[vif.VifIndex()]; ok {
		file.Close()
		delete(n.files, vif.VifIndex())
	}
}

func (n *NetlinkAgent) deleteAllTap() {
	for _, vif := range n.links {
		n.deleteTap(vif)
	}
}

func (n *NetlinkAgent) configTapHardwareAddr(vif *vswitch.VifInfo, ha net.HardwareAddr) {
	if tap, ok := n.taps[vif.String()]; ok {
		log.Printf("Netlink Agent: Configure MAC for %v to %v", vif, ha)
		// Set Hardware Addr
		if err := netlink.LinkSetHardwareAddr(tap, ha); err != nil {
			log.Fatalf("Netlink Agent: Setting TAP %v's MAC to %s failed: %v", vif, ha, err)
		}
	}
}

func (n *NetlinkAgent) configTapMTU(vif *vswitch.VifInfo, mtu vswitch.MTU) {
	if tap, ok := n.taps[vif.String()]; ok {
		log.Printf("Netlink Agent: Configure MTU for %v to %d", vif, mtu)
		// Set Hardware Addr
		if err := netlink.LinkSetMTU(tap, int(mtu)); err != nil {
			log.Fatalf("Netlink Agent: Setting TAP %v's MTU to %d failed: %v", vif, mtu, err)
		}
	}
}

func (n *NetlinkAgent) configTapLinkStatus(vif *vswitch.VifInfo, stat vswitch.LinkStatus) {
	if tap, ok := n.taps[vif.String()]; ok {
		log.Printf("Netlink Agent: Change link state of %v to %v", vif, stat)
		switch stat {
		case vswitch.LinkUp:
			if err := netlink.LinkSetUp(tap); err != nil {
				log.Fatalf("Netlink Agent: Bringing up Tap %s failed: %v", vif, err)
			}
		case vswitch.LinkDown:
			if err := netlink.LinkSetDown(tap); err != nil {
				log.Fatalf("Netlink Agent: Bringing down Tap %s failed: %v", vif, err)
			}
		}
	} else {
		log.Fatalf("Netlink Agent: Can't change link state of TAP %v", vif)
	}
}

func (n *NetlinkAgent) handleVIFNoti(t notifier.Type, vrf *vswitch.VrfInfo, vif *vswitch.VifInfo) {
	switch t {
	case notifier.Add:
		n.addTap(vrf, vif)
	case notifier.Delete:
		n.deleteTap(vif)
	default:
		// nop
	}
}

func (n *NetlinkAgent) handleIPAddr(t notifier.Type, vif *vswitch.VifInfo, ip vswitch.IPAddr) {
	if tap, ok := n.taps[vif.String()]; ok {
		addr := &netlink.Addr{IPNet: &net.IPNet{ip.IP, ip.Mask}}

		log.Printf("Netlink Agent: %s IP Address %s to TAP %s", t, ip, vif)

		switch t {
		case notifier.Add, notifier.Update:
			if err := netlink.AddrReplace(tap, addr); err != nil {
				log.Fatalf("Netlink Agent: Adding/Replacing IP Address %s to TAP %s failed: %v", ip, vif, err)
			}
		case notifier.Delete:
			if err := netlink.AddrDel(tap, addr); err != nil {
				log.Fatalf("Netlink Agent: Deleting IP Address %s to TAP %s failed: %v", ip, vif, err)
			}
		}
	} else {
		log.Fatalf("Netlink Agent: Can't set IP Address of TAP %v", vif)
	}
}

func (n *NetlinkAgent) handleRouteUpdate(ru netlink.RouteUpdate) {
	vi, ok := n.tables[ru.Table]
	if !ok {
		return
	}

	vif, ok := n.links[ru.LinkIndex]
	if !ok {
		return
	}

	entry := vswitch.Route{
		Dst:      ru.Dst,
		Src:      ru.Src,
		Gw:       ru.Gw,
		Metrics:  ru.Priority,
		VifIndex: vif.VifIndex(),
		Scope:    vswitch.RouteScope(ru.Scope),
	}
	if ru.Type == syscall.RTM_NEWROUTE {
		if !vi.AddEntry(entry) {
			log.Printf("Netlink Agent: Can't add a route entry for %v: %v", vi, entry)
		}
	} else {
		if !vi.DeleteEntry(entry) {
			log.Printf("Netlink Agent: Can't delete a route entry for %v: %v", vi, entry)
		}
	}
}

func (n *NetlinkAgent) handleNeighbourUpdate(nu NeighUpdate) {
	vif, ok := n.links[nu.LinkIndex]
	if !ok {
		return
	}

	// This is same as deleting
	if nu.HardwareAddr == nil {
		nu.Type = syscall.RTM_DELNEIGH
	}

	if nu.Type == syscall.RTM_NEWNEIGH {
		entry := vswitch.Neighbour{
			Dst:           nu.IP,
			LinkLocalAddr: nu.HardwareAddr,
			State:         stateToNudState[nu.State],
		}

		if !vif.AddEntry(entry) {
			log.Printf("Netlink Agent: Neigh entry add failed for %v: %v", vif, entry)
		}
	} else if nu.Type == syscall.RTM_DELNEIGH {
		if !vif.DeleteEntry(nu.IP) {
			log.Printf("Netlink Agent: Neigh entry delete failed for %v: %v", vif, nu.IP)
		}
	}
}

func (n *NetlinkAgent) listen() {
	for {
		select {
		case noti, ok := <-n.vswch:
			if !ok {
				return
			}
			log.Printf("Netlink Agent: VSW: %v\n", noti)

			if vif, ok := noti.Target.(*vswitch.VifInfo); ok {
				switch value := noti.Value.(type) {
				case nil:
					n.handleVIFNoti(noti.Type, nil, vif)

				case net.HardwareAddr:
					if noti.Type == notifier.Update {
						n.configTapHardwareAddr(vif, value)
					}

				case vswitch.MTU:
					if noti.Type == notifier.Update {
						n.configTapMTU(vif, value)
					}

				case vswitch.IPAddr:
					n.handleIPAddr(noti.Type, vif, value)

				case vswitch.LinkStatus:
					if noti.Type == notifier.Update {
						n.configTapLinkStatus(vif, value)
					}

				case vswitch.Neighbour:
					// Don't care. Came from me.

				default:
					log.Printf("Netlink Agent: Unexpectd value: %v\n", vif)
				}

			} else if vrf, ok := noti.Target.(*vswitch.VrfInfo); ok {
				switch vif := noti.Value.(type) {
				case nil:
					n.handleVRFNoti(noti.Type, vrf)

				case *vswitch.VifInfo:
					n.handleVIFNoti(noti.Type, vrf, vif)

				case vswitch.Route:
					// Don't care. This should have came out from me.

				default:
					log.Printf("Netlink Agent: Unexpectd value: %v\n", vif)
				}
			} else {
				log.Printf("Netlink Agent: Unexpectd target: %v\n", noti.Target)
			}

		case ru, ok := <-n.ruch:
			if !ok {
				return
			}
			log.Printf("Netlink Agent: RU: %v", ru)
			n.handleRouteUpdate(ru)

		case nu, ok := <-n.ndch:
			if !ok {
				return
			}
			log.Printf("Netlink Agent: NU: %v\n", nu)
			n.handleNeighbourUpdate(nu)
		}
	}
}

func GetTapFile(vifidx vswitch.VifIndex) (*os.File, bool) {
	netlinkInstance.filesMutex.RLock()
	defer netlinkInstance.filesMutex.RUnlock()
	file, ok := netlinkInstance.files[vifidx]
	return file, ok
}

func (n *NetlinkAgent) Start() bool {
	// Listen to changes on VIF/VRF
	n.vswch = vswitch.GetNotifier().Listen()

	if err := netlink.RouteSubscribe(n.ruch, n.nldone); err != nil {
		log.Printf("Netlink Agent: Can't receive route update: %v", err)
		return false
	}

	if err := NeighSubscribe(n.ndch, n.nldone); err != nil {
		log.Printf("Netlink Agent: Can't receive neighbour update: %v", err)
		return false
	}

	go func() {
		n.listen()
	}()

	return true
}

func (n *NetlinkAgent) Stop() {
	n.deleteAllVRF()
	n.deleteAllTap()
	vswitch.GetNotifier().Close(n.vswch)
	close(n.nldone)
}

func (n *NetlinkAgent) String() string {
	return "Netlink Agent"
}

func init() {
	netlinkInstance = &NetlinkAgent{
		ruch:   make(chan netlink.RouteUpdate),
		ndch:   make(chan NeighUpdate),
		nldone: make(chan struct{}),
		vrfs:   make(map[string]*netlink.Vrf),
		taps:   make(map[string]*netlink.Tuntap),
		links:  make(map[int]*vswitch.VifInfo),
		tables: make(map[int]*vswitch.VrfInfo),
		files:  make(map[vswitch.VifIndex]*os.File),
	}
	vswitch.RegisterAgent(netlinkInstance)
}
