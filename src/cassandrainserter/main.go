package main

import (
	"os"
	"github.com/caarlos0/env"
	"github.com/gocql/gocql"
	"log"
	"encoding/csv"
	"strconv"
	"fmt"
)

type config struct {
	CLUSTERS []string `env:"CLUSTERS" envSeparator:","`
	FILE string `env:"FILE" envDefault:"sample.csv"`
	KEYSPACE string `env:"KEYSPACE" envDefault:"example"`
	TABLE string `env:"TABLE" envDefault:"RS_SCORE_BY_ITEM"`
	USER string `env:"USER"`
	PASSWORD string `env:"PASSWORD"`

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

	cluster.Keyspace = cfg.KEYSPACE
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.USER,
		Password: cfg.PASSWORD,
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
	for _, record := range records {
		f, _ := strconv.ParseFloat(record[2], 32)
		if err := session.Query(csql, record[0], record[1], float32(f)).Exec(); err != nil {
			log.Fatal(err)
		}
	}
	defer session.Close()
}