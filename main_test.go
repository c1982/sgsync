package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
)

var permissions = []*ec2.IpPermission{
	&ec2.IpPermission{
		IpProtocol: aws.String("udp"),
		FromPort:   aws.Int64(4222),
		ToPort:     aws.Int64(4222),
		IpRanges: []*ec2.IpRange{
			&ec2.IpRange{
				CidrIp:      aws.String("10.10.10.1/32"),
				Description: aws.String("")},
			&ec2.IpRange{
				CidrIp:      aws.String("172.168.10.2/32"),
				Description: aws.String("")},
			&ec2.IpRange{
				CidrIp:      aws.String("4.2.2.2/32"),
				Description: aws.String("dns server")},
		},
	},
	&ec2.IpPermission{
		IpProtocol: aws.String("tcp"),
		FromPort:   aws.Int64(22),
		ToPort:     aws.Int64(22),
		IpRanges: []*ec2.IpRange{
			&ec2.IpRange{
				CidrIp:      aws.String("8.8.8.8/32"),
				Description: aws.String("")},
		},
	},
	&ec2.IpPermission{
		IpProtocol: aws.String("udp"),
		FromPort:   aws.Int64(4222),
		ToPort:     aws.Int64(4222),
		IpRanges: []*ec2.IpRange{
			&ec2.IpRange{
				CidrIp:      aws.String("8.8.8.8/32"),
				Description: aws.String("google dns server")},
		},
	},
	&ec2.IpPermission{
		IpProtocol: aws.String("udp"),
		FromPort:   aws.Int64(4222),
		ToPort:     aws.Int64(4222),
		IpRanges: []*ec2.IpRange{
			&ec2.IpRange{
				CidrIp:      aws.String("8.8.8.8/32"),
				Description: aws.String("")},
		},
	},
}

func Test_Deduplicate_Permissions(t *testing.T) {

	duplicateds := deduplicatePermissions(permissions)

	if len(duplicateds) != 2 {
		t.Errorf("permission not aggregated, got: %d, want: 2", len(duplicateds))
	}

	for _, iprange := range duplicateds {
		if *iprange.FromPort == 4222 && *iprange.ToPort == 4222 && *iprange.IpProtocol == "udp" {
			if len(iprange.IpRanges) != 4 {
				t.Errorf("4222 port rules ip, got: %d, want: %d.", len(iprange.IpRanges), 4)
			}
		}
	}
}

func Test_Deduplicate_IpRanges(t *testing.T) {

	ipranges := []*ec2.IpRange{
		&ec2.IpRange{
			CidrIp:      aws.String("10.10.10.1/32"),
			Description: aws.String("")},
		&ec2.IpRange{
			CidrIp:      aws.String("10.10.10.1/32"),
			Description: aws.String("")},
		&ec2.IpRange{
			CidrIp:      aws.String("4.2.2.2/32"),
			Description: aws.String("dns server")},
		&ec2.IpRange{
			CidrIp:      aws.String("4.2.2.2/32"),
			Description: aws.String("")},
	}

	duplicateds := deduplicateIpRanges(ipranges)

	if len(duplicateds) != 2 {
		t.Errorf("iprange not duplicated, got: %d, want %d", len(duplicateds), 2)
	}
}

func Test_Deduplicate_Ipv6Ranges(t *testing.T) {

	ipv6ranges := []*ec2.Ipv6Range{
		&ec2.Ipv6Range{
			CidrIpv6:    aws.String("::/0"),
			Description: aws.String("")},
		&ec2.Ipv6Range{
			CidrIpv6:    aws.String("::/0"),
			Description: aws.String("")},
		&ec2.Ipv6Range{
			CidrIpv6: aws.String("::/128")},
		&ec2.Ipv6Range{
			CidrIpv6:    aws.String("::/128"),
			Description: aws.String("")},
	}

	duplicateds := deduplicateIpv6Ranges(ipv6ranges)

	if len(duplicateds) != 2 {
		t.Errorf("ipv6range not duplicated, got: %d, want %d", len(duplicateds), 2)
	}
}
