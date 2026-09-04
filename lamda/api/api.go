package api

import (
	"context"
	"encoding/json"
	"fmt"
	"lamda/database"
	"lamda/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type ApiHandler struct {
	dbStore database.DB
}

func NewApiHandler(db database.DB) ApiHandler {
	return ApiHandler{dbStore: db}
}

func (api *ApiHandler) RegisterUserHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var userData types.RegisterUser
	err := json.Unmarshal([]byte(request.Body), &userData)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if userData.Username == "" || userData.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "request haves empty parameters",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request haves empty parameters")
	}

	user, err := types.NewUser(userData)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       "invalid user data",
				StatusCode: http.StatusBadRequest,
			},
			err
	}

	userExist, err := api.dbStore.ExistUser(ctx, user.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       "internal error server",
				StatusCode: http.StatusInternalServerError,
			},
			err
	}
	if userExist {
		return events.APIGatewayProxyResponse{
			Body:       "user already register",
			StatusCode: http.StatusConflict,
		}, nil
	}

	err = api.dbStore.InsertUser(ctx, user)

	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal error server",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	return events.APIGatewayProxyResponse{
		Body:       "successfully registered user",
		StatusCode: http.StatusOK,
	}, nil
}
func (api *ApiHandler) LoginUserHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if loginRequest.Username == "" || loginRequest.Password == "" {
		return events.APIGatewayProxyResponse{
			Body:       "request haves empty parameters",
			StatusCode: http.StatusBadRequest,
		}, fmt.Errorf("request haves empty parameters")
	}

	user, err := api.dbStore.GetUser(ctx, loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	err = types.ValidatePassword(user.PasswordHash, loginRequest.Password)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body:       "invalid user credentials",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	accessToken := types.CreateToken(user)
	fmt.Printf("token %s\n", accessToken)
	sucessMsg := fmt.Sprintf(`{"access_token": "%s"}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body:       sucessMsg,
		StatusCode: http.StatusOK,
	}, nil
}
