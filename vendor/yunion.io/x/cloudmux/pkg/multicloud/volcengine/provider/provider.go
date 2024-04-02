// Copyright 2023 Yunion
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

package volcengine

import (
	"context"
	"fmt"
	"strings"

	"yunion.io/x/pkg/errors"

	api "yunion.io/x/cloudmux/pkg/apis/compute"
	"yunion.io/x/cloudmux/pkg/cloudprovider"
	"yunion.io/x/cloudmux/pkg/multicloud/volcengine"
	"yunion.io/x/jsonutils"
)

type SVolcEngineProviderFactory struct {
	cloudprovider.SPublicCloudBaseProviderFactory
}

func (self *SVolcEngineProviderFactory) GetId() string {
	return volcengine.CLOUD_PROVIDER_VOLCENGINE
}

func (self *SVolcEngineProviderFactory) GetName() string {
	return volcengine.CLOUD_PROVIDER_VOLCENGINE_CN
}

func (self *SVolcEngineProviderFactory) ValidateCreateCloudaccountData(ctx context.Context, input cloudprovider.SCloudaccountCredential) (cloudprovider.SCloudaccount, error) {
	output := cloudprovider.SCloudaccount{}
	if len(input.AccessKeyId) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_id")
	}
	if len(input.AccessKeySecret) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_secret")
	}
	output.Account = input.AccessKeyId
	output.Secret = input.AccessKeySecret
	return output, nil
}

func (f *SVolcEngineProviderFactory) ValidateUpdateCloudaccountCredential(ctx context.Context, input cloudprovider.SCloudaccountCredential, cloudaccount string) (cloudprovider.SCloudaccount, error) {
	output := cloudprovider.SCloudaccount{}
	if len(input.AccessKeyId) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_id")
	}
	if len(input.AccessKeySecret) == 0 {
		return output, errors.Wrap(cloudprovider.ErrMissingParameter, "access_key_secret")
	}
	output = cloudprovider.SCloudaccount{
		Account: input.AccessKeyId,
		Secret:  input.AccessKeySecret,
	}
	return output, nil
}

func (factory *SVolcEngineProviderFactory) IsSupportSAMLAuth() bool {
	return true
}

func validateClientCloudenv(client *volcengine.SVolcEngineClient) error {
	regions := client.GetIRegions()
	if len(regions) == 0 {
		return nil
	}
	return nil
}

func parseAccount(account string) (accessKey string, projectId string) {
	segs := strings.Split(account, "::")
	if len(segs) == 2 {
		accessKey = segs[0]
		projectId = segs[1]
	} else {
		accessKey = account
		projectId = ""
	}

	return
}

func (self *SVolcEngineProviderFactory) GetProvider(cfg cloudprovider.ProviderConfig) (cloudprovider.ICloudProvider, error) {
	accessKey, accountId := parseAccount(cfg.Account)
	client, err := volcengine.NewVolcEngineClient(
		volcengine.NewVolcEngineClientConfig(
			accessKey,
			cfg.Secret,
		).CloudproviderConfig(cfg).AccountId(accountId),
	)
	if err != nil {
		return nil, err
	}

	err = validateClientCloudenv(client)
	if err != nil {
		return nil, errors.Wrap(err, "validateClientCloudenv")
	}

	return &SVolcEngineProvider{
		SBaseProvider: cloudprovider.NewBaseProvider(self),
		client:        client,
	}, nil
}

func (self *SVolcEngineProviderFactory) GetClientRC(info cloudprovider.SProviderInfo) (map[string]string, error) {
	accessKey, accountId := parseAccount(info.Account)
	return map[string]string{
		"VOLCENGINE_ACCESS_KEY": accessKey,
		"VOLCENGINE_SECRET_KEY": info.Secret,
		"VOLCENGINE_REGION":     volcengine.VOLCENGINE_DEFAULT_REGION,
		"VOLCENGINE_ACCOUNT_ID": accountId,
	}, nil
}

func init() {
	factory := SVolcEngineProviderFactory{}
	cloudprovider.RegisterFactory(&factory)
}

type SVolcEngineProvider struct {
	cloudprovider.SBaseProvider

	client *volcengine.SVolcEngineClient
}

func (self *SVolcEngineProvider) GetAccountId() string {
	return self.client.GetAccountId()
}

func (self *SVolcEngineProvider) GetSysInfo() (jsonutils.JSONObject, error) {
	regions := self.client.GetIRegions()
	info := jsonutils.NewDict()
	info.Add(jsonutils.NewInt(int64(len(regions))), "region_count")
	info.Add(jsonutils.NewString(volcengine.VOLCENGINE_API_VERSION), "api_version")
	return info, nil
}

func (self *SVolcEngineProvider) GetBalance() (*cloudprovider.SBalanceInfo, error) {
	ret := &cloudprovider.SBalanceInfo{Currency: "CNY", Status: api.CLOUD_PROVIDER_HEALTH_UNKNOWN}
	balance, err := self.client.QueryBalance()
	if err != nil {
		return ret, err
	}

	ret.Status = api.CLOUD_PROVIDER_HEALTH_NORMAL
	ret.Amount = balance.AvailableBalance

	if ret.Amount < 0 {
		ret.Status = api.CLOUD_PROVIDER_HEALTH_ARREARS
	} else if ret.Amount < 100 {
		ret.Status = api.CLOUD_PROVIDER_HEALTH_INSUFFICIENT
	}
	return ret, nil
}

