package statefulsets

import (
	"strings"

	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
)

// MetaDataDeploymentsTable is some metadata to be used to represent traffic
type MetaDataDeploymentsTable struct {
	Name            string `json:"name"`
	Namespace       string `json:"ns"`
	Ready           bool   `json:"ready"`
	Progressing     bool   `json:"progressing"`
	ReplicasDesired int32  `json:"replicas"`
	ReplicasUpdated int32  `json:"replicas_updated"`
	ReplicasCurrent int32  `json:"replicas_current"`
}

// AssembleDeploymentsTable preps some data about deployments
func AssembleDeploymentsTable(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) []MetaDataDeploymentsTable {
	var result []MetaDataDeploymentsTable

	data, err := getAllStatefulSets(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	for _, k8sObj := range data.Items {

		result = append(result, MetaDataDeploymentsTable{
			Name:            k8sObj.Name,
			Namespace:       k8sObj.Namespace,
			Ready:           isReady(k8sObj),
			Progressing:     isProgressing(k8sObj),
			ReplicasDesired: *k8sObj.Spec.Replicas,
			ReplicasUpdated: k8sObj.Status.UpdatedReplicas,
			ReplicasCurrent: k8sObj.Status.CurrentReplicas,
		})
	}

	return result
}

func isProgressing(k8sObject v1.StatefulSet) bool {
	for _, obj := range k8sObject.Status.Conditions {
		if strings.Contains(string(obj.Type), "Progressing") {
			if strings.Contains(string(obj.Status), "True") {
				return true
			}
		}
	}
	return true
}
