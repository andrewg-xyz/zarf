// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"
	"github.com/zarf-dev/zarf/src/pkg/logger"

	"github.com/defenseunicorns/pkg/oci"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	layout2 "github.com/zarf-dev/zarf/src/internal/packager2/layout"
	// "github.com/zarf-dev/zarf/src/pkg/layout"
	"oras.land/oras-go/v2/registry"
)

// PublishOpts declares the parameters to publish a package.
type PublishOpts struct {
	// Concurrency configures the zoci push concurrency if empty defaults to 3.
	Concurrency int
	// IsSkeleton flags whether a package path is a skeleton pkg.
	IsSkeleton bool
	// SigningKeyPath points to a signing key on the local disk.
	SigningKeyPath string
	// SigningKeyPassword holds a password to use the key at SigningKeyPath.
	SigningKeyPassword string
	// SkipSignatureValidation flags whether Publish should skip validating the signature.
	SkipSignatureValidation bool
	// WithPlainHTTP falls back to plain HTTP for the registry calls instead of TLS.
	WithPlainHTTP bool
}

// Publish takes a Path to the location of the build package, a ref to a registry, and a PublishOpts and uploads the
// package tarball, oci reference, or skeleton package to the registry.
func Publish(ctx context.Context, path string, ref registry.Reference, opts PublishOpts) error {
	l := logger.From(ctx)
	// TODO: determine if the source is an OCI reference and a zoci.CopyPackage() is required
	// TODO: can you copy to and from the same registry?

	// Validate inputs
	l.Debug("validating PublishOpts")
	if err := ref.ValidateRegistry(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}
	if path == "" {
		return fmt.Errorf("path must be specified")
	}

	// Load package layout
	l.Info("loading package", "path", path)
	var pkgLayout *layout2.PackageLayout
	var err error
	if opts.IsSkeleton {
		l.Debug("loading skeleton package", "path", path)
		pkgLayout, err = loadSkeleton(ctx, path, opts)
		if err != nil {
			return fmt.Errorf("unable to load skeleton: %w", err)
		}
	} else {
		l.Debug("loading package", "path", path)
		pkgLayout, err = loadPackage(ctx, path, opts.SkipSignatureValidation)
		if err != nil {
			return fmt.Errorf("unable to load package: %w", err)
		}
	}

	return pushToRemote(ctx, pkgLayout, ref, opts)
}

// loadSkeleton generates a skeleton package from a directory
func loadSkeleton(ctx context.Context, path string, opts PublishOpts) (*layout2.PackageLayout, error) {
	// Create skeleton buildpath
	createOpts := layout2.CreateOptions{
		SigningKeyPath:     opts.SigningKeyPath,
		SigningKeyPassword: opts.SigningKeyPassword,
		SetVariables:       map[string]string{},
	}
	buildPath, err := layout2.CreateSkeleton(ctx, path, createOpts)
	if err != nil {
		return nil, fmt.Errorf("unable to create skeleton: %w", err)
	}

	// Generate partial layout from buildpath
	layoutOpts := layout2.PackageLayoutOptions{
		SkipSignatureValidation: opts.SkipSignatureValidation,
		// TODO: define what IsPartial purpose is in code docs
		IsPartial: true,
	}
	return layout2.LoadFromDir(ctx, buildPath, layoutOpts)
}

// loadPackage loads an existing package's layout from tarball
func loadPackage(ctx context.Context, path string, isSkipValidation bool) (*layout2.PackageLayout, error) {
	layoutOpts := layout2.PackageLayoutOptions{
		SkipSignatureValidation: isSkipValidation,
	}
	return layout2.LoadFromTar(ctx, path, layoutOpts)
}

func pushToRemote(ctx context.Context, layout *layout2.PackageLayout, ref registry.Reference, opts PublishOpts) error {
	// Build Reference for remote from registry location and pkg
	r, err := layout2.ReferenceFromMetadata(ref.String(), layout.Pkg)
	if err != nil {
		return err
	}

	arch := layout.Pkg.Metadata.Architecture
	// Set platform
	p := ocispec.Platform{
		OS:           oci.MultiOS,
		Architecture: arch,
	}

	// Set up remote repo client
	rem, err := layout2.NewRemote(ctx, r, p, oci.WithPlainHTTP(opts.WithPlainHTTP))
	if err != nil {
		return fmt.Errorf("could not instantiate remote: %w", err)
	}

	logger.From(ctx).Info("pushing package to remote registry",
		"ref", ref,
		"architecture", arch)
	return rem.Push(ctx, layout, opts.Concurrency)
}
