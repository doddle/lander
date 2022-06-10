package nodes

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/doddle/lander/pkg/chart"
	v1 "k8s.io/api/core/v1"

	"github.com/patrickmn/go-cache"
	"github.com/withmandala/go-log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// hard limit cache for 10sec, expire at 10m
	pkgCache = cache.New(10*time.Second, 10*time.Minute)
)

// getAllNodes speaks to the cluster and attempt to pull all raw Nodes
func getAllNodes(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (*v1.NodeList, error) {
	cacheObj := "v1/NodeList"
	cached, found := pkgCache.Get(cacheObj)
	if found {
		return cached.(*v1.NodeList), nil
	}
	objList, err := clientSet.
		CoreV1().
		Nodes().
		List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	logger.Debugf("got all %s from k8s", cacheObj)
	pkgCache.Set(cacheObj, objList, cache.DefaultExpiration)
	return objList, err
}

func AssembleNodesPieChart(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
) (FinalPieChart, error) {
	var resultColors []string
	var resultLabels []string
	var resultSeries []int64
	var totalBad int64
	var totalGood int64
	data, err := getAllNodes(logger, clientSet)
	if err != nil {
		logger.Error(err)
	}
	for _, obj := range data.Items {
		if isReady(obj) && isSchedulable(obj) {
			totalGood++
		} else {
			totalBad++
		}
	}

	// colors: https://apexcharts.com/docs/options/colors/
	// vs https://vuetifyjs.com/en/styles/colors/#material-colors
	// green lighten-2
	colGood := "#81C784"
	// red lighten-2
	colBad := "#E57373"
	if totalBad > 0 {
		resultLabels = append(resultLabels, "UnHealthy")
		resultSeries = append(resultSeries, totalBad)
		resultColors = append(resultColors, colBad)
	}
	if totalGood > 0 {
		resultLabels = append(resultLabels, "Healthy")
		resultSeries = append(resultSeries, totalGood)
		resultColors = append(resultColors, colGood)
	}
	result := FinalPieChart{
		Total:  totalBad + totalGood,
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
				ID: "pie-nodes",
			},
			Labels: resultLabels,
		},
	}
	return result, err
}

func AssembleTable(
	logger *log.Logger,
	clientSet *kubernetes.Clientset,
	labelSlices []string,
) (NodeTable, error) {
	var result NodeTable

	var nodeStats []NodeStats //nolint:prealloc
	nodez, err := getAllNodes(logger, clientSet)
	if err != nil {
		logger.Fatal(err)
	}
	now := time.Now()
	for _, i := range nodez.Items {
		age := now.Sub(i.CreationTimestamp.Time)

		matchedLabels := map[string]string{}

		for _, desiredLabel := range labelSlices {
			humanKey := shortLabelName(desiredLabel)
			x := getLabelValue(i, desiredLabel)
			matchedLabels[humanKey] = x
		}

		// We're going to assume the node is this version (of k8s)
		kubeletVersion := i.Status.NodeInfo.KubeletVersion
		kernel := i.Status.NodeInfo.KernelVersion

		newNode := NodeStats{
			AgeSeconds:  intergerOnly(age.Seconds()),
			Ready:       isReady(i),
			Schedulable: isSchedulable(i),
			Name:        i.Name,
			Version:     kubeletVersion,
			LabelMap:    matchedLabels,
			Kernel:      kernel,
		}
		nodeStats = append(nodeStats, newNode)
	}

	standardHeaders := []TableHeaders{
		{
			Text:  "Name",
			Align: "start",
			Value: "name",
		},
		{Text: "Ready", Value: "ready"},
		{Text: "Schedulable", Value: "schedulable"},
		{Text: "Age", Value: "age"},
		{Text: "Version", Value: "version"},
		{Text: "Kernel", Value: "kernel"},
	}

	for _, customLabel := range labelSlices {
		humanKey := shortLabelName(customLabel)
		keyUpper := strings.ToUpper(humanKey)
		dataKey := fmt.Sprintf("labels.%s", humanKey)
		newHeader := TableHeaders{
			Text:  keyUpper,
			Value: dataKey,
		}
		standardHeaders = append(standardHeaders, newHeader)
	}

	result = NodeTable{
		Headers: standardHeaders,
		Nodes:   nodeStats,
	}

	// assemble headers
	return result, err
}

func shortLabelName(input string) string {
	// shortens labels like:
	//  - node.kubernetes.io/instance-type  => instance-type
	//  - node.kubernetes.io/instancegroup  => instancegroup
	//  - topology.kubernetes.io/zone       => zone
	// This basically grabs whatever comes after the final `/`
	splitString := strings.Split(input, "/")
	// wasn't split:
	if len(splitString) < 1 {
		return input
	}
	// return the last slice
	return splitString[len(splitString)-1]
}

func getLabelValue(node v1.Node, label string) string {
	for k, v := range node.Labels {
		if k == label {
			return v
		}
	}
	return ""
}

func intergerOnly(input float64) int {
	i, _ := math.Modf(input)
	return int(i)
}

func isReady(k8sObject v1.Node) bool {
	for _, conditions := range k8sObject.Status.Conditions {
		if strings.Contains(string(conditions.Type), "Ready") {
			if strings.Contains(string(conditions.Status), "True") {
				return true
			}
		}
	}
	return false
}

// simple check to see if the node is schedulable
// TODO: check for cordon (and other) taints one day?
func isSchedulable(k8sObject v1.Node) bool {
	return !(k8sObject.Spec.Unschedulable)
}
