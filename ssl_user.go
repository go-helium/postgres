// +build !windows

package postgres

import "os"

func userHomeDir() string { return os.Getenv("HOME") }
