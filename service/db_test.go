package service

import (
	"context"
	"mailer/model"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
)

type MockDynamoDBAPI struct {
	Output *dynamodb.ScanOutput
	Error  error
}

func (m *MockDynamoDBAPI) Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error) {
	return m.Output, m.Error
}

func TestScan(t *testing.T) {
	testUsers := []model.User{
		{
			Email:     "user1@example.com",
			Name:      "Alice",
			Age:       20,
			Residence: "Tokyo",
			Gender:    1,
		},
		{
			Email:     "user2@example.com",
			Name:      "Bob",
			Age:       22,
			Residence: "Osaka",
			Gender:    0,
		},
		{
			Email:     "user3@example.com",
			Name:      "Carol",
			Age:       24,
			Residence: "Hokkaido",
			Gender:    1,
		},
		{
			Email:     "user4@example.com",
			Name:      "David",
			Age:       26,
			Residence: "Fukuoka",
			Gender:    0,
		},
		{
			Email:     "user5@example.com",
			Name:      "Eve",
			Age:       28,
			Residence: "Saitama",
			Gender:    1,
		},
		{
			Email:     "user6@example.com",
			Name:      "Frank",
			Age:       30,
			Residence: "Chiba",
			Gender:    0,
		},
		{
			Email:     "user7@example.com",
			Name:      "Grace",
			Age:       32,
			Residence: "Nagoya",
			Gender:    1,
		},
		{
			Email:     "user8@example.com",
			Name:      "Henry",
			Age:       34,
			Residence: "Kyoto",
			Gender:    0,
		},
		{
			Email:     "user9@example.com",
			Name:      "Ivy",
			Age:       36,
			Residence: "Kobe",
			Gender:    1,
		},
		{
			Email:     "user10@example.com",
			Name:      "Jack",
			Age:       38,
			Residence: "Hiroshima",
			Gender:    0,
		},
	}

	av := make([]map[string]types.AttributeValue, 0)
	for _, v := range testUsers {
		u, err := attributevalue.MarshalMap(v)
		assert.NoError(t, err)
		av = append(av, u)
	}
	mock := &MockDynamoDBAPI{
		Output: &dynamodb.ScanOutput{
			Items: av,
		},
		Error: nil,
	}

	client := &dbClient{client: mock}
	users, err := client.Scan(context.TODO())
	assert.NoError(t, err)

	filter := model.Filter{
		MaxAge: 35,
		MinAge: 25,
		Gender: -1,
	}
	expected := []model.User{
		{
			Email:     "user4@example.com",
			Name:      "David",
			Age:       26,
			Residence: "Fukuoka",
			Gender:    0,
		},
		{
			Email:     "user5@example.com",
			Name:      "Eve",
			Age:       28,
			Residence: "Saitama",
			Gender:    1,
		},
		{
			Email:     "user6@example.com",
			Name:      "Frank",
			Age:       30,
			Residence: "Chiba",
			Gender:    0,
		},
		{
			Email:     "user7@example.com",
			Name:      "Grace",
			Age:       32,
			Residence: "Nagoya",
			Gender:    1,
		},
		{
			Email:     "user8@example.com",
			Name:      "Henry",
			Age:       34,
			Residence: "Kyoto",
			Gender:    0,
		},
	}

	result := client.Filter(users, filter)
	assert.Equal(t, result, expected)

}
