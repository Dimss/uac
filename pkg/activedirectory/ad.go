package activedirectory

import (
	"fmt"
	oauthv1 "github.com/openshift/api/oauth/v1"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"html"
	"net/http"

	//"github.com/uac/pkg/k8sclient"
	"gopkg.in/ldap.v3"
	"os"
	"strings"
)

func init() {
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
	// Read JSON configuration file
	viper.SetConfigFile("config.json")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}


func RegisterHandler() {
	http.HandleFunc("/bar1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello this is ad, %q", html.EscapeString(r.URL.Path))
	})
}


func SyncUserPermissions(oauthToken oauthv1.OAuthAccessToken) {
	logrus.Infof("Username: %s", oauthToken.UserName)
	userGroups := parseUserAdGroups(getUserADGroups(oauthToken.UserName))
	logrus.Infof("Users AD membership: %s", userGroups)
	//k8sclient.SetUserRbac(userGroups, oauthToken.UserName)
}

func getUserADGroups(user string) (userAdGroups []string) {
	adHost := viper.GetString("ad.host")
	adPort := viper.GetInt("ad.port")
	bindUser := viper.GetString("ad.bindUser")
	bindPass := viper.GetString("ad.bindPass")
	// Connect to AD
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", adHost, adPort))
	if err != nil {
		logrus.Fatal(err)
	}
	defer l.Close()
	// Bind with system user
	err = l.Bind(bindUser, bindPass)
	if err != nil {
		logrus.Fatal(err)
	}
	// Execute search
	sr, err := l.Search(getSearchRequest(user))
	if err != nil {
		logrus.Fatal(err)
	}
	// Parse search results
	if len(sr.Entries) > 0 {
		logrus.Warning("Getting user's AD groups")
		ldapEntry := sr.Entries[0]
		if len(ldapEntry.Attributes) > 0 {
			entryAttributes := ldapEntry.Attributes[0]
			for _, adGroup := range entryAttributes.Values {
				logrus.Infof("%s member of %s", user, adGroup)
				userAdGroups = append(userAdGroups, adGroup)
			}
			return
		} else {
			logrus.Warning("Empty user group list")
		}
	}
	return
}

func getSearchRequest(user string) (searchRequest *ldap.SearchRequest) {
	adBaseDn := viper.GetString("ad.baseDN")
	adQuery := fmt.Sprintf("(&(objectClass=user)(sAMAccountName=%s)(memberOf=*))", user)
	return ldap.NewSearchRequest(
		adBaseDn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		adQuery,
		[]string{"memberOf"},
		nil,
	)
}

func parseUserAdGroups(userGroups []string) (parsedUserGroups []UserGroups) {
	//parsedUserAdGroup := make(map[string]string)
	for _, userGroup := range userGroups {
		// split AD group string to slice by ','
		// example: 'CN=ocpns__capital-market,DC=ad,DC=lab'
		// becomes: ['CN=ocpns__capital-market','DC=ad','DC=lab']
		groupDn := strings.Split(userGroup, ",")
		if len(groupDn) > 0 {
			// First element in the slice contain AD group name
			groupName := strings.Split(groupDn[0], "=")
			// Split the group name by '='
			// Example: CN=ocpns__capital-market
			// becomes: ['CN','ocpns__capital-market']
			if len(groupName) == 2 {
				// Append AD group to result string array
				adGroupName := groupName[1]
				ocpNsName := parseAdGroupNameToOcpNsName(adGroupName)
				parsedUserGroups = append(parsedUserGroups, UserGroups{adGroupName, ocpNsName})
			} else {
				logrus.Warnf("Unexpected user group name %s", groupName)
			}
		} else {
			logrus.Warnf("Unexpected user group DN %s", groupDn)
		}
		fmt.Println(groupDn)
	}
	return
}

func parseAdGroupNameToOcpNsName(adGroupName string) (ocpNsName string) {
	// Get AD group name and parse it to OCP Project name
	// split AD group by '__' and return element at index 1
	// Example: ocpns__capital-market
	// becomes: ['ocpns','capital-market']
	ocpNs := strings.Split(adGroupName, "__")
	if len(ocpNs) == 2 {
		ocpNsName = ocpNs[1]
	} else {
		ocpNsName = ""
	}
	return
}
