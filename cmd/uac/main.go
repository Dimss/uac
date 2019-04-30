// https://github.com/Holmes89/hex-example/blob/hex/database/psql/ticket.go
// https://github.com/thockin/go-build-template
// https://www.youtube.com/watch?v=VQym87o91f8
// https://banzaicloud.com/blog/k8s-admission-webhooks/
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
	"path/filepath"
	"strings"
)

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

var dumpRuntimeConfigCmd = &cobra.Command{
	Use:   "dumpconfig",
	Short: "Dump all runtime configs",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Dumping runtime configs")
		logrus.Infof("http.crt: %s", viper.GetString("http.crt"))
		logrus.Infof("http.key: %s", viper.GetString("http.key"))
		logrus.Infof("ad.host: %s", viper.GetString("ad.host"))
		logrus.Infof("ad.port: %s", viper.GetString("ad.port"))
		logrus.Infof("ad.baseDN: %s", viper.GetString("ad.baseDN"))
		logrus.Infof("ad.bindUser: %s", viper.GetString("ad.bindUser"))
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringP("kubeconfig", "k", "", "Path to kubeconfig file, default to $home/.kube/config")
	rootCmd.PersistentFlags().StringP("configpath", "c", "", "Path to config directory with config.json file, default to . ")
	userSyncCmd.PersistentFlags().StringP("user", "u", "", "AD username for sync")
	if err := viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
		panic(err)
	}
	if err := viper.BindPFlag("configpath", rootCmd.PersistentFlags().Lookup("configpath")); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(webhookCmd)
	rootCmd.AddCommand(userSyncCmd)
	rootCmd.AddCommand(dumpRuntimeConfigCmd)
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

}

func initConfig() {
	viper.SetConfigType("json")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetEnvPrefix("UAC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	configPath := viper.GetString("configpath")
	logrus.Info(configPath)
	// If config flag is empty, assume config.json located in current directory
	if configPath != "" {
		viper.AddConfigPath(configPath)
	}
	// look for kubeconfig file, if not found, assume running inside OPC cluster
	kubeconfig := viper.GetString("kubeconfig")
	if kubeconfig == "" {
		// Check if kubeconfig file exists in user's HOME
		kubeconfig = filepath.Join(os.Getenv("HOME"), ".kube", "config")
		_, err := os.Stat(kubeconfig)
		if os.IsNotExist(err) {
			// The kubeconfig wasn't passed in and not found under user's home directory, assuming inClusterConfig mode
			logrus.Info("Unable to find kubeconfig, assuming running inside K8S cluster, gonna use inClusterConfig")
			viper.Set("kubeconfig", "useInClusterConfig")
		} else {
			// Use kubeconfig from user's home directory
			logrus.Info("Gonna use kubeconfig from user's HOME directory")
			viper.Set("kubeconfig", kubeconfig)
		}
	}
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Errorf("Unable to read config.json file, err: %s", err)
		os.Exit(1)
	}
}

func main() {

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
