package endpoints

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

var (
	// hard limit cache for 10sec, expire at 10m
	pkgCache = cache.New(10*time.Second, 10*time.Minute)
)

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
	return ingressList, nil
}
