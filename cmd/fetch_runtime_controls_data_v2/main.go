package main

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/log/logutil"
	"github.com/grokify/mogo/type/maputil"
)

func main() {
	_, err := config.LoadDotEnv([]string{".env"}, 1)
	logutil.FatalErr(err)

	clt, err := saviynt.NewClient(
		os.Getenv("SAVIYNT_BASE_URL"),
		saviynt.RelURLAPI,
		os.Getenv("SAVIYNT_USERNAME"),
		os.Getenv("SAVIYNT_PASSWORD"))
	logutil.FatalErr(err)

	attrs := map[string]any{}
	if attrsStr := os.Getenv("SAVIYNT_QUERY_ATTR"); len(attrsStr) > 0 {
		attrsVals, err := url.ParseQuery(attrsStr)
		logutil.FatalErr(err)
		attrsValsMap := maputil.MapStringSlice(attrsVals)
		attrs = attrsValsMap.FlattenAny(false, false)
		// attrs = MS3ToMSA(attrsVals)
	}

	resp, err := clt.FetchRuntimeControlsDataV2(
		os.Getenv("SAVIYNT_QUERY_NAME"),
		attrs,
		50, 0)
	logutil.FatalErr(err)

	b, err := io.ReadAll(resp.Body)
	logutil.FatalErr(err)

	if resp.StatusCode < 300 {
		b, err = jsonutil.IndentBytes(b, "", "  ")
		logutil.FatalErr(err)
	}
	fmt.Println(string(b))

	fmt.Println("DONE")
}

/*

func MS3ToMSS(m map[string][]string) map[string]string {
	mss := map[string]string{}
	for k, vs := range m {
		for _, v := range vs {
			mss[k] = v
		}
	}
	return mss
}

func MS3ToMSA(m map[string][]string) map[string]any {
	mss := map[string]any{}
	for k, vs := range m {
		for _, v := range vs {
			mss[k] = v
		}
	}
	return mss
}

*/
