package main

import (
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/option"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

const (
	project       = "ca-intern-201710-team02"
	instance      = "teamb-bigtable1"
	pathToKeyFile = "ca-intern-201710-team02-4d5815ebcb43.json"
	family        = "Log"
)

var (
	ctx = context.Background()
)

type Info struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	DeviceID  string  `json:"device_id"`
	SysName   string  `json:"sysname"`
	SysVer    string  `json:"sysver"`
}

func authenticate() (*bigtable.Client, error) {
	jsonKey, err := ioutil.ReadFile(pathToKeyFile)
	if err != nil {
		return nil, err
	}
	config, err := google.JWTConfigFromJSON(jsonKey, bigtable.Scope)
	if err != nil {
		return nil, err
	}
	client, err := bigtable.NewClient(ctx, project, instance, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func isDevelop() bool {
	return os.Getenv("DEV") == "1"
}

func openBigtable(tableName string) (table *bigtable.Table, err error) {
	var client *bigtable.Client
	if isDevelop() {
		client, err = authenticate()
	} else {
		client, err = bigtable.NewClient(ctx, project, instance)
	}
	if err != nil {
		log.Fatalln("エラー", err)
	}
	table = client.Open(tableName)
	return table, err
}

func write(table *bigtable.Table) (err error) {
	mut := bigtable.NewMutation()
	mut.Set(family, "lat", bigtable.Now(), []byte("1"))
	mut.Set(family, "lon", bigtable.Now(), []byte("2"))
	err = table.Apply(ctx, "2017102100000000#IDFA2", mut)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func main() {
	table, err := openBigtable("latlon-table")
	if err != nil {
		log.Fatalln(err)
	}

	err = write(table)
	if err != nil {
		log.Fatal(err)
	}
}
