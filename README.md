Getting started with Go for Data Processing Pipeline (DPP) CLI tool
===================================================================

[![OpenSSF
Scorecard](https://api.securityscorecards.dev/projects/github.com/data-engineering-helpers/dppctl/badge)](https://api.securityscorecards.dev/projects/github.com/data-engineering-helpers/dppctl)

# Overview
[This project](https://github.com/data-engineering-helpers/dppctl)
intends to develop and maintain a command-line (CLI) utility in Go
to help deploy data engineering pipelines on modern data stack (MDS).

Even though the members of the GitHub organization may be employed by
some companies, they speak on their personal behalf and do not represent
these companies.

# References
* [Data engineering pipeline deployment on the Modern Data Stack (MDS)](https://github.com/data-engineering-helpers/data-pipeline-deployment)
* [Architecture principles for data engineering pipelines on the Modern Data Stack (MDS)](https://github.com/data-engineering-helpers/architecture-principles)
* [Data Processing Pipeline (DPP) utility in Go (this repository)](https://github.com/data-engineering-helpers/dppctl)

## AWS SDK for Go
* [AWS - AWS SDK for Go v2](https://aws.github.io/aws-sdk-go-v2/docs/)
* [GO - AWS SDK for Go v2](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
* [AWS on GitHub - AWS SDK code examples in Go v2](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2)

# Getting started
* Check the latest versions/tags:
  https://github.com/data-engineering-helpers/dppctl/tags

* Import/download the module:
```bash
$ go get github.com/data-engineering-helpers/dppctl@vx.y.z
```

* Clone and edit the YAML deployment specification. For instance,
  for a deployment on AWS cloud:
```bash
$ cp depl/aws-dev-sample.yaml depl/aws-dev.yaml
$ vi depl/aws-dev.yaml
```

* Check the version of the `dppctl` utility:
```bash
$ dppctl -v
[dppctl] 0.0.x-alpha.x
```

* Launch the `dppctl` utility in checking mode (which is the default one):
```bash
$ dppctl -f depl/aws-dev.yaml
```

* Launch the `dppctl` utility in deployment mode:
```bash
$ dppctl -f depl/aws-dev.yaml -c deploy
```

# Publish the module
* Recompute the dependencies:
```bash
$ go mod tidy
```

* Check that the tests pass:
```bash
$ go test
```

* Tag the Git repository:
```bash
$ git commit -m "[Release][Go] v0.0.x-alpha.x"
$ git push
$ git tag -a v0.0.x-alpha.x -m "[Release][Go] v0.0.x-alpha.x"
$ git push --tags
```

* Publish the module:
```bash
$ GOPROXY=proxy.golang.org go list -m github.com/data-engineering-helpers/dppctl@v0.0.x-alpha.x
github.com/data-engineering-helpers/data-pipeline-deployment v0.0.x-alpha.x
```

# Troubleshooting

## AWS Airflow (MWAA)
As of beginning of 2023, apparently for security reasons, it does not seem
possible to target/use the Airflow API directly on
the AWS managed service (MWAA). One has to use instead the API backend
of the MWAA CLI. That is why the Go code of
[the corresponding `AWSAirflowCLI()` function](https://github.com/data-engineering-helpers/dppctl/blob/main/service/aws.go#AWSAirflowCLI)
is not straightforward.
Note that the use of the MWAA CLI API (through `curl`) is itself
convoluted, as detailed below.
	 
### References
* [Stack Overflow - Is it possible to access the Airflow API in AWS MWAA?](https://stackoverflow.com/questions/67884770/is-it-possible-to-access-the-airflow-api-in-aws-mwaa)
* [Apache Airflow - Airflow API reference guide](https://airflow.apache.org/docs/apache-airflow/stable/stable-rest-api-ref.html)
* [AWS - Amazon Managed Workflows for Apache Airflow (MWAA) User Guide](https://docs.aws.amazon.com/mwaa/index.html)
   + [AWS - Accessing the Apache Airflow UI](https://docs.aws.amazon.com/mwaa/latest/userguide/access-airflow-ui.html)
     - [AWS - Apache Airflow CLI command reference](https://docs.aws.amazon.com/mwaa/latest/userguide/airflow-cli-command-reference.html) 
* [GitHub - AWS - Sample code for MWAA](https://github.com/aws-samples/amazon-mwaa-examples)
  + [GitHub - AWS - Sample code for MWAA - Bash operator script](https://github.com/aws-samples/amazon-mwaa-examples/tree/main/dags/bash_operator_script)

### Listing the DAGs
* Configuration:
```bash
$ export MWAA_ENV="<the-MWAA-environment-name"
  export AWS_REGION="eu-west-1"
  export CLI_TOKEN
  export WEB_SERVER_HOSTNAME
```

* Create a CLI (command-line) token:
```bash
$ aws mwaa --region $AWS_REGION create-cli-token --name $MWAA_ENV
```
```javascript
{
    "CliToken": "someToken",
    "WebServerHostname": "<airflow-id>.$AWS_REGION.airflow.amazonaws.com"
}
```

* Copy/paste the web server hostname and the CLI token and save them
  as environment variables:
```bash
$ CLI_TOKEN="someToken"
  WEB_SERVER_HOSTNAME="<airflow-id>.$AWS_REGION.airflow.amazonaws.com"
```

* Note that the CLI token is very short-lived (valid for only one or two times)
  and the two operations (`aws mwaa create-cli-token` and
  `CLI_TOKEN="some-token"`) must be repeated every time before
  the following commands are perfomed

* Invoke an Airflow command through the API wrapping the MWAA CLI
  + Raw (not formatted) outpout:
```bash
$ curl -s --request POST "https://$WEB_SERVER_HOSTNAME/aws_mwaa/cli" --header "Authorization: Bearer $CLI_TOKEN" --header "Content-Type: text/plain" --data-raw "dags list -o json"|jq -r ".stdout" | base64 -d
```
```javascript
...
[{"dag_id": "dag_name", "filepath": "prefix/script.py", "owner": "airflow", "paused": "True"}, {"dag_id": ...}, ...]
```
  + CSV-formatted outpout (list of DAGs):
```bash
$ curl -s --request POST "https://$WEB_SERVER_HOSTNAME/aws_mwaa/cli" --header "Authorization: Bearer $CLI_TOKEN" --header "Content-Type: text/plain" --data-raw "dags list -o json"|jq -r ".stdout" | base64 -d | grep "^\[{\"dag_id\"" | jq -r ".[]|[.dag_id,.filepath,.owner,.paused]|@csv" | sed -e s/\"//g
```
```javascript
...
...
dag_name,prefix/script.py,airflow,True
...
```

