package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes/timestamp"
)

type PersonServiceServer interface {
	Create(*gin.Context, *CreateRequest) (*CreateResponse, error)
	Read(*gin.Context, *ReadRequest) (*ReadResponse, error)
	Update(*gin.Context, *UpdateRequest) (*UpdateResponse, error)
	Delete(*gin.Context, *DeleteRequest) (*DeleteResponse, error)
	ReadAll(*gin.Context, *ReadAllRequest) (*ReadAllResponse, error)
}

type Person struct {
	Id int64 `json:"id"`
	Name string `json:"name"`
	Birthdate *timestamp.Timestamp `json:"birthdate"`
}

type CreateRequest struct {
	Api string
	Person *Person
}

type CreateResponse struct {
	Api string
	Id int64
}

type ReadRequest struct {
	Api string
	Id int64
}

type ReadResponse struct {
	Api string
	Person *Person
}

type UpdateRequest struct {
	Api string
	Person *Person
}

type UpdateResponse struct {
	Api string
	Updated int64
}

type DeleteRequest struct {
	Api string
	Id int64
}

type DeleteResponse struct {
	Api string
	Deleted int64
}

type ReadAllRequest struct {
	Api string
}

type ReadAllResponse struct {
	Api string
	People []*Person
}