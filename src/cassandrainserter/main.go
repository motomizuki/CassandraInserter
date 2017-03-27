package main

import (
	"os"
	"github.com/caarlos0/env"
	"github.com/gocql/gocql"
	"log"
	"encoding/csv"
	"strconv"
	"fmt"
	"sync"
	"time"
)

type config struct {
	CLUSTERS []string `env:"CLUSTERS" envSeparator:","`
	FILE string `env:"FILE" envDefault:"sample.csv"`
	KEYSPACE string `env:"KEYSPACE" envDefault:"example"`
	TABLE string `env:"TABLE" envDefault:"RS_SCORE_BY_ITEM"`
	USER string `env:"USER"`
	PASSWORD string `env:"PASSWORD"`
	N int `env:"N" envDefault:"10"`
	N_CON int `env:"N_CON" envDefault:"2"`
	RETRY int `env:"RETRY" envDefault:"5"`
}

const csql_tmpl  = `INSERT INTO %s.%s (user_id, item_id, score) values (?, ?, ?)`
func main() {
	var cfg config
	var err error

	err = env.Parse(&cfg)
	if err != nil {
		panic(err)
	}
	log.Println(cfg)
	clusters := cfg.CLUSTERS
	cluster := gocql.NewCluster(clusters...)
	cluster.NumConns = cfg.N_CON
	cluster.ConnectTimeout = 5 * time.Second
	cluster.Keyspace = cfg.KEYSPACE
	if (cfg.USER != "" && cfg.PASSWORD != "") {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: cfg.USER,
			Password: cfg.PASSWORD,
		}
	}
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("create session error %s", err)
	}

	csv_file, err := os.Open(cfg.FILE)
	if err != nil {
		log.Fatalf("load error %s", err)
	}
	defer csv_file.Close()
	csql := fmt.Sprintf(csql_tmpl, cfg.KEYSPACE, cfg.TABLE)
	reader := csv.NewReader(csv_file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)

	}
	l := len(records)
	bs := l / cfg.N
	wg := sync.WaitGroup{}
	wg.Add(cfg.N + 1)
	for i := 0; i < cfg.N + 1; i++ {
		s := i * bs
		e := s + bs
		if e> l { e = l}
		go func (record_chunk [][]string) {
			for _, record := range record_chunk {
				f, _ := strconv.ParseFloat(record[2], 32)
				for t := 0; t < cfg.RETRY ; t++ {
					err := session.Query(csql, record[0], record[1], float32(f)).Exec()
					if err == nil {
						break
					}
					log.Printf("exec query error, %s", err)
					time.Sleep(1 * time.Second)
				}
			}
			wg.Done()
		}(records[s:e])
	}
	wg.Wait()
	defer session.Close()
}