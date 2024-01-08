package service

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
)

type MockS3API struct {
	Output *s3.GetObjectOutput
	Error  error
}

func (m *MockS3API) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	return m.Output, m.Error
}

func TestFetchTemplate(t *testing.T) {
	file, err := os.Open("../template.json")
	assert.NoError(t, err)

	mock := &MockS3API{
		Output: &s3.GetObjectOutput{
			Body: file,
		},
	}
	client := &s3Client{client: mock}

	template, err := client.FetchTemplate(context.TODO(), "", "")
	assert.NoError(t, err)
	assert.Equal(t, template.Subject, "おしらせ")
	assert.Equal(t, template.Body, "おしらせです")
}
