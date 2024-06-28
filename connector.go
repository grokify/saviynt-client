package saviynt

import (
	"net/url"
	"strings"
)

type ConnectionSaveTestOpts struct {
	ConnectionType string
	SaveConnection bool
	ConnectionName string
	Params         Values
}

func (opts ConnectionSaveTestOpts) Values() url.Values {
	v := opts.Params.Values()
	if v == nil {
		v = url.Values{}
	}
	if opts.ConnectionType != "" {
		v.Add("connectiontype", opts.ConnectionType)
	}
	if opts.ConnectionName != "" {
		v.Add("connectionName", opts.ConnectionName)
	}
	if opts.SaveConnection {
		v.Add("saveconnection", "Y")
	}
	return v
}

func (c Client) ConnectionSaveTest() {

}

type ConnectionParamsRESTOpts struct {
	ConnectionJSON string
	ImportUserJSON string
}

func (opts *ConnectionParamsRESTOpts) TrimSpace() {
	opts.ConnectionJSON = strings.TrimSpace(opts.ConnectionJSON)
}

type Values interface {
	Values() url.Values
}

func (opts *ConnectionParamsRESTOpts) Values() url.Values {
	v := url.Values{}
	if opts.ConnectionJSON != "" {
		v.Add("ConnectionJSON", opts.ConnectionJSON)
	}

	return v
}

/*
REST	ConnectionJSON,ImportUserJSON,ImportAccountEntJSON,STATUS_THRESHOLD_CONFIG,CreateAccountJSON,UpdateAccountJSON,EnableAccountJSON,DisableAccountJSON,AddAccessJSON,RemoveAccessJSON,UpdateUserJSON,ChangePassJSON,RemoveAccountJSON,TicketStatusJSON,CreateTicketJSON,ENDPOINTS_FILTER,PasswdPolicyJSON,ConfigJSON,AddFFIDAccessJSON,RemoveFFIDAccessJSON,MODIFYUSERDATAJSON,SendOtpJSON,ValidateOtpJSON,PAM_CONFIG
*/
