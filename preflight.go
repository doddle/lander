package main

import (
	"flag"
	"fmt"
	"github.com/digtux/lander/pkg/endpoints"
	"github.com/withmandala/go-log"
	"k8s.io/client-go/kubernetes"
	"os"
)

func envVarExists(key string) bool {
	_, exists := os.LookupEnv(key)
	return exists
}

// using flag.Visit, check if a flag was provided
// if not.. tell the user, print `-help` and bail
func checkRequredFlag() {
	required := []string{"host"}
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
	logger.Info("getting some initial data bootstrapped")
	clientSet, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		logger.Fatal(err)
	}
	_ = endpoints.ReallyAssemble(
		logger,
		clientSet,
		*flagHost,
	)
}
