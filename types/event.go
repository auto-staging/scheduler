package types

type Event struct {
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Action     string `json:"action"`
}
