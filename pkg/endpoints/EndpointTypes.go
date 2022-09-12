package endpoints

// Endpoint is for the metadata returned to the browser/frontend
type Endpoint struct {
	// these are not ac
	Address     string `json:"address"`     // combined hostname + /paths
	Class       string `json:"class"`       // name of the ingress path
	HTTPS       bool   `json:"https"`       // https
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
	Description string `json:"description"` // if we can match this to an app, we can propogate this
	Name        string `json:"name"`        // Application name
	Icon        string `json:"icon"`        // Appropriate file name for the Icon
}

// RouteMetaData is some metadata to be used to represent traffic
type RouteMetaData struct {
	Hostname    string `json:"hostname"`
	Path        string `json:"path"`
	Namespace   string `json:"ns"`
	Service     string `json:"svc"`
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
	Class       string `json:"class"`
}
