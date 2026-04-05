package helper

import "backend/dto"

type Response struct {
	Code     int           `json:"code"`
	Status   string        `json:"status"`
	Message  string        `json:"message"`
	Paginate *dto.Paginate `json:"paginate,omitempty"`
	Data     any           `json:"data,omitempty"`
}

func BuildResponse(params dto.ResponseParam) Response {
	status := "failed"
	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = "success"
	}

	return Response{
		Code:     params.StatusCode,
		Status:   status,
		Message:  params.Message,
		Paginate: params.Paginate,
		Data:     params.Data,
	}
}
