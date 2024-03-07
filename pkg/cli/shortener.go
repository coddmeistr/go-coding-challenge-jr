package cli

import (
	"challenge/pkg/proto"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	rootCmd.AddCommand(shortenerCommand)
	shortenerCommand.Flags().StringVarP(&url, "url", "u", "", "url that'll be shortened")
}

var url string
var shortenerCommand = &cobra.Command{
	Use:   "shortener",
	Short: "Shorten link",
	Long:  `gRPC call that'll shorten given link via Bilty API'`,
	Run: func(_ *cobra.Command, _ []string) {

		if url == "" {
			fmt.Println("url wasn't provided")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 20)
		defer cancel()
		conn, err := grpc.DialContext(ctx, address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			fmt.Printf("cannot connect to gRPC server: %v\n", err)
			return
		}

		client := proto.NewChallengeServiceClient(conn)
		shortened, err := client.MakeShortLink(context.Background(), &proto.Link{Data: url})
		if err != nil {
			fmt.Printf("cannot shorten link: %v\n", err)
			return
		}

		fmt.Printf("shortened link: %s\n", shortened.GetData())
	},
}
