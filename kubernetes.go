package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/withmandala/go-log"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func autoClientInit(logger *log.Logger) *rest.Config {
	// we will automatically decide if this is running inside the cluster or on someones laptop
	// if the ENV vars KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT exist
	// then we can assume this app is running inside a k8s cluster
	if envVarExists("KUBERNETES_SERVICE_HOST") && envVarExists("KUBERNETES_SERVICE_PORT") {
		config, err := rest.InClusterConfig()
		if err != nil {
			logger.Error(err)
		}
		return config
	}

	config, err := clientcmd.BuildConfigFromFlags("", findKubeConfig())
	if err != nil {
		logger.Fatal(err)
	}
	return config
}

// does a best guess to the location of your kubeconfig
func findKubeConfig() string {
	home := homedir.HomeDir()
	if envVarExists("KUBECONFIG") {
		kubeconfig := os.Getenv("KUBECONFIG")
		return kubeconfig
	}
	kubeconfig := fmt.Sprint(filepath.Join(home, ".kube", "config"))
	return kubeconfig
}
