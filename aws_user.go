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

	isOk, isErr := doesGroupExist(myConfig, myContext, groupName)
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

	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.AddUserToGroupInput{
		GroupName: &groupName,
		UserName:  &userName,
	}

	_, lErr := iamClient.AddUserToGroup(myContext, iamInput)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func doesUserExist(myConfig aws.Config, myContext context.Context, userName string) (bool, error) {

	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.ListUsersInput{}

	iamResp, lErr := iamClient.ListUsers(myContext, iamInput)

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

func CreateUser(myConfig aws.Config, myContext context.Context, userName string) (bool, error) {

	isOk, isErr := doesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr
	}

	if isOk {
		return true, nil
	}

	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.CreateUserInput{
		UserName: &userName,
	}

	_, lErr := iamClient.CreateUser(myContext, iamInput)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}
