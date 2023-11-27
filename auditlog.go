package saviynt

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

const AttrTimeFrame = "timeFrame" // from docs

// ExportColumns provides the column names in a standard audit log CSV or Excel export.
func AuditLogUIExportColumns() []string {
	return []string{
		"OBJECTTYPE",
		"OBJECTNAME",
		"ACTION",
		"ACCESSBY",
		"ACCESSTIME",
		"IPADDRESS",
		"ATTRIBUTENAME",
		"OLDVALUE",
		"NEWVALUE",
		"EVENTID",
		"MESSAGE",
	}
}

func AuditLogJobExportColumns() []string {
	return []string{
		"LOGINKEY",
		"LOGINTIME",
		"LOGOUTDATE",
		"COMMENTS",
		"OBJECTTYPE",
		"OBJECTNAME",
		"ACTION",
		"ACCESSBY",
		"IPADDRESS",
		"ATTRIBUTE",
		"OLDVALUE",
		"NEWVALUE",
		"EVENTID",
		"DETAIL",
		"ACCESS_URL",
		"EVENT_TIME",
		"QUERY_PARAM",
	}
}

func auditLogJobExportColumnsSQL() []string {
	return []string{
		"ua.LOGINKEY",
		"l.LOGINTIME",
		"l.LOGOUTDATE",
		"l.COMMENTS AS LOGIN_COMMENTS",
		"ua.TYPEOFACCESS AS OBJECTTYPE",
		"ua.OBJECTKEY AS OBJECTNAME",
		"ua.ActionType AS ACTION",
		"u.username AS ACCESSBY",
		"ua.IPADDRESS",
		"ua.OBJECT_ATTRIBUTE_NAME AS ATTRIBUTE",
		"ua.OLD_VALUE AS OLDVALUE",
		"ua.NEW_VALUE AS NEWVALUE",
		"ua.EVENT_ID AS EVENTID",
		"ua.DETAIL",
		"ua.ACCESS_URL",
		"ua.ACCESSTIME AS EVENT_TIME",
		"ua.QUERY_PARAM",
	}
}

// AnalyticsSQLSIEM represents the SIEM Integration query listed here: https://docs.saviyntcloud.com/bundle/EIC-Admin-v23x/page/Content/Chapter20-EIC-Integrations/Saviynt-SIEM-Integration.htm#Step
func AnalyticsSQLAuditLogSIEM() string {
	return `select ua.TYPEOFACCESS as 'Object Type',ua.ActionType as 'Action Taken',u.username as 'Accessed By', ua.IPADDRESS as 'IP Address',ua.ACCESSTIME as 'Event Time',ua.DETAIL as 'Message' from users u, userlogin_access ua, userlogins l where l.loginkey = ua.LOGINKEY and l.USERKEY = u.userkey and ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) and ua.Detail is not NULL`
}

// AuditLogSQLQueryUI represents a SQL query that very closely matches the CSV / XLSX download from the Audit Log UI. The primary difference is
// that the `MESSAGE` column value is wrapped in the API response.
func AnalyticsSQLAuditLogUI() string {
	return `SELECT ua.TYPEOFACCESS AS OBJECTTYPE, ua.OBJECTKEY AS OBJECTNAME, ua.ActionType AS ACTION, u.username AS ACCESSBY, ua.ACCESSTIME, ua.IPADDRESS, ua.OBJECT_ATTRIBUTE_NAME AS ATTRIBUTENAME, ua.OLD_VALUE AS OLDVALUE, ua.NEW_VALUE AS NEWVALUE, ua.EVENT_ID AS EVENTID, ua.DETAIL AS MESSAGE FROM users u, userlogin_access ua, userlogins l WHERE l.loginkey = ua.LOGINKEY AND l.USERKEY = u.userkey AND ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) AND ua.Detail is not NULL`
}

// AnalyticsSQLAuditLogJob provides a SQL query which returns the output of the EIC Archival job.
func AnalyticsSQLAuditLogJob() string {
	return "SELECT " + strings.Join(auditLogJobExportColumnsSQL(), ", ") + ` FROM users u, userlogin_access ua, userlogins l WHERE l.loginkey = ua.LOGINKEY AND l.USERKEY = u.userkey AND ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) AND ua.Detail is not NULL`
	// SELECT ua.LOGINKEY, l.LOGINTIME, l.LOGOUTDATE, l.COMMENTS AS LOGIN_COMMENTS, ua.TYPEOFACCESS AS OBJECTTYPE, ua.OBJECTKEY AS OBJECTNAME, ua.ActionType AS ACTION, u.username AS ACCESSBY, ua.IPADDRESS, ua.OBJECT_ATTRIBUTE_NAME AS ATTRIBUTE, ua.OLD_VALUE AS OLDVALUE, ua.NEW_VALUE AS NEWVALUE, ua.EVENT_ID AS EVENTID, ua.DETAIL, ua.ACCESS_URL, ua.ACCESSTIME AS EVENT_TIME, ua.QUERY_PARAM FROM users u, userlogin_access ua, userlogins l WHERE l.loginkey = ua.LOGINKEY AND l.USERKEY = u.userkey AND ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) AND ua.Detail is not NULL
}

func (c Client) FetchRuntimeControlsDataV2(name string, attrs map[string]any, limit, offset uint) (*http.Response, error) {
	// func (c Client) GetAuditLogRuntimeControlsData(name string, minutes, limit, offset uint) (*http.Response, error) {
	if limit == 0 {
		limit = 50
	}
	sreq := httpsimple.Request{
		Method:   http.MethodPost,
		URL:      urlutil.JoinAbsolute(c.BaseURL, RelURLECM, RelURLAPI, RelURLLoginRuntimeControlsData),
		BodyType: httpsimple.BodyTypeJSON,
		Body: AnalyticsRequest{
			AnalyticsName: name,
			Attributes:    attrs,
			Max:           strconv.Itoa(int(limit)),
			Offset:        strconv.Itoa(int(offset)),
		},
	}
	sclient := httpsimple.Client{
		BaseURL:    c.BaseURL,
		HTTPClient: c.HTTPClient}
	return sclient.Do(sreq)
}

type AnalyticsRequest struct {
	AnalyticsName string `json:"analyticsname"`
	Attributes    any    `json:"attributes,omitempty"`
	Max           string `json:"max,omitempty"`
	Offset        string `json:"offset,omitempty"`
}

type AnalyticsRequestAttributes struct {
	TimeFrame string `json:"timeFrame"`
}

/*
<html>
<head><title>404 Not Found</title></head>
<body>
<center><h1>404 Not Found</h1></center>
<hr><center>nginx</center>
</body>
</html>
*/
