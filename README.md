# go-difs

![banner](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/docs/difs.png)

[![](https://img.shields.io/badge/made%20by-DACC-yellowgreen.svg)](http://dacc.co)
[![Build Status](https://img.shields.io/badge/build-passing-green.svg)]()
[![Python](https://img.shields.io/badge/Golang-1.10.3-blue.svg)](https://golang.org/)

> DIPS implementation in Go

## What is DIFS?

![banner](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/docs/difs-blueprint.png)

DIFS is DACC IAM File System based on [IPFS](https://github.com/ipfs/ipfs) with IAM(Identity and Access Management) system.it is the most fundamental infrastructure and innovation of DACC architecture. It includes file sharing system, permission maps, migration engine from centralized storage to decentralized storage. More importantly, the whole system is modular and independent from the rest of DACC architecture, so any other public chain that needs modern decentralized storage with IAM capabilities can easily implement DACC file system.

For more information see -> http://dacc.co/whitepaper/Dacc.pdf

Please put all issues regarding:
  - DIFS _design_ in the [DIFS protocol repo issues](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues).
  - Go DIFS _implementation_ in [this repo](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues).

## Table of Contents

- [Security Issues](#security-issues)
- [Install](#install)
  - [System Requirements](#system-requirements)
  - [Install prebuilt packages](#install-prebuilt-packages)
  - [Build from Source](#build-from-source)
    - [Install Go](#install-go)
    - [Download and Compile IPFS](#download-and-compile-ipfs)
    - [Troubleshooting](#troubleshooting)
  - [Development Dependencies](#development-dependencies)
  - [Updating](#updating)
- [Usage](#usage)
- [Getting Started](#getting-started)
  - coming soon
- [Contributing](#contributing)
  - [Want to hack on DIFS?](#want-to-hack-on-difs)
  - [Want to read our code?](#want-to-read-our-code)
- [License](#license)

## Security Issues

The DIFS protocol and its implementations are still in preliminary stage. This means that there may be problems in our protocols, or there may be mistakes in our implementations. And DIFS is not production-ready yet. If you discover a security issue, please bring it to our attention! thx

If you find a vulnerability or bug, please submit the issue or send your report to security@dacc.co

## Install

The canonical download instructions for DIFS are over at: http://dacc.co/difs/docs/install/. It is **highly suggested** you follow those instructions if you are not interested in working on DIFS development.

### System Requirements

IPFS can run on most Linux, macOS, and Windows systems. We recommend running it on a machine with at least 2 GB of RAM (it’ll do fine with only one CPU core), but it should run fine with as little as 1 GB of RAM. On systems with less memory, it may not be completely stable.

### Install prebuilt packages

We host prebuilt binaries over at our [distributions page](https://ipfs.io/ipns/dist.ipfs.io#go-ipfs).

From there:
- Click the blue "Download go-ipfs" on the right side of the page.
- Open/extract the archive.
- Move `ipfs` to your path (`install.sh` can do it for you).

You can also download go-ipfs from this project's GitHub releases page if you are unable to access ipfs.io.

### From Linux package managers

- [Arch Linux](#arch-linux)
- [Nix](#nix)
- [Snap](#snap)

#### Arch Linux

In Arch Linux go-ipfs is available as
[go-ipfs](https://www.archlinux.org/packages/community/x86_64/go-ipfs/) package.

	$ sudo pacman -S go-ipfs

Development version of go-ipfs is also on AUR under
[go-ipfs-git](https://aur.archlinux.org/packages/go-ipfs-git/).
You can install it using your favourite AUR Helper or manually from AUR.

### Nix

For Linux and MacOSX you can use the purely functional package manager [Nix](https://nixos.org/nix/):

```
$ nix-env -i ipfs
```
You can also install the Package by using it's attribute name, which is also `ipfs`.

#### Snap

With snap, in any of the [supported Linux distributions](https://snapcraft.io/docs/core/install):

    $ sudo snap install ipfs

### Build from Source

#### Install Go

The build process for ipfs requires Go 1.10 or higher. If you don't have it: [Download Go 1.10+](https://golang.org/dl/).


You'll need to add Go's bin directories to your `$PATH` environment variable e.g., by adding these lines to your `/etc/profile` (for a system-wide installation) or `$HOME/.profile`:

```
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:$GOPATH/bin
```

(If you run into trouble, see the [Go install instructions](https://golang.org/doc/install)).

#### Download and Compile IPFS

```
$ go get -u -d github.com/Harold-the-Axeman/dacc-iam-filesystem

$ cd $GOPATH/src/github.com/Harold-the-Axeman/dacc-iam-filesystem
$ make install
```

If you are building on FreeBSD instead of `make install` use `gmake install`.

#### Building on less common systems

If your operating system isn't officially supported, but you still want to try
building ipfs anyways (it should work fine in most cases), you can do the
following instead of `make install`:

```
$ make install_unsupported
```

Note: This process may break if [gx](https://github.com/whyrusleeping/gx)
(used for dependency management) or any of its dependencies break as `go get`
will always select the latest code for every dependency, often resulting in
mismatched APIs.

#### Troubleshooting

* Separate [instructions are available for building on Windows](docs/windows.md).
* Also, [instructions for OpenBSD](docs/openbsd.md).
* `git` is required in order for `go get` to fetch all dependencies.
* Package managers often contain out-of-date `golang` packages.
  Ensure that `go version` reports at least 1.10. See above for how to install go.
* If you are interested in development, please install the development
dependencies as well.
* *WARNING: Older versions of OSX FUSE (for Mac OS X) can cause kernel panics when mounting!*
  We strongly recommend you use the [latest version of OSX FUSE](http://osxfuse.github.io/).
  (See https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues/177)
* For more details on setting up FUSE (so that you can mount the filesystem), see the docs folder.
* Shell command completion is available in `misc/completion/ipfs-completion.bash`. Read [docs/command-completion.md](docs/command-completion.md) to learn how to install it.
* See the [init examples](https://github.com/ipfs/website/tree/master/static/docs/examples/init) for how to connect IPFS to systemd or whatever init system your distro uses.

### Development Dependencies

If you make changes to the protocol buffers, you will need to install the [protoc compiler](https://github.com/google/protobuf).

### Updating

#### Updating using ipfs-update
IPFS has an updating tool that can be accessed through `ipfs update`. The tool is
not installed alongside IPFS in order to keep that logic independent of the main
codebase. To install `ipfs update`, [download it here](https://ipfs.io/ipns/dist.ipfs.io/#ipfs-update).

#### Downloading IPFS builds using IPFS
List the available versions of go-ipfs:
```
$ ipfs cat /ipns/dist.ipfs.io/go-ipfs/versions
```

Then, to view available builds for a version from the previous command ($VERSION):
```
$ ipfs ls /ipns/dist.ipfs.io/go-ipfs/$VERSION
```

To download a given build of a version:
```
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_darwin-386.tar.gz # darwin 32-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_darwin-amd64.tar.gz # darwin 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_freebsd-amd64.tar.gz # freebsd 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-386.tar.gz # linux 32-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-amd64.tar.gz # linux 64-bit build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_linux-arm.tar.gz # linux arm build
$ ipfs get /ipns/dist.ipfs.io/go-ipfs/$VERSION/go-ipfs_$VERSION_windows-amd64.zip # windows 64-bit build
```

## Usage

```
  ipfs - Global p2p merkle-dag filesystem.

  ipfs [<flags>] <command> [<arg>] ...

SUBCOMMANDS
  BASIC COMMANDS
    init          Initialize ipfs local configuration
    add <path>    Add a file to ipfs
    cat <ref>     Show ipfs object data
    get <ref>     Download ipfs objects
    ls <ref>      List links from an object
    refs <ref>    List hashes of links from an object

  DATA STRUCTURE COMMANDS
    block         Interact with raw blocks in the datastore
    object        Interact with raw dag nodes
    files         Interact with objects as if they were a unix filesystem

  ADVANCED COMMANDS
    daemon        Start a long-running daemon process
    mount         Mount an ipfs read-only mountpoint
    resolve       Resolve any type of name
    name          Publish or resolve IPNS names
    dns           Resolve DNS links
    pin           Pin objects to local storage
    repo          Manipulate an IPFS repository

  NETWORK COMMANDS
    id            Show info about ipfs peers
    bootstrap     Add or remove bootstrap peers
    swarm         Manage connections to the p2p network
    dht           Query the DHT for values or peers
    ping          Measure the latency of a connection
    diag          Print diagnostics

  TOOL COMMANDS
    config        Manage configuration
    version       Show ipfs version information
    update        Download and apply go-ipfs updates
    commands      List all available commands

  Use 'ipfs <command> --help' to learn more about each command.

  ipfs uses a repository in the local file system. By default, the repo is located
  at ~/.ipfs. To change the repo location, set the $IPFS_PATH environment variable:

    export IPFS_PATH=/path/to/ipfsrepo
```

## Getting Started

See also: http://ipfs.io/docs/getting-started/

To start using IPFS, you must first initialize IPFS's config files on your
system, this is done with `ipfs init`. See `ipfs init --help` for information on
the optional arguments it takes. After initialization is complete, you can use
`ipfs mount`, `ipfs add` and any of the other commands to explore!

### Some things to try

Basic proof of 'ipfs working' locally:

	echo "hello world" > hello
	ipfs add hello
	# This should output a hash string that looks something like:
	# QmT78zSuBmuS4z925WZfrqQ1qHaJ56DQaTfyMUF7F8ff5o
	ipfs cat <that hash>


### Docker usage

An IPFS docker image is hosted at [hub.docker.com/r/ipfs/go-ipfs](https://hub.docker.com/r/ipfs/go-ipfs/).
To make files visible inside the container you need to mount a host directory
with the `-v` option to docker. Choose a directory that you want to use to
import/export files from IPFS. You should also choose a directory to store
IPFS files that will persist when you restart the container.

    export ipfs_staging=</absolute/path/to/somewhere/>
    export ipfs_data=</absolute/path/to/somewhere_else/>

Start a container running ipfs and expose ports 4001, 5001 and 8080:

    docker run -d --name ipfs_host -v $ipfs_staging:/export -v $ipfs_data:/data/ipfs -p 4001:4001 -p 127.0.0.1:8080:8080 -p 127.0.0.1:5001:5001 ipfs/go-ipfs:latest

Watch the ipfs log:

    docker logs -f ipfs_host

Wait for ipfs to start. ipfs is running when you see:

    Gateway (readonly) server
    listening on /ip4/0.0.0.0/tcp/8080

You can now stop watching the log.

Run ipfs commands:

    docker exec ipfs_host ipfs <args...>

For example: connect to peers

    docker exec ipfs_host ipfs swarm peers

Add files:

    cp -r <something> $ipfs_staging
    docker exec ipfs_host ipfs add -r /export/<something>

Stop the running container:

    docker stop ipfs_host

### Troubleshooting

If you have previously installed IPFS before and you are running into
problems getting a newer version to work, try deleting (or backing up somewhere
else) your IPFS config directory (~/.ipfs by default) and rerunning `ipfs init`.
This will reinitialize the config file to its defaults and clear out the local
datastore of any bad entries.

Please direct general questions and help requests to our
[forum](https://discuss.ipfs.io) or our IRC channel (freenode #ipfs).

If you believe you've found a bug, check the [issues list](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/issues)
and, if you don't see your problem there, either come talk to us on IRC (freenode #ipfs) or
file an issue of your own!

## Contributing

Patches are welcome! If you would like to contribute, but don't know what to work on, check the issues list.

### Want to hack on DIFS?

[![](https://cdn.rawgit.com/jbenet/contribute-ipfs-gif/master/img/contribute.gif)](https://github.com/ipfs/community/blob/master/contributing.md)

### Want to read our code?

Some places to get you started. (WIP)

Main file: [cmd/ipfs/main.go](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/blob/master/cmd/ipfs/main.go) <br>
CLI Commands: [core/commands/](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/tree/master/core/commands) <br>
Bitswap (the data trading engine): [exchange/bitswap/](https://github.com/Harold-the-Axeman/dacc-iam-filesystem/tree/master/exchange/bitswap)

DHT: https://github.com/libp2p/go-libp2p-kad-dht <br>
PubSub: https://github.com/libp2p/go-floodsub <br>
libp2p: https://github.com/libp2p/go-libp2p

## License

MIT
