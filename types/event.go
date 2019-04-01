package types

// Event contains the event body used in the invokation of the Lambda
type Event struct {
	Operation  string `json:"operation"`
	Repository string `json:"repository"`
	Branch     string `json:"branch"`
	Action     string `json:"action"`
}
