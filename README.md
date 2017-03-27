### 使い方

```
CLUSTERS=127.0.0.1,127.0.0.2 FILE=/path/to/csv KEYSPACE=cassandra_keyspace TABLE=cassandra_table ./src/cassandrainserter/main
```

+ CLUSTERS: ips of cassandra (, separator)
+ FILE: file path of csv data
+ KEYSPACE: target key space of cassandra
+ TABLE: target table of cassandra
+ USER: username of cassandra password authentication
+ PASSWORD: password of cassandra password authentication
+ N: Number of data divisions.
+ N_CON: Number of connection of cassandra

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