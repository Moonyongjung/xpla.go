# Slashing module
## Usage
### (Tx) Unjail validator
```go
txbytes, err := xplac.Unjail().CreateAndSignTx()
res, err := xplac.Broadcast(txbytes)
```

### (Query) slashing params
```go
res, err := xplac.SlashingParams().Query()
```

### (Query) slashing signing infos
```go
// Query a validator's signing information
signingInfoMsg := types.SigningInfoMsg{
    ConsPubKey: `{"@type": "/cosmos.crypto.ed25519.PubKey","key": "6RBPm24ckoWhRt8mArcSCnEKvt0FMGvcaMwchfZ3ue8="}`,
}
res, err := xplac.SigningInfos(signingInfoMsg).Query()

// Query signing information of all validators
res, err := xplac.SigningInfos().Query()
```
