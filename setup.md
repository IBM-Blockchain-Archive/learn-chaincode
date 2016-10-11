# Chaincode Development Environment
The following is a list of dependencies that you should install in order to develop chaincode.

### Go
- [Go 1.6 install](https://golang.org/dl/#go1.6.3)
- [Go installation instructions](https://golang.org/doc/install)
- [Go documentation and tutorials](https://golang.org/doc/)

##### Why?
Go is the language that chaincode must be written in.  Go version 1.6 is required, as that is what the Hyperledger fabric is written in.  The CLI that is installed with Go has useful commands like `go build`, which allows you to verify that your chaincode actually compiles before you attempt to deploy it to a Hyperledger network.

##### Instruction?
Follow the installation instructions linked above.  You can verify that Go is installed properly by running the following commands.  Of course, the output of `go version` may change depending on your operating system.

```
$ go version
go version go1.6.3 windows/amd64

$ echo $GOPATH
C:\gopath
```

Your `GOPATH` does not need to match the one above.  What is important is that you have this variable set to an valid directory on your filesystem.  The installation instructions linked above should take you through the setup of this environment variable.

### Hyperledger fabric
##### Why?
In order to develop chaincode locally, you will need to 

### Postman

### Git

### Node.js

## IDE Suggestions
### Visual Studio Code
### Atom