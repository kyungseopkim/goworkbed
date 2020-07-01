package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type V2Vehicle map[string]struct{}

var exists = struct{}{}

func (v2 V2Vehicle) String() string {
	return v2.String()
}

func (v2 V2Vehicle) Contains(key string) bool {
	_, ok := v2[key]
	if ok {
		return true
	}
	return false
}

func (v2 V2Vehicle) FromS3() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1")},
	)

	bucketName := "java-library-repository"
	keyName := "resources/v2/v2vehicles.yaml"
	// Create S3 service client
	svc := s3.New(sess)
	req, resp := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})

	err = req.Send()
	if err != nil { // resp is now filled
		log.Fatalln(err)
	}

	objBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//v2s := V2Vehicles{}
	m := make(map[interface{}]interface{})
	err = yaml.Unmarshal(objBytes, &m)
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range m["vehicles"].([]string) {
		v2[v] = exists
	}
}
