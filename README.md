# dvf-blockchain-go

Original Python Code
https://github.com/dvf/blockchain

Learn Blockchains by Building One

- Original
  https://hackernoon.com/learn-blockchains-by-building-one-117428612f46

- 日本語翻訳
  https://qiita.com/hidehiro98/items/841ece65d896aeaa8a2a

```
go run main.go server.go blockchain.go

curl http://localhost:5000/mine
curl -X POST http://localhost:5000/transactions/new -d '{"sender":"d4ee26eee15148ee92c6cd394edd974e","recipient":"someone-other-address","amount":5}'
curl http://localhost:5000/chain
```

```
go run main.go server.go blockchain.go -addr=:5001

curl -X POST -H "Content-Type: application/json" -d '{
    "nodes": ["http://localhost:5001"]
}' "http://localhost:5000/nodes/register"
curl "http://localhost:5001/mine"
curl "http://localhost:5000/nodes/resolve"
```
