# How to test writing

```
$ go run main.go

$ curl -i -X POST -H "Content-Type: application/json" http://localhost:8080/collect -d @sample_data.json;
```