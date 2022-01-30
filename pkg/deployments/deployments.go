package deployments

import (
	"context"
	"strings"
	"time"

	"github.com/doddle/lander/pkg/chart"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// hard limit cache for 10sec, expire at 10m
	pkgCache = cache.New(10*time.Second, 10*time.Minute)
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
) (FinalPieChart, error) {
	var resultColors []string
	var resultLabels []string
	var resultSeries []int64
	var totalDown int64    // Conditions: anything else
	var totalHealthy int64 // Conditions: Available + Progressing
	var totalPending int64 // Conditions: Progressing (but containers available not 100%)

	// colors: https://apexcharts.com/docs/options/colors/
	// vs https://vuetifyjs.com/en/styles/colors/#material-colors
	//
	// try to stick to "$color lighten 2"
	colDown := "#E57373"
	colHealthy := "#81C784"
	colPending := "#FFB74D"

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, deployment := range data.Items {
		if isHealthy(deployment) {
			totalHealthy++
		} else {
			if isPending(deployment) {
				totalPending++
			} else {
				totalDown++
			}
		}
	}

	if totalDown > 0 {
		resultLabels = append(resultLabels, "Down")
		resultSeries = append(resultSeries, totalDown)
		resultColors = append(resultColors, colDown)
	}
	if totalHealthy > 0 {
		resultLabels = append(resultLabels, "Healthy")
		resultSeries = append(resultSeries, totalHealthy)
		resultColors = append(resultColors, colHealthy)
	}
	if totalPending > 0 {
		resultLabels = append(resultLabels, "Pending")
		resultSeries = append(resultSeries, totalPending)
		resultColors = append(resultColors, colPending)
	}

	result := FinalPieChart{
		Total:  totalDown + totalHealthy + totalPending,
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
				ID: "pie-deployments",
			},
			Labels: resultLabels,
		},
	}

	return result, err
}

// isReady is a meta kind of job
// isReady checks if a pod has a status condition of "Available==True"
func isReady(obj v1.Deployment) bool {
	for _, condition := range obj.Status.Conditions {
		if strings.Contains(string(condition.Type), "Available") {
			if strings.Contains(string(condition.Status), "True") {
				return true
			}
		}
	}
	return false
}

func isProgressing(k8sObject v1.Deployment) bool {
	for _, obj := range k8sObject.Status.Conditions {
		if strings.Contains(string(obj.Type), "Progressing") {
			if strings.Contains(string(obj.Status), "True") {
				return true
			}
		}
	}
	return true
}

func isPending(obj v1.Deployment) (result bool) {
	// at least one pod still working but stuff is looking bad
	if obj.Status.ReadyReplicas > 0 {
		if obj.Status.AvailableReplicas > 0 {
			return true
		}
	}
	return false
}

func isHealthy(obj v1.Deployment) (result bool) {
	// This is the most intense check.. basically if there is anything wrong, the object is not healthy
	if !isReady(obj) {
		return false
	}
	if !isProgressing(obj) {
		return false
	}
	// we want all of these to match replicas for the deployment to be classed as healthy
	//  updatedReplicas
	//  readyReplicas
	//  availableReplicas
	if obj.Status.UpdatedReplicas != obj.Status.Replicas {
		return false
	}
	if obj.Status.ReadyReplicas != obj.Status.Replicas {
		return false
	}
	if obj.Status.AvailableReplicas != obj.Status.Replicas {
		return false
	}
	return true
}
