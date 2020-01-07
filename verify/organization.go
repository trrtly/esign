package verify

import (
	"encoding/json"
	"fmt"

	"github.com/trrtly/esign/context"
	"github.com/trrtly/esign/util"
)

// Organization struct
type Organization struct {
	*context.Context
}

const (
	organizationBureau3URL = "%s/v2/identity/verify/organization/enterprise/bureau3Factors"
)

// NewOrganization init
func NewOrganization(ctx *context.Context) *Organization {
	return &Organization{ctx}
}

// ResponseOrganizationBureau3 企业3要素信息比对
type ResponseOrganizationBureau3 struct {
	context.CommonError
	Data struct {
		VerifyID string `json:"verifyId"`
	} `json:"data"`
}

// GetBureau3URL 企业3要素信息比对请求地址
func (org *Organization) GetBureau3URL() string {
	return fmt.Sprintf(organizationBureau3URL, org.GetDomain())
}

// Bureau3Factors 企业3要素信息比对
// http://open.esign.cn/docs/identity/信息比对/企业3要素信息比对.html
func (org *Organization) Bureau3Factors(name, orgCode, legalRepName string) (response *ResponseOrganizationBureau3, err error) {
	header, err := org.GetRequestHeader()
	if err != nil {
		return
	}
	request := map[string]string{
		"name":         name,
		"orgCode":      orgCode,
		"legalRepName": legalRepName,
	}
	resJSON, err := util.PostJSONWithHeader(org.GetBureau3URL(), request, header)
	if err != nil {
		return
	}
	response = &ResponseOrganizationBureau3{}
	err = json.Unmarshal(resJSON, &response)
	if err != nil {
		return
	}
	return
}
