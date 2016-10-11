# Chaincode Development Environment
The following is a list of dependencies that you should install in order to develop chaincode.

### Git
- [Git download page](https://git-scm.com/downloads)
- [Pro Git ebook](https://git-scm.com/book/en/v2)

Git is a great version control tool to be familiar with, both for chaincode development and software development in general.

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
Three different releases of the fabric are linked above.  The release you choose above should match the Hyperledger network you are deploying your chaincode on.  A list of known specific releases is included below

- [Blockchain service on Bluemix](https://new-console.ng.bluemix.net/catalog/services/blockchain/) - use the v0.5-developer-preview release

### Postman



### Node.js

## IDE Suggestions
### Visual Studio Code
##### W
### Atom