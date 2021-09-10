package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/digtux/lander/pkg/endpoints"
	"github.com/withmandala/go-log"
	"k8s.io/client-go/kubernetes"
)

func envVarExists(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// using flag.Visit, check if a flag was provided
// if not.. tell the user, print `-help` and bail
func checkRequredFlag() {
	required := []string{}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Printf("ERROR: missing required flag: '-%s'\n-------------\n", req)
			flag.PrintDefaults()
			os.Exit(2)
		}
	}
}

func onStartup(logger *log.Logger) {
	logger.Info("Startup: Prefetch endpoints & index icons")
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	_ = endpoints.ReallyAssemble(
		logger,
		clientSet,
		*flagLanderAnnotationBase,
	)
	files, err := ioutil.ReadDir(*flagAssetPath)
	if err != nil {
		logger.Fatal(err)
	}
	for _, f := range files {
		if !f.IsDir() {
			availableIcons = append(availableIcons, f.Name())
		}
	}
}
