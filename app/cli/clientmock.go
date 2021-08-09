package cli

import (
	"net/http"
	"net/http/httptest"
)

type apiClientMock struct {
	statusCode int
	body       string
	reqError   error
}

func newAPIClientMock(code int, body string, err error) *apiClientMock {
	return &apiClientMock{
		statusCode: code,
		body:       body,
		reqError:   err,
	}
}

func (m *apiClientMock) Do(req *http.Request) (*http.Response, error) {
	resp := httptest.NewRecorder()
	resp.WriteHeader(m.statusCode)
	_, err := resp.WriteString(m.body)
	if err != nil {
		return nil, err
	}
	if m.reqError != nil {
		return resp.Result(), m.reqError
	}
	return resp.Result(), nil
}
