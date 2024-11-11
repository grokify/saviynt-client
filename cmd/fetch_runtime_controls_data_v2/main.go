package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/go-saviynt/auditlog"
	"github.com/grokify/go-saviynt/auditlog/siem"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/type/maputil"
)

func main() {
	_, err := config.LoadDotEnv([]string{".env"}, 1)
	logutil.FatalErr(err, "load_dot_env")

	clt, err := saviynt.NewClient(
		context.Background(),
		os.Getenv(saviynt.EnvSaviyntServerURL),
		saviynt.RelURLAPI,
		os.Getenv(saviynt.EnvSaviyntUsername),
		os.Getenv(saviynt.EnvSaviyntPassword))
	logutil.FatalErr(err, "new_client")

	flagGetUser := true

	if flagGetUser {
		usr, _, resp, err := clt.UsersAPI.GetUserByUsername(os.Getenv(saviynt.EnvSaviyntUsername))
		logutil.FatalErr(err, "GetUserByUsername")

		b, err := io.ReadAll(resp.Body)
		logutil.FatalErr(err, "ReadAll")
		fmt.Println(string(b))
		fmtutil.MustPrintJSON(usr)
		scimUser, err := usr.UserDetails.SCIMUser()
		logutil.FatalErr(err)
		fmtutil.MustPrintJSON(scimUser)
		//panic("Z")
	}

	attrs := map[string]any{}
	if attrsStr := os.Getenv("SAVIYNT_QUERY_ATTR"); len(attrsStr) > 0 {
		attrsVals, err := url.ParseQuery(attrsStr)
		logutil.FatalErr(err, "parse_query")
		attrsValsMap := maputil.MapStringSlice(attrsVals)
		attrs = attrsValsMap.FlattenAny(false, false)
		// attrs = MS3ToMSA(attrsVals)
	}

	// DATE_FORMAT(ua.ACCESSTIME,  '%Y-%m-%dT%TZ')
	req, resp, err := clt.AnalyticsAPI.FetchRuntimeControlsDataV2(
		os.Getenv("SAVIYNT_QUERY_NAME"),
		"", "",
		attrs,
		9999, 0)
	fmtutil.MustPrintJSON(req)
	logutil.FatalErr(err, "fetchruntimecontrolsdatav2")
	fmt.Printf("status code (%d)\n", resp.StatusCode)
	if resp.StatusCode > 299 {
		fmt.Printf("status code error (%d)\n", resp.StatusCode)
	}

	if 1 == 0 {
		r, err := siem.ParseSIEMAuditResponse(resp.Body)
		logutil.FatalErr(err, "ParseSIEMAuditResponse")
		fmtutil.MustPrintJSON(r)
		fmt.Printf("REC_COUNT (%d)\n", len(r.Results))
	} else {
		b, err := io.ReadAll(resp.Body)
		logutil.FatalErr(err, "ioreadall")

		if resp.StatusCode < 300 {
			b, err = jsonutil.IndentBytes(b, "", "  ")
			logutil.FatalErr(err, "statuscode_lt_300")
		}
		flagWriteJSON := true
		if flagWriteJSON {
			err := os.WriteFile("output.json", b, 0600)
			logutil.FatalErr(err)
		}

		res, err := auditlog.ParseAnalyticsAuditLogArchivalAPIResponse(bytes.NewReader(b))
		logutil.FatalErr(err)
		fmtutil.MustPrintJSON(res)

		times := res.Results.EventTimes()
		fmtutil.MustPrintJSON(times)

		fmtutil.MustPrintJSON(times.Deltas())

		if times.IsSorted(true) {
			fmt.Printf("ORDERED ASC\n\n")
		}
		if times.IsSorted(false) {
			fmt.Printf("ORDERED DESC\n\n")
		}
	}

	fmt.Println("DONE")
}
