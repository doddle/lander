package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	"github.com/withmandala/go-log"
	//"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
)

// returns true if PLUGIN_DEBUG!=""
func newLogger(debug bool) *log.Logger {
	// check if debug enabled
	if debug {
		logger := log.New(os.Stdout).WithDebug().WithColor()
		logger.Debug("debug enabled")
		return logger
	} else {
		return log.New(os.Stdout).WithColor()
	}
}

func envVarExists(key string) bool {
	_, exists := os.LookupEnv(key)
	if exists {
		return true
	}
	return false
}

func autoClientInit(logger *log.Logger) *rest.Config {
	// we will automatically decide if this is running inside the cluster or on someones laptop
	// if the ENV vars KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT exist
	// then we can assume this app is running inside a k8s cluster
	if envVarExists("KUBERNETES_SERVICE_HOST") && envVarExists("KUBERNETES_SERVICE_PORT") {
		logger.Info("running in-cluster")
		config, err := rest.InClusterConfig()
		if err != nil {
			logger.Error(err)
		}
		return config
	}

	kubeconfig := findKubeConfig()
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		logger.Error(err)
	}
	return config
}

// does a best guess to the location of your kubeconfig
func findKubeConfig() string {
	home := homedir.HomeDir()
	if envVarExists("KUBECONFIG") {
		kubeconfig := os.Getenv("KUBECONFIG")
		return kubeconfig
	} else {
		kubeconfig := fmt.Sprintf(filepath.Join(home, ".kube", "config"))
		return kubeconfig
	}
}

func main() {

	logger := newLogger(true)

	config := autoClientInit(logger)

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Error(err)
	}


	for {
		ingressList, err := clientset.ExtensionsV1beta1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			// handle err
		}
		ingressCtrls := ingressList.Items
		if len(ingressCtrls) > 0 {
			for _, ingress := range ingressCtrls {
				logger.Infof("ingress %s exists in namespace %s\n", ingress.Name, ingress.Namespace)
			}
		} else {
			logger.Info("no ingress found")
		}
		time.Sleep(10 * time.Second)
	}

}
