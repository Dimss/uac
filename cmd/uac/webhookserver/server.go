package webhookserver

import (
	"crypto/tls"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/uac/pkg/activedirectory"
	"github.com/uac/pkg/k8sclient"
	"github.com/uac/pkg/oauthtokenwebhook"
	"log"
	"net/http"
)

func StartHttpRouter() {
	cert := viper.GetString("http.crt")
	key := viper.GetString("http.key")
	pair, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		logrus.Error("Failed to load key pair: %v", err)
	}
	// Buffered channel for AD users
	adUsersChan := make(chan string, 100)
	// Watch and process ad users when they pushed to adUsersChan channel
	go syncUsers(adUsersChan)
	// Handel admission webhook request
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		oauthtokenwebhook.WebHookHandler(w, r, adUsersChan)
	})
	// Create HTTPS server configuration
	s := &http.Server{
		Addr:      ":8080",
		TLSConfig: &tls.Config{Certificates: []tls.Certificate{pair}},
	}
	// Start HTTPS server
	log.Fatal(s.ListenAndServeTLS("", ""))
}

func syncUsers(adUsersChan chan string) {
	logrus.Info("Waiting for new ad users to process")
	for adUser := range adUsersChan {
		// Get parsed user's group from AD
		userGroups := activedirectory.GetUsersGroups(adUser)
		k8sclient.SetUserRbac(userGroups, adUser)
	}
}