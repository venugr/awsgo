package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func DoUser(myConfig aws.Config, myContext context.Context) {

	fmt.Println("Info: In DoUser()....")

	//*******************************************************************************
	log.Println()
	aliasName := "testaliasvenul"
	log.Println("Account Alias....Create.")
	isOk, isErr := CreateAccountAlias(myConfig, myContext, aliasName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create alias '%v': %v", aliasName, isErr)
	}

	log.Printf("Info: account alias '%v' is created.", aliasName)

	time.Sleep(10 * time.Second)

	log.Println()
	log.Println("Account Alias....Delete.")
	isOk, isErr = DeleteAccountAlias(myConfig, myContext, aliasName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete alias '%v': %v", aliasName, isErr)
	}

	log.Printf("Info: account alias '%v' is deleted.", aliasName)

	return

	//*******************************************************************************
	DisplayUsers(myConfig, myContext)
	userName := "testuser1"
	isOk, isErr = CreateUser(myConfig, myContext, userName)

	if isErr != nil {
		log.Fatalf("Error: User Create: %v\n", isErr)
	}

	if isOk {
		log.Printf("Info: User '%v' is created/already present.\n", userName)
	}
	DisplayUsers(myConfig, myContext)
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

	//*******************************************************************************
	newUserName := "DellUser1"
	isOk, isErr = ChangeUserName(myConfig, myContext, userName, newUserName)
	if isErr != nil || !isOk {
		if !strings.Contains(isErr.Error(), "EntityAlreadyExists") {
			log.Fatalf("Error: Change user name: %v\n", isErr)
		}
	}
	DisplayUsers(myConfig, myContext)

	//*******************************************************************************
	log.Println("\nDelete User....\n")
	DeleteAnUserFromGroup(myConfig, myContext, newUserName, groupName)
	isOk, isErr = DeleteUser(myConfig, myContext, newUserName)
	if isErr != nil || !isOk {
		log.Fatalf("Error: Delete user name: %v\n", isErr)
	}
	DisplayUsers(myConfig, myContext)

}

func DeleteAnUserFromGroup(myConfig aws.Config, myContext context.Context, userName string, groupName string) (bool, error) {
	_, isErr := DoesGroupExist(myConfig, myContext, groupName)
	if isErr != nil {
		log.Fatalf("Error: Unable to check Group '%v': %v\n", groupName, isErr)
	}

	_, lErr := GoDeleteUserFromGroup(myConfig, myContext, userName, groupName)
	if lErr != nil {
		return false, lErr
	}
	return true, nil

}

func GoDeleteUserFromGroup(myConfig aws.Config, myContext context.Context, userName string, groupName string) (*iam.RemoveUserFromGroupOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.RemoveUserFromGroupInput{
		GroupName: &groupName,
		UserName:  &userName,
	}

	return iamClient.RemoveUserFromGroup(myContext, iamInput)
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

func DisplayUsers(myConfig aws.Config, myContext context.Context) (bool, error) {
	iamResp, lErr := GoListUsers(myConfig, myContext)
	if lErr != nil {
		return false, lErr
	}

	log.Println("\n")
	log.Println("---------------------")
	log.Println("Users List")
	log.Println("---------------------")
	for _, user := range iamResp.Users {
		log.Println(*user.UserName)
	}
	log.Println("---------------------")
	log.Println()
	return true, nil
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

func GoUpdateUser(myConfig aws.Config, myContext context.Context, userName string, newUserName string) (*iam.UpdateUserOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.UpdateUserInput{
		UserName:    &userName,
		NewUserName: &newUserName,
	}

	return iamClient.UpdateUser(myContext, iamInput)
}

func ChangeUserName(myConfig aws.Config, myContext context.Context, userName string, newUserName string) (bool, error) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr
	}

	if !isOk {
		return false, nil
	}

	_, lErr := GoUpdateUser(myConfig, myContext, userName, newUserName)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoDeleteUser(myConfig aws.Config, myContext context.Context, userName string) (*iam.DeleteUserOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.DeleteUserInput{
		UserName: &userName,
	}

	return iamClient.DeleteUser(myContext, iamInput)
}

func DeleteUser(myConfig aws.Config, myContext context.Context, userName string) (bool, error) {
	isOk, isErr := DoesUserExist(myConfig, myContext, userName)
	if isErr != nil {
		return false, isErr
	}

	if !isOk {
		return false, nil
	}

	_, lErr := GoDeleteUser(myConfig, myContext, userName)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoCreateAccountAlias(myConfig aws.Config, myContext context.Context, aliasName string) (*iam.CreateAccountAliasOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.CreateAccountAliasInput{
		AccountAlias: &aliasName,
	}

	return iamClient.CreateAccountAlias(myContext, iamInput)

}

func CreateAccountAlias(myConfig aws.Config, myContext context.Context, aliasName string) (bool, error) {

	_, lErr := GoCreateAccountAlias(myConfig, myContext, aliasName)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoDeleteAccountAlias(myConfig aws.Config, myContext context.Context, aliasName string) (*iam.DeleteAccountAliasOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.DeleteAccountAliasInput{
		AccountAlias: &aliasName,
	}

	return iamClient.DeleteAccountAlias(myContext, iamInput)

}

func DeleteAccountAlias(myConfig aws.Config, myContext context.Context, aliasName string) (bool, error) {

	_, lErr := GoDeleteAccountAlias(myConfig, myContext, aliasName)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}

func GoListAccountAlias(myConfig aws.Config, myContext context.Context) (*iam.ListAccountAliasesOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.ListAccountAliasesInput{
		MaxItems: aws.Int32(int32(10)),
	}

	return iamClient.ListAccountAliases(myContext, iamInput)

	// aliasList := [:]

	// for _, alias := range iamResp.AccountAliases {
	// 	aliasList = append( aliasList, alias)
	// }

	//return iamClient.CreateAccountAlias(myContext, iamInput)

}

func GetAccountAliasList(myConfig aws.Config, myContext context.Context) (bool, error) {

	_, lErr := GoListAccountAlias(myConfig, myContext)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}
