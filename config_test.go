package main

import (
	"log"
	"testing"

	yaml "gopkg.in/yaml.v2"
)

var testdata = `interval: 1m
delete_destination_rules: true

source:
  aws_access_key_id: ACCESSKEY
  aws_secret_access_key: SECRETKEY
  region: eu-central-1
  group_id: sg-1

destinations:
  - 
    aws_access_key_id: ACCESSKEY1
    aws_secret_access_key: SECRETKEY1
    region: us-west-1
    group_ids: ["sg-1","sg-2"]
  - 
    aws_access_key_id: ACCESSKEY2
    aws_secret_access_key: SECRETKEY2
    region: eu-east-1
    group_ids: ["sg-1","sg-2","sg-3"]`

func Test_ParseConfig(t *testing.T) {

	cfg := SyncConfig{}

	err := yaml.Unmarshal([]byte(testdata), &cfg)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	if cfg.Interval != "1m" {
		t.Errorf("interval value not expected got: %s, want: %s", cfg.Interval, "1m")
	}

	if cfg.Source.AWSAccessKeyID == "ACCESSKEY1" {
		t.Errorf("source accesskeyid not expected got: %s, want: %s", cfg.Source.AWSAccessKeyID, "ACCESSKEY1")
	}

	if len(cfg.Destinations) != 2 {
		t.Errorf("destinations array size not expected got: %d, want: %d", len(cfg.Destinations), 2)
	}

	if len(cfg.Destinations[1].GroupIDs) != 3 {
		t.Errorf("destinations group array size not expected got: %d, want: %d", len(cfg.Destinations[1].GroupIDs), 3)
	}
}
