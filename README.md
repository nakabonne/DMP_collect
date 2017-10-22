# How to test writing

```
// ローカルの設定
$ DEV=1
$ export DEV
// 起動
$ go run main.go

$ curl -i -X POST -H "Content-Type: application/json" http://localhost:8080/collect -d @sample_data.json;
```
