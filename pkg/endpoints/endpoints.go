package endpoints

import (
	"github.com/doddle/lander/pkg/util"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes"
)

func ReallyAssemble(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	landerAnnotationRoot string,
) []Endpoint {
	var result []Endpoint

	ingressList, err := getIngressList(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	ingressObjects := ingressList.Items
	if len(ingressObjects) > 0 {
		for _, ingress := range ingressObjects {
			if !isAnnotatedForLander(ingress, landerAnnotationRoot) {
				continue
			}
			for _, rule := range ingress.Spec.Rules {
				for _, p := range rule.IngressRuleValue.HTTP.Paths {
					serviceDescription := ""
					if annotationKeyExists(ingress, landerAnnotationRoot+"/description") {
						serviceDescription = ingress.Annotations[landerAnnotationRoot+"/description"]
					}
					serviceName := p.Backend.ServiceName
					if annotationKeyExists(ingress, landerAnnotationRoot+"/name") {
						serviceName = ingress.Annotations[landerAnnotationRoot+"/name"]
					}
					serviceIcon := ""
					if annotationKeyExists(ingress, landerAnnotationRoot+"/icon") {
						serviceIcon = ingress.Annotations[landerAnnotationRoot+"/icon"]
					}

					// Strip out a trailing "/"
					uri := p.Path
					if p.Path == "/" {
						uri = ""
					}
					serviceUrl := "https://" + rule.Host + uri
					if annotationKeyExists(ingress, landerAnnotationRoot+"/url") {
						serviceUrl = ingress.Annotations[landerAnnotationRoot+"/url"]
					}
					result = append(result, Endpoint{
						Address:     serviceUrl,
						Https:       true,
						Oauth2proxy: getOauth2ProxyState(ingress),
						Class:       getIngressClass(logger, ingress),
						Description: serviceDescription,
						Name:        serviceName,
						Icon:        serviceIcon,
					})
				}
			}
		}
	}
	return result
}

// check if a key exists in an ingress annotation
func annotationKeyExists(ingress v1beta1.Ingress, key string) bool {
	_, exists := ingress.Annotations[key]
	return exists
}

func isAnnotatedForLander(ingress v1beta1.Ingress, annotationBase string) bool {
	return ingress.Annotations[annotationBase+"/show"] == "true"
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

// returns true/false if ingress Annotations contain what looks like oa2p
func getOauth2ProxyState(ingress v1beta1.Ingress) bool {
	if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-signin") {
		if annotationKeyExists(ingress, "nginx.ingress.kubernetes.io/auth-url") {
			return true
		}
	}
	return false
}
