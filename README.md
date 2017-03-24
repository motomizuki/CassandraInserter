### 使い方

```
CLUSTERS=127.0.0.1,127.0.0.2 FILE=/path/to/csv KEYSPACE=cassandra_keyspace TABLE=cassandra_table ./src/cassandrainserter/main
```

+ CLUSTERS: カサンドラのip　","区切りで複数指定可能
+ FILE: インサートするｃｓｖファイルのパス
+ KEYSPACE: カサンドラのkey_spaceの名前
+ TABLE: ｃｓｖをインサートするテーブル名



### ビルド方法
#### goのインストール
https://golang.org/doc/install 参照

#### direnvのインストール
https://github.com/direnv/direnv 参照

#### glideのインストール
https://github.com/Masterminds/glide　参照

#### ビルド
```
$ direnv allow
$ cd src
$ glide install
$ cd cassandrainserter
$ go build main.go
```