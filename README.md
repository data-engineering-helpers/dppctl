Getting started with Go for Data Processing Pipeline (DPP) CLI tool
===================================================================

# References

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


