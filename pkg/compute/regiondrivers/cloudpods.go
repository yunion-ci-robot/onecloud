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

package regiondrivers

import (
	"context"

	"yunion.io/x/pkg/utils"

	api "yunion.io/x/onecloud/pkg/apis/compute"
	"yunion.io/x/onecloud/pkg/compute/models"
	"yunion.io/x/onecloud/pkg/httperrors"
	"yunion.io/x/onecloud/pkg/mcclient"
)

type SCloudpodsRegionDriver struct {
	SManagedVirtualizationRegionDriver
}

func init() {
	driver := SCloudpodsRegionDriver{}
	models.RegisterRegionDriver(&driver)
}

func (self *SCloudpodsRegionDriver) GetProvider() string {
	return api.CLOUD_PROVIDER_CLOUDPODS
}

func (self *SCloudpodsRegionDriver) IsAllowSecurityGroupNameRepeat() bool {
	return false
}

func (self *SCloudpodsRegionDriver) GenerateSecurityGroupName(name string) string {
	return name
}

func (self *SCloudpodsRegionDriver) ValidateCreateVpcData(ctx context.Context, userCred mcclient.TokenCredential, input api.VpcCreateInput) (api.VpcCreateInput, error) {
	if !utils.IsInStringArray(input.CidrBlock, []string{"192.168.0.0/16", "10.0.0.0/8", "172.16.0.0/12"}) {
		return input, httperrors.NewInputParameterError("Invalid cidr_block, want 192.168.0.0/16|10.0.0.0/8|172.16.0.0/12, got %s", input.CidrBlock)
	}
	return input, nil
}
