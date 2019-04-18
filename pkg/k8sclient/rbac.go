package k8sclient

import (
	projectv1 "github.com/openshift/client-go/project/clientset/versioned/typed/project/v1"
	"github.com/sirupsen/logrus"
	"github.com/uac/pkg/activedirectory"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacV1 "k8s.io/client-go/kubernetes/typed/rbac/v1"
	apirbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/rest"
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
		return
	} else {
		logrus.Infof("User %s should have access to projects: %s", userName, userOcpProjects)
	}
	getAdminRoleBinding(userOcpProjects, userName)
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

func getExistingProjects() (projectsNames []string) {
	projectv1Client, err := projectv1.NewForConfig(getClientcmdConfigs())
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

func getAdminRoleBinding(userOcpProjects []string, userName string) {
	logrus.Info("Getting roles")
	rbacV1Client, err := rbacV1.NewForConfig(getClientcmdConfigs())
	if err != nil {
		panic(err.Error())
	}

	for _, ocpProject := range userOcpProjects {
		roleBinding, err := rbacV1Client.RoleBindings(ocpProject).Get("admin", metav1.GetOptions{})
		if err != nil {
			panic(err.Error())
		}
		// Create user Subject
		subject := apirbacv1.Subject{
			Kind:      "User",
			APIGroup:  "rbac.authorization.k8s.io",
			Name:      userName,
			Namespace: "",
		}
		// Append user Subject to RoleBinding subject slice
		roleBinding.Subjects = append(roleBinding.Subjects, subject)
		// Update RoleBinding
		roleBinding, err = rbacV1Client.RoleBindings(ocpProject).Update(roleBinding)
		if err != nil {
			panic(err.Error())
		}
		logrus.Infof("RoleBinding %s  in namespace %s for user %s was successfully updated",
			roleBinding.Name, ocpProject, userName)
	}

}

func getClientcmdConfigs() (config *rest.Config) {
	conf := "/Users/dima/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", conf)
	if err != nil {
		panic(err.Error())
	}
	return
}
