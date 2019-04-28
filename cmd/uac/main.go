// https://github.com/Holmes89/hex-example/blob/hex/database/psql/ticket.go
// https://github.com/thockin/go-build-template
// https://www.youtube.com/watch?v=VQym87o91f8
package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uac/cmd/uac/webhookserver"
	"github.com/uac/pkg/activedirectory"
	"github.com/uac/pkg/k8sclient"
	"os"
)

func init() {
	// Read JSON configuration file
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

var rootCmd = &cobra.Command{
	Use:   "uac",
	Short: "BNHP user access controller and permission sync manager between OCP clusters and AD",
}

var webhookCmd = &cobra.Command{
	Use:   "server",
	Short: "Start HTTP server for processing OCP OAuthaccessTokens dynamic admission webhooks",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting up webhook server...")
		webhookserver.StartHttpRouter()
	},
}

var userSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Synchronize user permission",
	Run: func(cmd *cobra.Command, args []string) {
		adUser, err := cmd.Flags().GetString("user")
		if err != nil {
			panic(err)
		}
		if adUser == "" {
			logrus.Error("Empty username, can't proceed!")
			os.Exit(1)
		}
		logrus.Infof("Gonna sync user: %s", adUser)
		userGroups := activedirectory.GetUsersGroups(adUser)
		k8sclient.SetUserRbac(userGroups, adUser)
	},
}

func main() {
	userSyncCmd.PersistentFlags().StringP("user", "u", "", "AD username for sync")
	rootCmd.AddCommand(webhookCmd)
	rootCmd.AddCommand(userSyncCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}