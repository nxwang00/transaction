package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/golang/gddo/httputil/header"

	"github.com/server/transaction/db"
	model "github.com/server/transaction/model"
)

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func decodeJSONBody(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
		if value != "application/json" {
			msg := "Content-Type header is not application/json"
			return &malformedRequest{status: http.StatusUnsupportedMediaType, msg: msg}
		}
	}

	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &malformedRequest{status: http.StatusBadRequest, msg: msg}

		case err.Error() == "http: request body too large":
			msg := "Request body must not be larger than 1MB"
			return &malformedRequest{status: http.StatusRequestEntityTooLarge, msg: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &malformedRequest{status: http.StatusBadRequest, msg: msg}
	}

	return nil
}

func httpSendResponse(w http.ResponseWriter, code int, resp interface{}, err error) {
	if code == 0 {
		// code not specified, determine based on error
		if err == nil {
			code = http.StatusOK
		} else {
			if strings.Contains(err.Error(), "unique constraint") {
				code = http.StatusConflict
			} else if strings.Contains(err.Error(), "Unauthorized") {
				code = http.StatusUnauthorized
			} else {
				code = http.StatusBadRequest
			}
		}
	}

	if resp == nil && err != nil {
		// build resp based on error message
		out := &Response{}
		out.Data.Status = code
		out.Code = ""
		out.Message = err.Error()

		resp = out
	}

	if isNil(resp) {
		w.Header().Set("Content-Length", "0")
		w.WriteHeader(code)
	} else {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(resp)
	}
}

func CreateTransaction(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Add Transaction ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var transaction model.TransactionReq
	var resp string
	var err error
	var code int
	err = decodeJSONBody(w, r, &transaction)
	if err == nil {
		resp, err = db.InsertTransaction(&transaction)
	}

	httpSendResponse(w, code, resp, err)
}

func ReadTransactions(w http.ResponseWriter, r *http.Request) {
	log.Printf("============== Get All Transactions with pagination and filter ===============\n")
	log.Printf("%s http://%s%s", r.Method, r.Host, r.RequestURI)

	var pageInfo = model.PageInfo{}
	var filterInfo = model.TransactionReq{}
	params := r.URL.Query()
	for k, v := range params {
		if k == "page_num" {
			pageNum, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Println("page number -> converting string to int is failed")
				return
			}
			pageInfo.Page_Number = pageNum
		}
		if k == "page_size" {
			pageSize, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Println("page size -> converting string to int is failed")
				return
			}
			pageInfo.Page_Size = pageSize
		}

		if k == "origin" {
			filterInfo.Origin = v[0]
		}
		if k == "user_id" {
			userID, err := strconv.Atoi(v[0])
			if err != nil {
				fmt.Println("User Id -> converting string to int is failed")
				return
			}
			filterInfo.User_ID = userID
		}
		if k == "amount" {
			filterInfo.Amount = "$" + v[0]
		}
		if k == "op_type" {
			filterInfo.Op_Type = v[0]
		}
		if k == "registered_at" {
			filterInfo.Registered_At = v[0]
		}
	}

	var err error
	var resp []*model.Transaction
	var code int

	resp = db.SelectTransactions(pageInfo, filterInfo)

	if len(resp) == 0 {
		httpSendResponse(w, code, "There ins`t nothing in DB now", err)
	}

	httpSendResponse(w, code, resp, err)
}