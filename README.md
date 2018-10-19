# projectmpl
Project boilerplate engine

[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/built-with-love.svg)](https://forthebadge.com)[![forthebadge](https://forthebadge.com/images/badges/uses-badges.svg)](https://forthebadge.com)

![Go Version](https://img.shields.io/badge/Go%20Version-latest-brightgreen.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/Depado/projectmpl)](https://goreportcard.com/report/github.com/Depado/projectmpl)
[![Build Status](https://drone.depado.eu/api/badges/Depado/projectmpl/status.svg)](https://drone.depado.eu/Depado/projectmpl)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/Depado/projectmpl/blob/master/LICENSE)
[![Say Thanks!](https://img.shields.io/badge/Say%20Thanks-!-1EAEDB.svg)](https://saythanks.io/to/Depado)

<!-- TOC -->

- [projectmpl](#projectmpl)
- [Usage](#usage)
    - [Keeping the template](#keeping-the-template)
- [Template Creation](#template-creation)
    - [The root `.projectmpl.yml` file](#the-root-projectmplyml-file)
    - [Variable declaration](#variable-declaration)
        - [Simple Input](#simple-input)
        - [Selection](#selection)
        - [Boolean/Confirmation](#booleanconfirmation)
        - [Other options and help](#other-options-and-help)
        - [Validation](#validation)
    - [Standard `.projectmpl.yml` files](#standard-projectmplyml-files)
    - [Per-file configuration](#per-file-configuration)

<!-- /TOC -->

# Usage

Projectmpl supports various provider to download the templates. It supports 
`git`, downloading an archive (`.zip`/`.tar.gz`/`.tar.xz`/...) from internet,
or using a local directory. 

## Keeping the template

When donwloading or cloning a template, `projectmpl` will create a temporary
directory and delete it once the operation completes. If you want to keep
the template (to play with it, or simply to keep a copy), make sure you pass
the `--template.keep` option. This option pairs well with `--template.output`
which defines where the template should be downloaded/cloned.

# Template Creation

## The root `.projectmpl.yml` file

To configure your template, place a `.projectmpl.yml` at the root of your template.
This is called the root configuration, and should contain some information about
your template such as its name, its version and a description.

It can also contain overrides for delimiters in the templates (defaults being
the go-style `{{ .var }}`) and variables.

```yaml
name: "Example Projectmpl"
version: "0.1.0"
description: "An example template to show how projectmpl works"
```

## Variable declaration

You can add a `variables` section to your root configuration (or to any 
`.projectmpl.yml` file, or directly inline in your template files, see below) to
define the variables you want your user to define. There are three types of
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

## Standard `.projectmpl.yml` files

If you place a `.projectmpl.yml` file in a sub-directory of your template, this
file will apply recursively to all the elements inside that directory and its 
own sub-directories, meaning that you can override some variables, add new ones, 
modify the delimiters, or completely ignore an entire directory.

For example you can completely ignore a director:

```
└──── change
    ├── override.go
    └── .projectmpl.yml
```

```yaml
ignore: true
```

In this case, the file `override.go` won't be rendered (but will simply be 
copied to the output directory). This would apply for every sub-directory, 
except if a directory contains a `.projectmpl.yml` telling otherwise, or a
file with an inline configuration. 

## Per-file configuration

You can also configure individual files by adding a front matter at the top
of the file (that will obviously be removed when rendered). 

Let's say I have a file that I don't want to render:

```
---
ignore: true
---
# This shouldn't be rendered at all !
```

You can even add per-file variables, or modify the delimiters. In fact, it's
like an inline `.projectmpl.yml` that applies to a single file.
