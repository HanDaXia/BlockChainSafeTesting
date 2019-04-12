// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !cgo

package x509

import "github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil"

func loadSystemRoots() (*cryptoutil.CertPool, error) {
	return execSecurityRoots()
}