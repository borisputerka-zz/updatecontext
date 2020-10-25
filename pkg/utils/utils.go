package utils

import (
	"bufio"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)

// AskForConfirmation function that return 1 if you type y/Y/yes or 0 otherwise
func AskForConfirmation() (bool, error) {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}

	switch strings.ToLower(response) {
	case "y\n", "yes\n", "\n":
		return true, nil
	}
	return false, nil
}

func GetKubernetesConfig(local bool, masterURL string) (config *restclient.Config, err error) {
	if local {
		configAccess := clientcmd.NewDefaultPathOptions()
		config, err = clientcmd.BuildConfigFromKubeconfigGetter(masterURL, configAccess.GetStartingConfig)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		config, err = restclient.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}
	return config, nil
}
