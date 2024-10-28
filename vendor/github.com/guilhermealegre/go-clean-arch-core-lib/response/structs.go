package response

// swagger:model Response
type Response struct {
	// Data
	Data interface{} `json:"data"`
	// Meta
	Meta Meta `json:"meta"`
	// Links
	Links Links `json:"links"`
}

// swagger:model Meta
type Meta struct {
	// Page Size
	PageSize int `json:"page_size"`
	// Length
	Length int `json:"length"`
	// MetaData
	MetaData map[string]interface{} `json:"meta_data,omitempty"`
}

// swagger:model Links
type Links struct {
	// First
	First *string `json:"first"`
	// Prev
	Prev *string `json:"prev"`
	// Self
	Self *string `json:"self"`
	// Next
	Next *string `json:"next"`
	// Last
	Last *string `json:"last"`
}

// swagger:model ErrorResponse
type ErrorResponse struct {
	//Status
	Status int `json:"status"`
	// Errors
	Errors []error `json:"errors"`
	// Meta
	Meta Meta `json:"meta,omitempty"`
}
