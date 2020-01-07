package verify

import (
	"encoding/json"
	"fmt"

	"github.com/trrtly/esign/context"
	"github.com/trrtly/esign/util"
)

const (
	individualBaseURL = "%s/v2/identity/verify/individual/base"
)

// Individual struct
type Individual struct {
	*context.Context
}

// NewIndividual init
func NewIndividual(ctx *context.Context) *Individual {
	return &Individual{ctx}
}

// ResponseIndividualBase 个人2要素信息比对返回值
type ResponseIndividualBase struct {
	context.CommonError
	Data struct {
		VerifyID string `json:"verifyId"`
	} `json:"data"`
}

// GetBaseURL 个人2要素信息比对请求地址
func (idl *Individual) GetBaseURL() string {
	return fmt.Sprintf(individualBaseURL, idl.GetDomain())
}

// Base 个人2要素信息比对
// http://open.esign.cn/docs/identity/信息比对/个人2要素信息比对.html
func (idl *Individual) Base(idno, name string) (response *ResponseIndividualBase, err error) {
	header, err := idl.GetRequestHeader()
	if err != nil {
		return
	}
	request := map[string]string{
		"idNo": idno,
		"name": name,
	}
	resJSON, err := util.PostJSONWithHeader(idl.GetBaseURL(), request, header)
	if err != nil {
		return
	}
	response = &ResponseIndividualBase{}
	err = json.Unmarshal(resJSON, &response)
	if err != nil {
		return
	}
	return
}
