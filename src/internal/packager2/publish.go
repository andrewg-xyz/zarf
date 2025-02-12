// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"
)

type PublishOpts struct {
	Path                    string
	Registry                string
	IsSkeleton              bool
	SigningKeyPath          string
	SigningKeyPassword      string
	SkipSignatureValidation bool
}

func NewPublishOpts(path string, registry string, skeleton bool, signingKeyPath string, signingKeyPassword string, skipSignaturevalidation bool) (PublishOpts, error) {

	if registry == "" {
		return PublishOpts{}, fmt.Errorf("registry must be specified")
	}

	opts := PublishOpts{
		Path:                    path,
		Registry:                registry,
		IsSkeleton:              skeleton,
		SigningKeyPath:          signingKeyPath,
		SigningKeyPassword:      signingKeyPassword,
		SkipSignatureValidation: skipSignaturevalidation,
	}

	return opts, nil
}

// Takes directory/tar file & OCI Registry

// TODO Dir points to a location on disk and registry is a URL.
func Publish(ctx context.Context, path string, registry string, opts PublishOpts) error {
	return nil
}
