package main

import (
	"bytes"
	"fmt"
	"html/template"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/joho/godotenv"
)

var (
	from    = "admin.sp@savvypool.com"
	to      = "harissownn@gmail.com"
	subject = "This is a test email from SavvyPool"
)

func init() {
	_ = godotenv.Load(".env")
}

func main() {
	t, err := template.ParseFiles("./test.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	// Render the email template with the OTP value
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000
	var tpl bytes.Buffer
	data := map[string]string{"OTP": strconv.Itoa(otp)}
	if err := t.Execute(&tpl, data); err != nil {
		fmt.Println(err)
		return
	}

	sess, err := session.NewSessionWithOptions(session.Options{
		Profile: "savvypool",
	})

	if err != nil {
		fmt.Println(err)
	}

	svc := ses.New(sess, &aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY"),
			os.Getenv("AWS_SECRET_KEY"),
			""),
	})
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: []*string{
				aws.String(to),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String("utf-8"),
					Data:    aws.String(tpl.String()),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String("utf-8"),
				Data:    aws.String(subject),
			},
		},
		Source: aws.String(from),
	}
	res, err := svc.SendEmail(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return
	}
	fmt.Println("Email has send to the destination : harissown@gmail.com")
	fmt.Println(res)
}
