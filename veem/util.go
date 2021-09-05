package veem

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func mustParseURL(u string) *url.URL {
	p, err := url.Parse(u)
	if err != nil {
		panic(err)
	}
	return p
}

func parseAPIError(body []byte) error {
	err := &APIError{}
	if merr := json.Unmarshal(body, err); merr != nil {
		return errors.New(string(body))
	}
	return err
}

func (c *client) newRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", c.apiURL.String(), endpoint)
	return http.NewRequest(method, url, body)
}

func (c *client) doWithAuth(req *http.Request, acceptType string) (io.ReadCloser, error) {
	if time.Now().Add(-time.Minute).After(c.token.ExpiresAt) {
		var err error
		c.token, err = c.getAccessToken()
		if err != nil {
			return nil, err
		}
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", strings.ToTitle(c.token.TokenType), c.token.AccessToken))
	if acceptType == "" {
		acceptType = "application/json"
	}
	req.Header.Add("Accept", acceptType)
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	return c.do(req)
}

func (c *client) doIntoWithAuth(req *http.Request, out interface{}) error {
	if time.Now().Add(-time.Minute).After(c.token.ExpiresAt) {
		var err error
		c.token, err = c.getAccessToken()
		if err != nil {
			return err
		}
	}
	if req.Header.Get("Content-Type") == "" {
		req.Header.Add("Content-Type", "application/json")
	}
	req.Header.Add("Authorization", fmt.Sprintf("%s %s", strings.ToTitle(c.token.TokenType), c.token.AccessToken))
	req.Header.Add("Accept", "application/json")
	return c.doInto(req, out)
}

func (c *client) do(req *http.Request) (io.ReadCloser, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, parseAPIError(body)
	}
	return res.Body, nil
}

func (c *client) doInto(req *http.Request, out interface{}) error {
	res, err := c.do(req)
	if err != nil {
		return err
	}
	defer res.Close()
	body, err := ioutil.ReadAll(res)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, out)
}
