package main

// Endpoint is a basic construct assembled from the k8s Ingress rules detected
// it is meant for creating a list of URLs and with some simple metadata based
// annotations, Paths and hostnames
type Endpoint struct {
	// these are not ac
	Address     string `json:"address"`     // combined hostname + /paths
	Class       string `json:"class"`       // name of the ingress path
	Https       bool   `json:"https"`       // https
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
}

// A slice of Endpoint objects
type EndpointList struct {
	Endpoints []Endpoint `json:"endpoints"`
}
