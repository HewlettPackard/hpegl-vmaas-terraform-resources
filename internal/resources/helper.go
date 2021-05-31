// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"
	"strings"
)

// Set header description for datastore
func setDsHeader(dsName, dsTitle string, examples ...string) string {
	ex := strings.Join(examples, ",")
	return f(dsHeadingDesc, dsName, dsTitle, dsTitle, ex)
}

// f for format
func f(format string, val ...interface{}) string {
	return fmt.Sprintf(format, val...)
}
