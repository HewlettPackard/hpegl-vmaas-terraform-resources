// (C) Copyright 2021 Hewlett Packard Enterprise Development LP

package resources

import "fmt"

const (
	generalNamedesc = "Name of the %s as it appears on GLC. If no %s is found with this name, an error will be returned"
	generalDDesc    = "Unique ID to identify a %s"
	dsHeadingDesc   = "Use this data source to get the %s. This data can be fetched under %s"
)

// f for format
func f(format string, val ...interface{}) string {
	return fmt.Sprintf(format, val...)
}
