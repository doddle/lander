package endpoints

import (
	"github.com/withmandala/go-log"
	"k8s.io/client-go/kubernetes"
)

func AssembleRouteMetaData(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) []RouteMetaData {
	var result []RouteMetaData

	ingressList, err := getIngressList(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	ingressObjects := ingressList.Items
	for _, ingress := range ingressObjects {
		for _, rule := range ingress.Spec.Rules {
			for _, p := range rule.IngressRuleValue.HTTP.Paths {
				path := p.Path
				result = append(result, RouteMetaData{
					Hostname:    rule.Host,
					Namespace:   ingress.Namespace,
					Oauth2proxy: getOauth2ProxyState(ingress),
					Path:        path,
					Service:     p.Backend.ServiceName,
				})
			}
		}
	}
	return result
}
