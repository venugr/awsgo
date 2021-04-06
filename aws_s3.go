package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func DoS3(myConfig aws.Config, myContext context.Context) {

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

	DoSleep(20, "deleting object...")

	//*******************************************************************************
	keyPutObj := "putobject.mp3"
	isOk, isErr = DoPutObject(myConfig, myContext, bucketName, keyPutObj)

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
func DoPutObject(myConfig aws.Config, myContext context.Context, bucketName string, key string) (bool, error) {

	s3Input := &s3.PutObjectInput{
		Bucket: &bucketName,
		Key:    &key,
	}

	_, isErr := GoPutObject(myConfig, myContext, s3Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}
