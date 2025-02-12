package packager2

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zarf-dev/zarf/src/test/testutil"
	"oras.land/oras-go/v2/registry"
)

func TestPublish(t *testing.T) {
	ctx := context.Background()

	// TODO add freeport
	registryURL := testutil.SetupInMemoryRegistry(ctx, t, 5000)
	ref := registry.Reference{
		Registry:   registryURL,
		Repository: "my-namespace",
	}

	tt := []struct {
		name   string
		opts   PublishOpts
		expect bool
	}{
		{
			name:   "Test empty publishopts",
			opts:   PublishOpts{},
			expect: true,
		},
		{
			name: "Test empty path",
			opts: PublishOpts{
				Path:     "",
				Registry: ref,
			},
			expect: true,
		},
		// {
		// 	name: "publish skeleton package",
		// 	opts: PublishOpts{
		// 		Path:                    "testdata/skeleton",
		// 		Registry:                ref,
		// 	},
		// 	expect: "",
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// TODO Make parallel
			// t.Parallel()
			err := Publish(context.Background(), tc.opts)
			if tc.expect {
				require.Error(t, err)
			}

			// TODO: Read manifest from registry

			// TODO: check sha of the resulting publish
			// err := oras.PackManifest()
			// require.NoError(t, err)
		})
	}
}
