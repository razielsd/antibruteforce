package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *clientAPI) GetWhitelist() ([]string, error) {
	req, err := createGetRequest(c.makeURL("/api/whitelist"))
	if err != nil {
		return nil, err
	}
	return c.GetBWlist(req)
}

func (c *clientAPI) GetBlacklist() ([]string, error) {
	req, err := createGetRequest(c.makeURL("/api/blacklist"))
	if err != nil {
		return nil, err
	}
	return c.GetBWlist(req)
}

func (c *clientAPI) GetBWlist(req *http.Request) ([]string, error) {
	client := c.getClient()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	if err != nil {
		return nil, err
	}
	ips := struct {
		Result []string `json:"result"`
	}{}
	err = json.Unmarshal(body, &ips)
	if err != nil {
		return nil, fmt.Errorf("%w: unable parse response", err)
	}
	return ips.Result, nil
}

func (c *clientAPI) appendBlacklist(clientIP string) error {
	return c.sendUpdate("/api/blacklist/add", "ip", clientIP)
}

func (c *clientAPI) appendWhitelist(clientIP string) error {
	return c.sendUpdate("/api/whitelist/add", "ip", clientIP)
}

func (c *clientAPI) removeBlacklist(clientIP string) error {
	return c.sendUpdate("/api/blacklist/remove", "ip", clientIP)
}

func (c *clientAPI) removeWhitelist(clientIP string) error {
	return c.sendUpdate("/api/whitelist/remove", "ip", clientIP)
}

func (c *clientAPI) dropBucketByLogin(key string) error {
	return c.sendUpdate("/api/bucket/drop/login", "key", key)
}

func (c *clientAPI) dropBucketByPasswd(key string) error {
	return c.sendUpdate("/api/bucket/drop/pwd", "key", key)
}

func (c *clientAPI) dropBucketByIP(key string) error {
	return c.sendUpdate("/api/bucket/drop/ip", "key", key)
}

func (c *clientAPI) sendUpdate(path, paramName, value string) error {
	client := c.getClient()
	req, err := c.createPostRequest(
		c.makeURL(path),
		map[string]string{paramName: value},
	)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	errResp, err := c.extractError(resp)
	if err != nil {
		return err
	}
	return errors.New(errResp.ErrMessage)
}
