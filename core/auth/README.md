# Auth module
## Usage
### (Query) auth params
```go
response, err := xplac.AuthParams().Query()
```

### (Query) account address
```go
queryAccAddressMsg := types.QueryAccAddressMsg{
    Address: "xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",
}
response, err := xplac.AccAddress(queryAccAddressMsg).Query()
```

### (Query) accounts
```go
response, err := xplac.Accounts().Query()
```

### (Query) Txs by events
```go
queryTxsByEventsMsg := types.QueryTxsByEventsMsg{
    Events: "transfer.recipient=xpla19w2r47nczglwlpfynqe5769cwkwq5fvmzu5pu7",
}
response, err := xplac.TxsByEvents(queryTxsByEventsMsg).Query()
```

### (Query) tx
```go
queryTxMsg := types.QueryTxMsg{
    Value: "B6BBBB649F19E8970EF274C0083FE945FD38AD8C524D68BB3FE3A20D72DF03C4",
}
response, err := xplac.Tx(queryTxMsg).Query()
```