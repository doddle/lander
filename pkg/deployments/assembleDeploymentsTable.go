package deployments

import (
	"github.com/withmandala/go-log"
	"k8s.io/client-go/kubernetes"
)


// DeploymentsMetaData is some metadata to be used to represent traffic
type DeploymentsMetaData struct {
	Name    string `json:"name"`
	Namespace   string `json:"ns"`
}

// AssembleDeploymentsTable preps some data about deployments
func AssembleDeploymentsTable(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) []DeploymentsMetaData {
	var result []DeploymentsMetaData

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	for _, deployment := range data.Items {

		//x := DeploymentsMetaData{
		//	Name:      deployment.Name,
		//	Namespace: deployment.Namespace,
		//}

		//if isHealthy(deployment) {
		//	totalHealthy++
		//} else {
		//	if isPending(deployment) {
		//		totalPending++
		//	} else {
		//		totalDown++
		//	}
		//}

		result = append(result, DeploymentsMetaData{
			Name: deployment.Name,
			Namespace: deployment.Namespace,
		})
	}

	return result
}
