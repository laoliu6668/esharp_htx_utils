package htx

import "encoding/json"

type ApiConfigModel struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	AccountId int64  `json:"account_id"`
	Uid       int64  `json:"uid"`
}

type ApiResponseV2 struct {
	// 定义响应结构体
	Code    int    `json:"code"`
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

func (r *ApiResponseV2) Success() bool {
	return r.Code == 200
}

type ApiResponseV1 struct {
	// 定义响应结构体
	Code    string `json:"err-code"`
	Status  string `json:"status"`
	Message string `json:"err-msg"`
}

type ApiResponseHBDM struct {
	// 定义响应结构体
	Code    int    `json:"err_code"`
	Status  string `json:"status"`
	Message string `json:"err_msg"`
}
type ApiResponseHBDMV3 struct {
	// 定义响应结构体
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (r *ApiResponseHBDMV3) Success() bool {
	return r.Code == 200
}
func (r *ApiResponseHBDM) Success() bool {
	return r.Status == "ok"
}

func (r *ApiResponseV1) Success() bool {
	return r.Status == "ok"
}

type ApiResponseIntData struct {
	ApiResponseV2
	Data int `json:"data"`
}

type SpotBalanceTicker struct {
	Symbol string      `json:"symbol"`
	Trade  json.Number `json:"trade"`
	Frozen json.Number `json:"frozen"`
}
