#How to write Chaincode

##What is covered?
You will be incrementally building up to a working chaincode that will be able to create generic assets.
Then you will interact with the chaincode via the network's API.

##What is chaincode?
Chaincode is a piece of code that lets you interact with a network's shared ledger.  Whenever you 'invoke' a transaction on the network, you are effectively calling a function in a piece of chaincode that reads and writes values to the ledger.


***

#Implementing Your First Chaincode

###Environment Setup
- Download and install GoLang for your OS - https://golang.org/dl/
- Add the Hyperledger shim code to your Go path by opening a command prompt/terminal and type:

	```
	go get github.com/hyperledger/fabric/core/chaincode/shim
	```

##GitHub Setup
The Bluemix IBM Blockchain service currently requires chaincode to be in a [GitHub](https://Github.com/) repository.
Therefore, you should register a GitHub account and setup Git locally on your computer. Need help? Vist [set up git](https://help.github.com/articles/set-up-git/) for more information.
- Create a new repo for this project named `learning_chaincode`
- Clone the repo to your local machine

###Download Chaincode
There is starting chaincode that you should [download](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/start/chaincode_start.go) and save to your project.
The finished chaincode we will build up to is also [available](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/finished/chaincode_finished.go).
Make sure it builds in your local environment:
- Open terminal/command prompt
- Browse to the folder that contains `chaincode_start.go` and type:

	```
	go build ./
	```
- It should return with no errors/text


###Implementing the chaincode interface
The first thing you need to do is implement the chaincode shim interface in your golang code.
The three main functions are **Init**, **Invoke**, and **Query**.
All three functions have the same prototype; they take in a function name and an array of strings.
The main difference between the functions is when they will be called.
We will be building up to a working chaincode to create generic assets.

###Dependencies
The `import` statement lists a few dependencies that you will need for your chaincode to build successfully.
- `fmt` - contains `Println` for debugging/logging.
- `errors` - standard go error format.
- `github.com/hyperledger/fabric/core/chaincode/shim` - the code that interfaces your golang code with a peer.

###Init()
Init is called when you first deploy your chaincode.
As the name implies, this function should be used to do any initialization your chaincode needs.
In our example, we use Init to configure the initial state of one variables on the ledger.

In your `chaincode.go` file,  change the `Init` function so that it stores the first element in the `args` argument to the key "hello_world".

```
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
```

This is done by using the shim function `stub.PutState`.
The first argument is the key as a string, and the second argument is the value as an array of bytes.
This function may return an error which our code inspects and returns if present.

###Invoke()
`Invoke` is called when you want to call chaincode functions to do real work.
Invocation transactions will be captured as blocks on the chain.
The structure of `Invoke` is simple.
It receives a `function` argument and based on this argument calls Go functions in the chaincode.

In your `chaincode.go` file, change the `Invoke` function so that it calls a generic write function.

```
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)

	return nil, errors.New("Received unknown function invocation")
}
```

Now that it’s looking for `write` let’s make that function somewhere in your `chaincode.go` file.

```
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the variable and value to set")
	}

	name = args[0]                            //rename for funsies
	value = args[1]
	err = stub.PutState(name, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
```

This `write` function should look similar to the `Init` change you just did.
One major difference is that you can now set the key and value for `PutState`.
This function allows you to store any key/value pair you want into the blockchain ledger.  

###Query()
As the name implies, `Query` is called whenever you query your chaincode state.
Queries do not result in blocks being added to the chain.
You will use `Query` to read the value of your chaincode state's key/value pairs.

In your `chaincode.go` file, change the `Query` function so that it calls a generic read function.

```
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {                            //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)

	return nil, errors.New("Received unknown function query")
}
```

Now that it’s looking for `read`, make that function somewhere in your `chaincode.go` file.

```
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var name, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the var to query")
	}

	name = args[0]
	valAsbytes, err := stub.GetState(name)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + name + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
```

This `read` function is using the complement to `PutState` called `GetState`.
This shim function just takes 1 string argument.
The argument is the name of the key to retrieve.
Next this function returns the value as an array of bytes back to `Query` who in turn sends it back to the REST handler.

### Main()
Finally, you need to create a short `main` function that will execute when each peer deploys their instance of the chaincode.
It just starts the chaincode and registers it with the peer.
You don’t need to add any code here beyond what is already in the example code.

```
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

#Interacting with Your First Chaincode
The fastest way to test your chaincode is to use the rest interface on your peers.
We’ve included a Swagger UI in the dashboard for your service instance that allows you to experiment with deploying chaincode without needing to write any additional code.

###Swagger API
The first step is to find the api swagger page.

1. Login to [Bluemix](https://console.ng.bluemix.net/login)
1. You probably landed on the Dashboard, but double check the top nav bar.  Click the "Dashboard" tab if you are not already there.
1. Also make sure you are in the same Bluemix "space" that contains your IBM Blockchain service. The space navigation is on the left.
1. There is a "Services" panel on this Bluemix dashboard near the bottom.  Look through your services and click your IBM Blockchain service square.
1. Now you should see a white page with the words "Welcome to the IBM Blockchain..." and there should be a teal "LAUNCH" button on the right, click it.
1. You are on the monitor page and you should see 2 tables, though the bottom one may be empty.
	- Noteworthy information on the network tab:
		- **Peer Logs** will be found in the top table. Find the row for peer 1 and then click the file-like icon in the last row.
			- It should have opened a new window. Congratulations you found your peer logs!
			- In addition to this static view there are live **streaming peer logs** in the "View Logs" tab near the top of the page.
		- **ChainCode Logs** will be found in the bottom table. There is one row for every chaincode, and they are labeled using the same chaincode hash that was returned to you when it was deployed. Find the chaincode ID you want, and then select the peer. Finally click the file-like icon.
			- It should have opened a new window. Congratulations you've found your peer's chaincode logs!
	- **Swagger Tab** is the one labeled **APIs**. Click it to see the API interactive documentation.
		- You are now on your swagger api page.

###Secure Enrollment
Calls to the `/chaincode` endpoint of the rest interface require a secure context ID.
This means that you must pass in a registered enrollID from the service credentials list in order for most REST calls to be accepted.
- Click the link "+ Network's Enroll IDs" to expand a list of enrollIDs and their secrets for your network.
- Expand the "Registrar" API section by clicking it
- Expand the `POST /registrar` section by clicking it
- Set the body's text field.  It should be JSON that contains an enrollID and secret from your list above. Example:


![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/registrar.PNG)

![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/register_response.PNG)


Now that you have enrollID set up, you can use this ID when deploying, invoking, and querying chaincode in the subsequent steps.

###Deploying the chaincode
In order to deploy chaincode through the rest interface, you will need to have the chaincode stored in a public git repository.
When you send a deploy request to a peer, you send it the url to your chaincode repository, as well as the parameters necessary to initialize the chaincode.

**Before you deploy** the code, make sure it builds locally!
- Open terminal/command prompt
- Browse to the folder that contains `chaincode_start.go` and type:

	```
	go build ./
	```
- It should return with no errors/text


- Expand the "Chaincode" API section by clicking it
- Expand the `POST /chaincode` section by clicking it
- Set the `DeploySpec` text field (make the other fields blank). Copy the example below but substitute in your chaincode repo path. Also use the same enrollID you used in the `/registrar` step.

	```
	{
		"jsonrpc": "2.0",
		"method": "deploy",
		"params": {
			"type": 1,
			"chaincodeID": {
				"path": "https://github.com/ibm-blockchain/learn-chaincode/finished"
			},
			"ctorMsg": {
				"function": "init",
				"args": [
					"hi there"
				]
			},
			"secureContext": "user_type1_191b8c2993"
		},
		"id": 1
	}
	```

The response should look like:

![Deploy Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/deploy_response.PNG)

The response for the deployment will contain an ID that is associated with this chaincode.
This is how you will reference the chaincode in any future invoke or query requests.

###Query
Next, let’s query the chaincode for the value of the `hello_world` key we set with the `Init` function.
- Expand the "Chaincode" API section by clicking it
- Expand the `POST /chaincode` section by clicking it
- Set the `QuerySpec` text field (make the other fields blank). Copy the example below but substitute in your chaincode name (the hashed ID from deploy). Also use the same enrollID you used in the `/registrar` step.

	```
	{
		"jsonrpc": "2.0",
		"method": "query",
		"params": {
			"type": 1,
			"chaincodeID": {
				"name": "CHAINCODE_HASH_HERE"
			},
			"ctorMsg": {
				"function": "read",
				"args": [
					"hello_world"
				]
			},
			"secureContext": "user_type1_xxxxxxxxx"
		},
		"id": 2
	}
	```
	
![Query Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/query_response.PNG)

Hopefully you see that the value of `hello_world` is "hi there".
This was set by the body of the deploy call you sent earlier.

###Invoke
Next, call your generic write function with `invoke`.
Change the value of `hello_world` to "go away".
- Expand the "Chaincode" API section by clicking it.
- Expand the `POST /chaincode` section by clicking it.
- Set the `InvokeSpec` text field (make the other fields blank). Copy the example below but substitute in your chaincode name (the hashed ID from deploy). Also use the same enrollID you used in the `/registrar` step.

	```
	{
		"jsonrpc": "2.0",
		"method": "invoke",
		"params": {
			"type": 1,
			"chaincodeID": {
				"name": "CHAINCODE_HASH_HERE"
			},
			"ctorMsg": {
				"function": "write",
				"args": [
					"go away"
				]
			},
			"secureContext": "user_type1_xxxxxxxxx"
		},
		"id": 3
	}
	```

![Invoke Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/invoke_response.PNG)

Now to test if it's stuck, just re-run the query above.

![Query2 Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/query2_response.PNG)


That’s all it takes to write basic chaincode.
