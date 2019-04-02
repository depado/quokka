<h1 align="center">Quokka</h1>
<div>
  <span align="left">
    <img src="/assets/mascot.png" alt="mascot">
  </span>
  <span align=right>

  [![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/uses-badges.svg)](https://forthebadge.com)

  ![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
  [![Go Report Card](https://goreportcard.com/badge/github.com/Depado/projectmpl)](https://goreportcard.com/report/github.com/Depado/projectmpl)
  [![Build Status](https://drone.depado.eu/api/badges/Depado/projectmpl/status.svg)](https://drone.depado.eu/Depado/projectmpl)
  [![codecov](https://codecov.io/gh/Depado/projectmpl/branch/master/graph/badge.svg)](https://codecov.io/gh/Depado/projectmpl)
  [![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/projectmpl/blob/master/LICENSE)
  [![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/Depado)
  
  </span>
</div>

<h2 align="center">Friendly Boilerplate Engine</h2>

<img align="center" src="/assets/projectmpl.gif">

- [Introduction](#introduction)
    - [Features](#features)
- [Installation](#installation)
    - [Download](#download)
    - [Build from source](#build-from-source)
- [Usage](#usage)
    - [Commands](#commands)
    - [Keeping the template](#keeping-the-template)
    - [Examples](#examples)
- [Template Creation](#template-creation)
    - [The root `.projectmpl.yml` file](#the-root-projectmplyml-file)
    - [Variable declaration](#variable-declaration)
        - [Simple Input](#simple-input)
        - [Selection](#selection)
        - [Boolean/Confirmation](#booleanconfirmation)
        - [Other options and help](#other-options-and-help)
        - [Validation](#validation)
        - [Sub Variables](#sub-variables)
    - [Standard `.projectmpl.yml` files](#standard-projectmplyml-files)
    - [Per-file configuration](#per-file-configuration)
    - [Conditional Rendering/Copy](#conditional-renderingcopy)
    - [After render commands](#after-render-commands)

# Introduction

Projectmpl is a boilerplate engine. It allows you to quickly use boilerplate 
templates and avoid copy-pasting chunks of code and snippets when you start a 
new project.

## Features

- **No external dependencies**  
  Projectmpl is written in Go and thus is compiled to a static binary. Download
  or build it and you're good to go.
- **Local or distant templates**  
  Either your template is a git repository, an archive stored on a distant
  server or a local directory, projectmpl knows how to handle these.
- **Sweet output and prompts**  
  Thanks to the wonderful [survey](https://github.com/AlecAivazis/survey) 
  library, the prompts are unified, can display an help text and support
  validation.
- **Clean configuration files**  
  Projectmpl uses YAML for its configuration file formats, making them clean
  and easy to read.
- **Powerful templating system**  
  Projectmpl uses [Go's template system](https://golang.org/pkg/text/template/)
  to render the boilerplate.
- **Configuration override**  
  Need a different behavior or additional variables in a specific directory? 
  Just add another `.projectmpl.yml` file in there. You can even overwrite
  variables. 
- **Conditional prompts (sub-variables)**  
  Each variable can have its own subset of variables which will only be
  prompted to the user if the parent variable is filled or set to true.
- **Customizable templates**  
  Projectmpl enables fine-grained control over what needs to be done when
  rendering the template. Just copy the file, ignore it, add conditionals based
  on what the user answered, change the template delimiters…
- **[On Hold] After render commands**  
  Projectmpl enables you to define commands to be run once the boilerplate has
  been rendered. _For security reasons, an explicit flag must be provided by the
  user for the commands to be executed_
  This feature is currently disabled for security reasons.

# Installation

## Download

You can grab the latest release from [the release page](https://github.com/Depado/projectmpl/releases).

## Build from source

```
$ go get -u github.com/Depado/projectmpl
$ cd $GOPATH/src/github.com/Depado/projectpml
$ dep ensure
$ make
```

Or directly install:

```
$ go get -u github.com/Depado/projectmpl
$ cd $GOPATH/src/github.com/Depado/projectpml
$ dep ensure
$ make install
```

# Usage

Projectmpl supports various provider to download the templates. It supports 
`git`, downloading an archive (`.zip`/`.tar.gz`/`.tar.xz`/...) from internet,
or using a local directory. 

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
  -d, --debug           debug mode
      --git.depth int   depth of git clone in case of git provider (default 1)
  -h, --help            help for qk
  -i, --input string    specify an input values file to automate template rendering
  -k, --keep            do not delete the template when operation is complete
  -o, --output string   specify the directory where the template should be downloaded or cloned
  -p, --path string     specify if the template is actually stored in a sub-directory of the downloaded file

Use "qk [command] --help" for more information about a command.
```

<!-- ## Commands

Some templates may define additional commands that will run once the template
has been rendered. If you wish to activate this behavior, you can pass the
`-c` or `--commands` flag. These commands can be anything, and may harm your
system so make sure you are ok with that.  -->

## Keeping the template

When downloading or cloning a template, `projectmpl` will create a temporary
directory and delete it once the operation completes. If you want to keep
the template (to play with it, or simply to keep a copy), make sure you pass
the `--keep` option. This option pairs well with the `--output` option which 
defines where the template should be downloaded/cloned.

## Examples

```sh
$ # Clone the repository and execute the template that is located in _example/license
$ qk git@github.com:Depado/projectmpl.git output --path _example/license
$ # Clone the template in a specific directory, render it in a specific directory and keep the template
$ qk git@github.com:Depado/projectmpl.git myamazingproject --path _example/cleanarch --keep --output "template"
$ # Reuse the downloaded template
$ projectmpl template/ myotherproject
```

# Template Creation

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
    └── .projectmpl.yml
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
like an inline `.projectmpl.yml` that applies to a single file.

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

<!-- ## After render commands

You can define some actions to be run once your template has been rendered.
You can only define those in the root configuration (not in sub-directory
configuration files). These actions can be configured to be run only when a
variable has been entered, just like the conditional rendering. Here is an
example of initializing the git repository when the template has been rendered:

```yml
after:
  - cmd: "git init"
    echo: "Intialized git repo"
    if: git
  - cmd: "git config core.hooksPath .githooks"
    echo: "Configured git hooks"
    if: git
variables:
  git:
    confirm: true
    prompt: "Initialize git repo and git hooks ?"
```

If the user answers yes to the question about git, then the repo will be 
initialized. You can also specify that you want the output of the command to be
displayed to the user using the `output: true`. `echo` is used to display a nice
message (instead of the command output).

**Note**: Due to the potential misbehavior of template creators, the user needs
to pass the `-c` or `--commands` to execute those commands. Otherwise the 
commands will be completely ignored. -->