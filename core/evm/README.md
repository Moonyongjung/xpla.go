# EVM module
## Usage
### (Tx) Send coin
```go
sendCoinMsg := types.SendCoinMsg{
    Amount: "10000",
    FromAddress: "0x6577385b5d959644ae31263208a88E921273C774",
    ToAddress: "0xF9AC4736D8034F2CB3BFF22A977CD8759934F090",
}

txbytes, err := xplac.EvmSendCoin(sendCoinMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Deploy solidity contract
```go
// Select input type of ABI and bytecode.
// It is also possible to input the entire abi and bytecode as a string type, but you can also enter a file path.
// ABI type must be json and bytecode file must be compiled file on Remix IDE.
deploySolContractMsg := types.DeploySolContractMsg{
    //ABI: `{ ABI json string type }`
    ABIJsonFilePath: "./abi.json",
    // Bytecode: "60806040523480156100......",
    BytecodeJsonFilePath: "./bytecode.json",
}

txbytes, err := xplac.DeploySolidityContract(deploySolContractMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Tx) Invoke(execute) solidity contract
```go
// When invoked, the arguments to be entered into the solidity contract are listed as []interface{}.
var args []interface{}
args = append(a, big.NewInt(2))

// Need contract address and invoke function name.
// Also, same as deployment, need ABI and bytecode.
invokeSolContractMsg := types.InvokeSolContractMsg{
    ContractAddress: "0xBe0AE9A424771C0D68D942A04994a97f928b0821",
    ContractFuncCallName: "store",
    Args: args,
    ABIJsonFilePath: "./abi.json",
    BytecodeJsonFilePath: "./bytecode.json",
}

txbytes, err := xplac.InvokeSolidityContract(invokeSolContractMsg).CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) Call solidity contract
```go
callSolContractMsg := types.CallSolContractMsg{
    ContractAddress: "0x80E123317190cAf36292A04776b0De020136526F",
    ContractFuncCallName: "retrieve",
    // Args: nil, // input params if needed to call
    ABIJsonFilePath: "./abi.json",
    BytecodeJsonFilePath: "./bytecode.json",
}

res, err := xplac.CallSolidityContract(callSolContractMsg).Query()
```

### (Query) Get transaction by hash
```go
getTransactionByHashMsg := types.GetTransactionByHashMsg {
    TxHash: "556c60576f9af3e4ae7d7fb28f8376e96803c4d9ff02eda6aacb86925f170d09",
}

res, err := xplac.GetTransactionByHash(getTransactionByHashMsg).Query()
```

### (Query) Get Block by hash or height
```go
// Query block by hash
getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg {
    BlockHash: "0xe083b9b3a8b5df69394f55d34cfdfa46e70743a812d7433aba0adf3b7fcecd21",
}

// Query block by height
getBlockByHashHeightMsg := types.GetBlockByHashHeightMsg {
    BlockHeight: "8",
}

res, err := xplac.GetBlockByHashOrHeight(getBlockByHashHeightMsg).Query()
```

### (Query) Account info
```go
// Query account info of user account or contract
// Response of query includes account address(Hex and Bech32), balances and etc. 
accountInfoMsg := types.AccountInfoMsg{
    Account: "0xCa8582862B82867C4Bb9E926682dD75820dE6013",
}

res, err := xplac.AccountInfo(accountInfoMsg).Query()
```

### (Query) Suggest gas price
```go
res, err := xplac.SuggestGasPrice().Query()
```

### (Query) ETH chain ID
```go
res, err := xplac.EthChainID().Query()
```

### (Query) Latest block number
```go
res, err := xplac.EthBlockNumber().Query()
```