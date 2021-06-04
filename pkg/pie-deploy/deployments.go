package pie_deploy

import (
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

import (
	"context"
	"github.com/withmandala/go-log"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

type DeploymentStats struct{
	Bad 	int `json:"bad"`
	Good 	int `json:"good"`
	Unknown int `json:"unknown"`
}

// getAllDeployments speaks to the cluster and attempt to pull all raw Deployments
func getAllDeployments(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	) (*v1.DeploymentList, error) {
	logger.Debug("getting deployment data")
	deploymentList, err := clientSet.
		AppsV1().
		Deployments("").
		List(
		context.TODO(),
		metav1.ListOptions{})

	if err != nil {
		return nil, err
	}
	return deploymentList, err
}


func AssembleDeploymentPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	) (DeploymentPieChart, error){
	colourBad := "rgba(225, 99, 132, 1)"
	colourBadBg := "rgba(255, 99, 132, 0.2)"
	colourGood := "rgba(65, 194, 93, 1)"
	colourGoodBg := "rgba(65, 194, 93, 0.2)"
	// colourUnknown := "rgba(235, 190, 54, 1)"
	// colourUnknownBg := "rgba(235, 190, 54, 1)"

	var totalBad int64
	var totalGood int64

	opts := ChartOptions{
			Legend: Legend{Display: true},
			Responsive: true,
			MaintainAspectRatio: true,
		}

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, deployment := range data.Items{
		logger.Info(deployment.Status)
		if isReady(deployment) {
			totalGood++
		} else {
			totalBad++
		}
	}
	var result = DeploymentPieChart{
		ChartOptions: opts,
		ChartData: ChartData{
			Labels: []string{
				"bad",
				"good",
			},
			Datasets: []Dataset{
				{
					BorderWidth: 0,
					BackgroundColor: []string{
						colourBadBg,
						colourGoodBg,
					},
					BorderColor: []string{
						colourBad,
						colourGood,
					},
					Data: []int64{
						totalBad,
						totalGood,
					},
				},
			},
		},
	}
	return result, err
}

func isReady(deployment v1.Deployment) bool{
	for _, obj := range deployment.Status.Conditions{
		if strings.Contains(string(obj.Type), "Available") {
			if strings.Contains(string(obj.Status), "True") {
				return true
			}
		}

	}
	return false
}