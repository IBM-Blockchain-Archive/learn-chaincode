# Chaincode Development Environment
The following is a list of dependencies that you should install in order to develop chaincode.

### Git
- [Git download page](https://git-scm.com/downloads)
- [Pro Git ebook](https://git-scm.com/book/en/v2)
- [Git Desktop (for those uncomfortable with git's CLI)](https://desktop.github.com/)

Git is a great version control tool to be familiar with, both for chaincode development and software development in general.  Also, git bash 

##### Instructions
After following the installation instructions above, you can verify that git is installed using the following command:

```
$ git version
git version 2.9.0.windows.1
```

### Go
- [Go 1.6 install](https://golang.org/dl/#go1.6.3)
- [Go installation instructions](https://golang.org/doc/install)
- [Go documentation and tutorials](https://golang.org/doc/)

Go is the language that chaincode must be written in.  Go version 1.6 is required, as that is what the Hyperledger fabric is written in.  The CLI that is installed with Go has useful commands like `go build`, which allows you to verify that your chaincode actually compiles before you attempt to deploy it to a Hyperledger network.

##### Instructions
Follow the installation instructions linked above.  You can verify that Go is installed properly by running the following commands.  Of course, the output of `go version` may change depending on your operating system.

```
$ go version
go version go1.6.3 windows/amd64

$ echo $GOPATH
C:\gopath
```

Your `GOPATH` does not need to match the one above.  What is important is that you have this variable set to an valid directory on your filesystem.  The installation instructions linked above should take you through the setup of this environment variable.

### Hyperledger fabric
- [v0.5-developer-preview Hyperledger fabric](https://github.com/hyperledger-archives/fabric/tree/v0.5-developer-preview)
- [v0.6-preview Hyperledger fabric](https://gerrit.hyperledger.org/r/gitweb?p=fabric.git;a=shortlog;h=refs/heads/v0.6)
- [master branch of the Hyperledger fabric](https://gerrit.hyperledger.org/r/gitweb?p=fabric.git;a=summary)

Any piece of chaincode that you write will need to import the chaincode shim from the Hyperledger fabric in order to be able to read and write data on the ledger.  To be able to compile chaincode locally, which you will be doing a lot, you will need to have the fabric code present in your `GOPATH`.

##### Instructions

Three different releases of the fabric are linked above.  The release you choose above should match the Hyperledger network you are deploying your chaincode on.  You will need to make sure that the fabric release you choose is stored under `$GOPATH/hyperledger/fabric`.

The instructions below should take you through the process of properly installing the v0.5 release on your `GOPATH`.
```

# Create the parent directories on your GOPATH
mkdir -p $GOPATH/github.com/hyperledger
cd $GOAPTH/github.com/hyperledger

# Clone the appropriate release codebase into $GOPATH/github.com/hyperledger/fabric
# Note that the v0.5 release is a branch of the repository
git clone -b v0.5-developer-preview https://github.com/hyperledger-archives/fabric.git
```

If the fabric is not installed properly on your `GOPATH`, you will see errors like the one below when building your chaincode:
```
$ go build .
chaincode_example02.go:27:2: cannot find package "github.com/hyperledger/fabric/core/chaincode/shim" in any of:
        C:\Go\src\github.com\hyperledger\fabric\core\chaincode\shim (from $GOROOT)
        C:\gopath\src\github.com\hyperledger\fabric\core\chaincode\shim (from $GOPATH)
```

A list of known specific releases is included below

- [Blockchain service on Bluemix](https://new-console.ng.bluemix.net/catalog/services/blockchain/) - use the v0.5-developer-preview release



### Postman
- [Home page](https://www.getpostman.com/)

Postman is a REST API testing tool.  The REST API, though it is deprecated, is an easy way to iterate on deploy your chaincode without



### Node.js
- [Download links](https://nodejs.org/en/download/)

Node.js is NOT necessary to develop chaincode, but most of our demos are built on Node.js, so it might be handy to go ahead and install it now.  Download the appropriate installation package and make sure the following commands work on your machine:
```
$ node -v
v4.4.7

$ npm -v
3.10.5
```

## IDE Suggestions
### Visual Studio Code
- [Download links](https://code.visualstudio.com/#alt-downloads)

Visual Studio Code is a free IDE that supports both Node.js and Go through plugins.  All of our demos and examples use either one or both of these languages.  It also has tab support, git integration, and debugging support.

### Atom
- [Home page](https://atom.io/)

Like VS Code, Atom has plugins to support any of the languages needed to develop chaincode or modify our examples.