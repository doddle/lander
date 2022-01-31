package deployments

import (
	"github.com/withmandala/go-log"
	"k8s.io/client-go/kubernetes"
)

// MetaDataDeploymentsTable is some metadata to be used to represent traffic
type MetaDataDeploymentsTable struct {
	Name              string `json:"name"`
	Namespace         string `json:"ns"`
	Ready             bool   `json:"ready"`
	Progressing       bool   `json:"progressing"`
	ReplicasDesired   int32  `json:"replicas"`
	ReplicasAvailable int32  `json:"replicas_available"`
}

// AssembleDeploymentsTable preps some data about deployments
func AssembleDeploymentsTable(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) []MetaDataDeploymentsTable {
	var result []MetaDataDeploymentsTable

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	for _, k8sObj := range data.Items {

		result = append(result, MetaDataDeploymentsTable{
			Name:              k8sObj.Name,
			Namespace:         k8sObj.Namespace,
			Ready:             isReady(k8sObj),
			Progressing:       isProgressing(k8sObj),
			ReplicasDesired:   *k8sObj.Spec.Replicas,
			ReplicasAvailable: k8sObj.Status.AvailableReplicas,
		})
	}

	return result
}
