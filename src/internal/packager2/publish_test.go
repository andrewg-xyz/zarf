package packager2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/defenseunicorns/pkg/oci"
	goyaml "github.com/goccy/go-yaml"
	"github.com/stretchr/testify/require"
	"github.com/zarf-dev/zarf/src/api/v1alpha1"
	"github.com/zarf-dev/zarf/src/internal/packager2/layout"
	"github.com/zarf-dev/zarf/src/pkg/zoci"
	"github.com/zarf-dev/zarf/src/test/testutil"
	"oras.land/oras-go/v2/registry"
)

func TestPublishError(t *testing.T) {
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
		expectErr error
	}{
		{
			name:      "Test empty publishopts",
			opts:      PublishOpts{},
			expectErr: errors.New("invalid registry"),
		},
		{
			name: "Test empty path",
			opts: PublishOpts{
				Path:     "",
				Registry: ref,
			},
			expectErr: errors.New("path must be specified"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// TODO Make parallel
			// t.Parallel()
			err := Publish(context.Background(), tc.opts)
			require.ErrorContains(t, err, tc.expectErr.Error())
		})
	}
}

func TestPublish(t *testing.T) {
	ctx := context.Background()

	// TODO add freeport
	registryURL := testutil.SetupInMemoryRegistry(ctx, t, 5000)
	ref := registry.Reference{
		Registry:   registryURL,
		Repository: "my-namespace",
	}

	tt := []struct {
		name string
		opts PublishOpts
	}{
		{
			name: "Publish skeleton package",
			opts: PublishOpts{
				Path:     "testdata/skeleton",
				Registry: ref,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			// TODO Make parallel
			// t.Parallel()
			err := Publish(context.Background(), tc.opts)
			require.NoError(t, err)

			// fmt.Println("efr"tc.opts.Registry.Reference)
			data, err := os.ReadFile(filepath.Join(tc.opts.Path, layout.ZarfYAML))
			require.NoError(t, err)

			var expectedPkg v1alpha1.ZarfPackage
			goyaml.Unmarshal(data, &expectedPkg)

			format := "%s/%s:%s"
			artifactURL := fmt.Sprintf(format, registryURL, expectedPkg.Metadata.Name, expectedPkg.Metadata.Version)

			fmt.Println("this is the artifact url", artifactURL)
			ref, err := registry.ParseReference(artifactURL)
			require.NoError(t, err)

			registry.ParseReference(artifactURL)
			rmt, err := zoci.NewRemote(ctx, ref.String(), zoci.PlatformForSkeleton(), oci.WithPlainHTTP(true))
			require.NoError(t, err)

			pkg, err := rmt.FetchZarfYAML(ctx)
			require.NoError(t, err)

			require.Equal(t, pkg, expectedPkg)

			// TODO: check sha of the resulting publish
			// err := oras.PackManifest()
			// require.NoError(t, err)
		})
	}
}
