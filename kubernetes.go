package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// takes a slice of endpoint and only returns ones containing a hostname
// excluded is a (optionaly comma seperated) list of endpoints to exclude (EG this landing page itself)
func onlyHostnamesContaining(input []Endpoint, host string) []Endpoint {
	var result []Endpoint
	for _, data := range input {
		if strings.Contains(data.Address, host) {
			result = append(result, data)
		}
	}
	return result
}

func getIngressEndpoints(logger *log.Logger) []Endpoint {
	cacheObj := "endpoints"
	var result []Endpoint

	cached, found := cacheShort.Get(cacheObj)
	if found {
		logger.Debugf("got all %s from cache", cacheObj)
		data := cached.([]Endpoint)
		logger.Info(data)
		return data
	}
	ingressList, err := getIngressList(logger)
	if err != nil {
		logger.Error(err)
	}
	ingressObjects := ingressList.Items
	// time.Sleep(5 * time.Second)
	if len(ingressObjects) > 0 {
		for _, ingress := range ingressObjects {
			for _, rule := range ingress.Spec.Rules {
				for _, p := range rule.IngressRuleValue.HTTP.Paths {
					serviceName := p.Backend.ServiceName
					guessed := guessApp(serviceName)
					// Strip out a trailing "/"
					uri := p.Path
					if p.Path == "/" {
						uri = ""
					}
					msg := Endpoint{
						Address:     "https://" + rule.Host + uri,
						Https:       true,
						Oauth2proxy: getOauth2ProxyState(logger, ingress),
						Class:       getIngressClass(logger, ingress),
						Icon:        guessed.Icon,
						Description: guessed.Desc,
						Name:        guessed.Name,
					}
					result = append(result, msg)
				}
			}
		}
		cacheShort.Set(cacheObj, result, cache.DefaultExpiration)
	}
	return result
}

// App is a generic definition of a known, we'll use these to attempt to guess the apps
// TODO: add some tags, slice of common service names and import the data from json/yaml maybe
type App struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
	Desc string `yaml:"desc"`
}

func genApps() (fallback App, index []App) {
	var result []App
	fallback = App{
		Name: "unknown",
		Icon: "/assets/link.png",
		Desc: "generic service",
	}

	x := App{
		Name: "grafana",
		Icon: "/assets/grafana.png",
		Desc: "View and create dashboards for prometheus metric data, can also view+stream logs",
	}
	result = append(result, x)

	x = App{
		Name: "prometheus",
		Icon: "/assets/prometheus.png",
		Desc: "Explore prometheus Alerts, AlertRules, Service discovery and run raw queries",
	}
	result = append(result, x)

	x = App{
		Name: "alertmanager",
		Icon: "/assets/alertmanager.png",
		Desc: "manage alerts and silences",
	}
	result = append(result, x)

	x = App{
		Name: "kibana",
		Icon: "/assets/kibana.png",
		Desc: "aggregate and explore log data+graphs",
	}
	result = append(result, x)
	return fallback, result
}

func guessApp(svc string) App {
	fallback, apps := genApps()
	for _, x := range apps {
		if strings.Contains(svc, x.Name) {
			return x
		}
	}
	return fallback
}

// returns true/false if ingress Annotations contain what looks like oa2p
func getOauth2ProxyState(logger *log.Logger, ingress v1beta1.Ingress) bool {
	if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-signin") {
		if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-url") {
			return true
		}
	}
	return false
}

// check if a key exists in an ingress annotation
func annotationKeyExists(ingress v1beta1.Ingress, key string) bool {
	for k := range ingress.Annotations {
		if k == key {
			return true
		}
	}
	return false
}

// attempts to return the ingress class (or an empty string)
func getIngressClass(logger *log.Logger, ingress v1beta1.Ingress) string {
	for k, v := range ingress.Annotations {
		if k == "kubernetes.io/ingress.class" {
			return v
		}
	}
	return ""
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
