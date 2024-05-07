// Copyright 2019 Yunion
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

package hostdrivers

import (
	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/compute/models"
)

type SOracleHostDriver struct {
	SManagedVirtualizationHostDriver
}

func init() {
	driver := SOracleHostDriver{}
	models.RegisterHostDriver(&driver)
}

func (self *SOracleHostDriver) GetHostType() string {
	return api.HOST_TYPE_ORACLE
}

func (self *SOracleHostDriver) GetHypervisor() string {
	return api.HYPERVISOR_ORACLE
}

func (self *SOracleHostDriver) GetProvider() string {
	return api.CLOUD_PROVIDER_ORACLE
}
