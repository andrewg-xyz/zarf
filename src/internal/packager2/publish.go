// SPDX-License-Identifier: Apache-2.0
// SPDX-FileCopyrightText: 2021-Present The Zarf Authors

package packager2

import "context"

type PublishOpts struct {
	SigningKeyPath          string
	SigningKeyPassword      string
	SkipSignatureValidation bool
}

func Publish(ctx context.Context, opts PublishOpts) error {
	return nil
}
