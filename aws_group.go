package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

func DoGroup(myConfig aws.Config, myContext context.Context) {

}

func doesGroupExist(myConfig aws.Config, myContext context.Context, groupName string) (bool, error) {

	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.ListGroupsInput{}

	iamResp, lErr := iamClient.ListGroups(myContext, iamInput)
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

func CreateGroup(myConfig aws.Config, myContext context.Context, groupName string) (bool, error) {

	iamClient := iam.NewFromConfig(myConfig)

	iamInput := &iam.CreateGroupInput{
		GroupName: &groupName,
	}

	_, lErr := iamClient.CreateGroup(myContext, iamInput)

	if lErr != nil {
		return false, lErr
	}

	return true, nil
}
