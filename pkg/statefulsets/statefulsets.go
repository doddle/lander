package statefulsets

import (
	"context"
	"time"

	"github.com/digtux/lander/pkg/chart"
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
	objList, err := clientSet.
		AppsV1().
		StatefulSets("").
		List(
			context.TODO(),
			metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("got all %s from k8s", cacheObj)
	pkgCache.Set(cacheObj, objList, cache.DefaultExpiration)
	return objList, err
}

// AssembleStatefulSetPieChart is used to assemble data to be returned to the API
func AssembleStatefulSetPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (FinalPieChart, error) {
	var resultColors []string
	var resultLabels []string
	var resultSeries []int64
	var totalBad int64
	var totalGood int64

	data, err := getAllStatefulSets(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, obj := range data.Items {
		if isReady(obj) {
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

	result := FinalPieChart{
		Total:  totalBad + totalGood,
		Series: resultSeries,
		ChartOpts: chart.Opts{
			Legend: chart.Legend{Show: true},
			PlotOpt: chart.PlotOpt{
				Pie: chart.PlotOptPie{
					ExpandOnClick: false,
					Size:          120,
				},
			},
			Colors: resultColors,
			Stroke: chart.Stroke{Width: 0},
			Chart: chart.Chart{
				ID: "pie-statefulsets",
			},
			Labels: resultLabels,
		},
	}
	return result, err
}

func isReady(k8sObject v1.StatefulSet) bool {
	// updated replicas should be the same as "ready replicas" and the current replica total
	// otherwise return false
	replicas := k8sObject.Status.Replicas
	replicasReady := k8sObject.Status.ReadyReplicas
	replicasUpdated := k8sObject.Status.UpdatedReplicas
	if replicas == replicasReady {
		if replicas == replicasUpdated {
			return true
		}
	}
	return false
}
