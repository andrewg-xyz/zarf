// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import (
	"context"
	"fmt"
	"strings"

	"github.com/defenseunicorns/pkg/helpers/v2"
	"github.com/zarf-dev/zarf/src/config"
	"github.com/zarf-dev/zarf/src/pkg/logger"
	"github.com/zarf-dev/zarf/src/pkg/zoci"

	"github.com/defenseunicorns/pkg/oci"
	layout2 "github.com/zarf-dev/zarf/src/internal/packager2/layout"

	// "github.com/zarf-dev/zarf/src/pkg/layout"
	"oras.land/oras-go/v2/registry"
)

// PublishOpts declares the parameters to publish a package.
type PublishOpts struct {
	// Concurrency configures the zoci push concurrency if empty defaults to 3.
	Concurrency int
	// SigningKeyPath points to a signing key on the local disk.
	SigningKeyPath string
	// SigningKeyPassword holds a password to use the key at SigningKeyPath.
	SigningKeyPassword string
	// SkipSignatureValidation flags whether Publish should skip validating the signature.
	SkipSignatureValidation bool
	// WithPlainHTTP falls back to plain HTTP for the registry calls instead of TLS.
	WithPlainHTTP bool
	// PublicKeyPath validates the create time signage of a package.
	PublicKeyPath string
	// Architecture is the architecture we are publishing to
	Architecture string
}

// Publish takes a Path to the location of the build package, a ref to a registry, and a PublishOpts and uploads the
// package tarball, oci reference, or skeleton package to the registry.
func Publish(ctx context.Context, path string, dst registry.Reference, opts PublishOpts) error {
	l := logger.From(ctx)
	// TODO: determine if the source is an OCI reference and a zoci.CopyPackage() is required
	// TODO: can you copy to and from the same registry?

	// Validate inputs
	l.Debug("validating PublishOpts")
	if err := dst.ValidateRegistry(); err != nil {
		return fmt.Errorf("invalid registry: %w", err)
	}
	if path == "" {
		return fmt.Errorf("path must be specified")
	}

	// If path is remote copy oci to oci
	if helpers.IsOCIURL(path) {

		// Build srcRef
		trimmed := strings.TrimPrefix(path, "oci://")
		srcRef, err := registry.ParseReference(trimmed)
		if err != nil {
			return fmt.Errorf("failed to parse path, path=%s: %w", path, err)
		}

		// Extract packageName from src and build dstUrl for dst remote
		parts := strings.Split(srcRef.String(), "/")
		packageName := parts[len(parts)-1]
		dstUrl := dst.String() + "/" + packageName

		// Build platform
		arch := config.GetArch(opts.Architecture)
		p := oci.PlatformForArch(arch)

		// Set up remote repo client
		src, err := zoci.NewRemote(ctx, srcRef.String(), p, oci.WithPlainHTTP(opts.WithPlainHTTP))
		if err != nil {
			return fmt.Errorf("could not instantiate remote: %w", err)
		}
		dstRem, err := zoci.NewRemote(ctx, dstUrl, p, oci.WithPlainHTTP(opts.WithPlainHTTP))
		if err != nil {
			return fmt.Errorf("could not instantiate remote: %w", err)
		}

		// Execute copy
		err = zoci.CopyPackage(ctx, src, dstRem, opts.Concurrency)
		if err != nil {
			return fmt.Errorf("could not copy package: %w", err)
		}

		return nil
	}

	// Load package layout
	l.Info("loading package", "path", path)
	layoutOpts := layout2.PackageLayoutOptions{
		PublicKeyPath:           opts.PublicKeyPath,
		SkipSignatureValidation: opts.SkipSignatureValidation,
	}
	pkgLayout, err := layout2.LoadFromTar(ctx, path, layoutOpts)
	if err != nil {
		return fmt.Errorf("unable to load package: %w", err)
	}

	return pushToRemote(ctx, pkgLayout, dst, opts.Concurrency, opts.WithPlainHTTP)
}

// PublishSkeletonOpts declares the parameters to publish a skeleton package.
type PublishSkeletonOpts struct {
	// Concurrency configures the zoci push concurrency if empty defaults to 3.
	Concurrency int
	// SigningKeyPath points to a signing key on the local disk.
	SigningKeyPath string
	// SigningKeyPassword holds a password to use the key at SigningKeyPath.
	SigningKeyPassword string
	// WithPlainHTTP falls back to plain HTTP for the registry calls instead of TLS.
	WithPlainHTTP bool
}

// PublishSkeleton takes a Path to the location of the build package, a ref to a registry, and a PublishOpts and uploads the
// package tarball, oci reference, or skeleton package to the registry.
func PublishSkeleton(ctx context.Context, path string, ref registry.Reference, opts PublishSkeletonOpts) error {
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
	l.Info("loading skeleton package", "path", path)
	// Create skeleton buildpath
	createOpts := layout2.CreateOptions{
		SigningKeyPath:     opts.SigningKeyPath,
		SigningKeyPassword: opts.SigningKeyPassword,
		SetVariables:       map[string]string{},
	}
	buildPath, err := layout2.CreateSkeleton(ctx, path, createOpts)
	if err != nil {
		return fmt.Errorf("unable to create skeleton: %w", err)
	}

	// Generate partial layout from buildpath
	layoutOpts := layout2.PackageLayoutOptions{
		SkipSignatureValidation: true,
		// TODO: define what IsPartial purpose is in code docs
		IsPartial: true,
	}
	pkgLayout, err := layout2.LoadFromDir(ctx, buildPath, layoutOpts)
	if err != nil {
		return fmt.Errorf("unable to load skeleton: %w", err)
	}

	return pushToRemote(ctx, pkgLayout, ref, opts.Concurrency, opts.WithPlainHTTP)
}

// pushToRemote pushes a package to a remote at ref.
func pushToRemote(ctx context.Context, layout *layout2.PackageLayout, ref registry.Reference, concurrency int, plainHTTP bool) error {
	// Build Reference for remote from registry location and pkg
	r, err := layout2.ReferenceFromMetadata(ref.String(), layout.Pkg)
	if err != nil {
		return err
	}

	arch := layout.Pkg.Metadata.Architecture
	// Set platform
	p := oci.PlatformForArch(arch)

	// Set up remote repo client
	rem, err := layout2.NewRemote(ctx, r, p, oci.WithPlainHTTP(plainHTTP))
	if err != nil {
		return fmt.Errorf("could not instantiate remote: %w", err)
	}

	logger.From(ctx).Info("pushing package to remote registry",
		"ref", ref,
		"architecture", arch)
	return rem.Push(ctx, layout, concurrency)
}
