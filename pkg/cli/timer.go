package cli

import (
	"challenge/pkg/proto"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

func init() {
	rootCmd.AddCommand(startTimerCommand)
	startTimerCommand.Flags().StringVarP(&name, "name", "n", "", "name of the timer")
	startTimerCommand.Flags().IntVarP(&freq, "freq", "f", 0, "frequency of the timer")
	startTimerCommand.Flags().IntVarP(&secs, "secs", "s", 0, "seconds of the timer")
}

var name string
var freq int
var secs int
var startTimerCommand = &cobra.Command{
	Use:   "timer",
	Short: "Start timer",
	Long:  `gRPC call that'll start new timer via timercheck.io and will provide updates with gRPC stream'`,
	Run: func(_ *cobra.Command, _ []string) {

		if name == "" {
			fmt.Println("timer name is empty")
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
		stream, err := client.StartTimer(context.Background(), &proto.Timer{Name: name, Frequency: int64(freq), Seconds: int64(secs)})
		if err != nil {
			fmt.Printf("cannot create or connect to timer: %v\n", err)
			return
		}

		for {
			ping, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("stream closed\n")
				return
			}
			if err != nil {
				fmt.Printf("stream error: %v\n", err)
				return
			}

			fmt.Printf("timer name: %s\n", ping.GetName())
			fmt.Printf("timer seconds left: %d\n", ping.GetSeconds())
			fmt.Printf("timer frequency: %d\n", ping.GetFrequency())
		}
	},
}
