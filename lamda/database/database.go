package database

import (
	"context"
	"fmt"
	appTypes "lamda/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	databaseStore *dynamodb.Client
}

const TABLE_NAME = "userTable"

type DB interface {
	ExistUser(ctx context.Context, username string) (bool, error)
	InsertUser(ctx context.Context, user appTypes.User) error
	GetUser(ctx context.Context, username string) (appTypes.User, error)
}

func NewDynamoDBClient(cfc aws.Config) DynamoDBClient {
	db := dynamodb.NewFromConfig(cfc)
	return DynamoDBClient{
		databaseStore: db,
	}
}

var _ DB = (*DynamoDBClient)(nil)

func (db DynamoDBClient) ExistUser(ctx context.Context, username string) (bool, error) {
	result, err := db.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return false, err
	}

	if result.Item == nil {
		return false, nil
	}

	return true, nil
}

func (db DynamoDBClient) InsertUser(ctx context.Context, user appTypes.User) error {
	item := &dynamodb.PutItemInput{
		TableName: aws.String(TABLE_NAME),
		Item: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: user.Username},
			"password": &types.AttributeValueMemberS{Value: user.PasswordHash},
		},
	}
	_, err := db.databaseStore.PutItem(ctx, item)
	return err
}

func (db DynamoDBClient) GetUser(ctx context.Context, username string) (appTypes.User, error) {

	var user appTypes.User

	result, err := db.databaseStore.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(TABLE_NAME),
		Key: map[string]types.AttributeValue{
			"username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return user, err
	}

	if result.Item == nil {
		return user, fmt.Errorf("user not found")
	}
	err = attributevalue.UnmarshalMap(result.Item, &user)
	if err != nil {
		return user, fmt.Errorf("error unmarshal data: %w", err)
	}
	return user, nil
}
