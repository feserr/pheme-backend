package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	fiberadapter "github.com/awslabs/aws-lambda-go-api-proxy/fiber"
	_ "github.com/feserr/pheme-backend/docs"
	"github.com/feserr/pheme-backend/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

// @title Pheme backend
// @version 1.0
// @description Pheme backend service.

// @contact.name Elias Serrano
// @contact.email feserr3@gmail.com

// @license.name BSD 2-Clause License
// @license.url http://opensource.org/licenses/BSD-2-Clause

var fiberLambda *fiberadapter.FiberLambda

// Handler will deal with Fiber working with Lambda
func Handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// If no name is provided in the HTTP request body, throw an error
	return fiberLambda.ProxyWithContext(ctx, req)
}

// @BasePath /api/
func init() {
	app := fiber.New(fiber.Config{
		Prefork:       true,
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		AppName:       "Pheme backend v1.0.0",
	})

	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	if len(os.Args[1:]) == 1 && os.Args[1] == "local" {
		err := app.Listen(fmt.Sprintf("%v:%v", os.Getenv("PHEME_HOST"), os.Getenv("PHEME_PORT")))
		if err != nil {
			panic(err.Error())
		}

		return
	}

	fiberLambda = fiberadapter.New(app)
}

func main() {
	if len(os.Args[1:]) < 1 || os.Args[1] != "local" {
		lambda.Start(Handler)
	}
}
