// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/pkg/oci"
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
	WithPlainHTTP           bool
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

	// TODO skeleton and flavors during publish
	// TODO Create skeleton locally
	cOpts := layout.CreateOptions{
		SigningKeyPath:     opts.SigningKeyPath,
		SigningKeyPassword: opts.SigningKeyPassword,
	}
	// TODO Resolve compiler errors
	buildPath, err := layout.CreateSkeleton(ctx, opts.Path, cOpts)
	if err != nil {
		return err
	}

	layoutOpt := layout.PackageLayoutOptions{
		SkipSignatureValidation: opts.SkipSignatureValidation,
		IsPartial:               true,
	}
	pkgLayout, err := layout.LoadFromTar(ctx, buildPath, layoutOpt)
	if err != nil {
		return err
	}

	// TODO can we convert from packager types to packager2 types

	// TODO Do publish to remote

	// TODO Resolve compiler errors
	rem, err := zoci.NewRemote(ctx, opts.Registry.String(), zoci.PlatformForSkeleton(), oci.WithPlainHTTP(opts.WithPlainHTTP))
	if err != nil {
		return err
	}
	// TODO(mkcp): Resolve compiler errors
	err = rem.PublishPackage(ctx, &pkgLayout.Pkg, paths, concurrency)
	if err != nil {
		return err
	}

	return nil
}
