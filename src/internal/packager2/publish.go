// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"

	"github.com/zarf-dev/zarf/src/internal/packager2/layout"
	"github.com/zarf-dev/zarf/src/pkg/zoci"
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

	// Validate inputs
	if err := opts.Registry.ValidateRegistry(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}

	if opts.Path == "" {
		return fmt.Errorf("path must be specified")
	}

	// TODO Create skeleton locally
	cOpts := layout.CreateOptions{
		Flavor:                  "",
		RegistryOverrides:       nil,
		SigningKeyPath:          "",
		SigningKeyPassword:      "",
		SetVariables:            nil,
		SkipSBOM:                false,
		DifferentialPackagePath: "",
	}
	// TODO Resolve compiler errors
	layout.CreateSkeleton(ctx, packagePath, cOpts)

	// TODO Do publish to remote
	// TODO Resolve compiler errors
	rem, err := zoci.NewRemote(ctx, url, platform, mods)
	if err != nil {
		return err
	}
	// TODO(mkcp): Resolve compiler errors
	err = rem.PublishPackage(ctx, pkg, paths, concurrency)
	if err != nil {
		return err
	}

	return nil
}
