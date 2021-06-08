package statefulsets

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
	// hard limit cache for 15sec, expire at 15m
	pkgCache = cache.New(15*time.Second, 15*time.Minute)
)

// StatefulSetStats is a simple slice/list of deployment pod numbers
type StatefulSetStats struct {
	Bad     int `json:"bad"`
	Good    int `json:"good"`
	Unknown int `json:"unknown"`
}

// getAllStatefulSets speaks to the cluster and attempt to pull all raw statefulSets
func getAllStatefulSets(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*v1.StatefulSetList, error) {
	cacheObj := "v1/StatefulSetList"
	cached, found := pkgCache.Get(cacheObj)
	if found {
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
	pkgCache.Set(cacheObj, deploymentList, cache.DefaultExpiration)
	return deploymentList, err
}

// AssembleStatefulSetPieChart is used to assemble data to be returned to the API
func AssembleStatefulSetPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (FinalResult, error) {
	var resultColors []string
	var resultLabels []string
	var resultSeries []int64
	var totalBad int64
	var totalGood int64

	data, err := getAllStatefulSets(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, deployment := range data.Items {
		if isReady(deployment) {
			totalGood++
		} else {
			totalBad++
		}
	}

	if totalBad > 0 {
		resultLabels = append(resultLabels, "Errored")
		resultSeries = append(resultSeries, totalBad)
		resultColors = append(resultColors, "#E57373")
	}
	if totalGood > 0 {
		resultLabels = append(resultLabels, "Healthy")
		resultSeries = append(resultSeries, totalGood)
		resultColors = append(resultColors, "#81C784")
	}

	result := FinalResult{
		Total:  totalBad + totalGood,
		Series: resultSeries,
		ChartOpts: ChartOpts{
			Colors: resultColors,
			Labels: resultLabels,
			Chart: Chart{
				ID: "pie-statefulsets",
			},
		},
	}
	return result, err
}

func isReady(k8sObject v1.StatefulSet) bool {
	current := k8sObject.Status.CurrentReplicas
	replicas := k8sObject.Status.Replicas
	readyReplicas := k8sObject.Status.ReadyReplicas
	if current == replicas {
		if readyReplicas == replicas {
			return true
		}
	}
	return false
}
