package bg

type Created struct {
	ID uint `json:"id"`
}

type General struct {
	Code    int `json:"code"`
	Message string `json:"message,omitempty"`
	Path string `json:"path,omitempty"`
}
