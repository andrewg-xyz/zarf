package packager2

import (
	"context"
	"os"
	"testing"

	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"github.com/zarf-dev/zarf/src/api/v1alpha1"
	"github.com/zarf-dev/zarf/src/pkg/zoci"
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
		name      string
		opts      PublishOpts
		expectErr bool
	}{
		{
			name:      "Test empty publishopts",
			opts:      PublishOpts{},
			expectErr: true,
		},
		{
			name: "Test empty path",
			opts: PublishOpts{
				Path:     "",
				Registry: ref,
			},
			expectErr: true,
		},
		{
			name: "Publish skeleton package",
			opts: PublishOpts{
				Path:     "testdata/skeleton",
				Registry: ref,
			},
			expectErr: true,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// TODO Make parallel
			// t.Parallel()
			err := Publish(context.Background(), tc.opts)
			if tc.expectErr {
				require.Error(t, err)
			}

			rmt, err := zoci.NewRemote(ctx, tc.opts.Registry.Reference, zoci.PlatformForSkeleton())
			require.NoError(t, err)

			pkg, err := rmt.FetchZarfYAML(ctx)
			require.NoError(t, err)

			data, err := os.ReadFile(tc.opts.Path)
			require.NoError(t, err)

			var expectedPkg v1alpha1.ZarfPackage
			goyaml.Unmarshal(data, &expectedPkg)

			require.Equal(t, pkg, expectedPkg)

			// TODO: check sha of the resulting publish
			// err := oras.PackManifest()
			// require.NoError(t, err)
		})
	}
}
