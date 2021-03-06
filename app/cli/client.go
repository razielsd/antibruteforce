package cli

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/razielsd/antibruteforce/app/api"
)

type clientAPI struct {
	Host       string
	HTTPClient httpClient
}

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

func newClientAPI(apiHost string) *clientAPI {
	client := &clientAPI{
		Host: apiHost,
	}
	client.HTTPClient = client.createDefaultHTTPClient()
	return client
}

func (c *clientAPI) createDefaultHTTPClient() httpClient {
	z := &http.Client{
		Transport: &http.Transport{
			MaxIdleConns:    100,
			IdleConnTimeout: 90 * time.Second,
		},
	}
	return z
}

func (c *clientAPI) getHTTPClient() httpClient {
	return c.HTTPClient
}

func (c *clientAPI) setHTTPClient(hc httpClient) {
	c.HTTPClient = hc
}

func (c *clientAPI) makeURL(path string) string {
	return "http://" + c.Host + path
}

func (c *clientAPI) createPostRequest(reqURL string, params map[string]string) (*http.Request, error) {
	data := url.Values{}
	for k, v := range params {
		data.Set(k, v)
	}
	req, err := http.NewRequestWithContext(
		context.Background(), "POST", reqURL, strings.NewReader(data.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func createGetRequest(reqURL string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "GET", reqURL, nil)
	if err != nil {
		return nil, err
	}
	return req, err
}

func (c *clientAPI) extractError(resp *http.Response) (*api.ErrorResponse, error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	errResp := &api.ErrorResponse{}
	err = json.Unmarshal(body, errResp)
	if err != nil {
		return nil, err
	}

	return errResp, nil
}
