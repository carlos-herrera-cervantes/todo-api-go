package models

import (
	"encoding/json"
	"net/http"
	"todo-api-fiber/locales"
)

type FailResponse struct {
	Status  bool   `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SuccessResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type SuccessResponseWithPager struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
	Pager  Pager       `json:"pager"`
}

// Returns an internal server error
func (response *FailResponse) SendInternalServerError(w http.ResponseWriter, r *http.Request) {
	lang := r.Header.Get("Accept-Language")
	message := locales.GetLocalizer(lang, "InternalServerError")

	w.WriteHeader(http.StatusInternalServerError)
	response.Code = "InternalServerError"
	response.Message = message

	json.NewEncoder(w).Encode(response)
}

// Returns an unprocessable entity
func (response *FailResponse) SendUnprocessableEntity(w http.ResponseWriter, r *http.Request) {
	response.Code = "UnprocessableEntity"
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(response)
}

// Returns a bad request response
func (response *FailResponse) SendBadRequest(w http.ResponseWriter, r *http.Request, key string) {
	response.Status = false
	response.Code = key

	lang := r.Header.Get("Accept-Language")
	message := locales.GetLocalizer(lang, key)
	response.Message = message

	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(response)
}

// Returns a forbidden response
func (response *FailResponse) SendForbidden(w http.ResponseWriter, r *http.Request, key string) {
	response.Status = false
	response.Code = key

	lang := r.Header.Get("Accept-Language")
	message := locales.GetLocalizer(lang, key)
	response.Message = message

	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(response)
}

// Returns a success response
func (response *SuccessResponse) SendOk(w http.ResponseWriter, r *http.Request) {
	response.Status = true
	json.NewEncoder(w).Encode(response)
}

// Returns a success response including a pager object
func (response *SuccessResponseWithPager) SendOk(w http.ResponseWriter, r *http.Request) {
	response.Status = true
	json.NewEncoder(w).Encode(response)
}

// Returns a created response
func (response *SuccessResponse) SendCreated(w http.ResponseWriter, r *http.Request) {
	response.Status = true
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Returns a no content response
func (response *SuccessResponse) SendNoContent(w http.ResponseWriter, r *http.Request) {
	response.Status = true
	w.WriteHeader(http.StatusNoContent)
	json.NewEncoder(w).Encode(response)
}
