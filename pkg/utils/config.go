package utils

import (
	"fmt"
	"regexp"

	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

type Config struct {
	cmdConfig    *api.Config
	configAccess *clientcmd.PathOptions
	restConfig   *rest.Config
}

func (c *Config) GetKubernetesClient() (*kubernetes.Clientset, error) {
	clientSet, err := kubernetes.NewForConfig(c.restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	return clientSet, nil
}

func (c *Config) Update() (err error) {
	return clientcmd.ModifyConfig(c.configAccess, *c.cmdConfig, true)
}

func (c *Config) ContextName(cluster string, namespace string) string {
	gkeRegexp := "^gke_([^_\\n]*)_([^_\\n]*)_([^_\\n]*)$"
	eksRegexp := "^arn:([^:\\n]*):([^:\\n]*):([^:\\n]*):([^:\\n]*):(([^:\\/\\n]*)[:\\/])?(.*)$"
	if matched, _ := regexp.MatchString(eksRegexp, cluster); matched {
		eksRe := regexp.MustCompile(eksRegexp)
		eksClusterSubmatch := eksRe.FindStringSubmatch(cluster)
		clusterName := eksClusterSubmatch[len(eksClusterSubmatch)-1]
		return fmt.Sprintf("%s/%s", namespace, clusterName)

	} else if matched, _ := regexp.MatchString(gkeRegexp, cluster); matched {
		gkeRe := regexp.MustCompile(gkeRegexp)
		gkeClusterSubmatch := gkeRe.FindStringSubmatch(cluster)
		clusterName := gkeClusterSubmatch[len(gkeClusterSubmatch)-1]
		return fmt.Sprintf("%s/%s", namespace, clusterName)
	} else {
		return fmt.Sprintf("%s/%s", namespace, cluster)
	}
}

func (c *Config) DeleteContexts(contexts map[string]*api.Context) error {
	for name := range contexts {
		delete(c.cmdConfig.Contexts, name)
	}
	return c.Update()
}

func (c *Config) AddContext(cluster string, namespace string) {
	newContext := *api.NewContext()
	newContext.Cluster = cluster
	newContext.Namespace = namespace
	newContext.AuthInfo = cluster

	contextName := c.ContextName(cluster, namespace)
	c.cmdConfig.Contexts[contextName] = &newContext
	c.cmdConfig.CurrentContext = contextName
}

func (c *Config) ListContexts() (map[string]*api.Context, string) {
	contexts := map[string]*api.Context{}
	currentContext := c.cmdConfig.Contexts[c.cmdConfig.CurrentContext]
	cluster := currentContext.Cluster
	for name, ctx := range c.cmdConfig.Contexts {
		if ctx.Cluster == cluster {
			contexts[name] = ctx
		}
	}
	return contexts, cluster
}

func NewConfig() (config *Config, err error) {
	configAccess := clientcmd.NewDefaultPathOptions()
	cmdConfig, err := configAccess.GetStartingConfig()
	if err != nil {
		return nil, err
	}

	configFlags := genericclioptions.NewConfigFlags(true)
	restConfig, err := configFlags.ToRESTConfig()
	if err != nil {
		return nil, err
	}

	return &Config{
		cmdConfig:    cmdConfig,
		configAccess: configAccess,
		restConfig:   restConfig,
	}, nil
}
