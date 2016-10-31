package bg

type Created struct {
	ID string
}

type General struct {
	Code    int
	Message string `json:"Message,omitempty"`
}
