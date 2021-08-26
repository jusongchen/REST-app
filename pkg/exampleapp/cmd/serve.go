package cmd

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/jusongchen/REST-app/pkg/exampleapp"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "serve command starts demoapp server",
	Long:  `serve starts an demoapp serve to provide DB stats service`,
	Run: func(cmd *cobra.Command, args []string) {
		exampleapp.PrintBanner()

		myapp, err := exampleapp.New()
		if err != nil {
			log.Errorf("app config error:%v", err)
			return
		}

		myapp.Run(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
