// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"

	"github.com/defenseunicorns/pkg/oci"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	layout2 "github.com/zarf-dev/zarf/src/internal/packager2/layout"
	// "github.com/zarf-dev/zarf/src/pkg/layout"
	"github.com/zarf-dev/zarf/src/pkg/zoci"
	"oras.land/oras-go/v2/registry"
)

type PublishOpts struct {
	// Concurrency configures the zoci push concurrency if empty defaults to 3.
	Concurrency             int
	Path                    string
	Registry                registry.Reference
	IsSkeleton              bool
	SigningKeyPath          string
	SigningKeyPassword      string
	SkipSignatureValidation bool
	WithPlainHTTP           bool
}

// Publish takes PublishOpts and uploads the package tarball, oci reference, or skeleton package to the registry.
func Publish(ctx context.Context, opts PublishOpts) error {
	// TODO: check linter for packager2
	// TODO: should we be using packager2 oci NewRemote and Push instead of zoci?
	// If so, do we need to implement a layout2.Copy function?
	var err error

	// Validate inputs
	if err := opts.Registry.ValidateRegistry(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}

	if opts.Path == "" {
		return fmt.Errorf("path must be specified")
	}

	// TODO: determine if the source is an OCI reference and a zoci.CopyPackage() is required
	// TODO: can you copy to and from the same registry?

	var pkgLayout *layout2.PackageLayout
	var platform ocispec.Platform

	layoutOpt := layout2.PackageLayoutOptions{
		SkipSignatureValidation: opts.SkipSignatureValidation,
	}

	if opts.IsSkeleton {
		cOpts := layout2.CreateOptions{
			SigningKeyPath:     opts.SigningKeyPath,
			SigningKeyPassword: opts.SigningKeyPassword,
			SetVariables:       map[string]string{},
		}

		buildPath, err := layout2.CreateSkeleton(ctx, opts.Path, cOpts)
		if err != nil {
			return fmt.Errorf("unable to create skeleton: %w", err)
		}

		// TODO: define what IsPartial purpose is in code docs
		layoutOpt.IsPartial = true
		pkgLayout, err = layout2.LoadFromDir(ctx, buildPath, layoutOpt)
		if err != nil {
			return fmt.Errorf("unable to load package: %w", err)
		}
		// cleanup zoci use
		platform = zoci.PlatformForSkeleton()
	} else {
		// publish a built package

		pkgLayout, err = layout2.LoadFromTar(ctx, opts.Path, layoutOpt)
		if err != nil {
			return err
		}
		platform = ocispec.Platform{OS: oci.MultiOS, Architecture: pkgLayout.Pkg.Metadata.Architecture}
	}

	ref, err := layout2.ReferenceFromMetadata(opts.Registry.String(), pkgLayout.Pkg)
	if err != nil {
		return err
	}

	// // TODO can we convert from packager types to packager2 types
	// ref, err := zoci.ReferenceFromMetadata(opts.Registry.String(), &pkgLayout.Pkg.Metadata, &pkgLayout.Pkg.Build)
	// if err != nil {
	// 	return fmt.Errorf("unable to create reference: %w", err)
	// }

	rem, err := layout2.NewRemote(ctx, ref, platform, oci.WithPlainHTTP(opts.WithPlainHTTP))

	// rem, err := zoci.NewRemote(ctx, ref, platform, oci.WithPlainHTTP(opts.WithPlainHTTP))

	if err != nil {
		return fmt.Errorf("could not instantiate remote: %w", err)
	}
	// layout1 := layout.New(pkgLayout.DirPath())

	err = rem.Push(ctx, pkgLayout, opts.Concurrency)

	// err = rem.PublishPackage(ctx, &pkgLayout.Pkg, layout1, opts.Concurrency)
	if err != nil {
		return fmt.Errorf("could not publish package: %w", err)
	}

	return nil
}
