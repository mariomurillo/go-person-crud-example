package http

import (
	"errors"
	"github.com/gin-gonic/gin"
	v2 "go-person-crud-example/pkg/api/v2"
	"log"
	"net/http"
)

type PersonHTTPHandler struct {
	server v2.PersonServiceServer
}

func NewPersonHttpHandler(server *v2.PersonServiceServer) *PersonHTTPHandler {
	return &PersonHTTPHandler{server:*server}
}

func (h *PersonHTTPHandler) Create(ctx *gin.Context)  {
	var request v2.CreateRequest

	if err := ctx.BindJSON(&request); err != nil {
		handleError(ctx, errors.New("GENERAL_ERROR"), 409, "Failed Bind Json CreateRequest")
		return
	}

	response, err := h.server.Create(ctx, &request)

	if err != nil {
		handleError(ctx, err, 409, "Failed creating Person")
		return
	}
	ctx.JSON(201, response)
}

func (h *PersonHTTPHandler) ReadAll(ctx *gin.Context) {
	var request v2.ReadAllRequest

	api := ctx.Param("api")

	if api == "" {
		handleError(ctx, errors.New("GENERAL_ERROR"), 409, "Failed Api is require")
		return
	}
	request.Api = api

	response, err := h.server.ReadAll(ctx, &request)

	if err != nil {
		handleError(ctx, err, 409, "Failed creating Person")
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *PersonHTTPHandler) Read(ctx *gin.Context) {
	var request v2.ReadRequest

	if err := ctx.BindJSON(&request); err != nil {
		handleError(ctx, errors.New("GENERAL_ERROR"), 409, "Failed Bind Json CreateRequest")
		return
	}

	response, err := h.server.Read(ctx, &request)

	if err != nil {
		handleError(ctx, err, 409, "Failed creating Person")
		return
	}
	ctx.JSON(200, response)
}

func (h *PersonHTTPHandler) Upddate(ctx *gin.Context) {
	var request v2.UpdateRequest

	if err := ctx.BindJSON(&request); err != nil {
		handleError(ctx, errors.New("GENERAL_ERROR"), 409, "Failed Bind Json CreateRequest")
		return
	}

	response, err := h.server.Update(ctx, &request)

	if err != nil {
		handleError(ctx, err, 409, "Failed creating Person")
		return
	}
	ctx.JSON(200, response)
}

func (h *PersonHTTPHandler) Delete(ctx *gin.Context) {
	var request v2.DeleteRequest

	if err := ctx.BindJSON(&request); err != nil {
		handleError(ctx, errors.New("GENERAL_ERROR"), 409, "Failed Bind Json CreateRequest")
		return
	}

	response, err := h.server.Delete(ctx, &request)

	if err != nil {
		handleError(ctx, err, 409, "Failed creating Person")
		return
	}
	ctx.JSON(200, response)
}

func handleError(ctx *gin.Context, err error, code int, message string) {
	log.Printf("Error: %s", err.Error())
	ctx.JSON(code, gin.H{
		"message": message,
		"status":  code,
		"error":   err.Error(),
	})
}
