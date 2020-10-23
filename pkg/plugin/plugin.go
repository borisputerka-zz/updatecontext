package plugin

import (
	"fmt"
	"strings"

	"github.com/borisputerka/updatecontext/pkg/logger"
	"github.com/borisputerka/updatecontext/pkg/utils"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func RunPlugin(configFlags *genericclioptions.ConfigFlags) error {
	logger := logger.NewLogger()

	configAccess := clientcmd.NewDefaultPathOptions()

	cmdConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		return errors.Wrap(err, "failed to get cmdConfig")
	}

	config, err := configFlags.ToRESTConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read kubeconfig")
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return errors.Wrap(err, "failed to create clientset")
	}

	namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to list namespaces")
	}

	currContext := cmdConfig.Contexts[cmdConfig.CurrentContext]
	cluster := currContext.Cluster
	currentContexts := listContexts(cmdConfig, cluster)

	var createdContexts []string
	for _, namespace := range namespaces.Items {
		contextName := fmt.Sprintf("%s/%s", namespace.Name, cluster)
		if _, ok := currentContexts[contextName]; !ok {
			addContext(cmdConfig, cluster, namespace.Name)
			createdContexts = append(createdContexts, contextName)
		}
		delete(currentContexts, contextName) // Non-existent context will be  in allContexts after all iterations
	}

	err = clientcmd.ModifyConfig(configAccess, *cmdConfig, true)
	if err != nil {
		return errors.Wrap(err, "failed to modify configAccess")
	}

	if len(createdContexts) > 0 {
		logger.Info("Contexts created: \n")
		logger.Info(strings.Join(createdContexts, "\n"))
	}

	if len(createdContexts) == 0 {
		logger.Info("Nothing to create \n")
	}

	if len(currentContexts) > 0 {
		err := deleteContexts(configAccess, cmdConfig, currentContexts)
		if err != nil {
			return errors.Wrap(err, "could not delete contexts")
		}
	}

	return nil
}

func listContexts(cmdConfig *api.Config, cluster string) map[string]*api.Context {
	contexts := map[string]*api.Context{}
	for name, ctx := range cmdConfig.Contexts {
		if ctx.Cluster == cluster {
			contexts[name] = ctx
		}
	}

	return contexts
}

func addContext(cmdConfig *api.Config, cluster string, namespace string) {
	newContext := *api.NewContext()
	newContext.Cluster = cluster
	newContext.Namespace = namespace
	newContext.AuthInfo = cluster

	contextName := fmt.Sprintf("%s/%s", namespace, cluster)
	cmdConfig.Contexts[contextName] = &newContext
	cmdConfig.CurrentContext = contextName
}

func deleteContexts(configAccess *clientcmd.PathOptions, cmdConfig *api.Config, contexts map[string]*api.Context) error {
	logger := logger.NewLogger()
	contextNames := ""
	for name := range contexts {
		contextNames += name + "\n"
	}

	logger.Info("Following contexts are not used anymore:")
	logger.Info(contextNames)
	logger.Info("Do you want to delete them? [Y/n]: ")
	confirmed, err := utils.AskForConfirmation()
	if err != nil {
		return err
	}

	if confirmed {
		for name := range contexts {
			delete(cmdConfig.Contexts, name)
		}

		err = clientcmd.ModifyConfig(configAccess, *cmdConfig, true)
		if err != nil {
			return err
		}

		logger.Info("Contexts deleted successfully")
	}
	return nil
}
