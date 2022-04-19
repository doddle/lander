package deployments

import (
	"github.com/withmandala/go-log"
	v1 "k8s.io/api/apps/v1"
	"k8s.io/client-go/kubernetes"
	"strings"
)

// MetaDataDeploymentsTable is some metadata to be used to represent traffic
type MetaDataDeploymentsTable struct {
	Name              string `json:"name"`
	Namespace         string `json:"ns"`
	Ready             bool   `json:"ready"`
	Progressing       bool   `json:"progressing"`
	ReplicasDesired   int32  `json:"replicas"`
	ReplicasAvailable int32  `json:"replicas_available"`
	Tag               string `json:"tag"`
}

type TagFilters struct {
	Registry string
	Name     string
}

// AssembleDeploymentsTable preps some data about deployments
func AssembleDeploymentsTable(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	filteredTags []TagFilters,
) []MetaDataDeploymentsTable {
	var result []MetaDataDeploymentsTable

	data, err := getAllDeployments(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}

	for _, k8sObj := range data.Items {
		tag := guessTheImportantTag(logger, k8sObj, filteredTags)

		result = append(result, MetaDataDeploymentsTable{
			Name:              k8sObj.Name,
			Namespace:         k8sObj.Namespace,
			Ready:             isReady(k8sObj),
			Progressing:       isProgressing(k8sObj),
			ReplicasDesired:   *k8sObj.Spec.Replicas,
			ReplicasAvailable: k8sObj.Status.AvailableReplicas,
			Tag:               tag,
		})
	}

	return result
}

func guessTheImportantTag(
	logger *log.Logger,
	deploy v1.Deployment,
	filteredTags []TagFilters,
) (result string) {
	fallback := ""
	for _, obj := range deploy.Spec.Template.Spec.Containers {
		logger.Debugf("deploy: %s, image: %s\n", deploy.Name, obj.Image)
		if matchImageToReg(obj.Image, filteredTags) {
			return parseMatchingTag(obj.Image, filteredTags)
		} else {
			// try to split the image by :
			imageSplit := strings.Split(obj.Image, ":")
			// if there was a : to split by, return that image
			if len(imageSplit) > 1 {
				fallback = imageSplit[1]
			} else {
				// otherwise the fallback tag must be the docker image :latest
				fallback = "latest"
			}
		}
	}
	return fallback
}

func parseMatchingTag(image string, filters []TagFilters) (result string) {
	splitData := strings.Split(image, ":")
	fullTag := splitData[1]
	return fullTag
}

func matchImageToReg(image string, filters []TagFilters) (result bool) {
	for _, i := range filters {
		if strings.Contains(image, i.Registry) {
			return true
		}

	}
	return false
}
