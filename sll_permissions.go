// +build !windows
// Copyright (c) 2011-2013, 'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany

package postgres

import "os"

// sslKeyPermissions checks the permissions on user-supplied ssl key files.
// The key file should have very little access.
//
// libpq does not check key file permissions on Windows.
func sslKeyPermissions(sslkey string) error {
	info, err := os.Stat(sslkey)
	if err != nil {
		return err
	}
	if info.Mode().Perm()&0077 != 0 {
		return ErrSSLKeyHasWorldPermissions
	}
	return nil
}
