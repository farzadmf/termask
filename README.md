> ***VERY IMPORTANT NOTE***:
> Something that may be concerning is that this tool would have access to your secret information and somehow steal them (store them somewhere, etc.)
> I GUARANTEE that's NEVER going to happen, and this tool will only mask the secrets and WILL NOT use them in any way possible without your knowledge

# termask

A utility to mask property values in the terminal

It supports different inputs:

- Terraform (v0.12)
- JSON

## Introduction

### Terraform

> DISCLAIMER: I created this tool while only using Azure provider, so it is, in theory, possible that there would be bugs
> when using other providers; I'll do my best to solve any issues opened for other scenarios

Although Terraform allows marking `output` variables as `sensitive`, at the time of this writing, it doesn't provide a way to mark arbitrary values as "secret"

Inspired by [tfmask](https://github.com/cloudposse/tfmask), this program allows masking property values (the ones in the form of
`"property" = "value"`) in the output of `terraform plan` and `terraform apply`

***NOTE***: it's worth noting that, for the moment, it only supports the `-no-color` option of Terraform

## Installation

You can use `go get` to download the tool (a proper executable will be available soon)

```bash
go get github.com/farzadmf/termask
```

## Usage

### Terraform

You can get help by running `termask --help`:

```text
NAME:
   termask - Mask values in the terminal

USAGE:
   termask [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --mode value, -m value      (tf|json) mode determines the type of the input
   --property value, -p value  property whose value we want to mask (can be specified multiple times)
   --ignore-case, -i           case insensitive match (default: false)
   --partial-match, -l         match if property partially contains the specified string (default: false)
   --help, -h                  show help (default: false)
```

**NOTE**: by default, any property that contains the word `password` will be masked (case insensitive), and the options
below don't change this; they just add to this.

As mentioned in the help, you can use `--property` (or `-p`) to specify properties whose values you want to mask
(this flag can be specified multiple times).

By default, matching is done case sensitive, you can disable that by specifying the `--ignore-case` (or `-i`) flag.

Also, by default, a property should match the specified string as a whole word, `--partial-match` (or `-l`) allows overriding that.

### Examples

Let's say you have the following line in your `terraform plan` output:

```text
+ resource "azurerm_resource_group" "rg" {
    + "name" = "my-secret-resource-group"
    ...
}
```

If you want to mask the `name` property, you can do this:

```bash
# Don't forget the '-no-color' switch
terraform plan -no-color | termask -m tf -p name
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
terraform plan -no-color | termask -m tf -p name -p location
```

Which will result in the following output:

```text
+ resource "azurerm_resource_group" "rg" {
    + location = "***"
    + name     = "***"
    ...
}
```

An example combining multiple options, given the following Terraform output:

```text
+ resource "azurerm_app_service" "main" {
    ...
    app_settings = {
      ~ "MyConnectionString" = "secret-connection-string" -> "new-secret"
    }
    ...
}
```

We use the following command:

```bash
terraform plan -no-color | termask -m tf -p connectionstring -i -l
```

And we get:

```text
+ resource "azurerm_app_service" "main" {
    ...
    app_settings = {
      ~ "MyConnectionString" = "***" -> "***"
    }
    ...
}
```

### JSON

We have a simialr concept for JSON input. Let's say you a file named `my.json` with the following content:

```json
{
  "password": "secret",
  "property": "value",
  "name": "John"
}
```

Since, by default, `password` is masked, if you run this:

```bash
cat my.json | termask -m json
```

you would see the following output:

```json
{
  "password": "***",
  "property": "value",
  "name": "John"
}
```

And, you can choose to mask other properties:

```bash
cat my.json | termask -m json -p name
```

Gives you:

```json
{
  "password": "***",
  "property": "value",
  "name": "***"
}
```
