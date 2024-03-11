package suits

import (
	"challenge/pkg/config"
	"challenge/pkg/proto"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Suite struct {
	*testing.T
	Client proto.ChallengeServiceClient
}

func NewDefault(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	// Load environment variables
	config.LoadEnvs("../.env")
	grpcAddress := viper.GetString("GRPC_HOST_PORT")
	require.NotEqual(t, "", grpcAddress)

	ctx, cancel := context.WithTimeout(context.Background(), 20)

	t.Cleanup(func() {
		t.Helper()
		cancel()
	})

	cc, err := grpc.DialContext(ctx, grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("grpc connection failed: %v", err)
	}

	return ctx, &Suite{
		T:      t,
		Client: proto.NewChallengeServiceClient(cc),
	}
}
