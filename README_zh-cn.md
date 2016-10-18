# ��α�д����
���̳���ʾ�˹��� [Hyperledger fabric](https://github.com/hyperledger/fabric) ����Ӧ�ó�������Ļ�����������ͺ������㽫�𲽴���һ����������ͨ���ʲ��Ĺ������롣
Ȼ���㽫ʹ������ API �����뽻�����Ķ�����ɱ��̳̺���Ӧ���ܹ���ȷ�ش��������⣺
- ʲô�����룿
- ���ʵ�����룿
- ʵ��������Ҫʲô������ϵ��
- ��Ҫ������ʲô��
- ��α����ҵ����룿
- ��δ��ݲ�ͬ��ֵ���ҵĲ�����
- ��ΰ�ȫ�����ҵ�������ע���û���
- ���ʹ�� REST API ���ҵ����뽻����


## ʲô�����룿
������һ�δ��룬�������� [Hyperledger fabric](https://github.com/hyperledger/fabric) �ڵ�������У�ʵ���������Ĺ����˱��Ľ�����

***

# ʵ����ĵ�һ������

## ��������
Ŀǰ��Hyperledger fabric ֧���� Go �б�д�����롣������Ҫ [Go 1.6](https://blog.golang.org/go1.6)�� ������Ѿ���װ�� Go 1.6 �����ú��˻������������������3������Ѿ������� Hyperledger fabric [��������](https://github.com/hyperledger/fabric/blob/master/docs/dev-setup/devenv.md)���������ȫ�����˲��֡�

1. ���� [Golang](https://golang.org/doc/install) Ϊ��Ĳ���ϵͳ���غͰ�װ Golang������������һ�ΰ�װ Go����Ӧ����ѭ���е�˵���������������Ƿ��Ѿ�[��ȷ��װ](https://golang.org/doc/install#testing)��Ӧ�ñ������Զ���λ�ð�װ Go��
2. ��������ʾ�����նˣ�������������� Hyperledger shim ������ӵ���Ĺ���Ŀ¼������Ŀ¼�� $GOPATH ����������ָ��λ�ã��У�

	```
	cd $GOPATH
	go get github.com/hyperledger/fabric/core/chaincode/shim
	```
3. [IBM Bluemix](https://console.ng.bluemix.net/) �� Blockchain ����ǰҪ�����뱣���� [GitHub](https://Github.com/) �ֿ��С���ˣ�����㻹û�� GitHub �ʻ�����Ӧ����[ע��һ��](http://github.com)��
4. ����㻹û���ڱ��ص��������� Git����Ӧ�ò���[����](https://help.github.com/articles/set-up-git/))���á�

## �����벿�� IBM Bluemix
1. Fork ���ֿ⵽��� Github �˻��У����Ϲ�����������Ȼ���� **Fork**����
2. ���� fork �Ĳֿ� clone ����� $GOPATH Ŀ¼�С�

	```
	cd $GOPATH
	mkdir -p src/github.com/<yourgithubid>/
	cd src/github.com/<yourgithubid>/
	git clone https://github.com/<yourgithubid>/learn-chaincode.git
	```
3. ��ע�⣬�ڱ��̳��У������ṩ��������ͬ�汾�����룺[Start](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/start/chaincode_start.go) - �㽫Ҫ�ڴ˻����Ͻ��п����Ŀ�����룬[Finished](https://github.com/IBM-Blockchain/learn-chaincode/blob/master/finished/chaincode_finished.go) - ����ɵ����롣
4. ȷ��������ı��ػ����й�����
	- ���ն˻�������ʾ��
	
	```
	cd $GOPATH/src/github.com/<yourgithubid>/learn-chaincode/start
	go build ./
	```
	- ��Ӧ��û�д����������������ȷ��������ȷ��ѭ���������裬����[���� Go ��װ]https://golang.org/doc/install#testing)��


## ʵ������ӿ�
����Ҫ���ĵ�һ����������� golang ������ʵ������ shim �ӿڡ�
����������Ҫ������ ** Init **��** Invoke ** �� ** Query **��
������������������ͬ��ԭ��; ���ǽ���һ����������һ���ַ������顣
������������Ҫ�����������Ǻ�ʱ�����á�
���ǽ���дһ�����Թ�����������������ͨ���ʲ���

## ����
`import` ����г��������ܹ��ɹ�������һЩ������ϵ��
- `fmt` - �������ڵ���/��־��¼�� `Println`��
- `errors` - ��׼ go �����ʽ��
- `github.com/hyperledger/fabric/core/chaincode/shim` - golang ������ڵ㽻���Ľӿڴ��롣

## Init()
���״β����������ʱ���� Init��
����˼�壬�˺�������������������г�ʼ��������
�����ǵ�ʾ���У�����ʹ�� Init ���˱��ϳ�ʼ��һ��������

����� `chaincode_start.go` �ļ��У��޸� `Init` �������Ա㽫 `args` �����еĵ�һ��Ԫ�ش洢���� ��hello_world�� �С�

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

����ͨ��ʹ�� shim ���� `stub.PutState` ����ɵġ�
��һ����������Ϊ�ַ����ļ����ڶ�����������Ϊ�ֽ������ֵ��
�˺������ܷ��ش�����δ�������Ƿ��д��󲢷��ء�

## Invoke()
������������뺯�����������Ĺ���ʱ����Ҫ���� `Invoke`��
���ý��׽���Ϊ���ϵ����鱻����
`Invoke` �Ľṹ�ܼ򵥡�
������һ�� `function` �����������ڴ˲������������е� Go ������

����� `chaincode_start.go` �ļ��У��޸� `Invoke` ��������������һ��ͨ�õ� `write` ������

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

���ڣ�������Ѱ�� `write`�������ǰ��������д����� `chaincode_start.go` �ļ���

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

��� `write` ����Ӧ�ÿ�������������ոն� `Init` ���޸ġ�
һ����Ҫ�������ǣ������ڿ������� `PutState` �ļ���ֵ��
�˹��������㽫�κ�����ļ�/ֵ�Դ洢���������˱��С�

## Query()
����˼�壬���ۺ�ʱ��ѯ����״̬��������� `Query`��
��ѯ���ᵼ�����鱻��ӵ����С�
�㽫ʹ�� `Query` ��ȡ����״̬�м�/ֵ�Ե�ֵ��

����� `chaincode_start.go` �ļ��У��޸� `Query` ��������������һ��ͨ�õ� `read` ������

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

���ڣ�������Ѱ�� `read`�������ǰѸú���д�� `chaincode_start.go` �ļ���

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

`read` ����ʹ���� `GetState`������������ `PutState` ǡ���෴��
shim ����ֻ��Ҫһ���ַ���������
������Ҫ�����ļ������ơ�
���������˺������ַ����鷵�ص� `Query`��Ȼ�������ظ� REST �����

## Main()
�������Ҫ����һ����̵� `main` ��������ÿ���ڵ㲿�������ʵ��ʱִ�С�
��ֻ���������벢����ע�ᵽ�ڵ��ϡ�
�㲻��ҪΪ�˺�������κδ��롣`chaincode_start.go` �� `chaincode_finished.go` ����һ�� `main` ��������λ���ļ��Ķ������ú���������ʾ��

```go
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
```

## ��Ҫ������
��������κ�ʱ�򱻿�ס����ʲô����ֻ��ȥ�鿴 `chaincode_finished.go` �ļ���ʹ�ø��ļ���������ڱ�д�� `chaincode_start.go` ������Ƿ���ȷ��


# ����ĵ�һ�����뽻��
���������������ķ�����ʹ�ýڵ��ϵ� REST �ӿڡ�
��������ķ���ʵ���� Dashboard �а�����һ�� Swagger UI���������㳢�Բ������룬�������д�κ��������롣

## Swagger API
��һ�����ҵ� api swagger ҳ�档

1. ��¼�� [IBM Bluemix](https://console.ng.bluemix.net/login)
1. ������Ѿ���¼���� Dashboard������ϸ��鶥��������������㻹����������� ��Dashboard�� ��ǩ��
1. ͬʱ����Ҫȷ������ͬһ������ IBM Blockchain ����� Bluemix ��space�� �С�space ��������ࡣ
1. �� Bluemix dashboard �ײ���һ�� ��Services�� ��塣���� IBM Blockchain �� service �鿴��ķ���
1. ������Ӧ�ÿ���һ������ ��Welcome to the IBM Blockchain...�� ��ҳ�棬�Ҳ���һ����ɫ�� ��LAUNCH�� ��ť������ð�ť��
1. �ڼ���ҳ���ϣ���Ӧ�ÿ���������������һ�����п����ǿյġ�
	- ����ѡ��ϵ�ֵ��ע�����Ϣ��
		- **Peer Logs** ��λ������ı��ҵ��ڵ� 1 ���ڵ��У�Ȼ�������һ���������ļ���ͼ�ꡣ
			- ��Ӧ�û��һ���´��ڡ���ϲ���ҵ�����Ľڵ���־��
			- ���˴˾�̬��ͼ�����п���ҳ�涥�� **View Logs** ѡ��е�ʵʱ���ڵ���־��
		- ** ChainCode Logs ** ��λ������ı�ÿ�����붼��һ�У�����ʹ�������ϣ����ǣ������ϣ���ڲ�������ʱ���صĹ�ϣֵ���ҵ���������� ID��Ȼ��ѡ��ڵ㡣��󵥻������ļ���ͼ�ꡣ
			- ��Ӧ�û��һ���´��ڡ���ϲ���ҵ�����Ľڵ��������־��
	- ���Ϊ ** APIs ** ����** Swagger Tab **��������鿴 API �����ĵ���
		- �������� Swagger API ҳ�����ˡ�

## ��ȫע��
���� REST �ӿڵ� `/chaincode' ��Ҫһ����ȫ�������� ID��
����ζ�������ӷ���ƾ���б��д���ע��� enrollID���������ſ��Խ��ܴ���� REST ���á�
- ������� **+ Network's Enroll IDs** չ��������� enrollID ���������б�
- �򿪼��±�������һ��ƾ�ݡ����������õ����ǡ�
- ͨ������չ�� ��Registrar�� API ����
- ͨ������չ�� ��POST / registrar�� ����
- �������ĵ��ı��ֶΡ���Ӧ���� JSON ��ʽ���������б��ϵ�һ�� enrollID �����롣���磺


![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/registrar.PNG)

![Register Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/register_response.PNG)


�����û���յ� ��Login successful�� ��Ӧ���뷵�ز�ȷ��������ȷ������ enrollID �����롣�������������� enrollID��������ں��������в��𡢵��úͲ�ѯ����ʱʹ�ô� ID��

## �������루Deploying the chaincode��
Ϊ��ͨ�� REST �ӿڲ������룬����Ҫ������洢��һ�������� git �ֿ��С�
������ڵ㷢�Ͳ�������ʱ���Ὣ�������ֿ� URL ���͸��ڵ㣬ͬʱ��������ʼ����������Ĳ�����

**���㲿������֮ǰ**��ȷ�������ڱ��ع�����
- ���ն˻�������ʾ��
- ������� `chaincode_start.go` ���ļ��У����룺

	```
	go build ./
	```
- ��Ӧ�ò��᷵�ش���/�ı�


- ͨ������չ�� `Chaincode` API ����
- ͨ������չ�� `POST /chaincode` ����
- ���� `DeploySpec` �ı��ֶΣ�ʹ�����ֶ�Ϊ�գ���������������ӣ����滻�������ֿ�·������ʹ������ `/registrar` ������ʹ�õ���ͬ enrollID��
- `��path��`��Ӧ������ `��https://github.com/johndoe/learn-chaincode/finished��`��������� fork ����Ŀ�е�һ��Ŀ¼��`chaincode_finished.go` �ļ�λ�ڸ�Ŀ¼�С�

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

��ӦӦ��������ʾ��

![Deploy Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/deploy_response.PNG)

�������Ӧ���������������� ID��ID��һ�� 128 ���ַ�����ĸ���ֹ�ϣ��Ҳ����ļ��±��ϸ��ƴ� ID��������Ӧ����һ�� enrollID ƾ�ݺͱ�ʶ�������ļ��ܹ�ϣ��
�������ڵ��û��ѯ�������������뽫��ʹ�õ����ǡ�

## ��ѯ��Query��
�������������ǲ�ѯ����� `hello_world` ����ֵ��֮ǰ����ʹ���� `Init` ����Ϊ�������˳�ʼֵ��
- ͨ������չ�� `Chaincode` API ����
- ͨ������չ�� `POST /chaincode` ����
���� `QuerySpec` �ı��ֶΣ�ʹ�����ֶ�Ϊ�գ������������ʾ�������滻����������ƣ�����ʱ��ȡ�Ĺ�ϣ ID������ʹ������ `/registrar` ������ʹ�õ���ͬ enrollID��

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

ϣ���㿴���� `hello_world` ��ֵ�� ��hi there����
��ֵ������֮ǰ���ò������õġ�

## ���ã�Invoke��
��������ʹ�� `invoke` ������� `write` ������
�� `hello_world` ��ֵ����Ϊ ��go away����
- ͨ������չ�� `Chaincode` API ���֡�
- ͨ������չ�� `POST /chaincode` ���֡�
- ���� `InvokeSpec` �ı��ֶΣ�ʹ�����ֶ�Ϊ�գ��� ���������ʾ�������滻����������ƣ�����ʱ��ȡ�Ĺ�ϣ ID������ʹ������ `/registrar` ������ʹ�õ���ͬ enrollID��

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

����Ҫ�������Ƿ�ס��ֻ��������������Ĳ�ѯ������

![Query2 Example](https://raw.githubusercontent.com/IBM-Blockchain/learn-chaincode/master/imgs/query2_response.PNG)


����Ǳ�д������������Ҫ��ȫ�����ݡ�
