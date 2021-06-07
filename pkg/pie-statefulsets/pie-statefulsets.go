package pie_statefulset

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// cache for 30sec and expire obj @ 1m
	cacheStatefulSets = cache.New(30*time.Second, 1*time.Minute)
)

// StatefulSetStats is a simple slice/list of deployment pod numbers
type StatefulSetStats struct {
	Bad     int `json:"bad"`
	Good    int `json:"good"`
	Unknown int `json:"unknown"`
}

// getAllStatefulSets speaks to the cluster and attempt to pull all raw StatefulSets
func getAllStatefulSets(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*v1.StatefulSetList, error) {
	cacheObj := "statefulsets"
	cached, found := cacheStatefulSets.Get(cacheObj)
	if found {
		logger.Debugf("got all %s from cache", cacheObj)
		return cached.(*v1.StatefulSetList), nil
	}
	deploymentList, err := clientSet.
		AppsV1().
		StatefulSets("").
		List(
			context.TODO(),
			metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("got all %s from k8s", cacheObj)
	cacheStatefulSets.Set(cacheObj, deploymentList, cache.DefaultExpiration)
	return deploymentList, err
}

// AssembleStatefulSetPieChart is used to assemble data to be returned to the API
func AssembleStatefulSetPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (StatefulSetPieChart, error) {
	var totalBad int64
	var totalGood int64

	data, err := getAllStatefulSets(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, deployment := range data.Items {
		if isHappy(deployment) {
			totalGood++
		} else {
			totalBad++
		}
	}
	var result = StatefulSetPieChart{
		Series: []int64{totalBad, totalGood},
	}
	return result, err
}

func isHappy(k8sObject v1.StatefulSet) bool {
	current := k8sObject.Status.CurrentReplicas
	replicas := k8sObject.Status.Replicas
	readyReplicas := k8sObject.Status.ReadyReplicas
	if current == replicas {
		if readyReplicas == replicas{
			return true
		}
	}
	return false
}
