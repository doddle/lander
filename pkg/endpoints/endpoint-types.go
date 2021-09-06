package endpoints

// Endpoint is for the metadata returned to the browser/frontend
type Endpoint struct {
	// these are not ac
	Address     string `json:"address"`     // combined hostname + /paths
	Class       string `json:"class"`       // name of the ingress path
	Https       bool   `json:"https"`       // https
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
	Description string `json:"description"` // if we can match this to an app, we can propogate this
	Name        string `json:"name"`        // Application name
	Icon        string `json:"icon"`
}
