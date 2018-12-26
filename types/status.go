package types

// StatusUpdate struct is used for DynamoDB updates, because the update command requires all json keys to start with ":"
type StatusUpdate struct {
	Status string `json:":status"`
}
