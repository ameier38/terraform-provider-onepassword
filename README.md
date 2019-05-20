[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/ameier38/ameier38%2Fop-terraform%2Fop-terraform?key=eyJhbGciOiJIUzI1NiJ9.NWMzMjE0ODA3YTJkOGI3ZjkxMzVhZjlm.WFn4I6XuUDBfWsKEp6LIuG-IlDsT4JCDTjMzeH7kGu8&type=cf-1)]( https://g.codefresh.io/pipelines/op-terraform/builds?filter=trigger:build~Build;pipeline:5ce2933ab66ecb8654fe386b~op-terraform)
[![Go Report Card](https://goreportcard.com/badge/github.com/ameier38/op-terraform)](https://goreportcard.com/report/github.com/ameier38/op-terraform)

# op-terraform
Thin wrapper around the 1Password CLI for use in Terraform.

Based on [this blog post](https://medium.com/@JesseDearing/using-1password-values-in-terraform-71d2e3077380)
but uses go instead of bash so that it will work on Windows.

## Setup

Install go.
```
scoop install go
```

Clone the repository.
```
git clone https://github.com/ameier38/op-terraform.git
```

Build the binary.
```
cd op-terraform
go build
```

Add the binary to your path.
```
vim $PROFILE
```
```
$env:Path += ";C:\path\to\op-terraform"
```

Restart your shell and check the installation.
```
op-terraform --help
```

## Usage in shell
First sign in to 1Password.
```
iex $(op signin)
```
> See [1Password CLI docs](https://support.1password.com/command-line-getting-started/) for
setting up the 1Password CLI.

Get an item from 1Password.
```
echo '{"vaultName": "test-vault", "itemName": "test-item"}' | op-terraform
```

## Usage in Terraform
Add external data source.
```t
data "external" "op_test_item" {
  program = ["op-terraform"]
  query = {
    vaultName = "test-vault"
    itemName = "test-item"
  }
}
```

And then in a resource
```t
resource "kubernetes_secret" "test_secret" {
  metadata {
    name      = "test-secret"
    namespace = "default"
  }

  data {
    user = "${data.external.development_redshift.result.username}"
    password = "${data.external.development_redshift.result.password}"
  }
}
```

## Resources
- [1Password CLI](https://support.1password.com/command-line-getting-started/)
- [Using 1Password in Terraform](https://medium.com/@JesseDearing/using-1password-values-in-terraform-71d2e3077380)
- [golang variable naming](https://talks.golang.org/2014/names.slide#1)
- [golang testing](https://golang.org/pkg/testing/)
- [golang writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)
- [interfaces in testing](https://nathanleclaire.com/blog/2015/10/10/interfaces-and-composition-for-effective-unit-testing-in-golang/)
