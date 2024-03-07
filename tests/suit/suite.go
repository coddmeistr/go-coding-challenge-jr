package suits

import (
	"challenge/pkg/proto"
	"os"
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

	grpcAddress := os.Getenv("GRPC_HOST_PORT")
	if grpcAddress == "" {
		t.Fatalf("GRPC_HOST_PORT is not set")
	}

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
