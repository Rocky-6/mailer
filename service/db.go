package service

import (
	"context"
	"log"

	"mailer/model"
	"mailer/repository"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type DynamoDBAPI interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

type dbClient struct {
	client DynamoDBAPI
}

func NewDBClient(ctx context.Context) (repository.DBRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("ap-northeast-1"))
	if err != nil {
		return nil, err
	}
	return &dbClient{client: dynamodb.NewFromConfig(cfg)}, err
}

func (c *dbClient) Scan(ctx context.Context) ([]model.User, error) {
	var users []model.User
	var err error
	var response *dynamodb.ScanOutput
	projEx := expression.NamesList(
		expression.Name("mail_address"), expression.Name("name"), expression.Name("age"), expression.Name("residence"), expression.Name("gender"))
	expr, err := expression.NewBuilder().WithProjection(projEx).Build()
	if err != nil {
		log.Fatal(err)
	} else {
		response, err = c.client.Scan(ctx, &dynamodb.ScanInput{
			TableName:                 aws.String("mailer-users"),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			FilterExpression:          expr.Filter(),
			ProjectionExpression:      expr.Projection(),
		})
		if err != nil {
			log.Fatal(err)
		} else {
			err = attributevalue.UnmarshalListOfMaps(response.Items, &users)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
	return users, err
}

func (c *dbClient) Filter(users []model.User, filt model.Filter) []model.User {
	var filteredUsers []model.User
	for _, u := range users {
		if u.Age < filt.MinAge || u.Age > filt.MaxAge {
			continue
		}
		if u.Residence != filt.Residence && filt.Residence != "" {
			continue
		}
		if u.Gender != filt.Gender && filt.Gender != -1 {
			continue
		}
		filteredUsers = append(filteredUsers, u)
	}
	return filteredUsers
}
