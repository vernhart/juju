// Copyright 2014 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package vmware_test

import (
	jc "github.com/juju/testing/checkers"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/environs"
	"github.com/juju/juju/instance"
	"github.com/juju/juju/provider/common"
	"github.com/juju/juju/provider/vmware"
)

type environInstanceSuite struct {
	vmware.BaseSuite
}

var _ = gc.Suite(&environInstanceSuite{})

func (s *environInstanceSuite) SetUpTest(c *gc.C) {
	s.BaseSuite.SetUpTest(c)
}

func (s *environAvailzonesSuite) TestInstances(c *gc.C) {
	client := vmware.ExposeEnvFakeClient(s.Env)
	client.SetPropertyProxyHandler("FakeDatacenter", vmware.RetrieveDatacenterProperties)
	vmName1 := common.MachineFullName(s.Env, "1")
	vmName2 := common.MachineFullName(s.Env, "2")
	s.FakeInstancesWithResourcePool(client, vmware.InstRp{Inst: vmName1, Rp: "rp1"}, vmware.InstRp{Inst: vmName2, Rp: "rp2"})

	instances, err := s.Env.Instances([]instance.Id{instance.Id(vmName1), instance.Id(vmName2)})

	c.Assert(err, jc.ErrorIsNil)
	c.Assert(len(instances), gc.Equals, 2)
	c.Assert(string(instances[0].Id()), gc.Equals, vmName1)
	c.Assert(string(instances[1].Id()), gc.Equals, vmName2)
}

func (s *environAvailzonesSuite) TestInstancesReturnNoInstances(c *gc.C) {
	client := vmware.ExposeEnvFakeClient(s.Env)
	client.SetPropertyProxyHandler("FakeDatacenter", vmware.RetrieveDatacenterProperties)
	s.FakeInstancesWithResourcePool(client, vmware.InstRp{Inst: "Some name that don't match naming convention", Rp: "rp1"})

	_, err := s.Env.Instances([]instance.Id{instance.Id("Some other name")})

	c.Assert(err, gc.Equals, environs.ErrNoInstances)
}

func (s *environAvailzonesSuite) TestInstancesReturnPartialInstances(c *gc.C) {
	client := vmware.ExposeEnvFakeClient(s.Env)
	client.SetPropertyProxyHandler("FakeDatacenter", vmware.RetrieveDatacenterProperties)
	vmName1 := common.MachineFullName(s.Env, "1")
	vmName2 := common.MachineFullName(s.Env, "2")
	s.FakeInstancesWithResourcePool(client, vmware.InstRp{Inst: vmName1, Rp: "rp1"}, vmware.InstRp{Inst: "Some inst", Rp: "rp2"})

	_, err := s.Env.Instances([]instance.Id{instance.Id(vmName1), instance.Id(vmName2)})

	c.Assert(err, gc.Equals, environs.ErrPartialInstances)
}
