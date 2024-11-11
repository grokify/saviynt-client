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
	if svc.client == nil {
		return nil, ErrClientNotSet
	}
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
	// Username specifies the name of the user for whom you need to change the password.
	Username string `json:"username"`

	// Password specifies the new password.
	Password string `json:"password"`

	// ChangePasswordAssociatedAccounts specifies whether to create the corresponding Change Password task, if `true`. You can specify `false` to only update the user password. The default value is `true`.
	ChangePasswordAssociatedAccounts bool `json:"changePasswordAssociatedAccounts"`

	// Endpointspecifies a list of endpoints (comma-separated) to update the password for
	// if `ChangePasswordAssociatedAccounts` is set to `true`.
	Endpoint []string `json:"endpoint"`

	// SetARSTaskSourcespecifies whether to set the source column in the `arstasks`
	// table with the `changeOwnPasswordFromAPI` value. When the source column in `arstasks`
	// table is set to the `changeOwnPasswordFromAPI` value then the `pwdLastSet` parameter is
	// not set to `0`` for Active Directory (AD). The default value is `false`.  Note: If your
	// AD password has expired (pwdLastSet = 0), you are forced to choose a new password in EIC
	// on the next login.
	SetARSTaskSource bool `json:"setarstasksource"`

	// UpdateUserPassword specifies whether to update the user password when
	// `ChangePasswordAssociatedAccounts` is set to `true`. Setting this to `true` specifies
	// when `ChangePasswordAssociatedAccounts` is `true` updates the user password and creates
	// the corresponding Change Password task. When set to `false``, only the Change Password
	// task is created. The default value is true.
	UpdateUserPassword bool `json:"updateUserPassword"`

	// ValidateAgainstPolicy specifies whether the new password conforms to the USER scope
	// password policy. If you do not want to apply the password policy, then specify `N`.
	// The default value is `Y`.
	ValidateAgainstPolicy bool `json:"validateagainstpolicy"`
}

func (opts ChangePasswordOpts) Values() url.Values {
	v := url.Values{}
	v.Add("username", opts.Username)
	v.Add("password", opts.Password)
	v.Add("changePasswordAssociatedAccounts", strconvutil.Btoa(opts.ChangePasswordAssociatedAccounts))
	v.Add("updateUserPassword", strconvutil.Btoa(opts.UpdateUserPassword))
	v.Add("setarstasksource", strconvutil.Btoa(opts.SetARSTaskSource))
	if opts.ValidateAgainstPolicy {
		v.Add("validateagainstpolicy", "Y")
	} else {
		v.Add("validateagainstpolicy", "N")
	}
	eps := stringsutil.SliceCondenseSpace(opts.Endpoint, true, false)
	if len(eps) > 0 {
		v.Add("endpoints", strings.Join(eps, ","))
	}
	return v
}
