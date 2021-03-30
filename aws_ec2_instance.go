package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

func DoEc2Instances(myConfig aws.Config, myContext context.Context) (string, bool, error) {

	amiId := "ami-0533f2ba8a1995cf9"
	instanceId, isOK, isErr := RunEc2Instaces(myConfig, myContext, amiId)

	if isErr != nil {
		log.Printf("Error: unable to start EC2 instance: %v", isErr)
		return "", false, isErr

	}

	if !isOK {
		log.Printf("Error: unable to start EC2 instance: %v", isErr)
		return "", false, isErr
	}

	log.Printf("Info: Instance Id: %v", instanceId)

	return instanceId, true, nil

}

func GoRuninstace(myConfig aws.Config, myContext context.Context, params *ec2.RunInstancesInput) (*ec2.RunInstancesOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.RunInstancesInput{}
	}

	ec2Resp, isErr := ec2Client.RunInstances(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func RunEc2Instaces(myConfig aws.Config, myContext context.Context, amiId string) (string, bool, error) {

	ec2Input := &ec2.RunInstancesInput{
		MaxCount: 1,
		MinCount: 1,
		ImageId:  &amiId,
	}

	outResp, isErr := GoRuninstace(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.Instances[0].InstanceId, true, nil
}

func GoCreateSecurityGroup(myConfig aws.Config, myContext context.Context, params *ec2.CreateSecurityGroupInput) (*ec2.CreateSecurityGroupOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateSecurityGroupInput{}
	}

	ec2Resp, isErr := ec2Client.CreateSecurityGroup(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoCreateSecurityGroup(myConfig aws.Config, myContext context.Context, vpcId string, desc string, grpName string) (string, bool, error) {

	ec2Input := &ec2.CreateSecurityGroupInput{
		GroupName:   &grpName,
		VpcId:       &vpcId,
		Description: &desc,
	}

	outResp, isErr := GoCreateSecurityGroup(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.GroupId, true, nil
}

func GoDeleteSecurityGroup(myConfig aws.Config, myContext context.Context, params *ec2.DeleteSecurityGroupInput) (*ec2.DeleteSecurityGroupOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteSecurityGroupInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteSecurityGroup(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoDeleteSecurityGroupById(myConfig aws.Config, myContext context.Context, grpId string) (bool, error) {

	ec2Input := &ec2.DeleteSecurityGroupInput{
		GroupId: &grpId,
	}

	_, isErr := GoDeleteSecurityGroup(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func DoDeleteSecurityGroupByName(myConfig aws.Config, myContext context.Context, grpName string) (bool, error) {

	ec2Input := &ec2.DeleteSecurityGroupInput{
		GroupName: &grpName,
	}

	_, isErr := GoDeleteSecurityGroup(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoAuthorizeSecurityGroupIngress(myConfig aws.Config, myContext context.Context, params *ec2.AuthorizeSecurityGroupIngressInput) (*ec2.AuthorizeSecurityGroupIngressOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.AuthorizeSecurityGroupIngressInput{}
	}

	ec2Resp, isErr := ec2Client.AuthorizeSecurityGroupIngress(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoAuthorizeSecurityGroupIngress(myConfig aws.Config, myContext context.Context, sgId string, cidrIp string,
	fromPort int32, toPort int32, protocallName string) (bool, error) {

	ec2Input := &ec2.AuthorizeSecurityGroupIngressInput{
		CidrIp:     &cidrIp,
		FromPort:   fromPort,
		ToPort:     toPort,
		GroupId:    &sgId,
		IpProtocol: &protocallName,
	}

	_, isErr := GoAuthorizeSecurityGroupIngress(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoRevokeSecurityGroupIngress(myConfig aws.Config, myContext context.Context, params *ec2.RevokeSecurityGroupIngressInput) (*ec2.RevokeSecurityGroupIngressOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.RevokeSecurityGroupIngressInput{}
	}

	ec2Resp, isErr := ec2Client.RevokeSecurityGroupIngress(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoRevokeSecurityGroupIngress(myConfig aws.Config, myContext context.Context, sgId string, cidrIp string,
	fromPort int32, toPort int32, protocallName string) (bool, error) {

	ec2Input := &ec2.RevokeSecurityGroupIngressInput{
		CidrIp:     &cidrIp,
		FromPort:   fromPort,
		ToPort:     toPort,
		GroupId:    &sgId,
		IpProtocol: &protocallName,
	}

	_, isErr := GoRevokeSecurityGroupIngress(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}
