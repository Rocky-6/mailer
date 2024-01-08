package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"

	"mailer/model"
	"mailer/service"

	"github.com/aws/aws-sdk-go/aws/awserr"
)

func main() {
	ctx := context.Background()
	sesRepository, err := service.NewSESClient(ctx)
	if err != nil {
		panic(err)
	}
	s3Repository, err := service.NewS3Client(ctx)
	if err != nil {
		panic(err)
	}
	dbRepository, err := service.NewDBClient(ctx)
	if err != nil {
		panic(err)
	}
	bucket := os.Getenv("BUCKET")

	filterMode := flag.Bool("f", true, "filter use")
	sender := flag.String("s", os.Getenv("SENDER"), "sender")
	recipient := flag.String("r", os.Getenv("RECIPIENT"), "recipient")
	templatefile := flag.String("t", "template.json", "template file")
	charset := "UTF-8"
	minAge := flag.Int("minAge", 0, "min age filter")
	maxAge := flag.Int("maxAge", math.MaxInt, "max age filter")
	gender := flag.Int("gender", -1, "gender filter")
	residence := flag.String("residence", "", "residence filter")
	flag.Parse()

	template, err := s3Repository.FetchTemplate(ctx, bucket, *templatefile)
	if err != nil {
		log.Fatal(err)
	}
	messages := make([]*model.MailMessege, 0)

	if *filterMode {
		users, err := dbRepository.Scan(ctx)
		if err != nil {
			panic(err)
		}
		recipients := dbRepository.Filter(users, model.Filter{
			MinAge:    *minAge,
			MaxAge:    *maxAge,
			Gender:    *gender,
			Residence: *residence,
		})
		for _, r := range recipients {
			messages = append(messages, &model.MailMessege{
				Sender:    *sender,
				Recipient: r.Email,
				Charset:   charset,
				Param: &model.UserParam{
					Name:      r.Name,
					Age:       r.Age,
					Gender:    r.Gender,
					Residence: r.Residence,
				},
			})
		}
	} else {
		messages = append(messages, &model.MailMessege{
			Sender:    *sender,
			Recipient: *recipient,
			Charset:   charset,
		})
	}

	for _, messege := range messages {
		fillInPlaceholder(messege, template)

		inputEmail := sesRepository.AssembleEmail(messege)
		result, err := sesRepository.SendEmail(ctx, inputEmail)

		// Display error messages if they occur.
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				fmt.Println(aerr.Error())
			} else {
				// Print the error, cast err to awserr.Error to get the Code and
				// Message from an error.
				fmt.Println(err.Error())
			}
			return
		}

		fmt.Println("Email Sent to address: " + messege.Recipient)
		fmt.Println("MessegeId: " + *result.MessageId)
	}

}

func fillInPlaceholder(message *model.MailMessege, template *model.Template) {
	message.Subject = template.Subject
	message.Body = template.Body
	for k, v := range template.Params.(map[string]interface{}) {
		regex := regexp.MustCompile(`\$\{` + k + `\}`)
		if message.Param != nil {
			switch k {
			case "name":
				message.Body = regex.ReplaceAllLiteralString(message.Body, message.Param.Name)
				continue
			case "age":
				message.Body = regex.ReplaceAllLiteralString(message.Body, fmt.Sprintf("%d", message.Param.Age))
				continue
			case "gender":
				message.Body = regex.ReplaceAllLiteralString(message.Body, fmt.Sprintf("%d", message.Param.Gender))
				continue
			case "residence":
				message.Body = regex.ReplaceAllLiteralString(message.Body, message.Param.Residence)
				continue
			default:
			}
		}
		switch v := v.(type) {
		case string:
			message.Body = regex.ReplaceAllLiteralString(message.Body, v)
		case float64:
			message.Body = regex.ReplaceAllLiteralString(message.Body, fmt.Sprintf("%g", v))
		case int:
			message.Body = regex.ReplaceAllLiteralString(message.Body, fmt.Sprintf("%d", v))
		default:
			log.Fatal("Invalid value type for placeholders in template")
		}
	}
}
