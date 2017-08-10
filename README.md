## CLIg 

clig is a (do the needful) cli generator which generates a skeleton for your "go" based cli application. This tool creates a go seed cli project right out of the box. The current version only supports urfave/cli framework


### Overview

The project generates an cli skeletorn for a typical go cli app. You can quickly generate the boilerplate code for you go cli projects. CLIg helps you to generate associated subcommands and flags for the given application. Using CLIg is fairly simple (no fancy things), just give an yml file as an argument and you should be good to go.


### Getting Started

To get you started you can simply install the go application and input and yml file to generate the boilerplate code


#### Prerequistes

Firstly you need to have go installed and GOPATH set, it is necessary that your GOPATH is set or the boiler plate code is created in home directory. You can install go and set GOPATH [here](https://golang.org/doc/install#install)

To install CLIg, simply run

    $ go get github.com/sridhav/clig

    $ go install github.com/sridhav/clig

Make sure your PATH includes $GOPATH/bin directory so your commands can be easily used

    export PATH=$PATH:$GOPATH/bin


#### Generating the boilerplate code

Once all the pre-requistes are satisfied, use your favorite editor to create a yml file and generate your command structure. A sample is provided [here](https://github.com/sridhav/clig/blob/master/clig.yml.dist), (will add more samples). Once your yaml file is set use the following command to generate your boilerplate code

    clig <yml file location>

This generates your boiler plate code. The generated boiler code is stored in your $GOPATH/src/\<vcshost\>/\<author\>/\<appname\> folder.

The important configuration parameters. As these are required parameters in your yml, there are some defaults attached to it

`vcshost`   defaults to `github.com`

`author`    defaults to `local username` 

`appname`   defaults to `myapp`


#### Sample YML file

```
# app name required param
name: "cligs"

# author required
author: "sridhav"

# vcs host required
vcshost: "github.com"

# global flags
flags:
      - name: "hadoop"
        type: "Bool"

# commands
commands:
  - name: "doo"
    usage: "do the doo"
    description: "no really"
    flags:
      - name: "flag"
        type: "Bool"
    # subcommands for command doo (you can nest commands inside commands)
    commands:
      - name: "doodoo"
        usage: "do tht doodoo"
        description: "no no really"
        flags:
            - name: "flag"
              type: "Bool"
```

### Supported Platforms

CLIg is tested agains multiple versions of Go on Linux and agains the latest released version of Go on OS X. Windows should also have no problem but haven't tested it.


### Versus Other go cli generators

A big thanks to this other option for existing. This has paved me the way to make CLIg approachable.

[gcli](https://github.com/tcnksm/gcli) - has support for multiple framework and can generate go cli code with different boilerplates. gcli lacks the support of generating sub commands. CLIg just does the needful boilerplate including subcommands.
