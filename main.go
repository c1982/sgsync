package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
)

func main() {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	configPath := flag.String("config", "config.yaml", "--config=config.yaml")
	flag.Parse()

	if *configPath == "" {
		log.Error().Msg("Config parameter value cannot be empty. Please pass --config parameter value")
		return
	}

	cfg, err := NewSyncConfig(*configPath)

	if err != nil {
		log.Error().Str("config_file", *configPath).Err(err).Msg("configuration error")
		return
	}

	interval, err := time.ParseDuration(cfg.Interval)

	if err != nil {
		log.Error().Err(err).Msg("interval cannot be expected format")
		return
	}

	client := NewAWSClient(*cfg)

	log.Info().Msg("scheduling started")
	log.Info().Str("interval", cfg.Interval).Msg("it has been set interval")

	for range time.Tick(interval) {
		src, dsts, err := describeAWSSecurityGroups(client)

		if err != nil {
			log.Error().Err(err).Msg("cannot describe security groups from AWS")
			continue
		}

		err = executeSecurityGroupFunctions(src, dsts, client)

		if err != nil {
			log.Error().Err(err).Msg("operation error")
		}
	}
}

func describeAWSSecurityGroups(client *AWSClient) (src *ec2.SecurityGroup, dsts []*ec2.SecurityGroup, err error) {

	src, err = client.GetSourceSecurityGroup()

	if err != nil {
		return nil, nil, err
	}

	dsts, err = client.GetDestinationSecurityGroups()

	if err != nil {
		return nil, nil, err
	}

	return src, dsts, err
}

func executeSecurityGroupFunctions(src *ec2.SecurityGroup, destinations []*ec2.SecurityGroup, client *AWSClient) (err error) {

	sync := NewSyncGroup(src, destinations)
	ingress := sync.willbeAddedIngress()
	egress := sync.willbeAddedEgress()
	revokeingress := sync.willbeDeleteIngress()
	revokeegress := sync.willbeDeleteEgress()

	log.Debug().Str("ingress", fmt.Sprintf("%v", ingress)).Msg("ingress operations executing")
	err = client.AuthorizeIngress(ingress)
	if err != nil {
		log.Debug().Err(err).Msg("ingress operation error")
		return err
	}

	log.Debug().Str("egress", fmt.Sprintf("%v", egress)).Msg("egress operations executing")
	err = client.AuthorizeEgress(egress)
	if err != nil {
		log.Debug().Err(err).Msg("egress operation error")
		return err
	}

	log.Debug().Str("revokeingress", fmt.Sprintf("%v", revokeingress)).Msg("revokeingress operations executing")
	err = client.RevokeIngress(revokeingress)
	if err != nil {
		log.Debug().Err(err).Msg("revokeingress operation error")
		return err
	}

	log.Debug().Str("revokeegress", fmt.Sprintf("%v", revokeegress)).Msg("revokeegress operations executing")
	err = client.RevokeEgress(revokeegress)
	if err != nil {
		log.Debug().Err(err).Msg("revokeegress operation error")
	}
	return err
}
