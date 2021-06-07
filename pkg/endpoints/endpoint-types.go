package endpoints
// Endpoint is for the metadata returned to the browser/frontend
type Endpoint struct {
	// these are not ac
	Address     string `json:"address"`     // combined hostname + /paths
	Class       string `json:"class"`       // name of the ingress path
	Https       bool   `json:"https"`       // https
	Oauth2proxy bool   `json:"oauth2proxy"` // if its secured by an oauth2proxy
	Icon        string `json:"icon"`        // we will attempt to guess the ICON for endpoints
	Description string `json:"description"` // if we can match this to an app, we can propogate this
	Name        string `json:"name"`        // Application name
}

// App is a generic definition of a known, we'll use these to attempt to guess the apps
// TODO: add some tags, slice of common service names and import the data from json/yaml maybe
type App struct {
	Name string `yaml:"name"`
	Icon string `yaml:"icon"`
	Desc string `yaml:"desc"`
}
