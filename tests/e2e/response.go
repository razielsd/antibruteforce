// +build e2e

package e2e

import "encoding/json"

type SuccessResponse struct {
	Result interface{} `json:"result"`
	Code   int         `json:"code"`
}

type ErrorResponse struct {
	ErrMessage string `json:"error"`
	Code       int    `json:"code"`
}

type SuccessOK struct {
	Status string `json:"status"`
}

type AllowResult struct {
	CanAuth   bool `json:"can-auth"`
	Login     bool `json:"login-check"`
	Pwd       bool `json:"pwd-check"`
	IP        bool `json:"ip-check"`
	Whitelist bool `json:"whitelist-check"`
	Blacklist bool `json:"blacklist-check"`
}

func NewSuccessOK() SuccessOK {
	return SuccessOK{
		Status: "OK",
	}
}

func (s SuccessResponse) JSON() string {
	js, _ := json.Marshal(s)
	return string(js)
}

func (s ErrorResponse) JSON() string {
	js, _ := json.Marshal(s)
	return string(js)
}
