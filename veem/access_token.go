package veem

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	UserID      string `json:"user_id"`
	AccountID   string `json:"account_id"`
	Username    string `json:"user_name"`

	// Calculated at retrieval
	ExpiresAt time.Time
}

func (c *client) getAccessToken() (*AccessTokenResponse, error) {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	form.Add("scope", "all")
	req, err := c.newRequest(http.MethodPost, "oauth/token", strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprintf(
		"Basic %s", base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.opts.ClientID, c.opts.ClientSecret))),
	))
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	res := &AccessTokenResponse{}
	if err := c.doInto(req, res); err != nil {
		return nil, err
	}
	res.ExpiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Millisecond)
	return res, nil
}
