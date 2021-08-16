package endpoints

import (
	"context"
	"github.com/digtux/lander/pkg/util"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// hard limit cache for 15sec, expire at 15m
	pkgCache = cache.New(15*time.Second, 15*time.Minute)
)

func ReallyAssemble(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	landerAnnotation string,
) []Endpoint {
	//var result []Endpoint
	allEndpoints := gatherEndpointData(logger, clientSet, landerAnnotation)
	return allEndpoints

}

func gatherEndpointData(logger *log.Logger, clientSet *kubernetes.Clientset, landerAnnotation string) []Endpoint {
	var result []Endpoint

	ingressList, err := getIngressList(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	ingressObjects := ingressList.Items
	// time.Sleep(5 * time.Second)
	if len(ingressObjects) > 0 {
		for _, ingress := range ingressObjects {

			// We only want to show links to ingress objects with certain annotations
			if isAnnotatedForLander(ingress, landerAnnotation) {
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
	}
	return result
}

func guessApp(svc string) App {
	fallback, apps := genApps()
	for _, app := range apps {
		if strings.Contains(svc, app.Name) {
			return app
		}
	}
	return fallback
}

// check if a key exists in an ingress annotation
func annotationKeyExists(ingress v1beta1.Ingress, key string) bool {
	_, exists := ingress.Annotations[key]
	return exists
}

func isAnnotatedForLander(ingress v1beta1.Ingress, annotation string) bool {
	return ingress.Annotations[annotation] == "true"
}

// attempts to return the ingress class (or an empty string)
// TODO: upgrade to v1?
func getIngressClass(logger util.LoggerIFace, ingress v1beta1.Ingress) string {
	if val, ok := ingress.Annotations["kubernetes.io/ingress.class"]; ok {
		return val
	}
	logger.Warnf(
		"Unable to determine ingress class for: %s/%s",
		ingress.Namespace,
		ingress.Name)
	return ""
}

// Speaks to the cluster and attempt to pull an IngressList
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

func genApps() (fallback App, index []App) {
	fallback = App{
		Name: "unknown",
		Icon: "/assets/link.png",
		Desc: "generic service",
	}
	//TODO: handle error in function calling this
	var _ error
	index, _ = readAppsFromFile()
	return
}

func readAppsFromFile() (apps []App, err error) {
	type Schema struct {
		Apps []App `json:"apps"`
	}
	var absPath string
	if absPath, err = filepath.Abs("../../assets/apps.json"); err == nil {
		var body []byte
		if body, err = ioutil.ReadFile(absPath); err == nil {
			data := new(Schema)
			err = json.Unmarshal(body, data)
			apps = data.Apps
		}
	}
	return
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

//func filterIngressForHostname(input []Endpoint, host string) []Endpoint {
//	var result []Endpoint
//	for _, data := range input {
//		if strings.Contains(data.Address, host) {
//			result = append(result, data)
//		}
//	}
//	return result
//}
