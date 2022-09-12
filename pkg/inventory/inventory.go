package inventory

import (
	"context"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
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
	logger *log.Logger,
	config *rest.Config) *resourceObjList {
	apiGroupList, err := GetAPIGroupList(logger, config)
	if err != nil {
		logger.Fatal(err)
	}

	susResourceObjList, err := ProcessAPIGroupList(logger, config, apiGroupList)
	if err != nil {
		logger.Fatal(err)
	}
	return susResourceObjList
}

func splitGroupVersion(logger *log.Logger, s string) (string, string) {
	sep := "/"
	gv := strings.Split(s, sep)
	length := len(gv)
	switch length {
	case 2:
		return gv[0], gv[1]
	case 1:
		return "", gv[0]
	default:
		logger.Errorf("Couldn't unpack '%s' into group and version\n", s)
		//TODO: is it worth it to "just fail" and returning nothing?
		return "", ""
	}
}

func GetAPIGroupList(
	logger *log.Logger,
	config *rest.Config,
) (
	[]*metav1.APIResourceList,
	error,
) {
	clientKubernetes, err := kubernetes.NewForConfig(config)
	if err != nil {
		logger.Fatal(err)
	}

	apiGroupList, err := clientKubernetes.ServerPreferredResources()
	if err != nil {
		logger.Fatal(err)
	}

	if len(apiGroupList) == 0 {
		logger.Warn("ApiGroupList had no items. Likely something went wrong.")
	}
	return apiGroupList, nil
}

func checkGroupAPIResourceAnnotations(
	logger *log.Logger,
	config *rest.Config,
	resourceObjList *resourceObjList,
	g *metav1.APIResourceList,
) {
	// its OK to ignore these kinds of objects as they're driven/managed by others
	ignoredKinds := []string{
		"ControllerRevision",
		"ReplicaSet",
	}
	clientDynamic, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	group, version := splitGroupVersion(logger, g.GroupVersion)

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
			// this err could simply be because there are no resources of this object kind
			// TODO: evaluate the actual error
			logger.Warnf(
				"APIGROUP: '%s', APIGROUPVERSION: '%s', KIND: '%s', err: %s",
				group, version, a.Name, err,
			)
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
		}
	}
}

func ProcessAPIGroupList(
	logger *log.Logger,
	config *rest.Config,
	apiGroupList []*metav1.APIResourceList,
) (*resourceObjList, error) {
	var susResourceObjList resourceObjList

	cacheObj := "v1/apiGroups"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*resourceObjList), nil
	}

	for _, g := range apiGroupList {
		// k8s api returns `group/version`, this helper function helps to split/unpack these
		group, version := splitGroupVersion(logger, g.GroupVersion)

		logger.Debugf("discovered apiGroup: %s, v: %s\n", group, version)
		checkGroupAPIResourceAnnotations(logger, config, &susResourceObjList, g)
	}

	pkgCache.Set(cacheObj, &susResourceObjList, cache.DefaultExpiration)

	return &susResourceObjList, nil
}
