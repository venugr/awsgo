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

	userName := "testuser1"
	isOk, isErr := CreateUser(myConfig, myContext, userName)

	if isErr != nil {
		log.Fatalf("Error: User Create: %v\n", isErr)
	}

	if isOk {
		log.Printf("Info: User '%v' is created/already present.\n", userName)
	}

	groupName := "testgrp1"
	isOk, isErr = AddAnUserToGroup(myConfig, myContext, userName, groupName)
	if isErr != nil {
		log.Fatalf("Error: Unable to add user:'%v' to group:'%v'...%v\n", userName, groupName, isErr)
	}
	log.Printf("Info: User '%v' has been added to group '%v'\n", userName, groupName)

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
