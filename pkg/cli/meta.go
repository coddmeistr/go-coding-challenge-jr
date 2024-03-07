package cli

import (
	"challenge/pkg/proto"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func init() {
	rootCmd.AddCommand(metadataCommand)
	metadataCommand.Flags().StringVarP(&meta, "meta", "m", "", "metadata to encode")
}

var meta string
var metadataCommand = &cobra.Command{
	Use:   "metadata",
	Short: "Extract metadata",
	Long:  `Extract metadata from gRPC context'`,
	Run: func(_ *cobra.Command, _ []string) {

		ctx, cancel := context.WithTimeout(context.Background(), 20)
		defer cancel()
		conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("cannot connect to gRPC server: %v\n", err)
			return
		}

		client := proto.NewChallengeServiceClient(conn)
		data, err := client.ReadMetadata(metadata.NewOutgoingContext(context.Background(), metadata.Pairs("i-am-random-key", meta)), &proto.Placeholder{})
		if err != nil {
			fmt.Printf("cannot extract metadata: %v\n", err)
			return
		}

		fmt.Printf("extracted metadata: %s\n", data.GetData())
	},
}
