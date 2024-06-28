package restconnector

import (
	"net/http"
	"strings"

	"github.com/grokify/mogo/mime/mimeutil"
	"github.com/grokify/mogo/net/http/httputilmore"
)

type CallInfo struct {
	AccountIDPath          string            `json:"accountIdPath,omitempty"`
	ResponseColsToPropsMap map[string]string `json:"responseColsToPropsMap,omitempty"`
	Call                   []Call            `json:"call,omitempty"`
}

func (ci CallInfo) CallBodies() []string {
	var bodies []string
	for _, c := range ci.Call {
		if body := strings.TrimSpace(c.HTTPParams); body != "" {
			bodies = append(bodies, body)
		}
	}
	return bodies
}

type Call struct {
	Connection       string            `json:"connection"`
	HTTPContentType  string            `json:"httpContentType"`
	HTTPHeaders      map[string]string `json:"httpHeaders"`
	HTTPMethod       string            `json:"httpMethod"`
	HTTPParams       string            `json:"httpParams"`
	Name             string            `json:"name"`
	URL              string            `json:"url"`
	SuccessResponses Responses         `json:"successResponses"`
}

func (c Call) Header() http.Header {
	return httputilmore.NewHeadersMSS(c.HTTPHeaders)
}

func (c Call) IsJSON() bool {
	return mimeutil.IsType(httputilmore.ContentTypeAppJSON, c.HTTPContentType)
}

type Responses struct {
	StatusCodes []uint `json:"statusCode"`
}
