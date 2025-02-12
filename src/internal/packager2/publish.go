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

// Takes directory/tar file & OCI Registry

// TODO Dir points to a location on disk and registry is a URL.
func Publish(ctx context.Context, opts PublishOpts) error {
	if err := opts.Registry.ValidateRegistry(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}

	if opts.Path == "" {
		return fmt.Errorf("path must be specified")
	}
	return nil
}
