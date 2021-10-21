package main

import (
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/holmes89/thoth/internal/database"
	"github.com/holmes89/thoth/internal/handlers"
	"github.com/holmes89/thoth/internal/handlers/rest"
)

var (
	once   sync.Once
	router *handlers.APIRouter
)

func main() {
	once.Do(func() {
		conn := database.NewConnection()
		router = handlers.NewAPIRouter()
		rest.MakeGameHandler(router.Router, conn)
	})

	lambda.Start(router.HandleRequest)
}
