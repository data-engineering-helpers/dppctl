module github.com/data-engineering-helpers/dppctl

go 1.20

require (
	github.com/aws/aws-sdk-go-v2 v1.17.7
	github.com/aws/aws-sdk-go-v2/config v1.18.18
	github.com/aws/aws-sdk-go-v2/service/codeartifact v1.17.1
	github.com/aws/aws-sdk-go-v2/service/ecr v1.18.7
	github.com/aws/aws-sdk-go-v2/service/mwaa v1.14.5
	github.com/aws/aws-sdk-go-v2/service/s3 v1.30.6
	github.com/aws/aws-sdk-go-v2/service/sts v1.18.6
	github.com/aws/smithy-go v1.13.5
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream v1.4.10 // indirect
	github.com/aws/aws-sdk-go-v2/credentials v1.13.17 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.13.0 // indirect
	github.com/aws/aws-sdk-go-v2/internal/configsources v1.1.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/endpoints/v2 v2.4.25 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.3.31 // indirect
	github.com/aws/aws-sdk-go-v2/internal/v4a v1.0.22 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding v1.9.11 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/checksum v1.1.25 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/presigned-url v1.9.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/internal/s3shared v1.13.24 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.12.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.14.5 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
)

replace github.com/data-engineering-helpers/dppctl/service => ./service

replace github.com/data-engineering-helpers/dppctl/utilities => ./utilities

replace github.com/data-engineering-helpers/dppctl/workflow => ./workflow
