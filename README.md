# op-terraform
[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/ameier38/ameier38%2Fop-terraform%2Fop-terraform?key=eyJhbGciOiJIUzI1NiJ9.NWMzMjE0ODA3YTJkOGI3ZjkxMzVhZjlm.WFn4I6XuUDBfWsKEp6LIuG-IlDsT4JCDTjMzeH7kGu8&type=cf-1)]( https://g.codefresh.io/pipelines/op-terraform/builds?filter=trigger:build~Build;pipeline:5ce2933ab66ecb8654fe386b~op-terraform)
[![Go Report Card](https://goreportcard.com/badge/github.com/ameier38/op-terraform)](https://goreportcard.com/report/github.com/ameier38/op-terraform)

___
Thin wrapper around the 1Password CLI for use in Terraform.

Based on [this blog post](https://medium.com/@JesseDearing/using-1password-values-in-terraform-71d2e3077380)
but uses go instead of bash so that it will work on Windows.

## Setup

Install go.
```
scoop install go
```

Clone the repository.

Build the binary.
```
go build
```

Add the binary 

## Resources
- [1Password CLI](https://support.1password.com/command-line-getting-started/)
- [Using 1Password in Terraform](https://medium.com/@JesseDearing/using-1password-values-in-terraform-71d2e3077380)
- [golang variable naming](https://talks.golang.org/2014/names.slide#1)
- [golang testing](https://golang.org/pkg/testing/)
- [golang writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)
- [interfaces in testing](https://nathanleclaire.com/blog/2015/10/10/interfaces-and-composition-for-effective-unit-testing-in-golang/)
