package auditlog

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
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
		"ua.LOGINKEY AS LOGIN_KEY",
		"l.LOGINTIME AS LOGIN_TIME",
		"l.LOGOUTDATE AS LOGOUT_TIME",
		"l.COMMENTS AS LOGIN_COMMENTS",
		"ua.TYPEOFACCESS AS OBJECT_TYPE",
		"ua.OBJECTKEY AS OBJECT_KEY",
		"ua.ActionType AS ACTION",
		"u.username AS ACCESS_BY",
		"ua.IPADDRESS AS IP_ADDRESS",
		"ua.OBJECT_ATTRIBUTE_NAME AS OBJECT_ATTRIBUTE_NAME",
		"ua.OLD_VALUE AS OLD_VALUE",
		"ua.NEW_VALUE AS NEW_VALUE",
		"ua.ACCESSTIME AS EVENT_TIME",
		"ua.EVENT_ID AS EVENT_ID",
		"ua.DETAIL AS DETAIL", //  ua.DETAIL as 'Message'
		"ua.ACCESS_URL AS ACCESS_URL",
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

// AnalyticsSQLAuditLogArchival provides a SQL query which returns the output of the EIC Archival job.
func AnalyticsSQLAuditLogArchival() string {
	return "SELECT " + strings.Join(auditLogJobExportColumnsSQL(), ", ") + ` FROM users u, userlogin_access ua, userlogins l WHERE l.loginkey = ua.LOGINKEY AND l.USERKEY = u.userkey AND ua.AccessTime >= (NOW() - INTERVAL ${fromMinutesAgo} Minute) AND ua.Detail is not NULL`
	// SELECT ua.LOGINKEY, l.LOGINTIME, l.LOGOUTDATE, l.COMMENTS AS LOGIN_COMMENTS, ua.TYPEOFACCESS AS OBJECTTYPE, ua.OBJECTKEY AS OBJECTNAME, ua.ActionType AS ACTION, u.username AS ACCESSBY, ua.IPADDRESS, ua.OBJECT_ATTRIBUTE_NAME AS ATTRIBUTE, ua.OLD_VALUE AS OLDVALUE, ua.NEW_VALUE AS NEWVALUE, ua.EVENT_ID AS EVENTID, ua.DETAIL, ua.ACCESS_URL, ua.ACCESSTIME AS EVENT_TIME, ua.QUERY_PARAM FROM users u, userlogin_access ua, userlogins l WHERE l.loginkey = ua.LOGINKEY AND l.USERKEY = u.userkey AND ua.AccessTime >= (NOW() - INTERVAL ${timeFrame} Minute) AND ua.Detail is not NULL
}

type AuditEventSQL struct {
	AccessBy            string `json:"ACCESS_BY,omitempty"`
	AccessURL           string `json:"ACCESS_URL,omitempty"`
	Action              string `json:"ACTION,omitempty"`
	Detail              string `json:"DETAIL,omitempty"`
	EventID             string `json:"EVENT_ID,omitempty"`
	EventTime           string `json:"EVENT_TIME,omitempty"`
	IPAddress           string `json:"IP_ADDRESS,omitempty"`
	LoginComments       string `json:"LOGIN_COMMENTS,omitempty"`
	LoginKey            string `json:"LOGIN_KEY,omitempty"`
	LoginTime           string `json:"LOGIN_TIME,omitempty"`
	LogoutTime          string `json:"LOGOUT_TIME,omitempty"`
	NewValue            string `json:"NEW_VALUE,omitempty"`
	ObjectAttributeName string `json:"OBJECT_ATTRIBUTE_NAME,omitempty"`
	ObjectKey           string `json:"OBJECT_KEY,omitempty"`
	ObjectType          string `json:"OBJECT_TYPE,omitempty"`
	OldValue            string `json:"OLD_VALUE,omitempty"`
	QueryParam          string `json:"QUERY_PARAM,omitempty"`
}

