package plugin

import (
	"fmt"
	"strings"

	"github.com/borisputerka/updatecontext/pkg/logger"
	"github.com/borisputerka/updatecontext/pkg/utils"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd/api"
)

type ConfigFlags struct {
	Config *utils.Config
}

func (o *ConfigFlags) Complete() (err error) {
	o.Config, err = utils.NewConfig()
	if err != nil {
		return fmt.Errorf("failed to get config: %v", err)
	}
	return nil
}

// RunPlugin function generated contexts
func (o *ConfigFlags) RunPlugin() error {
	logger := logger.NewLogger()

	namespaces, err := o.listNamespaces()
	if err != nil {
		return errors.Wrap(err, "failed to list namespaces")
	}

	currentContexts, cluster := o.Config.ListContexts()

	var createdContexts []string
	for _, namespace := range namespaces {
		contextName := fmt.Sprintf("%s/%s", namespace, cluster)
		if _, ok := currentContexts[contextName]; !ok {
			o.Config.AddContext(cluster, namespace)
			createdContexts = append(createdContexts, contextName)
		}
		delete(currentContexts, contextName) // Non-existent context will be  in allContexts after all iterations
	}

	err = o.Config.Update()
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
		err := o.deleteContexts(currentContexts)
		if err != nil {
			return errors.Wrap(err, "could not delete contexts")
		}
	}

	return nil
}

func (o *ConfigFlags) listNamespaces() ([]string, error) {
	client, err := o.Config.GetKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("could get kubernetes client: %v", err)
	}

	namespaces, err := client.CoreV1().Namespaces().List(metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("could not get list of namespaces: %v", err)
	}

	var namespaceList []string
	excludedNsList := []string{"kube-system", "kube-public", "kube-node-lease"}
	for _, ns := range namespaces.Items {
		if !utils.StringInSlice(ns.Name, excludedNsList) {
			namespaceList = append(namespaceList, ns.Name)
		}
	}

	return namespaceList, nil
}

func (o *ConfigFlags) deleteContexts(contexts map[string]*api.Context) error {
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
		err := o.Config.DeleteContexts(contexts)
		if err != nil {
			return err
		}
		logger.Info("Contexts deleted successfully")
	}
	return nil
}
