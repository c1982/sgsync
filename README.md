# AWS Security Group Synchronization Tool

Of course, you can easily do these things using AWC's VPC Peering feature. This is the best practice. This tool includes a complete dirty-hack.

This tool monitors the resource AWS Security Group that you specify and synchronizes it to multiple Security Groups that you specify. It applies both inbound and outbound rules to target security groups.

### Download

* [https://github.com/c1982/sgsync/releases](https://github.com/c1982/sgsync/releases)

### Installation

* copy sgsync binary file to /usr/local/bin/sgsync
* create config file to /etc/sgsync/config.yaml
* create service file to /etc/systemd/system/sgsync.service

sgsync.service:
```ini
[Unit]
Description=SGSYNC sevice
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
ExecStart=/usr/local/bin/sgsync --config=/etc/sgsync/config.yaml

[Install]
WantedBy=multi-user.target
```

### Configuration

config.yaml:

```yaml
interval: 5m

source:
  aws_access_key_id: YOUR_AWS_ACCESS_KEY_ID
  aws_secret_access_key: YOUR_AWS_SECRET_ACCESS_KEY
  region: eu-central-1
  group_id: sg-00000000000000001

destinations:
  - 
    aws_access_key_id: YOUR_AWS_ACCESS_KEY_ID
    aws_secret_access_key: YOUR_AWS_SECRET_ACCESS_KEY
    region: us-west-2
    group_ids: ["sg-00000000000000002","sg-00000000000000003"]
  - 
    aws_access_key_id: YOUR_AWS_ACCESS_KEY_ID
    aws_secret_access_key: YOUR_AWS_SECRET_ACCESS_KEY
    region: ap-east-1
    group_ids: ["sg-00000000000000004"]
```


# AWS Policy

AWS kullanıcısına tanımlayabileceğiniz inline policy aşağıdaki gibidir.

AWS IAM Policy:

```json
{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Sid": "VisualEditor0",
            "Effect": "Allow",
            "Action": [
                "ec2:RevokeSecurityGroupIngress",
                "ec2:AuthorizeSecurityGroupEgress",
                "ec2:AuthorizeSecurityGroupIngress",
                "ec2:UpdateSecurityGroupRuleDescriptionsEgress",
                "ec2:DescribeSecurityGroupReferences",
                "ec2:CreateSecurityGroup",
                "ec2:RevokeSecurityGroupEgress",
                "ec2:DescribeSecurityGroups",
                "ec2:UpdateSecurityGroupRuleDescriptionsIngress",
                "ec2:DescribeStaleSecurityGroups"
            ],
            "Resource": "*"
        }
    ]
}
```


## Contact

Oğuzhan - [@c1982](https://twitter.com/c1982) - aspsrc@gmail.com