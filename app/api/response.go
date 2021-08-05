package api

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
