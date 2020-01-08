package main

import (
	"github.com/aws/aws-sdk-go/service/ec2"
)

type syncGroup struct {
	source       *ec2.SecurityGroup
	destinations []*ec2.SecurityGroup
}

func NewSyncGroup(source *ec2.SecurityGroup, securityGroups []*ec2.SecurityGroup) *syncGroup {
	return &syncGroup{
		source:       source,
		destinations: securityGroups,
	}
}

func (s *syncGroup) willbeAddedIngress() []*ec2.SecurityGroup {

	willbeAddedIngress := []*ec2.SecurityGroup{}

	for _, dst := range s.destinations {

		newsg := s.popsg(dst)

		for _, src := range s.source.IpPermissions {

			exists, dstprm := s.isPermissionExists(src, dst.IpPermissions)

			if !exists {
				newperm := s.popperm(src)
				newperm.SetIpRanges(src.IpRanges)
				newperm.SetIpv6Ranges(src.Ipv6Ranges)
				newsg.IpPermissions = append(newsg.IpPermissions, newperm)
			} else {

				updatesg := s.popsg(dst)
				updateperm := s.popperm(dstprm)

				//ipv4
				for _, srcrange := range src.IpRanges {
					rangexists, _ := s.isIpRangeExists(srcrange, dstprm.IpRanges)
					if !rangexists {
						updateperm.IpRanges = append(updateperm.IpRanges, srcrange)
					}
				}

				//ipv6
				for _, srcrange6 := range src.Ipv6Ranges {

					rangexist, _ := s.isIpv6RangeExists(srcrange6, dstprm.Ipv6Ranges)
					if !rangexist {
						updateperm.Ipv6Ranges = append(updateperm.Ipv6Ranges, srcrange6)
					}
				}

				if len(updateperm.Ipv6Ranges) > 0 || len(updateperm.IpRanges) > 0 {

					updatesg.SetIpPermissions([]*ec2.IpPermission{
						updateperm,
					})

					willbeAddedIngress = append(willbeAddedIngress, updatesg)
				}
			}
		}

		if len(newsg.IpPermissions) > 0 {
			willbeAddedIngress = append(willbeAddedIngress, newsg)
		}
	}

	return willbeAddedIngress
}

func (s *syncGroup) willbeAddedEgress() []*ec2.SecurityGroup {

	willbeAddedEgress := []*ec2.SecurityGroup{}

	for _, dst := range s.destinations {

		newsg := s.popsg(dst)

		for _, src := range s.source.IpPermissionsEgress {

			exists, dstprm := s.isPermissionExists(src, dst.IpPermissionsEgress)

			if !exists {
				newperm := s.popperm(src)
				newperm.SetIpRanges(src.IpRanges)
				newperm.SetIpv6Ranges(src.Ipv6Ranges)
				newsg.IpPermissionsEgress = append(newsg.IpPermissionsEgress, newperm)
			} else {

				updatesg := s.popsg(dst)
				updateperm := s.popperm(dstprm)

				//ipv4
				for _, srcrange := range src.IpRanges {
					rangexists, _ := s.isIpRangeExists(srcrange, dstprm.IpRanges)
					if !rangexists {
						updateperm.IpRanges = append(updateperm.IpRanges, srcrange)
					}
				}

				//ipv6
				for _, srcrange6 := range src.Ipv6Ranges {
					rangexist, _ := s.isIpv6RangeExists(srcrange6, dstprm.Ipv6Ranges)

					if !rangexist {
						updateperm.Ipv6Ranges = append(updateperm.Ipv6Ranges, srcrange6)
					}
				}

				if len(updateperm.Ipv6Ranges) > 0 || len(updateperm.IpRanges) > 0 {
					updatesg.SetIpPermissionsEgress([]*ec2.IpPermission{
						updateperm,
					})

					willbeAddedEgress = append(willbeAddedEgress, updatesg)
				}
			}
		}

		if len(newsg.IpPermissionsEgress) > 0 {
			willbeAddedEgress = append(willbeAddedEgress, newsg)
		}
	}

	return willbeAddedEgress
}

func (s *syncGroup) willbeDeleteIngress() []*ec2.SecurityGroup {

	willbeDeleteIngress := []*ec2.SecurityGroup{}

	for _, dst := range s.destinations {

		newsg := s.popsg(dst)

		for _, dstpermission := range dst.IpPermissions {

			newperm := s.popperm(dstpermission)
			exists, srcpermission := s.isPermissionExists(dstpermission, s.source.IpPermissions)

			if !exists {
				newperm.SetIpRanges(dstpermission.IpRanges)
				newperm.SetIpv6Ranges(dstpermission.Ipv6Ranges)
			} else {

				for _, dstrange := range dstpermission.IpRanges {

					rangexists, _ := s.isIpRangeExists(dstrange, srcpermission.IpRanges)
					if !rangexists {
						newperm.IpRanges = append(newperm.IpRanges, dstrange)
					}
				}

				for _, dstrange := range dstpermission.Ipv6Ranges {

					rangexists, _ := s.isIpv6RangeExists(dstrange, srcpermission.Ipv6Ranges)
					if !rangexists {
						newperm.Ipv6Ranges = append(newperm.Ipv6Ranges, dstrange)
					}
				}
			}

			if len(newperm.IpRanges) > 0 || len(newperm.Ipv6Ranges) > 0 {
				newsg.IpPermissions = append(newsg.IpPermissions, newperm)
			}
		}

		if len(newsg.IpPermissions) > 0 {
			willbeDeleteIngress = append(willbeDeleteIngress, newsg)
		}
	}

	return willbeDeleteIngress
}

