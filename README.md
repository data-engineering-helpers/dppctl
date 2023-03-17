Getting started with Go for Data Processing Pipeline (DPP) CLI tool
===================================================================

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

# Getting started
* Check the latest versions/tags:
  https://github.com/data-engineering-helpers/data-pipeline-deployment/tags

* Import/download the module:
```bash
$ go get github.com/data-engineering-helpers/data-pipeline-deployment@vx.y.z
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
$ GOPROXY=proxy.golang.org go list -m github.com/data-engineering-helpers/data-pipeline-deployment@v0.0.x-alpha.x
github.com/data-engineering-helpers/data-pipeline-deployment v0.0.x-alpha.x
```

# First time
* Create the `dppctl` module:
```bash
$ mkdir -p dppctl
$ pushd dppctl
$ go mod init github.com/data-engineering-helpers/data-pipeline-deployment/go/dppctl
$ go mod tidy
$ popd
```

* Create a checker:
```bash
$ mkdir -p tests
$ pushd tests
$ got mod init check-dppctl
$ go mod edit -replace github.com/data-engineering-helpers/data-pipeline-deployment/go/dppctl=../dppctl
$ go mod tidy
$ go run check-dppctl.go
$ go build check-dppctl.go
$ ./check-dppctl
$ popd
```

* Install the checker:
```bash
$ go list -f '{{.Target}}'
~/go/bin/check-dppctl
$ go install
$ ls -laFh ~/go/bin/check-dppctl
-rwxr-xr-x  1 user staff 2.1M Mar 17 16:23 /Users/DENIS/go/bin/check-dppctl*
```

* Run the AWS checker:
```bash
$ go run check-aws.go
$ go build check-aws.go
```


