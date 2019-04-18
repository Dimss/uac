package k8sclient

import (
	"fmt"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/sirupsen/logrus"
	"github.com/uac/pkg/activedirectory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func init() {
	// Init log
	logrus.SetOutput(os.Stdout)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func SetUserRbac(userGroups []activedirectory.UserGroups, userName string) {
	// Get all OCP projects
	ocpProjects := getExistingProjects()
	logrus.Infof("OCP Projects: %s", ocpProjects)
	// Find to which OCP projects user should have access
	userOcpProjects := matchUserGroupsToOcpNs(userGroups, ocpProjects)
	// If no patch is found, inform in log
	if len(userOcpProjects) == 0 {
		logrus.Errorf("No projects matched for user %s in OCP cluster. User groups: %s", userName, userGroups)

	} else {
		logrus.Infof("User %s should have access to %s", userName, userOcpProjects)
	}

}

func matchUserGroupsToOcpNs(userGroups []activedirectory.UserGroups, ocpProjects []string) (userOcpProjects []string) {
	for _, userGroup := range userGroups {
		for _, ocpProject := range ocpProjects {
			if userGroup.OcpNs == ocpProject {
				userOcpProjects = append(userOcpProjects, ocpProject)
			}
		}
	}
	return
}

func GetRunningPods() {
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))
}

func getExistingProjects() (projectsNames []string) {
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err.Error())
	}

	projectv1Client, err := projectv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	projects, err := projectv1Client.Projects().List(metav1.ListOptions{})
	logrus.Infof("Total projects in OCP cluster: %d", len(projects.Items))
	for _, project := range projects.Items {
		projectsNames = append(projectsNames, project.Name)
	}
	return
}
