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
		name string
		opts PublishOpts
	}{
		{
			name: "publish skeleton package",
			opts: PublishOpts{
				Path:     "testdata/skeleton",
				Registry: ref,
			},
		},
		// {
		// 	name:     "simple",
		// 	dir:      "testdata/simple",
		// 	registry: "",
		// 	opts:     PublishOpts{},
		// },
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Publish(context.Background(), tc.opts)
			require.NoError(t, err)
		})
	}
}
