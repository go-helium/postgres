// +build windows
// Copyright (c) 2011-2013, 'pq' Contributors Portions Copyright (C) 2011 Blake Mizerany

package postgres

// sslKeyPermissions checks the permissions on user-supplied ssl key files.
// The key file should have very little access.
//
// libpq does not check key file permissions on Windows.
func sslKeyPermissions(string) error { return nil }
