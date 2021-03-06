package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
)

func DoEip(myConfig aws.Config, myContext context.Context) {

	//*******************************************************************************
	allcId, pubIp, isOk, isErr := AllocateEip(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to allocate EIP: %v", isErr)
	}

	log.Printf("Info: allocated IP '%v' id '%v'", pubIp, allcId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 20 secs...releasing public ip allocation")
	// time.Sleep(20 * time.Second)
	DoSleep(20, "releasing public ip allocation")
	isOk, isErr = ReleaseEip(myConfig, myContext, allcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to Release EIP '%v(%v)': %v", pubIp, allcId, isErr)
	}

	log.Printf("Info: released allocated IP '%v' id '%v'", pubIp, allcId)

}

func DoEc2(myConfig aws.Config, myContext context.Context) {
	//*******************************************************************************
	avlZone := "us-east-1a"
	volSize := int32(8)
	volId, isOk, isErr := DoCreateVolume(myConfig, myContext, avlZone, volSize)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create volume: %v", isErr)
	}

	log.Printf("Info: Volume Id '%v' is created", volId)

	for true {
		volStatus, isOk, isErr := DoDescribeVolumeStatus(myConfig, myContext, volId)
		if isErr != nil || !isOk {
			log.Fatalf("Error: unable to create snapshot: %v", isErr)
		}

		log.Printf("Info: Voloue '%v' status '%v'", volId, volStatus)
		if volStatus == "available" {
			break
		}

		DoSleep(2, "volume not available...")
	}

	snpId, isOk, isErr := DoCreateSnapshot(myConfig, myContext, volId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create snapshot: %v", isErr)
	}

	log.Printf("Info: Snapshot Id '%v' is created", snpId)

	cpysnpId, isOk, isErr := DoCopySnapshot(myConfig, myContext, "us-east-1", "us-west-2", snpId, "TestCopy Snapshot")

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to copy snapshot: %v", isErr)
	}

	log.Printf("Info: Snapshot Id '%v' is copy", cpysnpId)
	DoSleep(20, "deleting snapshot...")

	isOk, isErr = DoDeleteSnapshot(myConfig, myContext, cpysnpId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete snapshot: %v", isErr)
	}

	log.Printf("Info: Snapshot Id '%v' is deleted", cpysnpId)

	DoSleep(20, "deleting snapshot...")

	isOk, isErr = DoDeleteSnapshot(myConfig, myContext, snpId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete snapshot: %v", isErr)
	}

	log.Printf("Info: Snapshot Id '%v' is deleted", snpId)

	DoSleep(20, "deleting volume...")

	isOk, isErr = DoDeleteVolume(myConfig, myContext, volId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete volume: %v", isErr)
	}

	log.Printf("Info: Volume Id: %v deleted", volId)

}

func DoEc2_2(myConfig aws.Config, myContext context.Context) {
	//*******************************************************************************
	vpcId, isOk, isErr := CreateEc2VPC(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create VPC: %v", isErr)
	}

	log.Printf("Info: VPC Id: %v", vpcId)

	//*******************************************************************************
	keyName := "testkeypair"
	keyId, keyFingerPrint, keyPem, isOk, isErr := DoCreateKeyPair(myConfig, myContext, keyName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to cretae key-pair: %v", isErr)
	}

	log.Printf("Info: keypair '%v' is created, id: %v", keyName, keyId)

	log.Printf("Info: KeyFingerPrint:\n%v", keyFingerPrint)
	log.Printf("Info: KeyPem:\n%v", keyPem)

	DoSleep(20, "check keypair...")

	//*******************************************************************************

	isOk, isErr = DoDeleteKeyPair(myConfig, myContext, keyId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete key-pair: %v", isErr)
	}

	log.Printf("Info: keypair '%v' is deleted, id: %v", keyName, keyId)

	DoSleep(10, "check keypair deleted...")

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...")
	// time.Sleep(10 * time.Second)
	isOk, isErr = DeleteEc2VPC(myConfig, myContext, vpcId)
	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete VPC'%v': %v", vpcId, isErr)
	}

	log.Printf("Info: VPC '%v' is deleted.", vpcId)

}

