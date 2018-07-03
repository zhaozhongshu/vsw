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

package receiver

import (
	"bytes"
	"testing"

	"github.com/lagopus/vsw/vswitch"
	"github.com/stretchr/testify/suite"
)

type PFKeyMessageTestSuit struct {
	suite.Suite
}

func Test_PFKeyMessageTestSuite(t *testing.T) {
	suite.Run(t, new(PFKeyMessageTestSuit))
}

func (s *PFKeyMessageTestSuit) TestparseSadbAddMsg() {
	b := []byte{
		0x02, 0x00, 0x01, 0x00, 0xc8, 0xca, 0xc0, 0x29, 0x20, 0x01, 0x03, 0x0c, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc5, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1d, 0x7d, 0x80, 0x58, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x20, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x20, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x07, 0x00, 0xff, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x08, 0x00, 0xa0, 0x00, 0x00, 0x00,
		0x82, 0x45, 0x2e, 0xb8, 0x73, 0x29, 0xe9, 0x11, 0x11, 0xf8, 0x2c, 0xba, 0x26, 0x86, 0x56, 0xf4,
		0xe1, 0x7a, 0xf6, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x09, 0x00, 0x80, 0x00, 0x00, 0x00,
		0x34, 0xf4, 0x6b, 0xf2, 0x07, 0x2d, 0x86, 0xb8, 0xf6, 0x74, 0xe8, 0x78, 0x57, 0xd2, 0x0e, 0xe6,
		0x02, 0x00, 0x13, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbAddMsg{}
	err := smsg.Parse(r)

	s.Assert().NoError(err)
	s.Assert().NotNil(smsg)
}

func (s *PFKeyMessageTestSuit) TestparseSadbAddMsgWithError() {
	b := []byte{
		0x04, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbAddMsg{}
	err := smsg.Parse(r)

	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbUpdateMsg() {
	b := []byte{
		0x02, 0x00, 0x01, 0x00, 0xc8, 0xca, 0xc0, 0x29, 0x20, 0x01, 0x03, 0x0c, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x10, 0x0e, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xc5, 0x0b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x1d, 0x7d, 0x80, 0x58, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x20, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x20, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x07, 0x00, 0xff, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x04, 0x00, 0x08, 0x00, 0xa0, 0x00, 0x00, 0x00,
		0x82, 0x45, 0x2e, 0xb8, 0x73, 0x29, 0xe9, 0x11, 0x11, 0xf8, 0x2c, 0xba, 0x26, 0x86, 0x56, 0xf4,
		0xe1, 0x7a, 0xf6, 0x5b, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x09, 0x00, 0x80, 0x00, 0x00, 0x00,
		0x34, 0xf4, 0x6b, 0xf2, 0x07, 0x2d, 0x86, 0xb8, 0xf6, 0x74, 0xe8, 0x78, 0x57, 0xd2, 0x0e, 0xe6,
		0x02, 0x00, 0x13, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbUpdateMsg{}
	err := smsg.Parse(r)

	s.Assert().NoError(err)
	s.Assert().NotNil(smsg)
}

func (s *PFKeyMessageTestSuit) TestparseSadbGetSPIMsgReq() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0xc0, 0xff, 0xff, 0xff, 0xcf, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbGetSPIMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbGetSPIMsgReqWithError() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbGetSPIMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbGetMsg() {
	b := []byte{
		0x02, 0x00, 0x01, 0x00, 0xc8, 0xca, 0xc0, 0x29, 0x20, 0x01, 0x03, 0x0c, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbGetMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbGetMsgWithError() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbGetMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbDeleteMsg() {
	b := []byte{
		0x02, 0x00, 0x01, 0x00, 0xc8, 0xca, 0xc0, 0x29, 0x20, 0x01, 0x03, 0x0c, 0x00, 0x00, 0x00, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbDeleteMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbDeleteMsgWithError() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbDeleteMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDAddMsg() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xea, 0x05, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDAddMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDAddMsg2() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0xff, 0x18, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xc0, 0xa8, 0xc8, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0xff, 0x18, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xc0, 0xa8, 0xc9, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x04, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x7e, 0xd5, 0x89, 0x58, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x08, 0x00, 0x12, 0x00, 0x02, 0x00, 0x02, 0x00, 0xb9, 0x01, 0x00, 0x00, 0x43, 0x0b, 0x00, 0x00,
		0x30, 0x00, 0x32, 0x00, 0x02, 0x03, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDAddMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDAddMsgWithError() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDAddMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDUpdateMsg() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xea, 0x05, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDUpdateMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDUpdateMsgWithError() {
	b := []byte{
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDUpdateMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDGetMsg() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02, 0x00, 0x0f, 0x00, 0x00, 0x00, 0x80, 0xea, 0x05, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDGetMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDGetMsgWithError() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDGetMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDDeleteMsg() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x80, 0xea, 0x05, 0x00,
		0x03, 0x00, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0d,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x02, 0x00, 0x00, 0x00, 0xac, 0x10, 0x01, 0x0c, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDDeleteMsg{}
	err := smsg.Parse(r)
	s.Assert().NotNil(smsg)
	s.Assert().NoError(err)
}

func (s *PFKeyMessageTestSuit) TestparseSadbXSPDDeleteMsgWithError() {
	b := []byte{
		0x02, 0x00, 0x12, 0x00, 0x00, 0x00, 0x02,
	}

	r := bytes.NewReader(b)
	smsg := sadbXSPDDeleteMsg{}
	err := smsg.Parse(r)
	s.Assert().Error(err)
}

func (s *PFKeyMessageTestSuit) TestNewMsgMux() {
	m := NewMsgMux()
	s.Assert().NotNil(m)
}

func (s *PFKeyMessageTestSuit) TestNewMsgMuxForVRF() {
	m1 := NewMsgMuxForVRF(1)
	s.Assert().NotNil(m1)
	m2 := NewMsgMuxForVRF(2)
	s.Assert().NotNil(m2)
	var i vswitch.VRFIndex = 1
	for _, v := range m1.MsgMux {
		s.Assert().Equal(i, vrfMap[v])
	}
	i = 2
	for _, v := range m2.MsgMux {
		s.Assert().Equal(i, vrfMap[v])
	}
	m1.Free()
	for _, v := range m1.MsgMux {
		_, ok := vrfMap[v]
		s.Assert().False(ok)
	}
}