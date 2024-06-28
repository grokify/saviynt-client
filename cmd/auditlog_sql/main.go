package main

import (
	"fmt"

	"github.com/grokify/go-saviynt/auditlog"
)

func main() {
	archival := auditlog.AnalyticsSQLAuditLogArchival()
	fmt.Println(archival)
}
