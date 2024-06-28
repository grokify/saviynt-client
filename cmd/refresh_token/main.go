package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/grokify/go-saviynt"
	"github.com/grokify/mogo/config"
	"github.com/grokify/mogo/log/logutil"
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

	fmt.Printf("A.TOK (%s)\n", clt.Token.AccessToken)
	fmt.Printf("R.TOK (%s)\n", clt.Token.RefreshToken)

	for i := range 2 {
		rtok, err := clt.GetTokenRefresh(context.Background())
		logutil.FatalErr(err)
		fmt.Printf("R[%d].TOK (%s)\n", i, rtok.AccessToken)
		fmt.Printf("R[%d].TOK (%s)\n", i, rtok.RefreshToken)
		time.Sleep(10 * time.Second)
	}

	fmt.Println("DONE")
}
