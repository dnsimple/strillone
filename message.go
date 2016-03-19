package main

import (
	"fmt"
)

func fmtURL(path string, a ...interface{}) string {
	return fmt.Sprintf(dnsimpleURL+path, a...)
}
