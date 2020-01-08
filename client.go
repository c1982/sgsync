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
		a.configuration.Source.AWSAccessKeyID,
		a.configuration.Source.AWSSectedAccessKey,
		[]string{a.configuration.Source.GroupID})

	if err != nil {
		return nil, err
	}

	return groups[0], nil
}

func (a *AWSClient) GetDestinationSecurityGroups() ([]*ec2.SecurityGroup, error) {

	sglist := []*ec2.SecurityGroup{}

	if len(a.configuration.Destinations) < 1 {
		return nil, fmt.Errorf("there are no destination security groups")
	}

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

	if len(rules) < 1 {
		return
	}

	err = a.executeIngress(rules, func(groupid *string, permissions []*ec2.IpPermission, svc *ec2.EC2) error {
		_, err := svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       groupid,
			IpPermissions: permissions,
		})

		return err
	})

	return err
}

func (a *AWSClient) AuthorizeEgress(rules []*ec2.SecurityGroup) (err error) {
	if len(rules) < 1 {
		return
	}

	err = a.executeEgress(rules, func(groupid *string, permissions []*ec2.IpPermission, svc *ec2.EC2) error {
		_, err := svc.AuthorizeSecurityGroupEgress(&ec2.AuthorizeSecurityGroupEgressInput{
			GroupId:       groupid,
			IpPermissions: permissions,
		})

		return err
	})

	return err
}

func (a *AWSClient) RevokeIngress(rules []*ec2.SecurityGroup) (err error) {

	if len(rules) < 1 {
		return
	}

	err = a.executeIngress(rules, func(groupid *string, permissions []*ec2.IpPermission, svc *ec2.EC2) error {
		_, err := svc.RevokeSecurityGroupIngress(&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       groupid,
			IpPermissions: permissions,
		})

		return err
	})

	return err
}

func (a *AWSClient) RevokeEgress(rules []*ec2.SecurityGroup) (err error) {

	if len(rules) < 1 {
		return
	}

	err = a.executeEgress(rules, func(groupid *string, permissions []*ec2.IpPermission, svc *ec2.EC2) error {
		_, err := svc.RevokeSecurityGroupEgress(&ec2.RevokeSecurityGroupEgressInput{
			GroupId:       groupid,
			IpPermissions: permissions,
		})

		return err
	})

	return err
}

func (a *AWSClient) executeIngress(rules []*ec2.SecurityGroup, action func(*string, []*ec2.IpPermission, *ec2.EC2) error) (err error) {

	for _, destination := range a.configuration.Destinations {

		svc := ec2.New(session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(destination.Region),
			Credentials: credentials.NewStaticCredentials(destination.AWSAccessKeyID, destination.AWSSectedAccessKey, "")},
		)))

		for _, rule := range rules {
			err := action(rule.GroupId, rule.IpPermissions, svc)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (a *AWSClient) executeEgress(rules []*ec2.SecurityGroup, action func(*string, []*ec2.IpPermission, *ec2.EC2) error) (err error) {

	for _, destination := range a.configuration.Destinations {

		svc := ec2.New(session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(destination.Region),
			Credentials: credentials.NewStaticCredentials(destination.AWSAccessKeyID, destination.AWSSectedAccessKey, "")},
		)))

		for _, rule := range rules {
			err := action(rule.GroupId, rule.IpPermissionsEgress, svc)
			if err != nil {
				break
			}
		}
	}

	return err
}

func (a *AWSClient) describeSecurityGroups(region, accesskey, secretkey string, groupIDs []string) ([]*ec2.SecurityGroup, error) {

	if len(groupIDs) < 1 {
		return nil, fmt.Errorf("security group id cannot be empty")
	}

	if region == "" {
		return nil, fmt.Errorf("region parameter cannot be empty")
	}

	if accesskey == "" {
		return nil, fmt.Errorf("access key cannot be empty for %s region", region)
	}

	if secretkey == "" {
		return nil, fmt.Errorf("secret key cannot be empty for %s region", region)
	}

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

	if len(output.SecurityGroups) == 0 {
		return nil, fmt.Errorf("security group (%v) cannot find on %s region\n",
			groupIDs,
			region)
	}

	return output.SecurityGroups, nil
}
