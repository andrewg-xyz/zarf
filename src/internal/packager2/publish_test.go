package packager2

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPublish(t *testing.T) {

	tt := []struct {
		name string
		dir  string
		opts PublishOpts
	}{
		{
			name: "simple",
			dir:  "testdata/simple",
			opts: PublishOpts{},
		},
		{
			name: "simple",
			dir:  "testdata/simple",
			opts: PublishOpts{},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := Publish(context.Background(), tc.opts)
			require.NoError(t, err)
		})
	}
}
