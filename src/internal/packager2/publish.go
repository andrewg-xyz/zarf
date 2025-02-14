// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/defenseunicorns/pkg/oci"
	layout2 "github.com/zarf-dev/zarf/src/internal/packager2/layout"
	"github.com/zarf-dev/zarf/src/pkg/layout"
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
	cOpts := layout2.CreateOptions{
		SigningKeyPath:     opts.SigningKeyPath,
		SigningKeyPassword: opts.SigningKeyPassword,
		SetVariables:       map[string]string{},
	}
	// TODO Resolve compiler errors
	buildPath, err := layout2.CreateSkeleton(ctx, opts.Path, cOpts)
	if err != nil {
		return fmt.Errorf("unable to create skeleton: %w", err)
	}

	layoutOpt := layout2.PackageLayoutOptions{
		SkipSignatureValidation: opts.SkipSignatureValidation,
		IsPartial:               true,
	}
	pkgLayout, err := layout2.LoadFromDir(ctx, buildPath, layoutOpt)
	if err != nil {
		return fmt.Errorf("unable to load package: %w", err)
	}

	// TODO can we convert from packager types to packager2 types
	rem, err := zoci.NewRemote(ctx, opts.Registry.String(), zoci.PlatformForSkeleton(), oci.WithPlainHTTP(opts.WithPlainHTTP))
	if err != nil {
		return fmt.Errorf("could not instantiate remote: %w", err)
	}
	layout1 := layout.New(pkgLayout.DirPath())

	spew.Dump(rem, pkgLayout, layout1)
	err = rem.PublishPackage(ctx, &pkgLayout.Pkg, layout1, 3)
	if err != nil {
		return fmt.Errorf("could not publish package: %w", err)
	}

	return nil
}
