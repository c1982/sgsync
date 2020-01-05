package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var ruleTypeAllTraffic = func() string { return "-1" }

func main() {

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("eu-central-1")},
	)

	if err != nil {
		log.Fatal(err)
	}

	svc := ec2.New(sess)

	input := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("group-id"),
				Values: aws.StringSlice([]string{"sg-1", "sg-2"}),
			},
		},
	}

	sgout, err := svc.DescribeSecurityGroups(input)

	if err != nil {
		log.Fatal(err)
	}

	in, _ := aggregatePermissions(sgout.SecurityGroups)
	in = deduplicatePermissions(in)
	//out = deduplicatePermissions(out)

	for _, sg := range sgout.SecurityGroups {

		ingress := &ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       sg.GroupId,
			IpPermissions: in,
		}

		r, err := svc.AuthorizeSecurityGroupIngress(ingress)

		if err != nil {
			log.Println(err)
		} else {
			fmt.Println(r.String())
		}
	}

}

func aggregatePermissions(securityGroupList []*ec2.SecurityGroup) (inbound []*ec2.IpPermission, outbound []*ec2.IpPermission) {

	inbound = []*ec2.IpPermission{}
	outbound = []*ec2.IpPermission{}

	for i := 0; i < len(securityGroupList); i++ {
		sg := securityGroupList[i]
		inbound = append(inbound, sg.IpPermissions...)
		outbound = append(outbound, sg.IpPermissionsEgress...)
	}

	return inbound, outbound
}

func deduplicatePermissions(permissions []*ec2.IpPermission) []*ec2.IpPermission {

	tmp_ := []*ec2.IpPermission{}

	for _, perm := range permissions {
		index, ok := func(p *ec2.IpPermission) (int, bool) {
			for i := 0; i < len(tmp_); i++ {
				t := tmp_[i]

				//probably egress permission
				if p.FromPort == nil && p.ToPort == nil && *p.IpProtocol == ruleTypeAllTraffic() {
					return i, true
				}

				if *p.FromPort == *t.FromPort && *p.ToPort == *t.ToPort &&
					*p.IpProtocol == *t.IpProtocol {
					return i, true
				}
			}
			return -1, false
		}(perm)

		if ok {
			tmp_[index].IpRanges = append(tmp_[index].IpRanges, perm.IpRanges...)
			tmp_[index].Ipv6Ranges = append(tmp_[index].Ipv6Ranges, perm.Ipv6Ranges...)
			tmp_[index].IpRanges = deduplicateIpRanges(tmp_[index].IpRanges)
			tmp_[index].Ipv6Ranges = deduplicateIpv6Ranges(tmp_[index].Ipv6Ranges)
		} else {
			tmp_ = append(tmp_, &ec2.IpPermission{
				FromPort:   perm.FromPort,
				ToPort:     perm.ToPort,
				IpRanges:   perm.IpRanges,
				Ipv6Ranges: perm.Ipv6Ranges,
				IpProtocol: perm.IpProtocol,
			})
		}
	}

	return tmp_
}

func deduplicateIpRanges(ranges []*ec2.IpRange) []*ec2.IpRange {

	for i := 0; i < len(ranges); i++ {
		for z := i + 1; z < len(ranges); z++ {
			if *ranges[i].CidrIp == *ranges[z].CidrIp {
				ranges = append(ranges[:z], ranges[z+1:]...)
				z--
			}
		}
	}

	return ranges
}

func deduplicateIpv6Ranges(ranges []*ec2.Ipv6Range) []*ec2.Ipv6Range {

	for i := 0; i < len(ranges); i++ {
		for z := i + 1; z < len(ranges); z++ {
			if *ranges[i].CidrIpv6 == *ranges[z].CidrIpv6 {
				ranges = append(ranges[:z], ranges[z+1:]...)
				z--
			}
		}
	}

	return ranges
}
