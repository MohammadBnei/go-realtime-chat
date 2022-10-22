package adapter

type ReqAddCompany struct {
	// Name of the company
	// in: string
	Name string `json:"name" validate:"required,min=2,max=100,alpha_space"`
	// Status of the company
	// in: int64
	Status int64 `json:"status" validate:"required"`
}

// swagger:parameters admin addCompany
type MessageInput struct {
	// - name: body
	//  in: body
	//  description: name and status
	//  schema:
	//  type: object
	//     "$ref": "#/definitions/ReqAddCompany"
	//  required: true
	Body ReqAddCompany `json:"body"`
}

// swagger:model Company
type Company struct {
	// Id of the company
	// in: int64
	Id int64 `json:"id"`
	// Name of the company
	// in: string
	Name string `json:"name"`
	// Status of the company
	// in: int64
	Status int64 `json:"status"`
}

// swagger:model CommonError
type CommonError struct {
	// Status of the error
	// in: int64
	Status int64 `json:"status"`
	// Message of the error
	// in: string
	Message string `json:"message"`
}
