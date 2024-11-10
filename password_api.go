package saviynt

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type PasswordService struct {
	client *Client
}

func NewPasswordService(client *Client) *PasswordService {
	return &PasswordService{client: client}
}

func (svc *PasswordService) ChangePassword(opts ChangePasswordOpts) (*http.Response, error) {
	return svc.client.SimpleClient.Do(
		httpsimple.Request{
			Method:   http.MethodPost,
			URL:      svc.client.BuildURL(RelURLPasswordChange),
			Body:     opts.Values(),
			BodyType: httpsimple.BodyTypeForm,
		},
	)
}

type ChangePasswordOpts struct {
	Username                         string   `json:"username"`
	Password                         string   `json:"password"`
	ChangePasswordAssociatedAccounts bool     `json:"changePasswordAssociatedAccounts"`
	ChangeUserPassword               bool     `json:"changeUserPassword"`
	Endpoints                        []string `json:"endpoints"`
	SetARSTaskSource                 bool     `json:"setarstasksource"`
	ValidateAgainstPolicy            bool     `json:"validateagainstpolicy"`
}

func (opts ChangePasswordOpts) Values() url.Values {
	v := url.Values{}
	v.Add("username", opts.Username)
	v.Add("password", opts.Password)
	v.Add("changePasswordAssociatedAccounts", strconvutil.Btoa(opts.ChangePasswordAssociatedAccounts))
	v.Add("changeUserPassword", strconvutil.Btoa(opts.ChangeUserPassword))
	v.Add("setarstasksource", strconvutil.Btoa(opts.SetARSTaskSource))
	if opts.ValidateAgainstPolicy {
		v.Add("validateagainstpolicy", "Y")
	} else {
		v.Add("validateagainstpolicy", "N")
	}
	eps := stringsutil.SliceCondenseSpace(opts.Endpoints, true, false)
	if len(eps) > 0 {
		v.Add("endpoints", strings.Join(eps, ","))
	}
	return v
}
