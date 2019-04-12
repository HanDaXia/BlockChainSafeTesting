// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build dragonfly freebsd linux netbsd openbsd solaris

package x509

import (
	"fmt"
	"github.com/HanDaXia/BlockChainSafeTesting/messagehub/cryptoutil"
	"os"
	"testing"
)

const (
	testDir     = "testdata"
	testDirCN   = "test-dir"
	testFile    = "test-file.crt"
	testFileCN  = "test-file"
	testMissing = "missing"
)

func TestEnvVars(t *testing.T) {
	testCases := []struct {
		name    string
		fileEnv string
		dirEnv  string
		files   []string
		dirs    []string
		cns     []string
	}{
		{
			// Environment variables override the default locations preventing fall through.
			name:    "override-defaults",
			fileEnv: testMissing,
			dirEnv:  testMissing,
			files:   []string{testFile},
			dirs:    []string{testDir},
			cns:     nil,
		},
		{
			// File environment overrides default file locations.
			name:    "file",
			fileEnv: testFile,
			dirEnv:  "",
			files:   nil,
			dirs:    nil,
			cns:     []string{testFileCN},
		},
		{
			// Directory environment overrides default directory locations.
			name:    "dir",
			fileEnv: "",
			dirEnv:  testDir,
			files:   nil,
			dirs:    nil,
			cns:     []string{testDirCN},
		},
		{
			// File & directory environment overrides both default locations.
			name:    "file+dir",
			fileEnv: testFile,
			dirEnv:  testDir,
			files:   nil,
			dirs:    nil,
			cns:     []string{testFileCN, testDirCN},
		},
		{
			// Environment variable empty / unset uses default locations.
			name:    "empty-fall-through",
			fileEnv: "",
			dirEnv:  "",
			files:   []string{testFile},
			dirs:    []string{testDir},
			cns:     []string{testFileCN, testDirCN},
		},
	}

	// Save old settings so we can restore before the test ends.
	origCertFiles, origCertDirectories := cryptoutil.certFiles, cryptoutil.certDirectories
	origFile, origDir := os.Getenv(cryptoutil.certFileEnv), os.Getenv(cryptoutil.certDirEnv)
	defer func() {
		cryptoutil.certFiles = origCertFiles
		cryptoutil.certDirectories = origCertDirectories
		os.Setenv(cryptoutil.certFileEnv, origFile)
		os.Setenv(cryptoutil.certDirEnv, origDir)
	}()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := os.Setenv(cryptoutil.certFileEnv, tc.fileEnv); err != nil {
				t.Fatalf("setenv %q failed: %v", cryptoutil.certFileEnv, err)
			}
			if err := os.Setenv(cryptoutil.certDirEnv, tc.dirEnv); err != nil {
				t.Fatalf("setenv %q failed: %v", cryptoutil.certDirEnv, err)
			}

			cryptoutil.certFiles, cryptoutil.certDirectories = tc.files, tc.dirs

			r, err := cryptoutil.loadSystemRoots()
			if err != nil {
				t.Fatal("unexpected failure:", err)
			}

			if r == nil {
				if tc.cns == nil {
					// Expected nil
					return
				}
				t.Fatal("nil roots")
			}

			// Verify that the returned certs match, otherwise report where the mismatch is.
			for i, cn := range tc.cns {
				if i >= len(r.certs) {
					t.Errorf("missing cert %v @ %v", cn, i)
				} else if r.certs[i].Subject.CommonName != cn {
					fmt.Printf("%#v\n", r.certs[0].Subject)
					t.Errorf("unexpected cert common name %q, want %q", r.certs[i].Subject.CommonName, cn)
				}
			}
			if len(r.certs) > len(tc.cns) {
				t.Errorf("got %v certs, which is more than %v wanted", len(r.certs), len(tc.cns))
			}
		})
	}
}
