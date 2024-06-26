package saviynt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/oauth2"
)

const (
	RelURLLogin                    = "/ECM/api/login"
	RelURLECM                      = "/ECM"
	RelURLAPI                      = "/api/v5"
	RelURLLoginRuntimeControlsData = "/fetchRuntimeControlsDataV2" // API at https://documenter.getpostman.com/view/23973797/2s9XxwutWR#b821cc21-ee7c-49e3-9433-989ed87b2b03
)

type Client struct {
	BaseURL      string
	Path         string
	HTTPClient   *http.Client
	SimpleClient *httpsimple.Client
}

func NewClient(ctx context.Context, baseURL, path, username, password string) (Client, error) {
	c := Client{
		BaseURL: baseURL,
		Path:    path,
	}
	tok, err := GetToken(ctx, baseURL, username, password)
	if err != nil {
		return c, err
	}
	httpClient := authutil.NewClientTokenOAuth2(tok)
	c.HTTPClient = httpClient
	simClient := httpsimple.NewClient(httpClient, baseURL)
	c.SimpleClient = &simClient
	return c, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c Client) GetToken(ctx context.Context, username, password string) (*oauth2.Token, error) {
	return GetToken(ctx, c.BaseURL, username, password)
}

func GetToken(ctx context.Context, baseURL, username, password string) (*oauth2.Token, error) {
	sreq := httpsimple.Request{
		URL:      urlutil.JoinAbsolute(baseURL, RelURLLogin),
		Method:   http.MethodPost,
		BodyType: httpsimple.BodyTypeJSON,
		Body: LoginRequest{
			Username: username,
			Password: password,
		},
	}
	if resp, err := httpsimple.Do(ctx, sreq); err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("login api status code is (%d)", resp.StatusCode)
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		tok := &oauth2.Token{}
		return tok, json.Unmarshal(b, tok)
	}
}

func (c Client) GetUserByUsername(username string) (*GetUserResponse, []byte, *http.Response, error) {
	if c.SimpleClient == nil {
		return nil, []byte{}, nil, errors.New("simple client cannot be nil")
	}
	sreq := httpsimple.Request{
		URL:      urlutil.JoinAbsolute(c.BaseURL, RelURLECM, RelURLAPI, "getUser"),
		Method:   http.MethodPost,
		BodyType: httpsimple.BodyTypeJSON,
		Body: map[string]string{
			"username": username,
		},
	}
	if resp, err := c.SimpleClient.Do(sreq); err != nil {
		return nil, []byte{}, resp, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, b, resp, err
	} else {
		apiResp := &GetUserResponse{}
		return apiResp, b, resp, json.Unmarshal(b, apiResp)
	}
}
