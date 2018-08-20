package imports

type ImportType struct {
	ID      string `json:"_id"`
	Type    string `json:"type"`
	Name    string `json:"name"`
	Version string `json:"availableVersion"`
	Status  string `json:"lastStatus"`
}
