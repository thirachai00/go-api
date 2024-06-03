package model

type JsonResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
