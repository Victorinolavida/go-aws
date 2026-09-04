package main

import (
	"context"
	"errors"
	"fmt"
	"lamda/app"
	"lamda/middleware"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Username string `json:"username"`
}

// HandleRequest take a payload and process it
func HandleRequest(event MyEvent) (string, error) {
	if event.Username == "" {
		return "", errors.New("username can not be empty")

	}
	return fmt.Sprintf("sucess called by - %s", event.Username), nil
}

func ProtectedHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       "this is a secret path",
		StatusCode: http.StatusOK,
	}, nil
}

func main() {
	appLambda := app.NewApp()
	// lambda.Start(appLambda.ApiHandler.RegisterUserHandler)
	lambda.Start(func(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.Path {
		case "/register":
			return appLambda.ApiHandler.RegisterUserHandler(ctx, request)
		case "/login":
			return appLambda.ApiHandler.LoginUserHandler(ctx, request)
		case "/protected":
			return middleware.ValidateJWTMiddleware(ProtectedHandler)(request)
		default:
			return events.APIGatewayProxyResponse{
				Body:       "not found",
				StatusCode: http.StatusNotFound,
			}, nil

		}
	})
}
