package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func DoS3(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	bucketName := "myaudiomp3"
	isOk, isErr := DoPutBucketTagging(myConfig, myContext, bucketName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to put tagging: %v", isErr)
	}

}

func DoS3_2(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	buckets, isOk, isErr := DoListBuckets(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to list buckets: %v", isErr)
	}

	log.Println("---------------------------------")
	log.Println("\t\tBuckets")
	log.Println("---------------------------------")
	cnt := 1
	for _, bucket := range strings.Split(buckets, ",") {
		if bucket == "" {
			continue
		}
		log.Printf("%v.%v", cnt, bucket)
		cnt++
	}

	if cnt == 1 {
		log.Printf("<No Buckets>")
	}

	log.Println("---------------------------------")

	//*******************************************************************************
	for _, bucketName := range strings.Split(buckets, ",") {

		if bucketName == "" {
			continue
		}

		objects, isOk, isErr := DoListObjects(myConfig, myContext, bucketName)

		if isErr != nil || !isOk {
			log.Println()
			log.Printf("Error: unable to list objects in bucket '%v': %v", bucketName, isErr)
			log.Println()
			continue
		}

		log.Println()
		log.Println("---------------------------------")
		log.Printf("\tObjects in '%v'\n", bucketName)
		log.Println("---------------------------------")
		cnt = 1
		for _, object := range strings.Split(objects, ",") {
			if object == "" {
				continue
			}
			log.Printf("%v.%v", cnt, object)
			cnt++
		}

		if cnt == 1 {
			log.Printf("<No Objects>")
		}

		log.Println("---------------------------------")
		log.Println()
	}

}

func DoS3_1(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	getBucket := "myaudiomp3"
	getKey := "London.mp3"

	objReader, isOk, isErr := DoGetObject(myConfig, myContext, getBucket, getKey)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to get object '%v': %v", getKey, isErr)
	}

	log.Printf("Info: body: %v", objReader)

	outFile, isErr := os.Create("/tmp/my_" + getKey)
	_, isErr = io.Copy(outFile, objReader)

	return

	//*******************************************************************************
	bucketName := "mymestbucketdellme"
	location, isOk, isErr := DoCreateBucket(myConfig, myContext, bucketName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create bucket '%v: %v", bucketName, isErr)
	}

	log.Printf("Info: bucket '%v' is created, and location '%v'", bucketName, location)

	//*******************************************************************************
	srcObject := "myaudiomp3/London.mp3"
	key := "test.mp3"
	isOk, isErr = DoCopyObject(myConfig, myContext, bucketName, srcObject, key)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to copy object '%v': %v", srcObject, isErr)
	}

	log.Printf("Info: object '%v' is copied to as '%v'", srcObject, key)

	//*******************************************************************************
	keyPutObj := "putobject.mp3"
	isOk, isErr = DoPutObject(myConfig, myContext, bucketName, keyPutObj, strings.NewReader("My text for Put Object"))

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to put object '%v': %v", keyPutObj, isErr)
	}

	log.Printf("Info: object '%v' is uploaded", keyPutObj)
	DoSleep(20, "deleting object...")

	//*******************************************************************************
	isOk, isErr = DoDeleteObject(myConfig, myContext, bucketName, keyPutObj)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete object '%v': %v", keyPutObj, isErr)
	}

	log.Printf("Info: object '%v' is deleted", keyPutObj)

	DoSleep(20, "deleting object...")

	//*******************************************************************************
	isOk, isErr = DoDeleteObject(myConfig, myContext, bucketName, key)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete object '%v': %v", key, isErr)
	}

	log.Printf("Info: object '%v' is deleted", key)

	DoSleep(20, "deleting bucket...")

	//*******************************************************************************
	isOk, isErr = DoDeleteBucket(myConfig, myContext, bucketName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete bucket '%v: %v", bucketName, isErr)
	}

	log.Printf("Info: bucket '%v' is deleted", bucketName)

}

