package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func DoUser(myConfig aws.Config, myContext context.Context) {

	fmt.Println("Info: In DoUser()....")

	//*******************************************************************************
	userName := "testuser1"
	isOk, isErr := CreateUser(myConfig, myContext, userName)

	if isErr != nil {
		log.Fatalf("Error: User Create: %v\n", isErr)
	}

	if isOk {
		log.Printf("Info: User '%v' is created/already present.\n", userName)
	}

	//*******************************************************************************
	groupName := "testgrp1"
	isOk, isErr = AddAnUserToGroup(myConfig, myContext, userName, groupName)
	if isErr != nil {
		log.Fatalf("Error: Unable to add user:'%v' to group:'%v'...%v\n", userName, groupName, isErr)
	}
	log.Printf("Info: User '%v' has been added to group '%v'\n", userName, groupName)

	//*******************************************************************************
	isOk, isErr, accKey, secId := CreateAccessKeysForUser(myConfig, myContext, userName)
	if isErr != nil {
		log.Fatalf("Error: Access key Create: %v\n", isErr)
	}

	if isOk {
		log.Printf("Info: User '%v' AccessKeys created.\n", userName)
		log.Printf("Info: AccessKey: %v\n", accKey)
		log.Printf("Info: SecrKeyId: %v", secId)

	}

	//*******************************************************************************
	isOk1, isErr1, accKeyMap := ListAccessKeysForUser(myConfig, myContext, userName)
	if isErr1 != nil {
		log.Fatalf("Error: Access key Create: %v\n", isErr1)
	}

	if isOk1 {
		log.Printf("Info: Keys and Status")
		for k, v := range accKeyMap {
			log.Printf("%v ==> %v", k, v)
		}
	}

	//*******************************************************************************
	if isOk && accKey != "" {

		// time.Sleep(200 * time.Second)
		isOk, isErr = DeleteAccessKeysForUser(myConfig, myContext, userName, accKey)
		if isErr != nil {
			log.Fatalf("Error: Access key Deleye: %v\n", isErr)
		}

		if isOk {
			log.Printf("Info: User '%v' AccessKey '%v' is deleted.\n", userName, accKey)
		}

	}
}

func AddAnUserToGroup(myConfig aws.Config, myContext context.Context, userName string, groupName string) (bool, error) {

	isOk, isErr := DoesGroupExist(myConfig, myContext, groupName)
	if isErr != nil {
		log.Fatalf("Error: Unable to check Group '%v': %v\n", groupName, isErr)
	}

	if !isOk {
		log.Printf("Info: Group '%v' not present.\n", groupName)
		log.Printf("Info: Creating group '%v'...\n", groupName)
		isOk, isErr = CreateGroup(myConfig, myContext, groupName)

		if isErr != nil {
			return false, isErr
		}

		if !isOk {
			log.Printf("Error: Unable to create group '%v'\n", groupName)
			return false, nil
		}
		log.Printf("Info: Created group '%v'.\n", groupName)

	}

	_, lErr := GoAddUserToGroup(myConfig, myContext, userName, groupName)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoAddUserToGroup(myConfig aws.Config, myContext context.Context, userName string, groupName string) (*iam.AddUserToGroupOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.AddUserToGroupInput{
		GroupName: &groupName,
		UserName:  &userName,
	}

	return iamClient.AddUserToGroup(myContext, iamInput)
}

func GoListUsers(myConfig aws.Config, myContext context.Context) (*iam.ListUsersOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.ListUsersInput{}

	return iamClient.ListUsers(myContext, iamInput)
}

func DoesUserExist(myConfig aws.Config, myContext context.Context, userName string) (bool, error) {

	iamResp, lErr := GoListUsers(myConfig, myContext)
	if lErr != nil {
		return false, lErr
	}

	for _, user := range iamResp.Users {

		if *user.UserName == userName {
			return true, nil
		}
	}

	return false, nil
}

func GoCreateUser(myConfig aws.Config, myContext context.Context, userName string) (*iam.CreateUserOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.CreateUserInput{
		UserName: &userName,
	}

	return iamClient.CreateUser(myContext, iamInput)

}

func CreateUser(myConfig aws.Config, myContext context.Context, userName string) (bool, error) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr
	}

	if isOk {
		return true, nil
	}

	_, lErr := GoCreateUser(myConfig, myContext, userName)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoCreateAccessKey(myConfig aws.Config, myContext context.Context, userName string) (*iam.CreateAccessKeyOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.CreateAccessKeyInput{
		UserName: &userName,
	}

	return iamClient.CreateAccessKey(myContext, iamInput)
}

func CreateAccessKeysForUser(myConfig aws.Config, myContext context.Context, userName string) (bool, error, string, string) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr, "", ""
	}

	if !isOk {
		return false, nil, "", ""
	}

	iamResp, lErr := GoCreateAccessKey(myConfig, myContext, userName)

	if lErr != nil {
		return false, lErr, "", ""
	}

	return true, nil, *iamResp.AccessKey.AccessKeyId, *iamResp.AccessKey.SecretAccessKey
}

func GoDeleteAccessKey(myConfig aws.Config, myContext context.Context, userName string, accKey string) (*iam.DeleteAccessKeyOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.DeleteAccessKeyInput{
		UserName:    &userName,
		AccessKeyId: &accKey,
	}

	return iamClient.DeleteAccessKey(myContext, iamInput)
}

func DeleteAccessKeysForUser(myConfig aws.Config, myContext context.Context, userName string, accKey string) (bool, error) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr
	}

	if !isOk {
		return false, nil
	}

	_, lErr := GoDeleteAccessKey(myConfig, myContext, userName, accKey)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoListAccessKeys(myConfig aws.Config, myContext context.Context, userName string) (*iam.ListAccessKeysOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.ListAccessKeysInput{
		UserName: &userName,
		MaxItems: aws.Int32(int32(10)),
	}

	return iamClient.ListAccessKeys(myContext, iamInput)
}

func ListAccessKeysForUser(myConfig aws.Config, myContext context.Context, userName string) (bool, error, map[string]string) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr, nil
	}

	if !isOk {
		return false, nil, nil
	}

	iamResp, lErr := GoListAccessKeys(myConfig, myContext, userName)

	if lErr != nil {
		return false, lErr, nil
	}

	accKeyStatus := make(map[string]string)

	for _, key := range iamResp.AccessKeyMetadata {
		accKeyStatus[*key.AccessKeyId] = string(key.Status)
	}
	return true, nil, accKeyStatus
}
