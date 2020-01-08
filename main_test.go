package main

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-sdk-go/service/ec2"
)

var sourceGroup = []*ec2.SecurityGroup{
	(&ec2.SecurityGroup{
		Description: aws.String("source"),
		GroupId:     aws.String("gs-001"),
		GroupName:   aws.String("source group"),
		OwnerId:     aws.String("owner"),
		VpcId:       aws.String("wpc"),
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{
				FromPort:   aws.Int64(22),
				ToPort:     aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("10.5.5.1/32"),
					}),
				},
			}),
		},
		IpPermissionsEgress: []*ec2.IpPermission{
			(&ec2.IpPermission{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("0.0.0.0/0"),
					}),
				},
			}),
			(&ec2.IpPermission{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("5.5.5.2/32"),
					}),
				},
			}),
		},
	}),
}

var destinations = []*ec2.SecurityGroup{
	(&ec2.SecurityGroup{
		Description: aws.String("destination 2"),
		GroupId:     aws.String("gs-002"),
		GroupName:   aws.String("destination group 2"),
		OwnerId:     aws.String("owner-id"),
		VpcId:       aws.String("vpc-id"),
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{
				FromPort:   aws.Int64(22),
				ToPort:     aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("10.5.5.1/32"),
					}),
				},
			}),
			(&ec2.IpPermission{
				FromPort:   aws.Int64(22),
				ToPort:     aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("4.2.2.1/32"),
					}),
				},
			}),
		},
	}),
	(&ec2.SecurityGroup{
		Description: aws.String("destination 3"),
		GroupId:     aws.String("gs-003"),
		GroupName:   aws.String("destination group 3"),
		OwnerId:     aws.String("owner-id"),
		VpcId:       aws.String("vpc-id"),
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{
				FromPort:   aws.Int64(22),
				ToPort:     aws.Int64(22),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("8.8.8.8/32"),
					}),
				},
			}),
		},
	}),
	(&ec2.SecurityGroup{
		Description: aws.String("destination 4"),
		GroupId:     aws.String("gs-004"),
		GroupName:   aws.String("destination group 4"),
		OwnerId:     aws.String("owner-id"),
		VpcId:       aws.String("vpc-id"),
		IpPermissions: []*ec2.IpPermission{
			(&ec2.IpPermission{
				FromPort:   aws.Int64(8080),
				ToPort:     aws.Int64(8080),
				IpProtocol: aws.String("tcp"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("0.0.0.0/0"),
					}),
				},
			}),
		},
		IpPermissionsEgress: []*ec2.IpPermission{
			(&ec2.IpPermission{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("0.0.0.0/0"),
					}),
				},
			}),
		},
	}),
}

var dstengress = []*ec2.SecurityGroup{
	(&ec2.SecurityGroup{
		Description: aws.String("destination 4"),
		GroupId:     aws.String("gs-004"),
		GroupName:   aws.String("destination group 4"),
		OwnerId:     aws.String("owner-id"),
		VpcId:       aws.String("vpc-id"),
		IpPermissionsEgress: []*ec2.IpPermission{
			(&ec2.IpPermission{
				IpProtocol: aws.String("-1"),
				IpRanges: []*ec2.IpRange{
					(&ec2.IpRange{
						CidrIp: aws.String("0.0.0.0/0"),
					}),
				},
			}),
		},
	}),
}

func Test_WillbeAddedIngress(t *testing.T) {

	s := NewSync(sourceGroup[0], destinations)
	ingress := s.willbeAddedIngress()

	if len(ingress) != 2 {
		t.Errorf("ingress rules cannot prepare. got: %d, want: %d", len(ingress), 2)
	}

	fmt.Printf("%v", ingress)
}

func Test_WillbeAddedEgress(t *testing.T) {

	s := NewSync(sourceGroup[0], dstengress)
	egress := s.willbeAddedEgress()

	if len(egress) != 1 {
		t.Errorf("egress value not expected got: %d, want: %d", len(egress), 1)
	}

	fmt.Printf("%v", egress)
}

func Test_WillBeDeleteIngress(t *testing.T) {

	s := NewSync(sourceGroup[0], destinations)
	deleteds := s.willbeDeleteIngress()

	if len(deleteds) > 0 {
		t.Errorf("deleted value not expected got: %d, want: %d", len(deleteds), 1)
	}

	fmt.Printf("%v", deleteds)
}