func (s AuditEventSQL) Event() (*AuditEvent, error) {
	evt := &AuditEvent{
		AccessBy:            s.AccessBy,
		AccessURL:           s.AccessURL,
		Action:              s.Action,
		Detail:              s.Detail,
		EventID:             s.EventID,
		IPAddress:           s.IPAddress,
		LoginComments:       s.LoginComments,
		LoginKey:            s.LoginKey,
		NewValue:            s.NewValue,
		ObjectAttributeName: s.ObjectAttributeName,
		ObjectKey:           s.ObjectKey,
		ObjectType:          s.ObjectType,
		OldValue:            s.OldValue,
		QueryParam:          s.QueryParam,
	}
	if strings.TrimSpace(s.LoginTime) != "" {
		if loginTime, err := time.Parse(timeutil.SQLTimestamp, s.LoginTime); err != nil {
			return evt, err
		} else {
			evt.LoginTime = &loginTime
		}
	}
	if strings.TrimSpace(s.LogoutTime) != "" {
		if logoutTime, err := time.Parse(timeutil.SQLTimestamp, s.LogoutTime); err != nil {
			return evt, err
		} else {
			evt.LogoutTime = &logoutTime
		}
	}
	if eventTime, err := time.Parse(timeutil.SQLTimestamp, s.EventTime); err != nil {
		return evt, err
	} else {
		evt.EventTime = eventTime
	}
	detail := strings.TrimSpace(s.Detail)
	if strings.Index(detail, "{") == 0 && stringsutil.ReverseIndex(detail, "}") == 0 {
		d := saviynt.UserLoginAccessDetail{}
		if err := json.Unmarshal([]byte(detail), &d); err != nil {
			return evt, err
		} else {
			evt.Message = strings.TrimSpace(d.Message)
			evt.Data = d.Data
			evt.ObjectName = d.ObjectName
		}
	} else {
		evt.Message = strings.TrimSpace(s.Detail)
	}
	return evt, nil
}

func AuditEventsParseMaps(m []map[string]string) (AuditEvents, error) {
	var evts []AuditEvent
	for _, mi := range m {
		evt, err := AuditEventParseMap(mi)
		if err != nil {
			return evts, err
		}
		evts = append(evts, *evt)
	}
	return evts, nil
}

func AuditEventParseMap(m map[string]string) (*AuditEvent, error) {
	j, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	evtSQL := AuditEventSQL{}
	err = json.Unmarshal(j, &evtSQL)
	if err != nil {
		return nil, err
	}
	return evtSQL.Event()
}

type AuditEvents []AuditEvent

func (e AuditEvents) EventTimes() timeutil.Times {
	var times timeutil.Times
	for _, ei := range e {
		times = append(times, ei.EventTime)
	}
	return times
}

type AuditEvent struct {
	AccessBy            string     `json:"accessBy,omitempty"`
	AccessURL           string     `json:"accessURL,omitempty"`
	Action              string     `json:"action,omitempty"`
	Data                string     `json:"data,omitempty"`
	Detail              string     `json:"detail,omitempty"`
	EventID             string     `json:"eventID,omitempty"`
	EventTime           time.Time  `json:"eventTime,omitempty"`
	IPAddress           string     `json:"ipAddress,omitempty"`
	LoginComments       string     `json:"loginComments,omitempty"`
	LoginKey            string     `json:"loginKey,omitempty"`
	LoginTime           *time.Time `json:"loginTime,omitempty"`
	LogoutTime          *time.Time `json:"logoutTime,omitempty"`
	Message             string     `json:"message,omitempty"`
	NewValue            string     `json:"newValue,omitempty"`
	ObjectAttributeName string     `json:"objectAttributeName,omitempty"`
	ObjectKey           string     `json:"objectKey,omitempty"`
	ObjectName          string     `json:"objectName,omitempty"`
	ObjectType          string     `json:"objectType,omitempty"`
	OldValue            string     `json:"oldValue,omitempty"`
	QueryParam          string     `json:"queryParam,omitempty"`
}

type AnalyticsAuditResponse struct {
	DisplayCount int         `json:"displaycount"`
	Msg          string      `json:"msg"`
	TotalCount   int         `json:"totalcount"`
	Results      AuditEvents `json:"results"`
}

// ParseAnalyticsAuditLogArchivalAPIResponse parses an API response, e.g. `*http.Response.Body` that is
// associated with the SQL query defined by `AnalyticsSQLAuditLogArchival()`.`
func ParseAnalyticsAuditLogArchivalAPIResponse(r io.Reader) (*AnalyticsAuditResponse, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	resGeneric := &saviynt.AnalyticsResponse{}
	err = json.Unmarshal(b, resGeneric)
	if err != nil {
		return nil, err
	}
	resAudit := &AnalyticsAuditResponse{
		DisplayCount: resGeneric.DisplayCount,
		Msg:          resGeneric.Msg,
		TotalCount:   resGeneric.TotalCount,
	}
	evts, err := AuditEventsParseMaps(resGeneric.Results)
	if err != nil {
		return nil, err
	}
	resAudit.Results = evts
	return resAudit, nil
}
