package common

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mattermost-extend/configuration"
	"mattermost-extend/helper"
	"net/http"
)

type (
	successResource struct {
		Success bool            `json:"success"`
		Data    successResponse `json:"data"`
	}
	successResponse struct {
		Information interface{} `json:"information"`
		Message     string      `json:"message"`
		HttpStatus  int         `json:"status"`
	}
	appError struct {
		Message    interface{} `json:"message"`
		HttpStatus int         `json:"status"`
	}
	errorResource struct {
		Success bool     `json:"success"`
		Data    appError `json:"data"`
	}
)

func DisplayAppSuccessResponse(w http.ResponseWriter, data interface{}, message string) {
	response := successResponse{
		Information: data,
		Message:     message,
		HttpStatus:  http.StatusOK,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	if j, err := json.Marshal(successResource{Success: true, Data: response}); err == nil {
		w.Write(j)
	}
}

func DisplayAppErrorResponse(w http.ResponseWriter, errorMessage interface{}, code int) {
	errObj := appError{
		Message:    errorMessage,
		HttpStatus: code,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if j, err := json.Marshal(errorResource{Success: false, Data: errObj}); err == nil {
		w.Write(j)
	}
}

func CheckRequestToken(r *http.Request) ([]byte, error) {
	request, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	tokenChecker := helper.TokenChecker{}
	err = json.Unmarshal(request, &tokenChecker)
	if err != nil {
		return nil, err
	}

	if !helper.Contains(tokenChecker.Token, configuration.ChatWithMeToken) {
		return nil, errors.New("chatwithme tokens are not the same in MM and CRM")
	}
	return request, err
}
