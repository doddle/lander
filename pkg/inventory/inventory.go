package inventory

import (
	"context"
	"fmt"
	"github.com/patrickmn/go-cache"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
	"strings"
	"time"
)

type resourceObj struct {
	Resource        string
	Kind            string
	Namespace       string `json:",omitempty"`
	APIGroup        string
	APIGroupVersion string
}

type resourceObjList []resourceObj

var (
	// Create a cache with a default expiration time of 60 s, and which
	// purges expired items every 15 minutes
	pkgCache = cache.New(60*time.Second, 15*time.Minute)
)

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func AssembleFluxIgnored(
//clientSet *kubernetes.Clientset,
) *resourceObjList {
	config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}

	apiGroupList, err := GetAPIGroupList(config)
	if err != nil {
		panic(err)
	}

	susResourceObjList, err := ProcessApiGroupList(config, apiGroupList)
	if err != nil {
		panic(err)
	}

	return susResourceObjList
}

func splitGroupVersion(s string) (string, string) {
	sep := "/"
	gv := strings.Split(s, sep)
	if len(gv) == 2 {
		return gv[0], gv[1]
	} else if len(gv) == 1 {
		return "", gv[0]
	} else {
		fmt.Printf("Couldn't unpack %s into group and version", s)
		//TODO: I don't know how to just "fail" and don't return anything
		return "", ""
	}
}

func GetAPIGroupList(config *rest.Config) ([]*metav1.APIResourceList, error) {
	clientKubernetes, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	apiGroupList, err := clientKubernetes.ServerPreferredResources()
	if err != nil {
		panic(err)
	}

	if len(apiGroupList) == 0 {
		fmt.Println("ApiGroupList had no items. Likely something went wrong.")
	}
	return apiGroupList, nil
}

func checkGroupAPIResourceAnnotations(config *rest.Config, resourceObjList *resourceObjList, g *metav1.APIResourceList) {
	ignoredKinds := []string{
		"ControllerRevision",
		"ReplicaSet",
	}
	clientDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	group, version := splitGroupVersion(g.GroupVersion)

	// unpack the info that we need into our custom struct
	for _, a := range g.APIResources {
		objectKind := schema.GroupVersionResource{
			Group:    group,
			Version:  version,
			Resource: a.Name,
		}

		// we use the dynamic client instead of plain k8s to be able to query CRDs
		rs, err := clientDynamic.Resource(objectKind).Namespace("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("K8s Dynamic client failed to get resources for: APIGROUP: '%s', APIGROUPVERSION: '%s', KIND: '%s'", group, version, a.Name)
			continue
		}
		// go over all resources and check for the annotation of interest
		for _, r := range rs.Items {
			ns := r.GetNamespace()
			name := r.GetName()
			kind := r.GetKind()
			annotations := r.GetAnnotations()

			if !contains(ignoredKinds, kind) {
				if len(annotations) > 0 {
					for k, v := range annotations {
						if strings.Contains(k, "flux.weave.works/ignore") {
							if strings.Contains(v, "true") {
								//log.Infow("WARNING!",
								//	"resource", name,
								//	"kind", kind,
								//	"namespace", ns,
								//	"apiGroup", group,
								//	"apiGroupVersion", version,
								//)
								*resourceObjList = append(*resourceObjList, resourceObj{
									Resource:        name,
									Kind:            kind,
									Namespace:       ns,
									APIGroup:        group,
									APIGroupVersion: version,
								})
							}
						}
					}
				}
			}
			//} else {
			//	log.Infof("sus: no annotations for namespace: %s, resource: %s, kind: %s", r.GetNamespace(), r.GetName(), r.GetKind())
			//}
		}
	}
}

func ProcessApiGroupList(config *rest.Config, apiGroupList []*metav1.APIResourceList) (*resourceObjList, error) {
	var susResourceObjList resourceObjList

	cacheObj := "v1/apiGroups"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*resourceObjList), nil
	}

	for _, g := range apiGroupList {
		// k8s api returns `group/version`, this helper function helps to split/unpack these
		group, version := splitGroupVersion(g.GroupVersion)

		if len(g.APIResources) == 0 {
			fmt.Println("discovered resources-less api",
				"apiGroup", group,
				"apiGroupVersion", version,
			)
			continue
		}

		fmt.Println("discovered api",
			"apiGroup", group,
			"apiGroupVersion", version,
		)

		checkGroupAPIResourceAnnotations(config, &susResourceObjList, g)
	}

	pkgCache.Set(cacheObj, &susResourceObjList, cache.DefaultExpiration)

	return &susResourceObjList, nil
}
