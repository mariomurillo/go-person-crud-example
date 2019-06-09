package v1

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang/protobuf/ptypes"
	"go-person-crud-example/pkg/api/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

const apiVersion = "v1"

type PersonaServiceServer struct {
	db *sql.DB
}

func NewPersonaServiceServer(db *sql.DB) v1.PersonaServiceServer {
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

func (s *PersonaServiceServer) connect(ctx context.Context) (*sql.Conn, error) {
	c, err := s.db.Conn(ctx)

	if err != nil {
		return nil, status.Error(codes.Unknown, "failed to connect to database-> "+err.Error())
	}
	return c, nil
}

func (s *PersonaServiceServer) Create(ctx context.Context, request *v1.CreateRequest) (*v1.CreateResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var birthdate time.Time

	if request.Persona.Birthdate != nil {
		birthdate, err = ptypes.Timestamp(request.Persona.Birthdate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Birthdate field has invalid format-> "+err.Error())
		}
	}

	result, err := conn.ExecContext(ctx, "INSERT `mydatabase`.`Person` (`name`, `birthdate`) VALUES (?, ?)",
		request.Persona.Name, birthdate)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to insert into Person-> "+err.Error())
	}

	id, err := result.LastInsertId()

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve id for created Person-> "+err.Error())
	}

	return &v1.CreateResponse{
		Api: apiVersion,
		Id: id,
	}, nil
}

func (s *PersonaServiceServer) Read(ctx context.Context, request *v1.ReadRequest) (*v1.ReadResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "SELECT `id`, `name`, `birthdate` FROM mydatabase.Person WHERE `id`=?",
		request.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from Person-> "+err.Error())
	}

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, status.Error(codes.Unknown, "Failed to select from Person-> "+err.Error())
		}
		return nil, status.Error(codes.NotFound, fmt.Sprintf("ToDo with ID='%d' is not found", request.Id))
	}

	var p v1.Person
	var birthdate time.Time

	if err := rows.Scan(&p.Id, &p.Name, &birthdate); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve field values from Person row-> "+err.Error())
	}
	p.Birthdate, err = ptypes.TimestampProto(birthdate)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Birthdate field has invalid format-> "+err.Error())
	}
	if rows.Next() {
		return nil, status.Error(codes.Unknown, fmt.Sprintf("Found multiple Person rows with ID='%d'",
			request.Id))
	}
	return &v1.ReadResponse{
		Api: apiVersion,
		Persona: &p,
	}, nil
}

func (s *PersonaServiceServer) Update(ctx context.Context, request *v1.UpdateRequest) (*v1.UpdateResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var birthdate time.Time

	if request.Persona.Birthdate != nil {
		birthdate, err = ptypes.Timestamp(request.Persona.Birthdate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Birthdate field has invalid format-> "+err.Error())
		}
	}

	result, err := conn.ExecContext(ctx, "UPDATE `mydatabase`.`Person` SET `name` = ?, `birthdate` = ? WHERE `id` = ?",
		request.Persona.Name, birthdate, request.Persona.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to update into Person-> "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve rows affected value-> "+err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Person with ID='%d' is not found",
			request.Persona.Id))
	}
	return &v1.UpdateResponse{
		Api: apiVersion,
		Updated: rowsAffected,
	}, nil
}

func (s *PersonaServiceServer) Delete(ctx context.Context, request *v1.DeleteRequest) (*v1.DeleteResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	result, err := conn.ExecContext(ctx, "DELETE FROM `mydatabase`.`Person` WHERE `id` = ?",
		request.Id)

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to delete into Person-> "+err.Error())
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve rows affected value-> "+err.Error())
	}
	if rowsAffected == 0 {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("Person with ID='%d' is not found",
			request.Id))
	}
	return &v1.DeleteResponse{
		Api: apiVersion,
		Deleted: rowsAffected,
	}, nil
}

func (s *PersonaServiceServer) ReadAll(ctx context.Context, request *v1.ReadAllRequest) (*v1.ReadAllResponse, error) {
	if err := s.checkApi(request.Api); err != nil {
		return nil, err
	}
	conn, err := s.connect(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, "SELECT `id`, `name`, `birthdate` FROM mydatabase.Person")

	if err != nil {
		return nil, status.Error(codes.Unknown, "Failed to select from Person-> "+err.Error())
	}

	var birthdate time.Time
	people := []*v1.Person{}

	for rows.Next() {
		person := new(v1.Person)

		if err := rows.Scan(&person.Id, &person.Name, &birthdate); err != nil {
			return nil, status.Error(codes.Unknown, "Failed to retrieve field values from Person row-> "+err.Error())
		}
		person.Birthdate, err = ptypes.TimestampProto(birthdate)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "Birthdate field has invalid format-> "+err.Error())
		}

		people = append(people, person)
	}

	if err := rows.Err(); err != nil {
		return nil, status.Error(codes.Unknown, "Failed to retrieve data from Person-> "+err.Error())
	}
	return &v1.ReadAllResponse{
		Api: apiVersion,
		People: people,
	}, nil
}