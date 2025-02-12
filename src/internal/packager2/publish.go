// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import "context"

type PublishOpts struct {
	SigningKeyPath          string
	SigningKeyPassword      string
	SkipSignatureValidation bool
}

// Takes directory/tar file & OCI Registry

// TODO Dir points to a location on disk and registry is a URL.
func Publish(ctx context.Context, dir string, registry string, opts PublishOpts) error {
	return nil
}
