package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	v2 "go-person-crud-example/pkg/api/v2"
	"go-person-crud-example/pkg/protocol/http"
	service "go-person-crud-example/pkg/service/v2"
)

type Config struct {
	HTTPPort string
	DatastoreDBHost string
	DatastoreDBUser string
	DatastoreDBPassword string
	DatastoreDBSchema string
}

func RunServer() error {
	cfg := Config{
		HTTPPort: "9091",
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

	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("Failed to open database: %v", err)
	}
	defer db.Close()

	db.AutoMigrate(&v2.Person{})

	server := service.NewPersonaServiceServer(db)
	handler := http.NewPersonHttpHandler(&server)

	engine := gin.Default()

	group := engine.Group("/api/service/person")

	group.POST("/:api", handler.Create)
	group.GET("/:api", handler.ReadAll)
	group.GET("/:api/:id", handler.Read)
	group.PUT("/:api/:id", handler.Upddate)
	group.DELETE("/:api/:id", handler.Delete)

	return engine.Run(":"+cfg.HTTPPort)
}