/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package factory

import (
	"errors"
	"github.com/hyperledger/fabric/bccsp"
	"github.com/hyperledger/fabric/bccsp/gm"
)

const (
	// PluginFactoryName is the factory name for BCCSP plugins
	GMFactoryName = "SW"
)

// PluginFactory is the factory for BCCSP plugins
type GMFactory struct{}

// Name returns the name of this factory
func (f *GMFactory) Name() string {
	return GMFactoryName
}

// Get returns an instance of BCCSP using Opts.
func (f *GMFactory) Get(config *FactoryOpts) (bccsp.BCCSP, error) {
	// Validate arguments
	if config == nil || config.SwOpts == nil {
		return nil, errors.New("Invalid config. It must not be nil.")
	}

	swOpts := config.SwOpts

	var keystore string
	if swOpts.FileKeystore != nil && swOpts.Ephemeral != true {
		keystore = swOpts.FileKeystore.KeyStorePath
	}

	return gm.New(keystore)
}