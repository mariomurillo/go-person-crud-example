package v2

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/ptypes"
	"github.com/jinzhu/gorm"
	v2 "go-person-crud-example/pkg/api/v2"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
	"time"
)

const apiVersion = "v2"

type PersonaServiceServer struct {
	db *gorm.DB
}

func NewPersonaServiceServer(db *gorm.DB) v2.PersonServiceServer {
	return &PersonaServiceServer{db:db}
}

func (s *PersonaServiceServer) checkApi(api string) error {
	if len(api) > 0 {
		if apiVersion != api {
			return status.Errorf(codes.Unimplemented,
				"unsupported API version: service implements API version '%s', but asked for '%s'",
				apiVersion, api)
		}
	}
	return nil
}

func (s *PersonaServiceServer) Create(ctx *gin.Context, request *v2.CreateRequest) (*v2.CreateResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}

	s.db.Save(request.Person)

	response := v2.CreateResponse{
		Api: apiVersion,
		Id:  request.Person.Id,
	}

	ctx.JSON(http.StatusCreated, response)

	return &response, nil
}

func (s *PersonaServiceServer) Read(ctx *gin.Context, request *v2.ReadRequest) (*v2.ReadResponse, error) {
	panic("implement me")
}

func (s *PersonaServiceServer) Update(ctx *gin.Context, request *v2.UpdateRequest) (*v2.UpdateResponse, error) {
	panic("implement me")
}

func (s *PersonaServiceServer) Delete(ctx *gin.Context, request *v2.DeleteRequest) (*v2.DeleteResponse, error) {
	panic("implement me")
}

func (s *PersonaServiceServer) ReadAll(ctx *gin.Context, request *v2.ReadAllRequest) (*v2.ReadAllResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}

	var people []*v2.Person

	db := s.db.DB()
	rows, err := db.QueryContext(ctx, "SELECT `id`, `name`, `birthdate` FROM mydatabase.Person")
	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from Person-> "+err.Error())
	}
	var birthdate time.Time
	for rows.Next() {
		person := new(v2.Person)

		if err := rows.Scan(&person.Id, &person.Name, &birthdate); err != nil {
			return nil, status.Error(codes.Unknown, "Failed to retrieve field values from Person row-> "+err.Error())
		}
		person.Birthdate, err = ptypes.TimestampProto(birthdate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Birthdate field has invalid format-> "+err.Error())
		}

		people = append(people, person)
	}

	response := v2.ReadAllResponse{
		Api:    apiVersion,
		People: people,
	}

	return &response, nil
}