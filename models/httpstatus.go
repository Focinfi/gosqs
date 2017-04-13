package models

type HTTPStatusMeta struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// HTTPStatus contains Info of a failed request
type HTTPStatus struct {
	HTTPStatusMeta
	Data interface{} `json:"data"`
}
