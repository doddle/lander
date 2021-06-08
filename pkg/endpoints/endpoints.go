package endpoints

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
	"time"
)

var (
	// hard limit cache for 15sec, expire at 15m
	pkgCache = cache.New(15*time.Second, 15*time.Minute)
)

func ReallyAssemble(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	hostName string,
) []Endpoint {
	var result []Endpoint
	allEndpoints := gatherEndpointData(logger, clientSet)

	// lets filter them for only ones matching the hostname of the context
	containsHostname := filterIngressForHostname(allEndpoints, hostName)

	for _, i := range containsHostname {
		// exclude hostnames that are not identical to flagHost
		// input here will be like: https://example.com, http://example.com/foo
		split := strings.Split(i.Address, "/")

		// len 3 includes only addresses with something after the "/" .. eg:
		// - example.com         <- not match
		// - example.com/foo     <- would match
		// this should exclude the lander itself (designed to be on https://cluster.example.com/)
		if len(split) > 3 {
			got := split[2]
			want := hostName
			if strings.Compare(got, want) == 0 {
				result = append(result, i)
			}
		}
	}
	return result
}

func gatherEndpointData(logger *log.Logger, clientSet *kubernetes.Clientset) []Endpoint {
	var result []Endpoint

	ingressList, err := getIngressList(logger, clientSet)
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
						Oauth2proxy: getOauth2ProxyState(ingress),
						Class:       getIngressClass(logger, ingress),
						Icon:        guessed.Icon,
						Description: guessed.Desc,
						Name:        guessed.Name,
					}
					result = append(result, msg)
				}
			}
		}
	}
	return result
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
// TODO: upgrade to v1?
func getIngressClass(logger *log.Logger, ingress v1beta1.Ingress) string {
	for k, v := range ingress.Annotations {
		if k == "kubernetes.io/ingress.class" {
			return v
		}
	}
	logger.Warnf(
		"Unable to determine ingress class for: %s/%s",
		ingress.Namespace,
		ingress.Name)
	return ""
}

// Speaks to the cluster and attempt to pull an IngressListgg
func getIngressList(logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*v1beta1.IngressList, error) {
	cacheObj := "v1beta/ingress"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*v1beta1.IngressList), nil
	}

	// find ALL Ingressess
	ingressList, err := clientSet.
		ExtensionsV1beta1().
		Ingresses("").
		List(
			context.TODO(),
			v1.ListOptions{},
		)
	if err != nil {
		return nil, err
	}
	logger.Debugf("got all %s from k8s", cacheObj)
	pkgCache.Set(cacheObj, ingressList, cache.DefaultExpiration)
	return ingressList, err
}

// TODO: read from a json file or similar
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

// returns true/false if ingress Annotations contain what looks like oa2p
func getOauth2ProxyState(ingress v1beta1.Ingress) bool {
	if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-signin") {
		if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-url") {
			return true
		}
	}
	return false
}

func filterIngressForHostname(input []Endpoint, host string) []Endpoint {
	var result []Endpoint
	for _, data := range input {
		if strings.Contains(data.Address, host) {
			result = append(result, data)
		}
	}
	return result
}
