/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello Lyon")
		conf := &config{
			secure: viper.GetBool("secure"),
			port:   viper.GetInt32("port"),
			cert:   viper.GetString("cert"),
			key:    viper.GetString("key"),
		}
		switch cmd.Flag("type").Value.String() {
		case "grpc":
			serveGrpc(conf)
		case "html":
			serveHtml(conf)
		case "rest":
			serveRest(conf)
		case "all":
			serveAll(conf)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	serveCmd.PersistentFlags().BoolP("secure", "s", false, "Secure mode (default: false)")
	viper.BindPFlag("secure", serveCmd.PersistentFlags().Lookup("secure"))
	serveCmd.PersistentFlags().String("cert", "", "Certificate")
	viper.BindPFlag("cert", serveCmd.PersistentFlags().Lookup("cert"))
	serveCmd.PersistentFlags().String("key", "", "Key")
	viper.BindPFlag("key", serveCmd.PersistentFlags().Lookup("key"))

	serveCmd.MarkFlagsRequiredTogether("secure", "cert", "key")

	serveCmd.PersistentFlags().Int16P("port", "p", 4000, "Server port")
	viper.BindPFlag("port", serveCmd.PersistentFlags().Lookup("port"))

	serveCmd.PersistentFlags().StringP("type", "t", "grpc", "Serve type")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
