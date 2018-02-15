package main

import (
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
	"gopkg.in/urfave/cli.v1"
	"bytes"
	"fmt"
)

func main() {
	app := cli.NewApp()
	app.Name = "aws-s3-peek"
	app.Version = "0.1"
	app.Usage = "Preview s3 objects"
	app.Flags = flags()

	app.Action = func(c *cli.Context) error {
		bucket := c.String("bucket")
		key := c.String("key")

		if bucket == "" {
			fmt.Println("Bucket is missing")
			cli.ShowAppHelpAndExit(c, 1)
		}

		if key == "" {
			fmt.Println("Key is missing")
			cli.ShowAppHelpAndExit(c, 1)
		}

		var s3Cli = s3.New(session.Must(session.NewSessionWithOptions(session.Options{
			Config:            aws.Config{Region: aws.String(region())},
			Profile:           profile(c),
			SharedConfigState: session.SharedConfigEnable,
		})))

		output, e := s3Cli.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			Range:  aws.String(c.String("range")),
		})

		if e != nil {
			panic(e)
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(output.Body)
		fmt.Println(buf.String())

		return nil
	}

	app.Run(os.Args)
}

func flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "p, profile",
			Usage: "The aws profile to use",
		},
		cli.StringFlag{
			Name:  "b, bucket",
			Usage: "The s3 bucket name",
		},
		cli.StringFlag{
			Name:  "k, key",
			Usage: "The s3 object key",
		},
		cli.StringFlag{
			Name:  "r, range",
			Value: "bytes=0-1024",
			Usage: "The byte range to retrieve",
		},
	}
}

func region() string {
	return getEnv("AWS_DEFAULT_REGION", "eu-west-1")
}

func profile(c *cli.Context) string {
	profile := c.String("profile")
	if profile != "" {
		return profile
	}
	return getEnv("AWS_DEFAULT_PROFILE", "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
