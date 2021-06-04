package deployments

import (
	"context"
	"github.com/withmandala/go-log"
	"k8s.io/api/extensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)



// GetDeployments speaks to the cluster and attempt to pull all raw Deployments
func GetDeployments(logger *log.Logger, clientset *kubernetes.Clientset) (*v1beta1.IngressList, error) {

	ingressList, err := clientset.ExtensionsV1beta1().Ingresses("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return ingressList, err
}
