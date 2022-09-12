package endpoints

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	networkingv1 "k8s.io/api/networking/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// hard limit cache for 10sec, expire at 10m
	pkgCache = cache.New(10*time.Second, 10*time.Minute)
)

// Speaks to the cluster and attempt to pull an IngressList using standard networking/v1
func getIngressListV1(logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*networkingv1.IngressList, error) {
	cacheObj := "networkingv1/ingress"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*networkingv1.IngressList), nil
	}

	// find ALL Ingresses (Yes we know we can watch them, but this is ok)
	ingressList, err := clientSet.
		NetworkingV1().
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
	return ingressList, nil
}
