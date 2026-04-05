package dto

type ResponseParam struct {
	StatusCode int
	Message    string
	Paginate   *Paginate
	Data       any
}
