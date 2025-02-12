// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"

	"oras.land/oras-go/v2/registry"
)

type PublishOpts struct {
	Path                    string
	Registry                registry.Reference
	IsSkeleton              bool
	SigningKeyPath          string
	SigningKeyPassword      string
	SkipSignatureValidation bool
}

func NewPublishOpts(path string, registry registry.Reference, skeleton bool, signingKeyPath string, signingKeyPassword string, skipSignatureValidation bool) (PublishOpts, error) {

	if err := registry.ValidateRegistry(); err != nil {
		return PublishOpts{}, fmt.Errorf("invalid registry: %w", err)
	}

	if path == "" {
		return PublishOpts{}, fmt.Errorf("path must be specified")
	}

	opts := PublishOpts{
		Path:                    path,
		Registry:                registry,
		IsSkeleton:              skeleton,
		SigningKeyPath:          signingKeyPath,
		SigningKeyPassword:      signingKeyPassword,
		SkipSignatureValidation: skipSignatureValidation,
	}

	return opts, nil
}

// Takes directory/tar file & OCI Registry

// TODO Dir points to a location on disk and registry is a URL.
func Publish(ctx context.Context, opts PublishOpts) error {
	return nil
}
