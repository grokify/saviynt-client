package saviynt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/goauth/authutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/http/httputilmore"
	"github.com/grokify/mogo/net/urlutil"
	"golang.org/x/oauth2"
)

const (
	RelURLLogin                    = "/ECM/api/login"
	RelOAuthAccessToken            = "/ECM/oauth/access_token"
	RelURLECM                      = "/ECM"
	RelURLAPI                      = "/api/v5"
	RelURLLoginRuntimeControlsData = "/fetchRuntimeControlsDataV2" // API at https://documenter.getpostman.com/view/23973797/2s9XxwutWR#b821cc21-ee7c-49e3-9433-989ed87b2b03
)

type Client struct {
	BaseURL      string
	Path         string
	HTTPClient   *http.Client
	SimpleClient *httpsimple.Client
	Token        *oauth2.Token
}

func NewClient(ctx context.Context, baseURL, path, username, password string) (Client, error) {
	c := Client{
		BaseURL: baseURL,
		Path:    path,
	}
	tok, err := GetTokenPassword(ctx, baseURL, username, password, false)
	if err != nil {
		return c, err
	}
	c.Token = tok
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

func (c Client) GetTokenPassword(ctx context.Context, username, password string) (*oauth2.Token, error) {
	return GetTokenPassword(ctx, c.BaseURL, username, password, true)
}

func GetTokenPassword(ctx context.Context, baseURL, username, password string, useBasicAuth bool) (*oauth2.Token, error) {
	var sreq httpsimple.Request
	if useBasicAuth {
		hval, err := authutil.BasicAuthHeader(username, password)
		if err != nil {
			return nil, err
		}
		sreq = httpsimple.Request{
			URL:    urlutil.JoinAbsolute(baseURL, RelURLLogin),
			Method: http.MethodPost,
			Headers: http.Header{
				httputilmore.HeaderAuthorization: []string{hval},
			},
		}
	} else {
		sreq = httpsimple.Request{
			URL:      urlutil.JoinAbsolute(baseURL, RelURLLogin),
			Method:   http.MethodPost,
			BodyType: httpsimple.BodyTypeJSON,
			Body: LoginRequest{
				Username: username,
				Password: password,
			},
		}
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

func (c Client) GetTokenRefresh(ctx context.Context) (*oauth2.Token, error) {
	if c.Token == nil {
		return nil, errors.New("oauth2.Token cannot be nil")
	} else if strings.TrimSpace(c.Token.AccessToken) == "" {
		return nil, errors.New("oauth2.Token.AccessToken cannot be empty")
	}
	return GetTokenRefresh(ctx, c.BaseURL, c.Token.AccessToken)
}

func GetTokenRefresh(ctx context.Context, baseURL, refreshToken string) (*oauth2.Token, error) {
	sreq := httpsimple.Request{
		URL:      urlutil.JoinAbsolute(baseURL, RelOAuthAccessToken),
		Method:   http.MethodPost,
		BodyType: httpsimple.BodyTypeForm,
		Body: url.Values{
			"grant_type":    []string{"refresh_token"},
			"refresh_token": []string{refreshToken}},
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
