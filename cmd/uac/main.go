// https://github.com/Holmes89/hex-example/blob/hex/database/psql/ticket.go
// https://github.com/thockin/go-build-template
// https://www.youtube.com/watch?v=VQym87o91f8
package main

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	oauthTokenWebHook "github.com/uac/pkg/oauthtokenwebhook"
	"net/http"
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
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
}

func startHttpRouter() {
	cert := viper.GetString("http.cert")
	key := viper.GetString("http.key")
	pair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Error("Failed to load key pair: %v", err)
	}
	// Register webhook handler
	http.HandleFunc("/", oauthTokenWebHook.WebHookHandler)
	// Create HTTPS server configuration
	s := &http.Server{
		Addr:      ":8080",
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	}
	// Start HTTPS server
	log.Fatal(s.ListenAndServeTLS("", ""))
}

func main() {
	log.Info("Starting up...")
	startHttpRouter()

}
