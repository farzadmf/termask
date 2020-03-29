# tfmask

A utility to mask property values in terraform output to the console

This project is heavily inspired by [tfmask from cloudposse](https://github.com/cloudposse/tfmask), with a few differences:

- It's tailored towards Terraform 0.12
- It doesn't use environment variables to determine what to mask; instead, it accepts CLI arguments
- It *may* support less scenarios, but it doesn't have to stay that way 🙂

## Introduction

Although Terraform allows marking `output` variables as `sensitive`, at the time of this writing, it doesn't provide a way to mark arbitrary values as "secret"

Inspired by [tfmask](https://github.com/cloudposse/tfmask), this is a go program to allow masking property values (the ones in the form of
`"property" = "value"`) in the output of `terraform plan` and `terraform apply`

***NOTE***: it's worth noting that, for the moment, it only supports the `-no-color` option of Terraform

## How to use

### Installation

You can use `go get` to download the tool (a proper executable will be available soon)

```bash
go get github.com/farzadmf/tfmask
```

### Usage

You can get help by running `tfmask --help`:

```bash
NAME:
   tfmask - Mask Terraform property values

USAGE:
   tfmask [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --property value, -p value  property to mask (can be specified multiple times)
   --ignore-case, -i           case insensitive match (default: false)
   --help, -h                  show help (default: false)
```

**NOTE**: by default, any property that contains the word `password` will be masked (case insensitive), and the options
below don't change this; they just add to this.

As mentioned in the help, you can use `--property` (or `-p`) to specify properties whose values you want to mask
(this flag can be specified multiple times).

By default, matching is done case sensitive, you can disable that by specifying the `--ignore-case` (or `-i`) flag.

### Examples

Let's say you have the following line in your `tf plan` output:

```text
+ resource "azurerm_resource_group" "rg" {
    + "name" = "my-secret-resource-group"
    ...
}
```

If you want to mask the `name` property, you can do this:

```bash
tf plan | tfmask -p name
```

And the output will be:

```text
+ resource "azurerm_resource_group" "rg" {
  + "name" = "***"
  ...
}
```

You can also mask multiple properties; let's say you have the following output:

```text
+ resource "azurerm_resource_group" "rg" {
    + location = "eastus"
    + name     = "mysecretrgname"
    ...
}
```

And you want to mask `name` and `location`:

```bash
tf plan | tfmask -p name -p location
```

Which will result in the following output:

```text
+ resource "azurerm_resource_group" "rg" {
    + location = "***"
    + name     = "***"
    ...
}
```