func DoEc2_1(myConfig aws.Config, myContext context.Context) {
	//*******************************************************************************
	vpcId, isOk, isErr := CreateEc2VPC(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create VPC: %v", isErr)
	}

	log.Printf("Info: VPC Id: %v", vpcId)

	//*******************************************************************************
	sgId, isOk, isErr := DoCreateSecurityGroup(myConfig, myContext, vpcId, "TestGroup AWS SDK GO-V2", "testgroup")

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create security group: %v", isErr)
	}

	log.Printf("Info: Sg Id '%v' is created.", sgId)
	DoSleep(10, "check security group.")

	//*******************************************************************************
	cidrIp := "149.24.0.0/16"
	toPort := int32(22)
	protocallName := "tcp"
	fromPort := int32(22)
	isOk, isErr = DoAuthorizeSecurityGroupIngress(myConfig, myContext, sgId, cidrIp, fromPort, toPort, protocallName)
	if isErr != nil || !isOk {
		log.Printf("Error: unable to add ingress '%v' to  security group '%v': %v", cidrIp, sgId, isErr)
	}

	log.Printf("Info: Ingress '%v' is added to sgid '%v'.", cidrIp, sgId)
	DoSleep(20, "Check Ingress entry")

	//*******************************************************************************
	isOk, isErr = DoRevokeSecurityGroupIngress(myConfig, myContext, sgId, cidrIp, fromPort, toPort, protocallName)
	if isErr != nil || !isOk {
		log.Printf("Error: unable to revoke ingress '%v' from  security group '%v': %v", cidrIp, sgId, isErr)
	}

	log.Printf("Info: Ingress '%v' is revoked from sgid '%v'.", cidrIp, sgId)
	DoSleep(20, "Check Ingress entry")

	//*******************************************************************************

	isOk, isErr = DoDeleteSecurityGroupById(myConfig, myContext, sgId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete security group: %v", isErr)
	}

	log.Printf("Info: SG Id '%v' is deleted.", sgId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...")
	// time.Sleep(10 * time.Second)
	DoSleep(10, "")

	isOk, isErr = DeleteEc2VPC(myConfig, myContext, vpcId)
	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete VPC'%v': %v", vpcId, isErr)
	}

	log.Printf("Info: VPC '%v' is deleted.", vpcId)

}

