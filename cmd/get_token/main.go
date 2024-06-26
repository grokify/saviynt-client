package main

import (
	"context"
	"fmt"
	"os"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/fmt/fmtutil"
	"github.com/grokify/mogo/log/logutil"
)

func main() {
	_, err := config.LoadDotEnv([]string{".env"}, 1)
	logutil.FatalErr(err)

	tok, err := saviynt.GetToken(
		context.Background(),
		os.Getenv("SAVIYNT_BASE_URL"),
		os.Getenv("SAVIYNT_USERNAME"),
		os.Getenv("SAVIYNT_PASSWORD"))
	logutil.FatalErr(err)

	fmtutil.MustPrintJSON(tok)

	fmt.Println("DONE")
}
