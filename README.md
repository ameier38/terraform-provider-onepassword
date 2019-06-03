[![Codefresh build status]( https://g.codefresh.io/api/badges/pipeline/ameier38/ameier38%2Fterraform-provider-onepassword%2Fterraform-provider-onepassword?branch=master&key=eyJhbGciOiJIUzI1NiJ9.NWMzMjE0ODA3YTJkOGI3ZjkxMzVhZjlm.WFn4I6XuUDBfWsKEp6LIuG-IlDsT4JCDTjMzeH7kGu8&type=cf-1)]( https://g.codefresh.io/pipelines/terraform-provider-onepassword/builds?repoOwner=ameier38&repoName=terraform-provider-onepassword&serviceName=ameier38%2Fterraform-provider-onepassword&filter=trigger:build~Build;branch:master;pipeline:5cf5099a3d8de566d41eda11~terraform-provider-onepassword)
[![Go Report Card](https://goreportcard.com/badge/github.com/ameier38/terraform-provider-onepassword)](https://goreportcard.com/report/github.com/ameier38/terraform-provider-onepassword)

# :lock: 1Password Terraform Provider :unlock:
Terraform data source (read: READ ONLY) provider for 1Password.

> This provider __does not__ create resources in 1Password. It requires
the user to have a 1Password account and to have created items in 1Password
prior to using. If you are interested in managing 1Password items through
terraform, please check out https://github.com/anasinnyk/terraform-provider-1password.

## Usage
First install terraform. On Window's you can use `scoop`.
```
scoop install terraform
```

Create a variables file `terraform.tfvars` and make sure it is added to 
your `.gitignore` so you do not expose your 1Password credentials.
> Learn more about terraform variables [here](https://www.terraform.io/docs/configuration/variables.html).

Add the following variables to the `terraform.tfvars` files.
```tf
# terraform.tfvars

password = "<your 1Password password>"
secret_key = "<your 1Password secret key>"
```

Next, create a terraform file `providers.tf` and add the provider.
```tf
# providers.tf

variable "onepassword_password" {
  description = "Login password for 1Password"
}

variable "onepassword_secret_key" {
  description = "Login secret key for 1Password"
}

provider "onepassword" {
	email = "test@testing.com"
	password = "${var.onepassword_password}"
	secret_key = "${var.onepassword_secret_key}"
	subdomain = "test" # Optional. Comes from <subdomain>.1password.com.
}
```

Then start using the provider to pull secrets. Let's say we
wanted to create a Kubernetes secret for our Redshift cluster.
We could create an item in 1Password called 'Redshift' in a
vault called 'Development'.

![redshift-item](./images/redshift-item.png)

We could then use the item in terraform to create a Kubernetes secret.
```tf
data "onepassword_item" "dev_redshift" {
  vault = "Development"
  item = "Redshift"
}

resource "kubernetes_secret" "redshift" {
  metadata {
    name      = "redshift"
    namespace = "default"
  }

  data {
    "redshift-user" = "${data.onepassword_item.dev_redshift.result["username"]}"
    "redshift-password" = "${data.onepassword_item.dev_redshift.result["password"]}"
    "redshift-host" = "${data.onepassword_item.dev_redshift.result["server"]}"
    "redshift-database" = "${data.onepassword_item.dev_redshift.result["database"]}"
  }
}
```
> Read more about creating Kubernetes secrets in terraform
[here](https://www.terraform.io/docs/providers/kubernetes/r/secret.html).

### Caveats
As of right now, you cannot create a schema with a nested map, which makes it difficult
to use a 1Password item with multiple sections. The workaround right now is to __only use
the default section with a blank section name__. 

![default-section](./images/default-section.png)

If you define other sections or add a section name to the default section, __you will
not be able to access the fields__.

`Elem` documentation ([link](https://godoc.org/github.com/hashicorp/terraform/helper/schema#Schema)):
```
// Elem represents the element type. For a TypeMap, it must be a *Schema
// with a Type that is one of the primitives: TypeString, TypeBool,
// TypeInt, or TypeFloat.
```

## Developement

### Setup
Install go. On Window's you can use `scoop`.
```
scoop install go
```

Clone the repository.
```
git clone https://github.com/ameier38/terraform-provider-onepassword.git
```

### Testing
Run acceptance tests.
```
$env:TF_ACC="true"
go test ./onepassword
```
> On macOS or Linux, set the environment variable with `export TF_ACC=true`.

## Contributing
If you find a :bug: please [create an issue](https://github.com/ameier38/terraform-provider-onepassword/issues)
and I will try to help resolve. If you would like to improve the library, feel free to
[open a pull request](https://github.com/ameier38/terraform-provider-onepassword/pulls).

## Resources
- [1Password CLI](https://support.1password.com/command-line-getting-started/)
- [Using 1Password in Terraform](https://medium.com/@JesseDearing/using-1password-values-in-terraform-71d2e3077380)
- [golang variable naming](https://talks.golang.org/2014/names.slide#1)
- [golang testing](https://golang.org/pkg/testing/)
- [golang writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)
- [interfaces in testing](https://nathanleclaire.com/blog/2015/10/10/interfaces-and-composition-for-effective-unit-testing-in-golang/)
- [terraform-provider-external](https://github.com/terraform-providers/terraform-provider-external)
- [Creating a Terraform Provider Part 1](https://medium.com/spaceapetech/creating-a-terraform-provider-part-1-ed12884e06d7)
- [Creating a Terraform Provider Part 2](https://medium.com/spaceapetech/creating-a-terraform-provider-part-2-1346f89f082c)
- [Terraform Schemas](https://www.terraform.io/docs/extend/schemas/index.html)
- [GitHub issue templates](https://github.com/stevemao/github-issue-templates)