func GoCreateBucket(myConfig aws.Config, myContext context.Context, params *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.CreateBucketInput{}
	}

	s3Resp, isErr := s3Client.CreateBucket(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoCreateBucket(myConfig aws.Config, myContext context.Context, bucketName string) (string, bool, error) {

	s3Input := &s3.CreateBucketInput{
		Bucket: &bucketName,
	}

	outResp, isErr := GoCreateBucket(myConfig, myContext, s3Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.Location, true, nil
}

func GoDeleteBucket(myConfig aws.Config, myContext context.Context, params *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.DeleteBucketInput{}
	}

	s3Resp, isErr := s3Client.DeleteBucket(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoDeleteBucket(myConfig aws.Config, myContext context.Context, bucketName string) (bool, error) {

	s3Input := &s3.DeleteBucketInput{
		Bucket: &bucketName,
	}

	_, isErr := GoDeleteBucket(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoCopyObject(myConfig aws.Config, myContext context.Context, params *s3.CopyObjectInput) (*s3.CopyObjectOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.CopyObjectInput{}
	}

	s3Resp, isErr := s3Client.CopyObject(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoCopyObject(myConfig aws.Config, myContext context.Context, bucketName string, srcObject string, key string) (bool, error) {

	s3Input := &s3.CopyObjectInput{
		Bucket:     &bucketName,
		CopySource: &srcObject,
		Key:        &key,
	}

	_, isErr := GoCopyObject(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoDeleteObject(myConfig aws.Config, myContext context.Context, params *s3.DeleteObjectInput) (*s3.DeleteObjectOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.DeleteObjectInput{}
	}

	s3Resp, isErr := s3Client.DeleteObject(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoDeleteObject(myConfig aws.Config, myContext context.Context, bucketName string, key string) (bool, error) {

	s3Input := &s3.DeleteObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	_, isErr := GoDeleteObject(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoPutObject(myConfig aws.Config, myContext context.Context, params *s3.PutObjectInput) (*s3.PutObjectOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.PutObjectInput{}
	}

	s3Resp, isErr := s3Client.PutObject(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoPutObject(myConfig aws.Config, myContext context.Context, bucketName string, key string, srcData io.Reader) (bool, error) {

	s3Input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   srcData,
	}

	_, isErr := GoPutObject(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoGetObject(myConfig aws.Config, myContext context.Context, params *s3.GetObjectInput) (*s3.GetObjectOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.GetObjectInput{}
	}

	s3Resp, isErr := s3Client.GetObject(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoGetObject(myConfig aws.Config, myContext context.Context, bucketName string, key string) (io.ReadCloser, bool, error) {

	s3Input := &s3.GetObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	outResp, isErr := GoGetObject(myConfig, myContext, s3Input)

	if isErr != nil {
		return nil, false, isErr
	}

	return outResp.Body, true, nil
}

func GoListBuckets(myConfig aws.Config, myContext context.Context, params *s3.ListBucketsInput) (*s3.ListBucketsOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.ListBucketsInput{}
	}

	s3Resp, isErr := s3Client.ListBuckets(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoListBuckets(myConfig aws.Config, myContext context.Context) (string, bool, error) {

	s3Input := &s3.ListBucketsInput{}

	outResp, isErr := GoListBuckets(myConfig, myContext, s3Input)

	if isErr != nil {
		return "", false, isErr
	}

	buckets := ""

	for _, bucket := range outResp.Buckets {

		buckets += *bucket.Name + ","
	}

	return buckets, true, nil
}

func GoListObjects(myConfig aws.Config, myContext context.Context, params *s3.ListObjectsInput) (*s3.ListObjectsOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.ListObjectsInput{}
	}

	s3Resp, isErr := s3Client.ListObjects(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoListObjects(myConfig aws.Config, myContext context.Context, bucketName string) (string, bool, error) {

	s3Input := &s3.ListObjectsInput{
		Bucket: &bucketName,
	}

	outResp, isErr := GoListObjects(myConfig, myContext, s3Input)

	if isErr != nil {
		return "", false, isErr
	}

	objects := ""

	for _, content := range outResp.Contents {
		objects += *content.Key + "(size:" + fmt.Sprint(content.Size) + "),"
	}

	return objects, true, nil
}

func GoPutBucketTagging(myConfig aws.Config, myContext context.Context, params *s3.PutBucketTaggingInput) (*s3.PutBucketTaggingOutput, error) {

	s3Client := s3.NewFromConfig(myConfig)

	if params == nil {
		params = &s3.PutBucketTaggingInput{}
	}

	s3Resp, isErr := s3Client.PutBucketTagging(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return s3Resp, nil

}
func DoPutBucketTagging(myConfig aws.Config, myContext context.Context, bucketName string) (bool, error) {

	s3Input := &s3.PutBucketTaggingInput{
		Bucket: &bucketName,
		Tagging: &types.Tagging{
			TagSet: []types.Tag{
				{
					Key:   aws.String(string("testTagName")),
					Value: aws.String(string("Name1")),
				},
				{
					Key:   aws.String(string("Name")),
					Value: aws.String(string("Audio")),
				},
				{
					Key:   aws.String(string("Type")),
					Value: aws.String(string("MP3")),
				},
				{
					Key:   aws.String(string("DND")),
					Value: aws.String(string("FALSE")),
				},
			},
		},
	}

	_, isErr := GoPutBucketTagging(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}
