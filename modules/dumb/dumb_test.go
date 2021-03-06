//
// Copyright 2017-2019 Nippon Telegraph and Telephone Corporation.
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

package dumb

import (
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/lagopus/vsw/dpdk"
	"github.com/lagopus/vsw/modules/testvif"
	"github.com/lagopus/vsw/vswitch"
)

var tx_chan chan *dpdk.Mbuf
var rx_chan chan *dpdk.Mbuf
var vif_mac net.HardwareAddr
var pool *dpdk.MemPool

func send(t *testing.T, mbuf *dpdk.Mbuf) bool {
	eh := mbuf.EtherHdr()
	src_ha := eh.SrcAddr()
	dst_ha := eh.DstAddr()

	//
	t.Logf("Sending: %s -> %s\n", src_ha, dst_ha)

	// send
	rx_chan <- mbuf

	// recv
	rmbuf := <-tx_chan
	reh := rmbuf.EtherHdr()
	md := (*vswitch.Metadata)(rmbuf.Metadata())

	t.Logf("Rcv'd: src=%s, dst=%s, vif=%d\n", reh.SrcAddr(), reh.DstAddr(), md.InVIF())

	return bytes.Compare(reh.SrcAddr(), src_ha) == 0 &&
		bytes.Compare(reh.DstAddr(), dst_ha) == 0 &&
		md.InVIF() == 1
}

func TestNormalFlow(t *testing.T) {
	t.Logf("testvif mac address: %s\n", vif_mac)

	src_ha, _ := net.ParseMAC("11:22:33:44:55:66")
	dst_ha, _ := net.ParseMAC("aa:bb:cc:dd:ee:ff")

	mbuf := pool.AllocMbuf()
	eh := mbuf.EtherHdr()
	eh.SetSrcAddr(src_ha)
	eh.SetDstAddr(dst_ha)

	if !send(t, mbuf) {
		t.Errorf("Unexpected packet metadata recv'd")
	}

	// send to testvif
	eh.SetDstAddr(vif_mac)

	if !send(t, mbuf) {
		t.Errorf("Unexpected packet metadata recv'd")
	}
}

func TestMain(m *testing.M) {
	// Initialize vswitch core
	vswitch.Init("../../vsw.conf")
	pool = vswitch.GetDpdkResource().Mempool

	//
	// Setup Vswitch
	//

	// Create Instances
	tv0, _ := vswitch.NewInterface("testvif", "tv0", nil)
	tv0_0, _ := tv0.NewVIF(0)
	dumb, _ := vswitch.NewTestModule("dumb", "dumb0", nil)
	testif, _ := tv0.Instance().(*testvif.TestIF)

	// Connect Instances
	dumb.AddVIF(tv0_0)

	// Get Channels
	tx_chan = testif.TxChan()
	rx_chan = testif.RxChan()
	vif_mac = tv0.MACAddress()

	// Enable Modules
	tv0.Enable()
	tv0_0.Enable()
	if err := dumb.Enable(); err != nil {
		fmt.Printf("Can't enable dumb module: %v\n", err)
		os.Exit(1)
	}

	// Execute test
	flag.Parse()
	rc := m.Run()

	// Teardown
	tv0.Disable()
	tv0_0.Disable()
	dumb.Disable()

	// Done
	os.Exit(rc)
}
