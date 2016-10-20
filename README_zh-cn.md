# 如何编写链码
本教程演示了构建 [Hyperledger fabric](https://github.com/hyperledger/fabric) 链码应用程序所需的基本构建区块和函数。你将逐步创建一个可以用于通用资产的工作链码。
然后，你将使用网络 API 与链码交互。阅读并完成本教程后，你应该能够明确回答以下问题：
- 什么是链码？
- 如何实现链码？
- 实现链码需要什么依赖关系？
- 主要功能是什么？
- 如何编译我的链码？
- 如何传递不同的值到我的参数？
- 如何安全地在我的网络上注册用户？
- 如何使用 REST API 与我的链码交互？


## 什么是链码？
链码是一段代码，它被部署到 [Hyperledger fabric](https://github.com/hyperledger/fabric) 节点的网络中，实现与该网络的共享账本的交互。

***

# 实现你的第一个链码

## 环境配置
目前，Hyperledger fabric 支持在 Go 中编写的链码。我们需要 [Go 1.6](https://blog.golang.org/go1.6)。 如果你已经安装了 Go 1.6 并配置好了环境，则可以跳到步骤3。如果已经配置了 Hyperledger fabric [开发环境](https://github.com/hyperledger/fabric/blob/master/docs/dev-setup/devenv.md)，则可以完全跳过此部分。

1. 参照 [Golang](https://golang.org/doc/install) 为你的操作系统下载和安装 Golang。如果这是你第一次安装 Go，你应该遵循所有的说明，包括测试它是否已经[正确安装](https://golang.org/doc/install#testing)，应该避免在自定义位置安装 Go。
2. 打开命令提示符或终端，输入以下命令，把 Hyperledger shim 代码添加到你的工作目录（工作目录即 $GOPATH 环境变量所指的位置）中：

	```
	cd $GOPATH
	go get github.com/hyperledger/fabric/core/chaincode/shim
	```
3. [IBM Bluemix](https://console.ng.bluemix.net/) 的 Blockchain 服务当前要求链码保存在 [GitHub](https://Github.com/) 仓库中。因此，如果你还没有 GitHub 帐户，你应该先[注册一个](http://github.com)。
4. 如果你还没有在本地电脑上设置 Git，你应该参照[这里](https://help.github.com/articles/set-up-git/))设置。

## 将链码部署到 IBM Bluemix
1. Fork 本仓库到你的 Github 账户中（向上滚动到顶部，然后点击 **Fork**）。
2. 把你 fork 的仓库 clone 到你的 $GOPATH 目录中。

	```
	cd $GOPATH
	mkdir -p src/github.com/<yourgithubid>/
	cd src/github.com/<yourgithubid>/
	git clone https://github.com/<yourgithubid>/learn-chaincode.git
	```
3. 请注意，在本教程中，我们提供了两个不同版本的链码：[Start](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/start/chaincode_start.go) - 你将要在此基础上进行开发的框架链码，[Finished](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/finished/chaincode_finished.go) - 已完成的链码。
4. 确保它在你的本地环境中构建：
	- 打开终端或命令提示符
	
	```
	cd $GOPATH/src/github.com/<yourgithubid>/learn-chaincode/start
	go build ./
	```
	- 它应该没有错误，如果发生错误，请确保你已正确遵循了上述步骤，包括[测试 Go 安装]https://golang.org/doc/install#testing)。


## 实现链码接口
你需要做的第一件事是在你的 golang 代码中实现链码 shim 接口。
它的三个主要功能是 ** Init **，** Invoke ** 和 ** Query **。
所有三个函数具有相同的原型; 它们接受一个函数名和一个字符串数组。
三个函数的主要区别在于它们何时被调用。
我们将编写一个可以工作的链码用来创建通用资产。

## 依赖
`import` 语句列出了链码能够成功构建的一些依赖关系。
- `fmt` - 包含用于调试/日志记录的 `Println`。
- `errors` - 标准 go 错误格式。
- `github.com/hyperledger/fabric/core/chaincode/shim` - golang 代码与节点交互的接口代码。

## Init()
当首次部署你的链码时调用 Init。
顾名思义，此函数用于链码所需的所有初始化工作。
在我们的示例中，我们使用 Init 在账本上初始化一个变量。

在你的 `chaincode_start.go` 文件中，修改 `Init` 函数，以便将 `args` 参数中的第一个元素存储到键 “hello_world” 中。

```go
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

这是通过使用 shim 函数 `stub.PutState` 来完成的。
第一个参数是作为字符串的键，第二个参数是作为字节数组的值。
此函数可能返回错误，这段代码会检查是否有错误并返回。

## Invoke()
当你想调用链码函数来做真正的工作时，就要调用 `Invoke`。
调用交易将作为链上的区块被捕获。
`Invoke` 的结构很简单。
它接收一个 `function` 参数，并基于此参数调用链码中的 Go 函数。

在你的 `chaincode_start.go` 文件中，修改 `Invoke` 函数，让它调用一个通用的 `write` 函数。

```go
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

现在，它正在寻找 `write`，让我们把这个函数写入你的 `chaincode_start.go` 文件。

```go
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, value string
	var err error
	fmt.Println("running write()")

	if len(args) != 2 {
		return nil, errors.New("Incorrect number of arguments. Expecting 2. name of the key and value to set")
	}

	key = args[0]                            //rename for fun
	value = args[1]
	err = stub.PutState(key, []byte(value))  //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}
```

这个 `write` 函数应该看起来类似于你刚刚对 `Init` 的修改。
一个主要的区别是，你现在可以设置 `PutState` 的键和值。
此功能允许你将任何所需的键/值对存储在区块链账本中。

## Query()
顾名思义，无论何时查询链码状态，都会调用 `Query`。
查询不会导致区块被添加到链中。
你将使用 `Query` 读取链码状态中键/值对的值。

在你的 `chaincode_start.go` 文件中，修改 `Query` 函数，让它调用一个通用的 `read` 函数。

```go
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

现在，它正在寻找 `read`，让我们把该函数写入 `chaincode_start.go` 文件。

```go
func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
```

`read` 函数使用了 `GetState`，它的作用与 `PutState` 恰好相反。
shim 函数只需要一个字符串参数。
参数是要检索的键的名称。
接下来，此函数将字符数组返回到 `Query`，然后将它返回给 REST 句柄。

## Main()
最后，你需要创建一个简短的 `main` 函数，当每个节点部署链码的实例时执行。
它只是启动链码并将其注册到节点上。
你不需要为此函数添加任何代码。`chaincode_start.go` 和 `chaincode_finished.go` 都有一个 `main` 函数，它位于文件的顶部。该函数如下所示：

```go
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

## 需要帮助？
如果你在任何时候被卡住或有什么困惑，只需去查看 `chaincode_finished.go` 文件。使用该文件检查您正在编写的 `chaincode_start.go` 代码段是否正确。


# 与你的第一个链码交互
测试你的链码的最快的方法是使用节点上的 REST 接口。
我们在你的服务实例的 Dashboard 中包含了一个 Swagger UI，它允许你尝试部署链码，而无需编写任何其他代码。

## Swagger API
第一步，找到 api swagger 页面。

1. 登录到 [IBM Bluemix](https://console.ng.bluemix.net/login)
1. 你可能已经登录到了 Dashboard，请仔细检查顶部导航栏。如果你还不在那里，请点击 “Dashboard” 标签。
1. 同时，还要确保你在同一个包含 IBM Blockchain 服务的 Bluemix “space” 中。space 导航在左侧。
1. 在 Bluemix dashboard 底部有一个 “Services” 面板。单击 IBM Blockchain 的 service 查看你的服务。
1. 现在你应该看到一个带有 “Welcome to the IBM Blockchain...” 的页面，右侧有一个蓝色的 “LAUNCH” 按钮，点击该按钮。
1. 在监视页面上，你应该看到两个表，最下面一个表有可能是空的。
	- 网络选项卡上的值得注意的信息：
		- **Peer Logs** 是位于上面的表。找到节点 1 所在的行，然后点击最后一行中类似文件的图标。
			- 它应该会打开一个新窗口。恭喜你找到了你的节点日志！
			- 除了此静态视图，还有靠近页面顶部 **View Logs** 选项卡中的实时流节点日志。
		- ** ChainCode Logs ** 是位于下面的表。每个链码都有一行，它们使用链码哈希来标记，链码哈希是在部署链码时返回的哈希值。找到所需的链码 ID，然后选择节点。最后单击类似文件的图标。
			- 它应该会打开一个新窗口。恭喜你找到了你的节点的链码日志！
	- 标记为 ** APIs ** 的是** Swagger Tab **。点击它查看 API 交互文档。
		- 你现在在 Swagger API 页面上了。

## 安全注册
调用 REST 接口的 `/chaincode' 需要一个安全的上下文 ID。
这意味着你必须从服务凭据列表中传入注册的 enrollID，服务器才可以接受大多数 REST 调用。
- 点击链接 **+ Network's Enroll IDs** 展开你的网络 enrollID 及其密码列表。
- 打开记事本并复制一组凭据。接下来会用到它们。
- 通过单击展开 “Registrar” API 部分
- 通过单击展开 “POST / registrar” 部分
- 设置正文的文本字段。它应该是 JSON 格式，包含了列表上的一个 enrollID 和密码。例如：


![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/registrar.PNG)

![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/register_response.PNG)


如果你没有收到 “Login successful” 响应，请返回并确保您已正确复制了 enrollID 和密码。现在你已设置了 enrollID，你可以在后续步骤中部署、调用和查询链码时使用此 ID。

## 部署链码（Deploying the chaincode）
为了通过 REST 接口部署链码，你需要将链码存储在一个公开的 git 仓库中。
当你向节点发送部署请求时，会将你的链码仓库 URL 发送给节点，同时还包括初始化链码所需的参数。

**在你部署链码之前**，确保它能在本地构建！
- 打开终端或命令提示符
- 进入包含 `chaincode_start.go` 的文件夹，输入：

	```
	go build ./
	```
- 它应该不会返回错误/文本


- 通过单击展开 `Chaincode` API 部分
- 通过单机展开 `POST /chaincode` 部分
- 设置 `DeploySpec` 文本字段（使其他字段为空）。复制下面的例子，但替换你的链码仓库路径。仍使用你在 `/registrar` 步骤中使用的相同 enrollID。
- `“path”`：应该类似 `“https://github.com/johndoe/learn-chaincode/finished”`。它是你的 fork 的项目中的一个目录，`chaincode_finished.go` 文件位于该目录中。

	```js
	{
		"jsonrpc": "2.0",
		"method": "deploy",
		"params": {
			"type": 1,
			"chaincodeID": {
				"path": "https://github.com/johndoe/learn-chaincode/finished"
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

响应应该如下所示：

![Deploy Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/deploy_response.PNG)

部署的响应包含与该链码关联的 ID。ID是一个 128 个字符的字母数字哈希。也在你的记事本上复制此 ID。你现在应该有一组 enrollID 凭据和标识你的链码的加密哈希。
接下来在调用或查询请求中引用链码将会使用到它们。

## 查询（Query）
接下来，让我们查询链码的 `hello_world` 键的值，之前我们使用了 `Init` 函数为它设置了初始值。
- 通过单击展开 `Chaincode` API 部分
- 通过单击展开 `POST /chaincode` 部分
设置 `QuerySpec` 文本字段（使其他字段为空）。复制下面的示例，但替换你的链码名称（部署时获取的哈希 ID）。仍使用你在 `/registrar` 步骤中使用的相同 enrollID。

	```js
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

希望你看到的 `hello_world` 的值是 “hi there”。
该值是由你之前调用部署设置的。

## 调用（Invoke）
接下来，使用 `invoke` 调用你的 `write` 函数。
将 `hello_world` 的值更改为 “go away”。
- 通过单击展开 `Chaincode` API 部分。
- 通过单击展开 `POST /chaincode` 部分。
- 设置 `InvokeSpec` 文本字段（使其他字段为空）。 复制下面的示例，但替换你的链码名称（部署时获取的哈希 ID）。仍使用你在 `/registrar` 步骤中使用的相同 enrollID。

	```js
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
					"hello_world",
					"go away"
				]
			},
			"secureContext": "user_type1_xxxxxxxxx"
		},
		"id": 3
	}
	```

![Invoke Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/invoke_response.PNG)

现在要测试它是否卡住，只需重新运行上面的查询操作。

![Query2 Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/query2_response.PNG)


这就是编写基本链码所需要的全部内容。
