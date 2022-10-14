package cmd

import (
	"chat/client"

	"github.com/spf13/cobra"
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "run the client",
	Run: func(cmd *cobra.Command, args []string) {
		channelName, _ := cmd.Flags().GetString("channel")
		sendersName, _ := cmd.Flags().GetString("sender")
		interval, _ := cmd.Flags().GetInt("interval")
		client.Run(channelName, sendersName, interval)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.PersistentFlags().String("channel", "default", "specify which channel to connect")
	clientCmd.PersistentFlags().String("sender", "default", "specify the sendersName")
	clientCmd.PersistentFlags().Int16("interval", 0, "specify the interval to resend the message, 0 is a one-time message")
}
