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

func GoCreateKeyPair(myConfig aws.Config, myContext context.Context, params *ec2.CreateKeyPairInput) (*ec2.CreateKeyPairOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateKeyPairInput{}
	}

	ec2Resp, isErr := ec2Client.CreateKeyPair(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoCreateKeyPair(myConfig aws.Config, myContext context.Context, keyName string) (string, string, string, bool, error) {

	ec2Input := &ec2.CreateKeyPairInput{
		KeyName: &keyName,
	}

	outResp, isErr := GoCreateKeyPair(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", "", "", false, isErr
	}

	return *outResp.KeyPairId, *outResp.KeyFingerprint, *outResp.KeyMaterial, true, nil
}

func GoDeleteKeyPair(myConfig aws.Config, myContext context.Context, params *ec2.DeleteKeyPairInput) (*ec2.DeleteKeyPairOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteKeyPairInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteKeyPair(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoDeleteKeyPair(myConfig aws.Config, myContext context.Context, keyId string) (bool, error) {

	ec2Input := &ec2.DeleteKeyPairInput{
		KeyPairId: &keyId,
	}

	_, isErr := GoDeleteKeyPair(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoCreateVolume(myConfig aws.Config, myContext context.Context, params *ec2.CreateVolumeInput) (*ec2.CreateVolumeOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateVolumeInput{}
	}

	ec2Resp, isErr := ec2Client.CreateVolume(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoCreateVolume(myConfig aws.Config, myContext context.Context, avlZone string, size int32) (string, bool, error) {

	ec2Input := &ec2.CreateVolumeInput{
		AvailabilityZone: &avlZone,
		Size:             size,
	}

	outResp, isErr := GoCreateVolume(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.VolumeId, true, nil
}

func GoDeleteVolume(myConfig aws.Config, myContext context.Context, params *ec2.DeleteVolumeInput) (*ec2.DeleteVolumeOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteVolumeInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteVolume(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoDeleteVolume(myConfig aws.Config, myContext context.Context, volId string) (bool, error) {

	ec2Input := &ec2.DeleteVolumeInput{
		VolumeId: &volId,
	}

	_, isErr := GoDeleteVolume(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoCreateSnapshot(myConfig aws.Config, myContext context.Context, params *ec2.CreateSnapshotInput) (*ec2.CreateSnapshotOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CreateSnapshotInput{}
	}

	ec2Resp, isErr := ec2Client.CreateSnapshot(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoCreateSnapshot(myConfig aws.Config, myContext context.Context, volId string) (string, bool, error) {

	ec2Input := &ec2.CreateSnapshotInput{
		VolumeId: &volId,
	}

	outResp, isErr := GoCreateSnapshot(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.SnapshotId, true, nil
}

func GoCopySnapshot(myConfig aws.Config, myContext context.Context, params *ec2.CopySnapshotInput) (*ec2.CopySnapshotOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.CopySnapshotInput{}
	}

	ec2Resp, isErr := ec2Client.CopySnapshot(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoCopySnapshot(myConfig aws.Config, myContext context.Context, srcRegion string, destRegion string, srcSnapId string, desc string) (string, bool, error) {

	ec2Input := &ec2.CopySnapshotInput{
		SourceRegion:     &srcRegion,
		SourceSnapshotId: &srcSnapId,
		Description:      &desc,
	}

	outResp, isErr := GoCopySnapshot(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	return *outResp.SnapshotId, true, nil
}

func GoDeleteSnapshot(myConfig aws.Config, myContext context.Context, params *ec2.DeleteSnapshotInput) (*ec2.DeleteSnapshotOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DeleteSnapshotInput{}
	}

	ec2Resp, isErr := ec2Client.DeleteSnapshot(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoDeleteSnapshot(myConfig aws.Config, myContext context.Context, snapId string) (bool, error) {

	ec2Input := &ec2.DeleteSnapshotInput{
		SnapshotId: &snapId,
	}

	_, isErr := GoDeleteSnapshot(myConfig, myContext, ec2Input)

	if isErr != nil {
		return false, isErr
	}

	return true, nil
}

func GoDescribeVolumes(myConfig aws.Config, myContext context.Context, params *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {

	ec2Client := ec2.NewFromConfig(myConfig)

	if params == nil {
		params = &ec2.DescribeVolumesInput{}
	}

	ec2Resp, isErr := ec2Client.DescribeVolumes(myContext, params)

	if isErr != nil {
		return nil, isErr
	}

	return ec2Resp, nil

}

func DoDescribeVolumeStatus(myConfig aws.Config, myContext context.Context, volId string) (string, bool, error) {

	ec2Input := &ec2.DescribeVolumesInput{}

	outResp, isErr := GoDescribeVolumes(myConfig, myContext, ec2Input)

	if isErr != nil {
		return "", false, isErr
	}

	volStatus := ""

	for _, volume := range outResp.Volumes {

		if volId == *volume.VolumeId {
			volStatus = string(volume.State)
		}
	}
	return volStatus, true, nil
}
