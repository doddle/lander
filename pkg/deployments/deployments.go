package deployments

import (
	"context"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// hard limit cache for 15sec, expire at 15m
	pkgCache = cache.New(15 * time.Second, 15 * time.Minute)
)


// DeploymentStats is a simple slice/list of deployment pod numbers
type DeploymentStats struct {
	Bad     int `json:"bad"`
	Good    int `json:"good"`
	Unknown int `json:"unknown"`
}

// getAllDeployments speaks to the cluster and attempt to pull all raw Deployments
func getAllDeployments(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*v1.DeploymentList, error) {
	cacheObj := "v1/DeploymentList"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*v1.DeploymentList), nil
	}
	deploymentList, err := clientSet.
		AppsV1().
		Deployments("").
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

// AssembleDeploymentPieChart is used to assemble data to be returned to the API
func AssembleDeploymentPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (FinalResult, error) {
	var resultColors []string
	var resultLabels []string
	var resultSeries []int64
	var totalBad int64
	var totalGood int64

	data, err := getAllDeployments(logger, clientSet)
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

	// colors: https://apexcharts.com/docs/options/colors/
	// vs https://vuetifyjs.com/en/styles/colors/#material-colors
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
		Total: totalBad + totalGood,
		Series: resultSeries,
		ChartOpts: ChartOpts{
			Colors: resultColors,
			Chart: Chart{
				ID: "pie-deployments",
				DropShadow: DropShadow{Effect: false},
			},
			Labels: resultLabels,
		},
	}
	return result, err
}

// isReady is a meta kind of job
// isReady checks if a pod has a status condition of "Available==True"
// TODO: possibly check for "Progressing" also?
func isReady(deployment v1.Deployment) bool {
	for _, obj := range deployment.Status.Conditions {
		if strings.Contains(string(obj.Type), "Available") {
			if strings.Contains(string(obj.Status), "True") {
				return true
			}
		}
	}
	return false
}
