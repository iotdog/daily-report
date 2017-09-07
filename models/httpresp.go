package models

type CommonResponse struct {
	Code int    `json:"code" description:"response code"`
	Msg  string `json:"message" description:"description of response code"`
}
