package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func DoGroup(myConfig aws.Config, myContext context.Context) {

}

func GoListGroups(myConfig aws.Config, myContext context.Context) (*iam.ListGroupsOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.ListGroupsInput{}

	return iamClient.ListGroups(myContext, iamInput)
}

func DoesGroupExist(myConfig aws.Config, myContext context.Context, groupName string) (bool, error) {

	iamResp, lErr := GoListGroups(myConfig, myContext)
	if lErr != nil {
		return false, lErr
	}

	for _, group := range iamResp.Groups {
		if *group.GroupName == groupName {
			return true, nil
		}
	}

	return false, nil
}

func GoCreateGroup(myConfig aws.Config, myContext context.Context, groupName string) (*iam.CreateGroupOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.CreateGroupInput{
		GroupName: &groupName,
	}

	return iamClient.CreateGroup(myContext, iamInput)
}

func CreateGroup(myConfig aws.Config, myContext context.Context, groupName string) (bool, error) {

	_, lErr := GoCreateGroup(myConfig, myContext, groupName)
	if lErr != nil {
		return false, lErr
	}

	return true, nil
}
