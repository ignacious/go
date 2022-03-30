
# Version
I used `go1.18 darwin/amd64`

# Commands

## In project root:
`go get .`

`go run .`

# API

## [GET] assets/:assetId

### Request
Pass in a hashed token as a url param
`http://localhost:8080/assets/0xdac17f958d2ee523a2206206994597c13d831ec7`

### Response
```
{
    "id": "0xdac17f958d2ee523a2206206994597c13d831ec7",
    "volume": 308191887.8126221,
    "pools": [
        "0x0cbe2f86e2fd90040ebb557b99f83400bf8f3717",
        "0x11b815efb8f581194ae79006d24e0d814b7697f6",
        "0x1dee9d7b7cfd8febf38982bc8ab715ec8c3050d1",
        "0x1e8f1568b598908785064809ebf5745004ce3962",
        "0x3416cf6c708da44db2624d63ea0aaef7113527c6",
        "0x4773e2c1c0b400a16dfec4ca6e305141859a5542",
        "0x4d1ad4a9e61bc0e5529d64f38199ccfca56f5a42",
        "0x4e68ccd3e89f51c3074ca5072bbac773960dfa36",
        "0x55ec9256077a311256b2daf81f70c0992d9fbd66",
        "0x56534741cd8b152df6d48adf7ac51f75169a83b2",
        ...
    ]
}
```

#### query string params
`http://localhost:8080/assets/0xdac17f958d2ee523a2206206994597c13d831ec7?poolCreatedAtStart=1624163307&poolCreatedAtEnd=1624163505`

**poolCreatedAtStart** : Sets a minimum `createdAtTimestamp` to filter. See `createdAtTimestamp_gte`.

**poolCreatedAtEnd** : Sets a maximum `createdAtTimestamp` to filter. See `createdAtTimestamp_lte`.

## [GET] blocks/:blockNumber

### Request
Pass in a block number as url param
http://localhost:8080/blocks/12533989

### Response
Lists all of the swaps and assets associated with the block number.
```
{
    "id": 12533989,
    "swaps": [
        "0x0000a8b55e0ea1bbd6323cdd5f0b993f486c1a936656b71f662aa082a21dc9eb#491",
        "0x0001b5fbc5172f64b12c5ebcf05db2ab378b4df1cb22b57056725c07996b8fdb#15990",
        "0x00025ab9cf525801fa7e17410d23bfcdf7d8fc64e063efacbb62878b32488f18#36451",
        ...
    "assets": [
        "0x2f109021afe75b949429fe30523ee7c0d5b27207",
        "0x50de6856358cc35f3a9a57eaaa34bd4cb707d2cd",
        "0x4a220e6096b25eadb88358cb44068a3248254675"
        ...
    ]
```