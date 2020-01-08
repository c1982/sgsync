package main

type SyncConfig struct {
	Interval               string              `yaml:"interval"`
	DeleteDestinationRules bool                `yaml:"delete_destination_rules"`
	Source                 SourceConfig        `yaml:"source"`
	Destinations           []DestinationConfig `yaml:"destinations"`
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
