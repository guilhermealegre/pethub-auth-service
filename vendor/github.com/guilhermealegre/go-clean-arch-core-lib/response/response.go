package response

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/fatih/structs"
	"github.com/guilhermealegre/go-clean-arch-core-lib/errors"
	"github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
	msg "github.com/guilhermealegre/go-clean-arch-core-lib/pagination"
)

// IsEmpty verify if the error response is empty
func (e *ErrorResponse) IsEmpty() bool {
	return len(e.Errors) == 0
}

// AddError add an error to the list
func (e *ErrorResponse) AddError(err error) {
	if err != nil {
		switch err := err.(type) {
		case errors.ErrorDetails:
			e.Errors = append(e.Errors, err)
		case errors.ErrorDetailsList:
			for _, er := range err {
				e.Errors = append(e.Errors, er)
			}
		default:
			e.Errors = append(e.Errors, ErrorGeneric.Formats(err.Error()))
		}
	}
}

// SetMetadata sets the metadata
func (e *ErrorResponse) SetMetadata(metadata interface{}) {
	mBytes, _ := json.Marshal(metadata)
	_ = json.Unmarshal(mBytes, &e.Meta.MetaData)
}

// GetResponse get response
func GetResponse(data interface{}, pagination *msg.Pagination, meta interface{}, errs ...error) (int, interface{}) {
	if len(errs) > 0 && errs[0] != nil {
		errorResponse := ErrorResponse{}
		for _, e := range errs {
			if errorDetail, ok := e.(errors.ErrorDetails); ok {
				if errorDetail.StatusCode > errorResponse.Status {
					errorResponse.Status = errorDetail.StatusCode
					errorDetail.StatusCode = 0 // clean the status code to not be sent
					e = errorDetail
				}
			}

			errorResponse.AddError(e)
		}

		if meta != nil {
			errorResponse.SetMetadata(meta)
		}

		if errorResponse.Status == 0 {
			errorResponse.Status = http.StatusBadRequest
		}

		return errorResponse.Status, errorResponse
	} else {
		response := &Response{}
		response.Data = data
		response.CalculateLength()

		if pagination != nil {
			response.SetLinks(pagination)
		}

		if meta != nil {
			response.SetMetadata(meta)
		}

		response.SetDefaults()

		return http.StatusOK, response
	}
}

// CalculateLength calculate response length
func (r *Response) CalculateLength() {
	switch reflect.TypeOf(r.Data).Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.String, reflect.Slice:
		r.Meta.Length = reflect.ValueOf(r.Data).Len()
	case reflect.Struct:
		r.Meta.Length = len(structs.Map(r.Data))
	}
}

// SetLinks set the links for pagination
func (r *Response) SetLinks(pages *pagination.Pagination) {
	r.Links.First = pages.ToString(pages.First())
	r.Links.Prev = pages.ToString(pages.Prev())
	r.Links.Self = pages.ToString(pages)
	r.Links.Next = pages.ToString(pages.Next())
	r.Links.Last = pages.ToString(pages.Last())
}

// SetDefaults set default response
func (r *Response) SetDefaults() {
	if r.Data == nil || (reflect.TypeOf(r.Data).Kind() == reflect.Slice && reflect.ValueOf(r.Data).Len() == 0) {
		r.Data = make([]interface{}, 0)
	}
}

// SetMetadata sets the metadata
func (r *Response) SetMetadata(metadata interface{}) {
	mBytes, _ := json.Marshal(metadata)
	_ = json.Unmarshal(mBytes, &r.Meta.MetaData)
}
