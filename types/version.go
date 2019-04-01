package types

// SingleComponentVersion is the implementation of the TowerAPI SingleComponentVersion schema
type SingleComponentVersion struct {
	Name       string `json:"name"`
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	Branch     string `json:"branch"`
	BuildTime  string `json:"buildTime"`
}
