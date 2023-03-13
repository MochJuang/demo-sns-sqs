module publisher-01

go 1.20

require github.com/aws/aws-sdk-go-v2 v1.16.11

require (
	github.com/aws/aws-sdk-go-v2/config v1.17.1 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.12.14 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.12.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.18 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.12 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.19 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.12 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.11.17 // indirect
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.13 // indirect
	github.com/aws/smithy-go v1.12.1 // indirect
)

require (
	github.com/aws/aws-sdk-go-v2/service/sns v1.17.13
	github.com/aws/aws-sdk-go-v2/service/sqs v1.19.4
	gitlab.com/myhelpers/aws_helpers v0.0.0-00010101000000-000000000000
)

replace gitlab.com/myhelpers/aws_helpers => ../aws_helpers
