package main

import "log"

func main() {

	cfg := SyncConfig{}
	client := NewAWSClient(cfg)

	src, err := client.GetSourceSecurityGroup()

	if err != nil {
		log.Fatal(err)
	}

	destinations, err := client.GetDestinationSecurityGroups()

	if err != nil {
		log.Fatal(err)
	}

	sync := NewSyncGroup(src, destinations)
	ingress := sync.willbeAddedIngress()
	egress := sync.willbeAddedEgress()
	revokeingress := sync.willbeDeleteIngress()

	err = client.AuthorizeIngress(ingress)
	if err != nil {

	}

	err = client.AuthorizeEgress(egress)
	if err != nil {

	}

	err = client.RevokeIngress(revokeingress)
	if err != nil {

	}
}
