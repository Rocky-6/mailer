package main

import (
	"fmt"
	"mailer/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFillInPlaceholder(t *testing.T) {
	template := &model.Template{
		Subject: `おしらせ`,
		Body:    `${name}さんへ${topic}のおしらせです`,
		Params: map[string]interface{}{
			"name":  "name",
			"topic": "お得な商品",
		},
	}

	users := []model.User{
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

	messages := make([]*model.MailMessege, 0)
	for _, u := range users {
		messages = append(messages, &model.MailMessege{
			Sender:    "xxx",
			Recipient: u.Email,
			Charset:   "UTF-8",
			Param: &model.UserParam{
				Name:      u.Name,
				Age:       u.Age,
				Gender:    u.Gender,
				Residence: u.Residence,
			},
		})
	}

	for _, messege := range messages {
		fillInPlaceholder(messege, template)
		assert.Equal(t, fmt.Sprintf("%sさんへお得な商品のおしらせです", messege.Param.Name), messege.Body)
	}
}
