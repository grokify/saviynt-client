package main

import (
	"fmt"
	"io"
	"os"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/log/logutil"
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

	resp, err := clt.GetAuditLogRuntimeControlsData(os.Getenv("SAVIYNT_QUERY_NAME"), 100000, 50, 0)
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
