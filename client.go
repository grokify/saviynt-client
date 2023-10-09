package saviynt

import (
	"encoding/json"
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
	BaseURL    string
	Path       string
	HTTPClient *http.Client
}

func NewClient(baseURL, path, username, password string) (Client, error) {
	c := Client{
		BaseURL: baseURL,
		Path:    path,
	}
	tok, err := GetToken(baseURL, username, password)
	if err != nil {
		return c, err
	}
	httpClient := authutil.NewClientTokenOAuth2(tok)
	c.HTTPClient = httpClient
	return c, nil
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (c Client) GetToken(username, password string) (*oauth2.Token, error) {
	return GetToken(c.BaseURL, username, password)
}

func GetToken(baseURL, username, password string) (*oauth2.Token, error) {
	sreq := httpsimple.SimpleRequest{
		URL:      urlutil.JoinAbsolute(baseURL, RelURLLogin),
		Method:   http.MethodPost,
		BodyType: httpsimple.BodyTypeJSON,
		Body: LoginRequest{
			Username: username,
			Password: password,
		},
	}
	resp, err := httpsimple.Do(sreq)
	if err != nil {
		return nil, err
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tok := &oauth2.Token{}
	return tok, json.Unmarshal(b, tok)
}
