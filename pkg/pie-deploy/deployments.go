package pie_deploy

import (
	"context"
	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"strings"
	"time"
)

var (
	// cache for 30sec and expire obj @ 1m
	cacheDeployments = cache.New(30*time.Second, 1*time.Minute)
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
	cacheObj := "deployments"
	cached, found := cacheDeployments.Get(cacheObj)
	if found {
		logger.Info("got all deployments from cache")
		return cached.(*v1.DeploymentList) , nil
	} else {
		deploymentList, err := clientSet.
			AppsV1().
			Deployments("").
			List(
				context.TODO(),
				metav1.ListOptions{})

		if err != nil {
			return nil, err
		}
		logger.Info("got all deployments from k8s")
		cacheDeployments.Set(cacheObj, deploymentList, cache.DefaultExpiration)
		return deploymentList, err
	}
}


func AssembleDeploymentPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	) (DeploymentPieChart, error){
	var totalBad int64
	var totalGood int64

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, deployment := range data.Items{
		if isReady(deployment) {
			totalGood++
		} else {
			totalBad++
		}
	}
	var result = DeploymentPieChart{
		Series: []int64{totalBad, totalGood},
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