package main

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
)

type StatementEntry struct {
	Effect   string
	Action   []string
	Resource string
}

// PolicyDocument is our definition of our policies to be uploaded to AWS Identity and Access Management (IAM).
type PolicyDocument struct {
	Version   string
	Statement []StatementEntry
}

func DoPolicy(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	policyName := "myTestPolicy3"
	isOk, isErr, policyArn := CreatePolicy(myConfig, myContext, policyName)
	if isErr != nil {
		log.Fatalf("Error: Policy create '%v': %v", policyName, isErr)
	}

	if !isOk {
		log.Fatalf("Error: policy '%v' is generated.", policyName)
	}

	log.Printf("Info: Policy ARN '%v'", policyArn)

	//*******************************************************************************
	//policyArn := "arn:aws:iam::764524628794:policy/myTestPolicy1"
	isOk, isErr, policyDoc := GetPolicy(myConfig, myContext, policyArn)
	if isErr != nil {
		log.Fatalf("Error: Policy get '%v': %v", policyArn, isErr)
	}

	if !isOk {
		log.Fatalf("Error: policy '%v' is generated.", policyArn)
	}

	log.Printf("Info: Policy Doc:'%v'", policyDoc)

	//*******************************************************************************
	userName := "testuser"
	isOk, isErr = AttachPolicyToAnUser(myConfig, myContext, policyArn, userName)
	if isErr != nil {
		log.Fatalf("Error: unable to attach policy to user: %v", isErr)
	}

	if !isOk {
		log.Fatalf("Error: policy '%v' is not attached.", policyArn)
	}

	log.Printf("Info: Policy is attached to user '%v'", userName)

	//*******************************************************************************

	log.Println("Waiting for 60 sec....")
	time.Sleep(60 * time.Second)

	isOk, isErr = DettachPolicyToAnUser(myConfig, myContext, policyArn, userName)
	if isErr != nil {
		log.Fatalf("Error: unable to dettach policy from user: %v", isErr)
	}

	if !isOk {
		log.Fatalf("Error: policy '%v' is not dettached.", policyArn)
	}

	log.Printf("Info: Policy is dettached from user '%v'", userName)

	//*******************************************************************************
	//policyArn := "arn:aws:iam::764524628794:policy/myTestPolicy1"
	isOk, isErr = DeletePolicy(myConfig, myContext, policyArn)
	if isErr != nil {
		log.Fatalf("Error: Policy delete '%v': %v", policyName, isErr)
	}

	if !isOk {
		log.Fatalf("Error: policy '%v' is not delete.", policyName)
	}

	log.Printf("Info: Policy is delete:'%v'", policyName)

}

func CreatePolicyDoc() ([]byte, error) {

	policy := PolicyDocument{
		Version: "2012-10-17",
		Statement: []StatementEntry{
			{
				Effect: "Allow",
				Action: []string{
					"logs:CreateLogGroup", // Allow for creating log groups
				},
				Resource: "arn:aws:s3:::myaudiomp3:*",
			},
			{
				Effect: "Allow",
				// Allows for DeleteItem, GetItem, PutItem, Scan, and UpdateItem
				Action: []string{
					"dynamodb:DeleteItem",
					"dynamodb:GetItem",
					"dynamodb:PutItem",
					"dynamodb:Scan",
					"dynamodb:UpdateItem",
				},
				Resource: "arn:aws:dynamodb:us-east-1:764524628794:table/InfoDevsUsers:*",
			},
		},
	}

	b, err := json.Marshal(&policy)

	return b, err

}

func GoCreatePolicy(myConfig aws.Config, myContext context.Context, policyName string, policyDoc []byte) (*iam.CreatePolicyOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.CreatePolicyInput{
		PolicyDocument: aws.String(string(policyDoc)),
		PolicyName:     &policyName,
		Description:    aws.String(string("my test Policy thru AWS SDK GO-V2")),
	}

	return iamClient.CreatePolicy(myContext, iamInput)
}

func CreatePolicy(myConfig aws.Config, myContext context.Context, policyName string) (bool, error, string) {

	policyDoc, isErr := CreatePolicyDoc()

	if isErr != nil {
		log.Fatalf("Error: policy doc %v", isErr)
	}

	iamResp, isErr := GoCreatePolicy(myConfig, myContext, policyName, policyDoc)

	if isErr != nil {
		return false, isErr, ""
	}

	return true, nil, *iamResp.Policy.Arn
}

func GoDeletePolicy(myConfig aws.Config, myContext context.Context, policyArn string) (*iam.DeletePolicyOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.DeletePolicyInput{
		PolicyArn: &policyArn,
	}

	return iamClient.DeletePolicy(myContext, iamInput)
}

func DeletePolicy(myConfig aws.Config, myContext context.Context, policyArn string) (bool, error) {

	_, isErr := GoDeletePolicy(myConfig, myContext, policyArn)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoGetPolicy(myConfig aws.Config, myContext context.Context, policyArn string) (*iam.GetPolicyOutput, error) {

	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.GetPolicyInput{
		PolicyArn: &policyArn,
	}

	return iamClient.GetPolicy(myContext, iamInput)
}

func GetPolicy(myConfig aws.Config, myContext context.Context, policyArn string) (bool, error, string) {

	description := ""

	iamResp, isErr := GoGetPolicy(myConfig, myContext, policyArn)

	if isErr != nil {
		return false, isErr, ""
	}

	if iamResp.Policy == nil {
		description = "Policy Nil"
	} else {
		if iamResp.Policy.Description == nil {
			description = "Description Nil"
		} else {
			description = *iamResp.Policy.Description
		}

	}

	return true, nil, description
}

func GoAttachUserPolicy(myConfig aws.Config, myContext context.Context, policyArn string, userName string) (*iam.AttachUserPolicyOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.AttachUserPolicyInput{
		PolicyArn: &policyArn,
		UserName:  &userName,
	}

	return iamClient.AttachUserPolicy(myContext, iamInput)
}

func AttachPolicyToAnUser(myConfig aws.Config, myContext context.Context, policyArn string, userName string) (bool, error) {

	isOk, isErr := CreateUser(myConfig, myContext, userName)

	if isErr != nil {
		return false, isErr
	}

	if !isOk {
		return false, nil
	}

	_, isErr = GoAttachUserPolicy(myConfig, myContext, policyArn, userName)
	if isErr != nil {
		return false, isErr
	}

	return true, nil

}

func GoDettachUserPolicy(myConfig aws.Config, myContext context.Context, policyArn string, userName string) (*iam.DetachUserPolicyOutput, error) {
	iamClient := iam.NewFromConfig(myConfig)
	iamInput := &iam.DetachUserPolicyInput{
		PolicyArn: &policyArn,
		UserName:  &userName,
	}

	return iamClient.DetachUserPolicy(myContext, iamInput)
}

func DettachPolicyToAnUser(myConfig aws.Config, myContext context.Context, policyArn string, userName string) (bool, error) {

	isOk, isErr := DoesUserExist(myConfig, myContext, userName)

	if isErr != nil {
		return false, isErr
	}

	if !isOk {
		return false, nil
	}

	_, isErr = GoDettachUserPolicy(myConfig, myContext, policyArn, userName)
	if isErr != nil {
		return false, isErr
	}

	return true, nil

}
