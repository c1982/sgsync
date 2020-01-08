package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type AWSClient struct {
	configuration SyncConfig
}

func NewAWSClient(cfg SyncConfig) *AWSClient {
	return &AWSClient{
		configuration: cfg,
	}
}

func (a *AWSClient) GetSourceSecurityGroup() (*ec2.SecurityGroup, error) {

	groups, err := a.describeSecurityGroups(a.configuration.Source.Region,
		a.configuration.Source.AWSAccessKeyID, a.configuration.Source.AWSSectedAccessKey,
		[]string{a.configuration.Source.Region})

	if err != nil {
		return nil, err
	}

	return groups[0], nil
}

func (a *AWSClient) describeSecurityGroups(region, accesskey, secretkey string, groupIDs []string) ([]*ec2.SecurityGroup, error) {

	svc := ec2.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accesskey,
			secretkey, "")},
	)))

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("group-id"),
				Values: aws.StringSlice(groupIDs),
			},
		},
	}

	output, err := svc.DescribeSecurityGroups(input)

	if err != nil {
		return nil, err
	}

	if len(output.SecurityGroups) > 0 {
		return nil, fmt.Errorf("security group (%v) cannot find on %s region\n",
			groupIDs,
			region)
	}

	return output.SecurityGroups, nil
}

func (a *AWSClient) GetDestinationSecurityGroups() ([]*ec2.SecurityGroup, error) {

	sglist := []*ec2.SecurityGroup{}

	for _, dst := range a.configuration.Destinations {

		if len(dst.GroupIDs) == 0 {
			continue
		}

		groups, err := a.describeSecurityGroups(dst.Region, dst.AWSAccessKeyID, dst.AWSSectedAccessKey, dst.GroupIDs)

		if err != nil {
			return nil, err
		}

		sglist = append(sglist, groups...)
	}

	return sglist, nil
}

func (a *AWSClient) AuthorizeIngress(rules []*ec2.SecurityGroup) (err error) {

	for _, destination := range a.configuration.Destinations {

		svc := ec2.New(session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(destination.Region),
			Credentials: credentials.NewStaticCredentials(destination.AWSAccessKeyID, destination.AWSSectedAccessKey, "")},
		)))

		for _, rule := range rules {

			input := &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId:       rule.GroupId,
				IpPermissions: rule.IpPermissions,
			}

			_, err = svc.AuthorizeSecurityGroupIngress(input)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (a *AWSClient) AuthorizeEgress(rules []*ec2.SecurityGroup) (err error) {
	return nil
}

func (a *AWSClient) RevokeIngress(rules []*ec2.SecurityGroup) (err error) {
	return nil
}

func (a *AWSClient) RevokeEgress(rules []*ec2.SecurityGroup) (err error) {
	return nil
}
