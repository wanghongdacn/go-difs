# go-difs

![banner](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/docs/difs.png)

[![](https://img.shields.io/badge/made%20by-DACC-yellowgreen.svg)](http://dacc.co)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()
[![Golang](https://img.shields.io/badge/Golang-1.10%2B-blue.svg)](https://golang.org/)

> DIFS implementation in Go

## What is DIFS?

![banner](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/docs/difs-blueprint.png)

DIFS is DACC IAM File System based on [difs](https://github.com/difs/difs) with IAM(Identity and Access Management) system.it is the most fundamental infrastructure and innovation of DACC architecture. It includes file sharing system, permission maps, migration engine from centralized storage to decentralized storage. More importantly, the whole system is modular and independent from the rest of DACC architecture, so any other public chain that needs modern decentralized storage with IAM capabilities can easily implement DACC file system.

For more information see -> http://dacc.co/whitepaper/Dacc.pdf

Please put all issues regarding:
  - DIFS _design_ in the [DIFS protocol repo issues](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues).
  - Go DIFS _implementation_ in [this repo](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues).

## Table of Contents

- [Security Issues](#security-issues)
- [Install](#install)
  - [System Requirements](#system-requirements)
  - [Build from Source](#build-from-source)
    - [Install Go](#install-go)
    - [Download and Compile difs](#download-and-compile-difs)
    - [Troubleshooting](#troubleshooting)
  - [Development Dependencies](#development-dependencies)
- [Usage](#usage)
- [Getting Started](#getting-started)
  - Coming soon (WIP)
- [Contributing](#contributing)
  - [Want to read our code?](#want-to-read-our-code)
- [License](#license)

## Security Issues

The DIFS protocol and its implementations are still in preliminary stage. This means that there may be problems in our protocols, or there may be mistakes in our implementations. And DIFS is not production-ready yet. If you discover a security issue, please bring it to our attention! thx

If you find a vulnerability or bug, please submit the issue or send your report to security@dacc.co

## Install

The canonical download instructions for DIFS are over at: http://dacc.co/difs/docs/install/. It is **highly suggested** you follow those instructions if you are not interested in working on DIFS development.

### System Requirements

DIFS can run on most Linux, macOS, and Windows systems. We recommend running it on a machine with at least 2 GB of RAM (it’ll do fine with only one CPU core), but it should run fine with as little as 1 GB of RAM. On systems with less memory, it may not be completely stable.

### Build from Source

#### Install Go

The build process for difs requires Go 1.10 or higher. If you don't have it: [Download Go 1.10+](https://golang.org/dl/).


You'll need to add Go's bin directories to your `$PATH` environment variable e.g., by adding these lines to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

(If you run into trouble, see the [Go install instructions](https://golang.org/doc/install)).

#### Download and Compile DIFS

```
$ go get -u -d github.com/Harold-the-Axeman/dacc-iam-filesystem

$ cd $GOPATH/src/github.com/Harold-the-Axeman/dacc-iam-filesystem
$ make install
```

If you are building on FreeBSD instead of `make install` use `gmake install`.

#### Building on less common systems

If your operating system isn't officially supported, but you still want to try
building difs anyways (it should work fine in most cases), you can do the
following instead of `make install`:

```
$ make install_unsupported
```

Note: This process may break if [gx](https://github.com/whyrusleeping/gx)
(used for dependency management) or any of its dependencies break as `go get`
will always select the latest code for every dependency, often resulting in
mismatched APIs.

#### Troubleshooting

Being collected (WIP)

### Development Dependencies

If you make changes to the protocol buffers, you will need to install the [protoc compiler](https://github.com/google/protobuf).

## Usage

```
  difs - Global p2p merkle-dag filesystem.

  difs [<flags>] <command> [<arg>] ...

SUBCOMMANDS
  BASIC COMMANDS
    init          Initialize difs local configuration
    add <path>    Add a file to difs
    cat <ref>     Show difs object data
    get <ref>     Download difs objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

  DATA STRUCTURE COMMANDS
    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    files         Interact with objects as if they were a unix filesystem

  ADVANCED COMMANDS
    daemon        Start a long-running daemon process
    mount         Mount an difs read-only mountpoint
    resolve       Resolve any type of name
    name          Publish or resolve IPNS names
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo          Manipulate an difs repository

  NETWORK COMMANDS
    id            Show info about difs peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the DHT for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

  TOOL COMMANDS
    config        Manage configuration
    version       Show difs version information
    update        Download and apply go-difs updates
    commands      List all available commands

  Use 'difs <command> --help' to learn more about each command.

  difs uses a repository in the local file system. By default, the repo is located
  at ~/.difs. To change the repo location, set the $DIFS_PATH environment variable:

    export DIFS_PATH=/path/to/difsrepo
```

## Getting Started

Coming soon (WIP)

## Contributing

Patches are welcome! If you would like to contribute, but don't know what to work on, check the [whitepaper](http://dacc.co/whitepaper/Dacc.pdf) and issues list.

### Want to read our code?

Some places to get you started. (WIP)

Main file: [cmd/difs/main.go](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/cmd/difs/main.go) <br>
CLI Commands: [core/commands/](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/tree/master/core/commands) <br>
Bitswap (the data trading engine): [exchange/bitswap/](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/tree/master/exchange/bitswap)

DHT: https://github.com/libp2p/go-libp2p-kad-dht <br>
PubSub: https://github.com/libp2p/go-floodsub <br>
libp2p: https://github.com/libp2p/go-libp2p

## License

MIT
