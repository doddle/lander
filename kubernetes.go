package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
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

func getEndpoints(logger *log.Logger) []string {

	cacheObj := "endpoints"
	fakeResult := []string{}

	cached, found := cacheShort.Get(cacheObj)

	if found {
		logger.Debug("getEndpoints serving from cache")
		return cached.([]string)
	} else {
		ingressList, err := getIngressList(logger)
		if err != nil {
			logger.Error(err)
		}
		ingressObjects := ingressList.Items
		time.Sleep(5 * time.Second)
		if len(ingressObjects) > 0 {
			for _, ingress := range ingressObjects {
				for _, rule := range ingress.Spec.Rules {
					for _, p := range rule.IngressRuleValue.HTTP.Paths {
						uri := p.Path
						if p.Path == "/" {
							uri = ""
						}
						msg := fmt.Sprintf("https://%s%s", rule.Host, uri)
						fakeResult = append(fakeResult, msg)
						logger.Debug(msg)

					}
				}
				// for k, v := range ingress.Annotations {
				// 	logger.Debugf("%s = %s", k, v)
				// }
			}
		}
		cacheShort.Set(cacheObj, fakeResult, cache.DefaultExpiration)
		return fakeResult
	}
}

// Speaks to the cluster and attempt to pull an IngressList
func getIngressList(logger *log.Logger) (*v1beta1.IngressList, error) {

	config := autoClientInit(logger)
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	ingressList, err := clientset.ExtensionsV1beta1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingressList, err
}
