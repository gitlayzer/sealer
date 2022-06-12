// Copyright © 2021 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package apply

import (
	"strconv"

	"github.com/sealerio/sealer/apply/applydriver"

	"github.com/sealerio/sealer/common"
	v2 "github.com/sealerio/sealer/types/api/v2"
	"github.com/sealerio/sealer/utils/net"
)

type ClusterArgs struct {
	cluster   *v2.Cluster
	imageName string
	runArgs   *Args
	hosts     []v2.Host
}

func PreProcessIPList(joinArgs *Args) error {
	if err := net.AssemblyIPList(&joinArgs.Masters); err != nil {
		return err
	}
	if err := net.AssemblyIPList(&joinArgs.Nodes); err != nil {
		return err
	}
	return nil
}

func (c *ClusterArgs) SetClusterArgs() error {
	c.cluster.APIVersion = common.APIVersion
	c.cluster.Kind = common.Cluster
	c.cluster.Name = c.runArgs.ClusterName
	c.cluster.Spec.Image = c.imageName
	c.cluster.Spec.SSH.User = c.runArgs.User
	c.cluster.Spec.SSH.Pk = c.runArgs.Pk
	c.cluster.Spec.SSH.PkPasswd = c.runArgs.PkPassword
	c.cluster.Spec.SSH.Port = strconv.Itoa(int(c.runArgs.Port))
	c.cluster.Spec.Env = append(c.cluster.Spec.Env, c.runArgs.CustomEnv...)
	c.cluster.Spec.CMDArgs = append(c.cluster.Spec.CMDArgs, c.runArgs.CMDArgs...)
	if c.runArgs.Password != "" {
		c.cluster.Spec.SSH.Passwd = c.runArgs.Password
	}
	err := PreProcessIPList(c.runArgs)
	if err != nil {
		return err
	}
	if net.IsIPList(c.runArgs.Masters) && (net.IsIPList(c.runArgs.Nodes) || c.runArgs.Nodes == "") {
		c.hosts = []v2.Host{}
		c.cluster.Spec.Hosts = c.hosts
	} else {
		ip, err := net.GetLocalDefaultIP()
		if err != nil {
			return err
		}
		c.cluster.Spec.Hosts = []v2.Host{
			{
				IPS:   []string{ip},
				Roles: []string{common.MASTER},
			},
		}
	}
	return err
}

func NewApplierFromArgs(imageName string, runArgs *Args) (applydriver.Interface, error) {
	c := &ClusterArgs{
		cluster:   &v2.Cluster{},
		imageName: imageName,
		runArgs:   runArgs,
	}
	if err := c.SetClusterArgs(); err != nil {
		return nil, err
	}
	return NewApplier(c.cluster, nil)
}
