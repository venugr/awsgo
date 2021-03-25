package main

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func DoEc2(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	vpcId, isOk, isErr := CreateEc2VPC(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create VPC: %v", isErr)
	}

	log.Printf("Info: VPC Id: %v", vpcId)

	//*******************************************************************************
	igwId, isOk, isErr := CreateIgw(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create IGW: %v", isErr)
	}

	log.Printf("Info: IGW Id: %v", igwId)

	//*******************************************************************************
	isOk, isErr = AttachIgwToVpc(myConfig, myContext, igwId, vpcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to attach igw '%v' to  VPC '%v': %v", igwId, vpcId, isErr)
	}

	log.Printf("Info: attached igw '%v' to VPC '%v'", igwId, vpcId)

	//*******************************************************************************
	log.Println("Info: Sleep for 20 secs...")
	time.Sleep(20 * time.Second)
	isOk, isErr = DetachIgwFromVpc(myConfig, myContext, igwId, vpcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to detach igw '%v' from  VPC '%v': %v", igwId, vpcId, isErr)
	}

	log.Printf("Info: detached igw '%v' from VPC '%v'", igwId, vpcId)

	//*******************************************************************************
	log.Println("Info: Sleep for 20 secs...")
	time.Sleep(20 * time.Second)
	isOk, isErr = DeleteIgw(myConfig, myContext, igwId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete IGW '%v: %v", igwId, isErr)
	}

	log.Printf("Info: deleted IGW Id: %v", igwId)

	//*******************************************************************************
	log.Println("Info: Sleep for 20 secs...")
	time.Sleep(20 * time.Second)

	isOk, isErr = DeleteEc2VPC(myConfig, myContext, vpcId)
	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete VPC'%v': %v", vpcId, isErr)
	}

	log.Printf("Info: VPC '%v' is deleted.", vpcId)

}

func GoCreateVpc(myConfig aws.Config, myContext context.Context, params *ec2.CreateVpcInput) (*ec2.CreateVpcOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateVpcInput{
			CidrBlock: aws.String(string("10.0.0.0/16")),
			TagSpecifications: []types.TagSpecification{
				types.TagSpecification{
					ResourceType: "vpc",
					Tags: []types.Tag{
						{
							Key:   aws.String("Name"),
							Value: aws.String("VPC-10.0.0.0/16-CIDR"),
						},
					},
				},
			},
		}

	}

	ec2Resp, isErr := ec2Client.CreateVpc(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateEc2VPC(myConfig aws.Config, myContext context.Context) (string, bool, error) {

	ec2Input := &ec2.CreateVpcInput{
		CidrBlock: aws.String(string("10.0.0.0/16")),
		// AmazonProvidedIpv6CidrBlock: true,
		TagSpecifications: []types.TagSpecification{
			types.TagSpecification{
				ResourceType: "vpc",
				Tags: []types.Tag{
					{
						Key:   aws.String("Name"),
						Value: aws.String("MyTestVPC"),
					},
					{
						Key:   aws.String("DND"),
						Value: aws.String("FALSE"),
					},
				},
			},
		},
	}

	outResp, isErr := GoCreateVpc(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	// log.Printf("Info: VPC is created: Id: %v", *outResp.Vpc.VpcId)

	return *outResp.Vpc.VpcId, true, nil

}

func GoDeleteVpc(myConfig aws.Config, myContext context.Context, params *ec2.DeleteVpcInput) (*ec2.DeleteVpcOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteVpcInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteVpc(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteEc2VPC(myConfig aws.Config, myContext context.Context, vpcId string) (bool, error) {

	ec2Input := &ec2.DeleteVpcInput{
		VpcId: &vpcId,
	}

	_, isErr := GoDeleteVpc(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil

}

func GoCreateInternetGatewayInput(myConfig aws.Config, myContext context.Context, params *ec2.CreateInternetGatewayInput) (*ec2.CreateInternetGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateInternetGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.CreateInternetGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateIgw(myConfig aws.Config, myContext context.Context) (string, bool, error) {

	ec2Input := &ec2.CreateInternetGatewayInput{}

	outResp, isErr := GoCreateInternetGatewayInput(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.InternetGateway.InternetGatewayId, true, nil

}

func GoDeleteInternetGatewayInput(myConfig aws.Config, myContext context.Context, params *ec2.DeleteInternetGatewayInput) (*ec2.DeleteInternetGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteInternetGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteInternetGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteIgw(myConfig aws.Config, myContext context.Context, igwId string) (bool, error) {

	ec2Input := &ec2.DeleteInternetGatewayInput{
		InternetGatewayId: &igwId,
	}

	_, isErr := GoDeleteInternetGatewayInput(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil

}

func GoAttachInternetGateway(myConfig aws.Config, myContext context.Context, params *ec2.AttachInternetGatewayInput) (*ec2.AttachInternetGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.AttachInternetGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.AttachInternetGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func AttachIgwToVpc(myConfig aws.Config, myContext context.Context, igwId string, vpcId string) (bool, error) {

	ec2Input := &ec2.AttachInternetGatewayInput{
		VpcId:             &vpcId,
		InternetGatewayId: &igwId,
	}

	_, isErr := GoAttachInternetGateway(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil

}

func GoDetachInternetGateway(myConfig aws.Config, myContext context.Context, params *ec2.DetachInternetGatewayInput) (*ec2.DetachInternetGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DetachInternetGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.DetachInternetGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DetachIgwFromVpc(myConfig aws.Config, myContext context.Context, igwId string, vpcId string) (bool, error) {

	ec2Input := &ec2.DetachInternetGatewayInput{
		VpcId:             &vpcId,
		InternetGatewayId: &igwId,
	}

	_, isErr := GoDetachInternetGateway(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil

}
