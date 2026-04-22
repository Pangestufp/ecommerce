package dto

type ResponseParam struct {
	StatusCode int       `json:"status_code"`
	Message    string    `json:"message"`
	Paginate   *Paginate `json:"paginate,omitempty"`
	Data       any       `json:"data"`
}
