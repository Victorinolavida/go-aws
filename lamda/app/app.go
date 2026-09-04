package app

import (
	"context"
	"fmt"
	"lamda/api"
	"lamda/database"

	"github.com/aws/aws-sdk-go-v2/config"
)

type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1")) //TODO: check region
	if err != nil {
		fmt.Printf("unable to load SDK config, error: %v", err)
		panic(1)
	}
	db := database.NewDynamoDBClient(cfg)
	handler := api.NewApiHandler(db)
	return App{ApiHandler: handler}
}