func (s *syncGroup) willbeDeleteEgress() []*ec2.SecurityGroup {

	willbeDeleteEgress := []*ec2.SecurityGroup{}

	for _, dst := range s.destinations {

		newsg := s.popsg(dst)

		for _, dstpermission := range dst.IpPermissionsEgress {

			newperm := s.popperm(dstpermission)
			exists, srcpermission := s.isPermissionExists(dstpermission, s.source.IpPermissionsEgress)

			if !exists {
				newperm.SetIpRanges(dstpermission.IpRanges)
				newperm.SetIpv6Ranges(dstpermission.Ipv6Ranges)
			} else {

				for _, dstrange := range dstpermission.IpRanges {

					rangexists, _ := s.isIpRangeExists(dstrange, srcpermission.IpRanges)
					if !rangexists {
						newperm.IpRanges = append(newperm.IpRanges, dstrange)
					}
				}

				for _, dstrange := range dstpermission.Ipv6Ranges {

					rangexists, _ := s.isIpv6RangeExists(dstrange, srcpermission.Ipv6Ranges)
					if !rangexists {
						newperm.Ipv6Ranges = append(newperm.Ipv6Ranges, dstrange)
					}
				}
			}

			if len(newperm.IpRanges) > 0 || len(newperm.Ipv6Ranges) > 0 {
				newsg.IpPermissionsEgress = append(newsg.IpPermissionsEgress, newperm)
			}
		}

		if len(newsg.IpPermissionsEgress) > 0 {
			willbeDeleteEgress = append(willbeDeleteEgress, newsg)
		}
	}

	return willbeDeleteEgress
}

func (s *syncGroup) isPermissionExists(src *ec2.IpPermission, destinations []*ec2.IpPermission) (bool, *ec2.IpPermission) {

	for _, dst := range destinations {
		if dst.FromPort != nil && dst.ToPort != nil && src.FromPort != nil && src.ToPort != nil {
			if *src.FromPort == *dst.FromPort && *src.ToPort == *dst.ToPort && *src.IpProtocol == *dst.IpProtocol {
				return true, dst
			}
		}

		if dst.FromPort == nil && dst.ToPort == nil && src.FromPort == nil && src.ToPort == nil {
			if *src.IpProtocol == *dst.IpProtocol {
				return true, dst
			}
		}
	}

	return false, nil
}

func (s *syncGroup) isIpRangeExists(src *ec2.IpRange, destinations []*ec2.IpRange) (bool, *ec2.IpRange) {

	for _, dst := range destinations {
		if *dst.CidrIp == *src.CidrIp {
			return true, dst
		}
	}

	return false, nil
}

func (s *syncGroup) isIpv6RangeExists(src *ec2.Ipv6Range, destinations []*ec2.Ipv6Range) (bool, *ec2.Ipv6Range) {

	for _, dst := range destinations {
		if *dst.CidrIpv6 == *src.CidrIpv6 {
			return true, dst
		}
	}

	return false, nil
}

func (s *syncGroup) popsg(sg *ec2.SecurityGroup) *ec2.SecurityGroup {

	newsg := &ec2.SecurityGroup{}

	if sg.GroupId != nil {
		newsg.SetGroupId(*sg.GroupId)
	}

	if sg.GroupName != nil {
		newsg.SetGroupName(*sg.GroupName)
	}

	if sg.Description != nil {
		newsg.SetDescription(*sg.Description)
	}

	if sg.OwnerId != nil {
		newsg.SetOwnerId(*sg.OwnerId)
	}

	if sg.VpcId != nil {
		newsg.SetVpcId(*sg.VpcId)
	}

	if sg.Tags != nil {
		newsg.SetTags(sg.Tags)
	}

	return newsg
}

func (s *syncGroup) popperm(perm *ec2.IpPermission) *ec2.IpPermission {

	newperm := &ec2.IpPermission{}

	if perm.FromPort != nil {
		newperm.SetFromPort(*perm.FromPort)
	}

	if perm.ToPort != nil {
		newperm.SetToPort(*perm.ToPort)
	}

	if perm.IpProtocol != nil {
		newperm.SetIpProtocol(*perm.IpProtocol)
	}

	return newperm
}
