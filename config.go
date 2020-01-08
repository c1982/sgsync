package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type SyncConfig struct {
	Interval               string              `yaml:"interval"`
	DeleteDestinationRules bool                `yaml:"delete_destination_rules"`
	Source                 SourceConfig        `yaml:"source"`
	Destinations           []DestinationConfig `yaml:"destinations"`
}

func NewSyncConfig(path string) (*SyncConfig, error) {
	cfg := &SyncConfig{}

	dat, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dat, cfg)
	return cfg, err
}

type SourceConfig struct {
	AWSAccessKeyID     string `yaml:"aws_access_key_id"`
	AWSSectedAccessKey string `yaml:"aws_secret_access_key"`
	Region             string `yaml:"region"`
	GroupID            string `yaml:"group_id"`
}

type DestinationConfig struct {
	AWSAccessKeyID     string   `yaml:"aws_access_key_id"`
	AWSSectedAccessKey string   `yaml:"aws_secret_access_key"`
	Region             string   `yaml:"region"`
	GroupIDs           []string `yaml:"group_ids"`
}