func DoEc2Old(myConfig aws.Config, myContext context.Context) {

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
	snetId, isOk, isErr := CreateSubnetInVPC(myConfig, myContext, vpcId, "10.0.0.0/24")

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create VPC: %v", isErr)
	}

	log.Printf("Info: subnet Id: '%v' is generated", snetId)

	//*******************************************************************************
	rtId, isOk, isErr := CreateRouteTableInVPC(myConfig, myContext, vpcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create route-table: %v", isErr)
	}

	log.Printf("Info: route-table Id: '%v' is generated", rtId)

	//*******************************************************************************
	ascId, isOk, isErr := AssociateRouteTableToSubnet(myConfig, myContext, rtId, snetId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to associate route-table '%v' to subnet '%v': %v", rtId, snetId, isErr)
	}

	log.Printf("Info: route-table Id: '%v' is associated to subnet '%v'", rtId, snetId)
	log.Printf("Info: association Id: '%v'", ascId)

	//*******************************************************************************
	destCidrBlock := "1.2.3.4/32"
	isOk, isErr = CreateRouteInRouteTable(myConfig, myContext, rtId, destCidrBlock, igwId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to add route '%v' in route-table '%v': %v", destCidrBlock, rtId, isErr)
	}

	log.Printf("Info: route '%v' is route-table '%v' is added.", destCidrBlock, rtId)

	//*******************************************************************************
	allcId, pubIp, isOk, isErr := AllocateEip(myConfig, myContext)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to allocate EIP: %v", isErr)
	}

	log.Printf("Info: allocated IP '%v' id '%v'", pubIp, allcId)

	//*******************************************************************************
	ngwId, isOk, isErr := CreateNatGatewayInSubnet(myConfig, myContext, allcId, snetId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create NatGateway: %v", isErr)
	}

	log.Printf("Info: NAT gateway '%v' is created.", ngwId)

	//*******************************************************************************
	// instId, isOk, isErr := DoEc2Instances(myConfig, myContext)

	// if isErr != nil || !isOk {
	// 	log.Fatalf("Error: unable to create NatGateway: %v", isErr)
	// }

	// log.Printf("Info: Instance Id '%v' is created.", instId)
	// DoSleep(60)

	//*******************************************************************************
	sgId, isOk, isErr := DoCreateSecurityGroup(myConfig, myContext, vpcId, "TestGroup AWS SDK GO-V2", "testgroup")

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create security group: %v", isErr)
	}

	log.Printf("Info: Sg Id '%v' is created.", sgId)
	DoSleep(10, "check security group.")

	//*******************************************************************************
	cidrIp := "149.24.0.0/16"
	toPort := int32(22)
	fromPort := int32(-1)
	isOk, isErr = DoAuthorizeSecurityGroupIngress(myConfig, myContext, sgId, cidrIp, fromPort, toPort, "tcp")
	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to add ingress '%v' to  security group '%v': %v", cidrIp, sgId, isErr)
	}

	log.Printf("Info: Ingress '%v' is added to sgid '%v'.", cidrIp, sgId)
	DoSleep(20, "Check Ingress entry")

	//*******************************************************************************
	grpName := "testgroup1"
	tsgId, isOk, isErr := DoCreateSecurityGroup(myConfig, myContext, vpcId, "TestGroup AWS SDK GO-V2", grpName)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create security group: %v", isErr)
	}

	log.Printf("Info: Sg Id '%v' is created.", tsgId)
	DoSleep(10, "")

	//*******************************************************************************

	isOk, isErr = DoDeleteSecurityGroupById(myConfig, myContext, tsgId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete security group: %v", isErr)
	}

	log.Printf("Info: Sg name '%v' is deleted.", grpName)
	DoSleep(10, "")

	//*******************************************************************************
	isOk, isErr = DoDeleteSecurityGroupById(myConfig, myContext, sgId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete security group: %v", isErr)
	}

	log.Printf("Info: Sg Id '%v' is deleted.", sgId)
	DoSleep(10, "")

	//*******************************************************************************
	// log.Println("Info: Sleep for 20 secs...deleting nat gateway")
	// time.Sleep(20 * time.Second)
	DoSleep(20, "deleting nat gateway")

	isOk, isErr = DeleteNatGatewayInSubnet(myConfig, myContext, ngwId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to create NatGateway: %v", isErr)
	}

	log.Printf("Info: NAT gateway '%v' is created.", ngwId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 20 secs...releasing public ip allocation")
	// time.Sleep(20 * time.Second)
	DoSleep(20, "deleting nat gateway")

	isOk, isErr = ReleaseEip(myConfig, myContext, allcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to Release EIP '%v(%v)': %v", pubIp, allcId, isErr)
	}

	log.Printf("Info: released allocated IP '%v' id '%v'", pubIp, allcId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 20 secs...deleting route in route-table")
	// time.Sleep(20 * time.Second)
	DoSleep(10, "deleting route in route-table")

	isOk, isErr = DeleteRouteInRouteTable(myConfig, myContext, rtId, destCidrBlock)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete route '%v' in route-table '%v': %v", destCidrBlock, rtId, isErr)
	}

	log.Printf("Info: route '%v' in route-table '%v' is deleted", destCidrBlock, rtId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...disassociation of route table")
	// time.Sleep(10 * time.Second)
	DoSleep(10, "disassociation of route table")

	isOk, isErr = DisassociateRouteTableFromSubnet(myConfig, myContext, ascId)
	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to disassociate route-table '%v': %v", rtId, isErr)
	}

	log.Printf("Info: route-table '%v' is disassociated", rtId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...deleting routetable")
	// time.Sleep(10 * time.Second)
	DoSleep(10, "deleting routetable")
	isOk, isErr = DeleteRouteTableInVPC(myConfig, myContext, rtId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete route-table '%v': %v", rtId, isErr)
	}

	log.Printf("Info: route-table Id: '%v' is deleted", rtId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 20 secs...")
	// time.Sleep(20 * time.Second)
	DoSleep(10, "")

	isOk, isErr = DetachIgwFromVpc(myConfig, myContext, igwId, vpcId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to detach igw '%v' from  VPC '%v': %v", igwId, vpcId, isErr)
	}

	log.Printf("Info: detached igw '%v' from VPC '%v'", igwId, vpcId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...")
	// time.Sleep(10 * time.Second)
	DoSleep(10, "")

	isOk, isErr = DeleteIgw(myConfig, myContext, igwId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete IGW '%v: %v", igwId, isErr)
	}

	log.Printf("Info: deleted IGW Id: %v", igwId)

	//*******************************************************************************
	isOk, isErr = DeleteSubnetInVPC(myConfig, myContext, snetId)

	if isErr != nil || !isOk {
		log.Fatalf("Error: unable to delete subnet '%v': %v", snetId, isErr)
	}

	log.Printf("Info: subnet '%v' is deleted", snetId)

	//*******************************************************************************
	// log.Println("Info: Sleep for 10 secs...")
	// time.Sleep(10 * time.Second)
	DoSleep(10, "")

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

func GoCreateSubnet(myConfig aws.Config, myContext context.Context, params *ec2.CreateSubnetInput) (*ec2.CreateSubnetOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateSubnetInput{}
	}

	ec2Resp, isErr := ec2Client.CreateSubnet(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateSubnetInVPC(myConfig aws.Config, myContext context.Context, vpcId string, cidrBlock string) (string, bool, error) {

	ec2Input := &ec2.CreateSubnetInput{
		VpcId:     &vpcId,
		CidrBlock: &cidrBlock,
	}

	outResp, isErr := GoCreateSubnet(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.Subnet.SubnetId, true, nil
}

func GoDeleteSubnet(myConfig aws.Config, myContext context.Context, params *ec2.DeleteSubnetInput) (*ec2.DeleteSubnetOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteSubnetInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteSubnet(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteSubnetInVPC(myConfig aws.Config, myContext context.Context, snetId string) (bool, error) {

	ec2Input := &ec2.DeleteSubnetInput{
		SubnetId: &snetId,
	}

	_, isErr := GoDeleteSubnet(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoCreateRouteTable(myConfig aws.Config, myContext context.Context, params *ec2.CreateRouteTableInput) (*ec2.CreateRouteTableOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateRouteTableInput{}
	}

	ec2Resp, isErr := ec2Client.CreateRouteTable(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateRouteTableInVPC(myConfig aws.Config, myContext context.Context, vpcId string) (string, bool, error) {

	ec2Input := &ec2.CreateRouteTableInput{
		VpcId: &vpcId,
	}

	outResp, isErr := GoCreateRouteTable(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.RouteTable.RouteTableId, true, nil
}

func GoCreateRoute(myConfig aws.Config, myContext context.Context, params *ec2.CreateRouteInput) (*ec2.CreateRouteOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateRouteInput{}
	}

	ec2Resp, isErr := ec2Client.CreateRoute(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateRouteInRouteTable(myConfig aws.Config, myContext context.Context, rtId string, destCidrBlock string, igwId string) (bool, error) {

	ec2Input := &ec2.CreateRouteInput{
		RouteTableId:         &rtId,
		DestinationCidrBlock: &destCidrBlock,
		GatewayId:            &igwId,
	}

	_, isErr := GoCreateRoute(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoDeleteRoute(myConfig aws.Config, myContext context.Context, params *ec2.DeleteRouteInput) (*ec2.DeleteRouteOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteRouteInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteRoute(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteRouteInRouteTable(myConfig aws.Config, myContext context.Context, rtId string, destCidrBlock string) (bool, error) {

	ec2Input := &ec2.DeleteRouteInput{
		RouteTableId:         &rtId,
		DestinationCidrBlock: &destCidrBlock,
	}

	_, isErr := GoDeleteRoute(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoDeleteRouteTable(myConfig aws.Config, myContext context.Context, params *ec2.DeleteRouteTableInput) (*ec2.DeleteRouteTableOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteRouteTableInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteRouteTable(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteRouteTableInVPC(myConfig aws.Config, myContext context.Context, rtId string) (bool, error) {

	ec2Input := &ec2.DeleteRouteTableInput{
		RouteTableId: &rtId,
	}

	_, isErr := GoDeleteRouteTable(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoAssociateRouteTable(myConfig aws.Config, myContext context.Context, params *ec2.AssociateRouteTableInput) (*ec2.AssociateRouteTableOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.AssociateRouteTableInput{}
	}

	ec2Resp, isErr := ec2Client.AssociateRouteTable(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func AssociateRouteTableToSubnet(myConfig aws.Config, myContext context.Context, rtId string, snetId string) (string, bool, error) {

	ec2Input := &ec2.AssociateRouteTableInput{
		RouteTableId: &rtId,
		SubnetId:     &snetId,
	}

	outResp, isErr := GoAssociateRouteTable(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.AssociationId, true, nil
}

func GoDisassociateRouteTable(myConfig aws.Config, myContext context.Context, params *ec2.DisassociateRouteTableInput) (*ec2.DisassociateRouteTableOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DisassociateRouteTableInput{}
	}

	ec2Resp, isErr := ec2Client.DisassociateRouteTable(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DisassociateRouteTableFromSubnet(myConfig aws.Config, myContext context.Context, ascId string) (bool, error) {

	ec2Input := &ec2.DisassociateRouteTableInput{
		AssociationId: &ascId,
	}

	_, isErr := GoDisassociateRouteTable(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoAllocateAddress(myConfig aws.Config, myContext context.Context, params *ec2.AllocateAddressInput) (*ec2.AllocateAddressOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.AllocateAddressInput{}
	}

	ec2Resp, isErr := ec2Client.AllocateAddress(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func AllocateEip(myConfig aws.Config, myContext context.Context) (string, string, bool, error) {

	ec2Input := &ec2.AllocateAddressInput{}

	outResp, isErr := GoAllocateAddress(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", "", false, isErr
	}

	return *outResp.AllocationId, *outResp.PublicIp, true, nil
}

func GoReleaseAddress(myConfig aws.Config, myContext context.Context, params *ec2.ReleaseAddressInput) (*ec2.ReleaseAddressOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.ReleaseAddressInput{}
	}

	ec2Resp, isErr := ec2Client.ReleaseAddress(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func ReleaseEip(myConfig aws.Config, myContext context.Context, allcId string) (bool, error) {

	ec2Input := &ec2.ReleaseAddressInput{
		AllocationId: &allcId,
	}

	_, isErr := GoReleaseAddress(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoCreateNatGateway(myConfig aws.Config, myContext context.Context, params *ec2.CreateNatGatewayInput) (*ec2.CreateNatGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateNatGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.CreateNatGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func CreateNatGatewayInSubnet(myConfig aws.Config, myContext context.Context, allcId string, snetId string) (string, bool, error) {

	ec2Input := &ec2.CreateNatGatewayInput{
		AllocationId: &allcId,
		SubnetId:     &snetId,
	}

	outResp, isErr := GoCreateNatGateway(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.NatGateway.NatGatewayId, true, nil
}

func GoDeleteNatGateway(myConfig aws.Config, myContext context.Context, params *ec2.DeleteNatGatewayInput) (*ec2.DeleteNatGatewayOutput, error) {
	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteNatGatewayInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteNatGateway(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil
}

func DeleteNatGatewayInSubnet(myConfig aws.Config, myContext context.Context, ngwId string) (bool, error) {

	ec2Input := &ec2.DeleteNatGatewayInput{
		NatGatewayId: &ngwId,
	}

	_, isErr := GoDeleteNatGateway(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}
