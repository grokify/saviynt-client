package saviynt

import (
	"net/http"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
)

type UsersService struct {
	client *Client
}

func NewUsersService(client *Client) *UsersService {
	return &UsersService{client: client}
}

func (svc *UsersService) UpdateUsers(matchField, matchValue string, attrs map[string]any) (*httpsimple.Request, *http.Response, error) {
	if attrs == nil {
		attrs = map[string]any{}
	}
	matchField = strings.TrimSpace(matchField)
	if matchField == "username" {
		attrs[matchField] = matchValue
	} else if matchField != "" {
		attrs[matchField] = matchValue
	}

	sreq := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      svc.client.BuildURL(RelURLUserUpdate),
		BodyType: httpsimple.BodyTypeJSON,
		Body:     attrs,
	}
	resp, err := svc.client.SimpleClient.Do(sreq)
	return &sreq, resp, err
}
