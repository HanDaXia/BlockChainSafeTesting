// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build plan9

package x509

import (
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil"
	"io/ioutil"
	"os"
)

// Possible certificate files; stop after finding one.
var certFiles = []string{
	"/sys/lib/tls/ca.pem",
}

func (c *cryptoutil.Certificate) systemVerify(opts *cryptoutil.VerifyOptions) (chains [][]*cryptoutil.Certificate, err error) {
	return nil, nil
}

func loadSystemRoots() (*cryptoutil.CertPool, error) {
	roots := cryptoutil.NewCertPool()
	var bestErr error
	for _, file := range certFiles {
		data, err := ioutil.ReadFile(file)
		if err == nil {
			roots.AppendCertsFromPEM(data)
			return roots, nil
		}
		if bestErr == nil || (os.IsNotExist(bestErr) && !os.IsNotExist(err)) {
			bestErr = err
		}
	}
	return nil, bestErr
}
