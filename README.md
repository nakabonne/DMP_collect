# DMP_collect

## Architecture

![result](https://github.com/ryonakao/DMP_collect/blob/master/media/ArchitectureB.png)

上記アーキテクチャの"CollectAPI"の部分。
AnswerAPIは→https://github.com/ryonakao/DMP_answer

## Tasks

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

当初は以下のようにジョブキューを噛ませる予定だったが、2000QPSなら使わなくても耐えることができた。

![result](https://github.com/ryonakao/DMP_collect/blob/master/media/architectureA.png)
