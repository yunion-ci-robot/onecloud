package aliyun

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"

	"yunion.io/x/jsonutils"
	"yunion.io/x/log"
)

func businessRequest(client *sdk.Client, apiName string, params map[string]string) (jsonutils.JSONObject, error) {
	return _jsonRequest(client, "business.aliyuncs.com", ALIYUN_BSS_API_VERSION, apiName, params)
}

func (self *SAliyunClient) businessRequest(apiName string, params map[string]string) (jsonutils.JSONObject, error) {
	cli, err := self.getDefaultClient()
	if err != nil {
		return nil, err
	}
	return businessRequest(cli, apiName, params)
}

type SAccountBalance struct {
	AvailableAmount     float64
	AvailableCashAmount float64
	CreditAmount        float64
	MybankCreditAmount  float64
	Currency            string
}

type SCashCoupon struct {
	ApplicableProducts  string
	ApplicableScenarios string
	Balance             float64
	CashCouponId        string
	CashCouponNo        string
	EffectiveTime       time.Time
	ExpiryTime          time.Time
	GrantedTime         time.Time
	NominalValue        float64
	Status              string
}

type SPrepaidCard struct {
	PrepaidCardId       string
	PrepaidCardNo       string
	GrantedTime         time.Time
	EffectiveTime       time.Time
	ExpiryTime          time.Time
	NominalValue        float64
	Balance             float64
	ApplicableProducts  string
	ApplicableScenarios string
}

func (self *SAliyunClient) QueryAccountBalance() (*SAccountBalance, error) {
	body, err := self.businessRequest("QueryAccountBalance", nil)
	if err != nil {
		log.Errorf("QueryAccountBalance fail %s", err)
		return nil, err
	}
	balance := SAccountBalance{}
	err = body.Unmarshal(&balance, "Data")
	if err != nil {
		log.Errorf("Unmarshal AccountBalance fail %s", err)
		return nil, err
	}
	return &balance, nil
}

func (self *SAliyunClient) QueryCashCoupons() ([]SCashCoupon, error) {
	params := make(map[string]string)
	params["EffectiveOrNot"] = "True"
	body, err := self.businessRequest("QueryCashCoupons", params)
	if err != nil {
		log.Errorf("QueryCashCoupons fail %s", err)
		return nil, err
	}
	coupons := make([]SCashCoupon, 0)
	err = body.Unmarshal(&coupons, "Data", "CashCoupon")
	if err != nil {
		log.Errorf("Unmarshal fail %s", err)
		return nil, err
	}
	return coupons, nil
}

func (self *SAliyunClient) QueryPrepaidCards() ([]SPrepaidCard, error) {
	params := make(map[string]string)
	params["EffectiveOrNot"] = "True"
	body, err := self.businessRequest("QueryPrepaidCards", params)
	if err != nil {
		log.Errorf("QueryPrepaidCards fail %s", err)
		return nil, err
	}
	cards := make([]SPrepaidCard, 0)
	err = body.Unmarshal(&cards, "Data", "PrepaidCard")
	if err != nil {
		log.Errorf("Unmarshal fail %s", err)
		return nil, err
	}
	return cards, nil
}
