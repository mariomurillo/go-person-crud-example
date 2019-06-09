package v1

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-person-crud-example/pkg/protocol/grpc"
	v1 "go-person-crud-example/pkg/service/v1"
)

type Config struct {

	GRPCPort string
	DatastoreDBHost string
	DatastoreDBUser string
	DatastoreDBPassword string
	DatastoreDBSchema string
}

func RunServer() error {
	ctx := context.Background()

	cfg := Config{
		GRPCPort: "9090",
		DatastoreDBHost: "localhost:3306",
		DatastoreDBSchema: "mydatabase",
		DatastoreDBUser: "root",
		DatastoreDBPassword: "12345",
	}

	param := "parseTime=true"

	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s",
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBHost,
		cfg.DatastoreDBSchema,
		param)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}
	defer db.Close()

	v1Api := v1.NewPersonaServiceServer(db)

	return grpc.RunServer(ctx, v1Api, cfg.GRPCPort)
}