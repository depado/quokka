<h1 align="center">Quokka</h1>
<h2 align="center">
  <img src="/assets/mascot.png" alt="mascot" height="200px">

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/github.com/Depado/quokka)](https://goreportcard.com/report/github.com/Depado/quokka)
  [![Build Status](https://drone.depa.do/api/badges/Depado/quokka/status.svg)](https://drone.depa.do/Depado/quokka)
  [![codecov](https://codecov.io/gh/Depado/quokka/branch/master/graph/badge.svg)](https://codecov.io/gh/Depado/quokka)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/quokka/blob/master/LICENSE)
  [![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/Depado)
</h2>

<h2 align="center">Friendly Boilerplate Engine</h2>

<img align="center" src="/assets/quokka.gif">

- [Introduction](#introduction)
    - [Features](#features)
- [Installation](#installation)
    - [Download](#download)
    - [Build from source](#build-from-source)
- [Usage](#usage)
    - [Keeping the template](#keeping-the-template)
    - [Input File](#input-file)
    - [Set](#set)
    - [Examples](#examples)
- [Template Creation](#template-creation)
    - [The root `.quokka.yml` file](#the-root-quokkayml-file)
    - [Variable declaration](#variable-declaration)
        - [Simple Input](#simple-input)
        - [Selection](#selection)
        - [Boolean/Confirmation](#booleanconfirmation)
        - [Other options and help](#other-options-and-help)
        - [Validation](#validation)
        - [Sub Variables](#sub-variables)
    - [Standard `.quokka.yml` files](#standard-quokkayml-files)
    - [Per-file configuration](#per-file-configuration)
    - [Conditional Rendering/Copy](#conditional-renderingcopy)

# Introduction

Quokka is a boilerplate engine. It allows you to quickly use boilerplate
templates and avoid copy-pasting chunks of code and snippets when you start a
new project. You can create templates for literally anything you want!

## Example Usages

- Generating your CI/CD configuration file
- Generating skeleton applications with your own best practices
- More

## Features

- **No external dependencies**
  Quokka is written in Go and thus is compiled to a static binary. Download
  or build it and you're good to go.
- **Local or distant templates**
  Quokka supports both git repositories and local files.
- **Sweet output and prompts**
  Thanks to the wonderful [survey](https://github.com/AlecAivazis/survey)
  library, the prompts are unified, can display an help text and support
  validation.
- **Clean configuration files**
  Quokka uses YAML for its configuration file formats, making them clean
  and easy to read.
- **Powerful templating system**
  Quokka uses [Go's template system](https://golang.org/pkg/text/template/)
  to render the boilerplate.
- **Configuration override**
  Need a different behavior or additional variables in a specific directory?
  Just add another `.quokka.yml` file in there. You can even overwrite
  variables.
- **Conditional prompts (sub-variables)**
  Each variable can have its own subset of variables which will only be
  prompted to the user if the parent variable is filled or set to true.
- **Customizable templates**
  Quokka enables fine-grained control over what needs to be done when
  rendering the template. Just copy the file, ignore it, add conditionals based
  on what the user answered, change the template delimiters…

# Installation

## Download

You can grab the latest release from [the release page](https://github.com/Depado/quokka/releases).

## Build from source

```
$ go get -u github.com/Depado/quokka
$ cd $GOPATH/src/github.com/Depado/quokka
$ make
```

Or directly install:

```
$ go get -u github.com/Depado/quokka
$ cd $GOPATH/src/github.com/Depado/quokka
$ make install
```

# Usage

Quokka has two ways of retrieving the templates. It supports
`git` or using a local directory.

```
Quokka (qk) is a template engine that enables to render local or distant
templates/boilerplates in a user friendly way. When given a URL/Git repository
or a path to a local Quokka template, quokka will ask for the required values
in an interactive way except if an inpute file is given to the CLI.

Usage:
  qk [template] [output] <options> [flags]
  qk [command]

Available Commands:
  help        Help about any command
  new         Create a new quokka template
  version     Show build and version

Flags:
      --debug           Enable or disable debug mode
      --git.depth int   depth of git clone in case of git provider (default 1)
  -h, --help            help for qk
  -i, --input string    specify an input values file to automate template rendering
  -k, --keep            do not delete the template when operation is complete
  -o, --output string   specify the directory where the template should be downloaded or cloned
  -p, --path string     specify if the template is actually stored in a sub-directory of the downloaded file
  -e, --set strings     specify values on the command line
  -y, --yes             Automatically accept

Use "qk [command] --help" for more information about a command.
```

## Keeping the template

When cloning a template, Quokka will create a temporary
directory and delete it once the operation completes. If you want to keep
the template (to play with it, or simply to keep a copy), make sure you pass
the `-k/--keep` option. This option pairs well with the `-o/--output` option
which defines where the template should be downloaded/cloned.

## Input file

The rendering of a Quokka template can be automated if the template was designed
with this in mind and if an input file is provided on the command line.

Since there is no clear way for specifying overriding values (for example a
variable that applies to a single file and overrides an already existing
variable in the root config), the input values will also fill the overriding
variables.

The format of the input file is also yaml. The following example demonstrates
how an input file could be used:

`.quokka.yml`
```yaml
name: "Quokka Template"
description: "New Quokka Template"
version: "0.1.0"
variables:
  slack:
    confirm: true
    prompt: "Add Slack integration?"
    variables:
      channel:
        required: true
      webhook:
        required: true
```

`input.yml`
```yaml
slack: true
slack_channel: "#mychan"
slack_webhook: "complexurl
```

If this input file is given to Quokka, it won't prompt for these three
variables, thus requiring no input from the user to render the template.

## Set

Additionally, you can provide Quokka with the `-e/--set` flag (multiple time if
you wish). This works the same way as the input file but has a higher priority,
meaning that if you pass both an input file and a `-e` flag that defines a
variable, the one passed on the command line will have a higher priority.

The `--set` flags work by providing it with a `key=value` style kind of string.
If we take the example above using the input file, we can effectively replace
the `slack_channel` variable by doing so:

```sh
$ qk template/ output -i input.yml --set "slack_channel=#anotherchan"
$ # Or
$ qk template/ output -i input.yml -e "slack_channel=#anotherchan"
```

## Examples

```sh
$ # Clone the repository and execute the template that is located in _example/license
$ qk git@github.com:Depado/quokka.git output --path _example/license
$ # Clone the template in a specific directory, render it in a specific directory and keep the template
$ qk git@github.com:Depado/quokka.git myamazingproject --path _example/cleanarch --keep --output "template"
$ # Reuse the downloaded template
$ qk template/ myotherproject
$ # Pass an input file to Quokka
$ qk template/ output -i in.yml
```

# Template Creation

## New command

If `quokka` is installed, simply run `quokka new <path>`. This will ask for
basic information such as the template name, description and version with some
sane defaults (version number for example is set to `0.1.0` by default).
You can also pass these values as flags on the command line.

This command will check if the output directory and a `.quokka.yml` file already
exist. This command is in charge of creating a new directory and creating the
initial `.quokka.yml` file with those basic information, helping you getting
started with Quokka template development.

<details><summary>Command Line Help</summary>

```
$ qk new --help
Create a new quokka template

Usage:
  qk new [output] <options> [flags]

Flags:
  -d, --description string   description of the new template
  -h, --help                 help for new
  -n, --name string          name of the new template
  -v, --version string       version of the new template

Global Flags:
      --debug   Enable or disable debug mode
  -y, --yes     Automatically accept
```
</details>

## The root `.quokka.yml` file

To configure your template, place a `.quokka.yml` at the root of your template.
This is called the root configuration, and should contain some information about
your template such as its name, its version and a description.

It can also contain overrides for delimiters in the templates (defaults being
the go-style `{{ .var }}`) and variables.

```yaml
name: "Example Quokka Template"
version: "0.1.0"
description: "An example template to show how quokka works"
```

## Variable declaration

You can add a `variables` section to your root configuration (or to any
`.quokka.yml` file, or directly inline in your template files, see below) to
define the variables you want your user to fill in. There are three types of
input you can use:

### Simple Input

If you just specify the name of your variable, it will result in a simple input.

```yaml
variables:
  name:
```

### Selection

```yaml
variables:
  license:
    values: ["MIT", "Apache License 2.0", "BSD 3", "FreeBSD", "GPL", "LGPL", "WTFPL", "None"]
```
This will result in a selection input where the user can choose one of the
provided choices.

### Boolean/Confirmation

```yaml
variables:
  test:
    confirm: true
```
If you're using the `confirm` keyword, it will generate a simple yes/no input.
The value you give that `confirm` key becomes the default value.

### Other options and help

You can also help your users by changing the prompt, adding a help text or
providing a default value:

```yaml
variables:
  license:
    values: ["MIT", "Apache License 2.0", "BSD 3", "FreeBSD", "GPL", "LGPL", "WTFPL", "None"]
    prompt: Which license do you want for your project?"
    help: "License file that will be added to your project"
    default: "MIT"
  name:
    default: amazingproject
    prompt: "What's the name of your project?"
    help: "Used to render the README file and various configuration files"
```

### Validation

You can mark any variable as required using the `required` keyword:

```yaml
variables:
  name:
    default: amazingproject
    prompt: "What's the name of your project?"
    help: "Used to render the README file and various configuration files"
    required: true
```

This will prevent the user from rendering your template with missing variables.
Note that if you specified a default value for an input, it becomes impossible
to not fill in that value. So the validator becomes obsolete.

### Sub Variables

It's not uncommon to ask for additional information when the user answered yes
or filled in a variable. Thus, each variable can have its own variables:

```yaml
variables:
  slack:
    confirm: true
    prompt: "Add Slack integration?"
    variables:
      channel:
        required: true
        prompt: "In which Slack channel should the result be posted?"
      webhook:
        required: true
        help: "See https://api.slack.com/incoming-webhooks for more information"
        prompt: "Provide the Slack webhook URL:"
```

In the example above we ask the user if he wants a Slack integration. If he
answers yes to that, then we'll ask him about the Slack channel and the webhook
URL. Otherwise we won't bother him with these details since they won't be used
in our template rendering.

The sub variables can be accessed in your templates with the form `.parent_sub`.
In this case, `.slack_channel` and `.slack_webhook`.

## Standard `.quokka.yml` files

If you place a `.quokka.yml` file in a sub-directory of your template, this
file will apply recursively to all the elements inside that directory and its
own sub-directories, meaning that you can override some variables, add new ones,
modify the delimiters, or completely ignore an entire directory.

For example you can completely ignore a directory:

```
└──── change
    ├── override.go
    └── .quokka.yml
```

```yaml
copy: true
```

In this case, the file `override.go` won't be rendered (but will simply be
copied to the output directory). This would apply for every sub-directory,
except if a directory contains a `.quokka.yml` telling otherwise, or a
file with an inline configuration. The `ignore` option can also be used to
completely ignore a file or a directory.

```yaml
ignore: true
```

## Per-file configuration

You can also configure individual files by adding a front matter at the top
of the file (that will obviously be removed when rendered).

Let's say I have a file that I don't want to render but simply copy to the
output directory:

```
---
copy: true
---
# This shouldn't be rendered at all !
```

You can even add per-file variables, or modify the delimiters. In fact, it's
like an inline `.quokka.yml` that applies to a single file.

Supported instructions are as follows:

- `if: condition`: Conditional rendering using an [expr](https://expr.medv.io/docs/Language-Definition)
  expression that must return a boolean value
- `copy: true`: Do not attempt to render the file and simply copy it to its
  destination
- `delimiters: ["[[", "]]"]`: Change the delimiters used for template rendering.
  This can be useful for files that already are templates or use extensively the
  `{}` chars
- `rename: newname`: Rename the file to this new name once rendered
- `ignore: true`: Completely ignore the file (no render, no copy)

## Conditional Rendering/Copy

You may want some files to not be copied or rendered according to what the user
answers to your prompt. You can use the `if` key (in a `.quokka.yml`
or inline in a file), with the name of one of your variables. For example if
you have a variable defined like this in your root config:

```yaml
variables:
  drone:
    prompt: "Do you want to add a Drone config file?"
    confirm: true
```

You can then add this at the top of the file:

```
---
if: drone
---
workspace:
  base: /go
...
```

This file will be rendered if, and only if, the user answered yes to that
question. Note that `if` and `copy` can work together if you just want
to copy the file and not render it.