func (self *SVolcEngineProvider) GetBucketCannedAcls(regionId string) []string {
	return []string{
		string(cloudprovider.ACLPrivate),
		string(cloudprovider.ACLPublicRead),
		string(cloudprovider.ACLPublicReadWrite),
	}
}

func (self *SVolcEngineProvider) GetObjectCannedAcls(regionId string) []string {
	return []string{
		string(cloudprovider.ACLPrivate),
		string(cloudprovider.ACLPublicRead),
		string(cloudprovider.ACLPublicReadWrite),
	}
}

func (self *SVolcEngineProvider) GetStorageClasses(regionId string) []string {
	return []string{
		"STANDARD", "IA", "ARCHIVE_FR", "INTELLIGENT_TIERING", "COLD_ARCHIVE",
	}
}

func (self *SVolcEngineProvider) GetCapabilities() []string {
	return self.client.GetCapabilities()
}

func (self *SVolcEngineProvider) GetIProjects() ([]cloudprovider.ICloudProject, error) {
	return self.client.GetIProjects()
}

func (self *SVolcEngineProvider) GetIRegionById(extId string) (cloudprovider.ICloudRegion, error) {
	return self.client.GetIRegionById(extId)
}

func (self *SVolcEngineProvider) GetIRegions() []cloudprovider.ICloudRegion {
	return self.client.GetIRegions()
}

func (self *SVolcEngineProvider) GetSubAccounts() ([]cloudprovider.SSubAccount, error) {
	return self.client.GetSubAccounts()
}

func (self *SVolcEngineProvider) GetVersion() string {
	return volcengine.VOLCENGINE_API_VERSION
}

func (self *SVolcEngineProvider) GetMetrics(opts *cloudprovider.MetricListOptions) ([]cloudprovider.MetricValues, error) {
	return self.client.GetMetrics(opts)
}

func (self *SVolcEngineProvider) CreateICloudSAMLProvider(opts *cloudprovider.SAMLProviderCreateOptions) (cloudprovider.ICloudSAMLProvider, error) {
	sp, err := self.client.CreateSAMLProvider(opts.Name, opts.Metadata.String(), opts.Desc)
	if err != nil {
		return nil, errors.Wrapf(err, "CreateSAMLProvider")
	}
	return sp, nil
}

func (self *SVolcEngineProvider) GetICloudSAMLProviders() ([]cloudprovider.ICloudSAMLProvider, error) {
	return self.client.GetICloudSAMLProviders()
}

func (self *SVolcEngineProvider) GetICloudroles() ([]cloudprovider.ICloudrole, error) {
	return self.client.GetICloudroles()
}

func (self *SVolcEngineProvider) GetICloudroleById(id string) (cloudprovider.ICloudrole, error) {
	role, err := self.client.GetRole(id)
	if err != nil {
		return nil, errors.Wrapf(err, "GetRole(%s)", id)
	}
	return role, nil
}

func (self *SVolcEngineProvider) CreateICloudrole(opts *cloudprovider.SRoleCreateOptions) (cloudprovider.ICloudrole, error) {
	stetement := fmt.Sprintf(`{"Statement":[{"Effect":"Allow","Action":["sts:AssumeRoleWithSAML"],"Principal":{"Federated":["trn:iam::%s:saml-provider/%s"]}}]}`, self.client.GetAccountId(), opts.SAMLProvider)
	role, err := self.client.CreateRole(opts.Name, stetement, opts.Desc)
	if err != nil {
		return nil, errors.Wrapf(err, "CreateRole")
	}
	return role, nil
}

func (self *SVolcEngineProvider) CreateIClouduser(conf *cloudprovider.SClouduserCreateConfig) (cloudprovider.IClouduser, error) {
	return self.client.CreateIClouduser(conf)
}

func (self *SVolcEngineProvider) GetICloudusers() ([]cloudprovider.IClouduser, error) {
	return self.client.GetICloudusers()
}

func (self *SVolcEngineProvider) GetIClouduserByName(name string) (cloudprovider.IClouduser, error) {
	return self.client.GetIClouduserByName(name)
}

func (self *SVolcEngineProvider) GetICloudgroups() ([]cloudprovider.ICloudgroup, error) {
	return self.client.GetICloudgroups()
}

func (self *SVolcEngineProvider) CreateICloudgroup(name, desc string) (cloudprovider.ICloudgroup, error) {
	return self.client.CreateICloudgroup(name, desc)
}

func (self *SVolcEngineProvider) GetICloudgroupByName(name string) (cloudprovider.ICloudgroup, error) {
	return self.client.GetICloudgroupByName(name)
}

func (self *SVolcEngineProvider) GetISystemCloudpolicies() ([]cloudprovider.ICloudpolicy, error) {
	return self.client.GetISystemCloudpolicies()
}

func (self *SVolcEngineProvider) GetICustomCloudpolicies() ([]cloudprovider.ICloudpolicy, error) {
	return self.client.GetICustomCloudpolicies()
}

func (self *SVolcEngineProvider) CreateICloudpolicy(opts *cloudprovider.SCloudpolicyCreateOptions) (cloudprovider.ICloudpolicy, error) {
	return nil, cloudprovider.ErrNotImplemented
}

func (self *SVolcEngineProvider) GetSamlEntityId() string {
	return cloudprovider.SAML_ENTITY_ID_VOLC_ENGINE
}
