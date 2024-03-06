package challenge_server

import (
	"challenge/pkg/proto"
	"context"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
	"testing"
)

func TestReadMetadata_TestCases(t *testing.T) {
	tc := []struct {
		name           string
		hasMetadata    bool
		metadataString string
		wantString     string
		wantError      bool
	}{
		{
			name:           "ok",
			hasMetadata:    true,
			metadataString: "some_random_metadata",
			wantString:     "some_random_metadata",
			wantError:      false,
		},
		{
			name:           "no metadata",
			hasMetadata:    false,
			metadataString: "",
			wantString:     "",
			wantError:      true,
		},
		{
			name:           "empty metadata",
			hasMetadata:    true,
			metadataString: "",
			wantString:     "",
			wantError:      false,
		},
	}

	caller := &server{}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			if tt.hasMetadata {
				ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(metadataKey, tt.metadataString))
			}

			got, err := caller.ReadMetadata(ctx, &proto.Placeholder{})
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantString, got.Data)
			}
		})
	}
}
