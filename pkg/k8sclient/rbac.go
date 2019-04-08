package k8sclient

import (
	"fmt"
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/sirupsen/logrus"
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

func SetUserRbac(userGroups []string, userName string) {
	ocpProjects := getExistingProjects()
	logrus.Infof("OCP Projects: %s", ocpProjects)


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
