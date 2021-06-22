// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import (
	"fmt"
)

// f for format
func f(format string, val ...interface{}) string {
	return fmt.Sprintf(format, val...)
}
