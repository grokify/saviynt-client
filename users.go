package saviynt

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/grokify/goauth/scim"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/strconv/strconvutil"
)

type GetUserResponse struct {
	Msg          string      `json:"msg"`
	DisplayCount string      `json:"displaycount"`
	TotalCount   string      `json:"totalcount"`
	ErrorCode    string      `json:"errorCode"`
	UserDetails  UserDetails `json:"userdetails"`
}

type UserDetails struct {
	City      string `json:"city"`
	Email     string `json:"email"`
	Enabled   string `json:"enabled"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	State     string `json:"state"`
	Username  string `json:"username"`
}

func (ud *UserDetails) SCIMUser() (scim.User, error) {
	usr := scim.User{
		Active: strconvutil.MustParseBool(ud.Enabled),
		Name: scim.Name{
			GivenName:  strings.TrimSpace(ud.FirstName),
			FamilyName: strings.TrimSpace(ud.LastName),
		},
		UserName: strings.TrimSpace(ud.Username),
	}
	ud.City = strings.TrimSpace(ud.City)
	ud.State = strings.TrimSpace(ud.State)
	if ud.City != "" || ud.State != "" {
		usr.Addresses = []scim.Address{
			{
				Locality: ud.City,
				Region:   ud.State},
		}
	}
	if em := strings.TrimSpace(ud.Email); em != "" {
		if err := usr.AddEmail(em, true); err != nil {
			return usr, err
		}
	}
	return usr, nil
}

func (svc *UsersService) GetUserByUsername(username string) (*GetUserResponse, []byte, *http.Response, error) {
	if svc.client.SimpleClient == nil {
		return nil, []byte{}, nil, errors.New("simple client cannot be nil")
	}
	sreq := httpsimple.Request{
		Method: http.MethodPost,
		// URL:      urlutil.JoinAbsolute(c.BaseURL, RelURLECM, RelURLAPI, "getUser"),
		URL:      svc.client.BuildURL(RelURLUserGet),
		BodyType: httpsimple.BodyTypeJSON,
		Body: map[string]string{
			"username": username,
		},
	}
	if resp, err := svc.client.SimpleClient.Do(sreq); err != nil {
		return nil, []byte{}, resp, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, b, resp, err
	} else {
		apiResp := &GetUserResponse{}
		return apiResp, b, resp, json.Unmarshal(b, apiResp)
	}
}
