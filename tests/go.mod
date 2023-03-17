module check-dppctl

go 1.20

replace github.com/data-engineering-helpers/dppctl => ../dppctl

require github.com/data-engineering-helpers/dppctl v0.0.1-alpha.1

require (
	golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c // indirect
	rsc.io/quote v1.5.2 // indirect
	rsc.io/sampler v1.3.0 // indirect
)
