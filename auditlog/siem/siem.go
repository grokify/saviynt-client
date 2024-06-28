package siem

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

type SIEMAuditResponse struct {
	DisplayCount int              `json:"displaycount"`
	TotalCount   int              `json:"totalcount"`
	Message      string           `json:"msg"`
	Results      []SIEMAuditEvent `json:"results"`
	ErrorCode    string           `json:"errorcode"`
}

func ParseSIEMAuditResponse(r io.Reader) (*SIEMAuditResponse, error) {
	resp := &SIEMAuditResponse{}
	if b, err := io.ReadAll(r); err != nil {
		return nil, err
	} else {
		return resp, json.Unmarshal(b, resp)
	}
}

type SIEMAuditEvent struct {
	ActionTaken string `json:"Action Taken"`
	IPAddress   string `json:"IP Address"`
	EventTime   string `json:"Event Time"`
	// EventTime   SQLTimestampReader `json:"Event Time"`
	Message    string `json:"Message"`
	ObjectType string `json:"Object Type"`
	AccessedBy string `json:"Accessed By"`
}

type SQLTimestampReader struct {
	time.Time
}

func (t *SQLTimestampReader) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" {
		t.Time = time.Time{}
		return
	}
	t.Time, err = time.Parse(timeutil.ISO9075, s)
	return
}
