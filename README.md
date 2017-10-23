# team_b_collection

## Architecture

![result](https://github.com/ryonakao/DMP_collect/blob/master/media/ArchitectureB.png)

上記アーキテクチャの"CollectAPI"の部分。

このAPIは以下のタスクを行います

- 位置情報ログを受け取る
- CloudBigtableへの書き込み
- Stack driver loggingへログを渡す

# How to test writing

```
// ローカルの設定
$ DEV=1
$ export DEV
// 起動
$ go run main.go

$ curl -i -X POST -H "Content-Type: application/json" http://localhost:8080/collect -d @sample_data.json;
```

# Digression

当初は以下のアーキテクチャの予定だったが時間なかった

![result](https://github.com/ryonakao/DMP_collect/blob/master/media/architectureA.png)
