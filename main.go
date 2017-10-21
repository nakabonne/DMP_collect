package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"google.golang.org/api/option"

	"cloud.google.com/go/bigtable"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

const (
	project  = "ca-intern-201710-team02"
	instance = "teamb-bigtable1"
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

func authenticate() *bigtable.Client {
	pathToKeyFile := "ca-intern-201710-team02-4d5815ebcb43.json"
	jsonKey, err := ioutil.ReadFile(pathToKeyFile)
	if err != nil {
		fmt.Println(err)
	}
	config, err := google.JWTConfigFromJSON(jsonKey, bigtable.Scope)
	if err != nil {
		fmt.Println(err)
	}
	client, err := bigtable.NewClient(ctx, project, instance, option.WithTokenSource(config.TokenSource(ctx)))
	if err != nil {
		fmt.Println(err)
	}
	return client
}

func isDevelop() bool {
	return os.Getenv("DEV") == "1"
}

func main() {
	var client *bigtable.Client
	var err error
	if isDevelop() {
		client = authenticate()
	} else {
		client, err = bigtable.NewClient(ctx, project, instance)
		if err != nil {
			log.Fatalln("エラー", err)
		}
	}
	table := client.Open("latlon-table")
	fmt.Println(table)

	mut := bigtable.NewMutation()
	mut.Set("links", "maps.google.com", bigtable.Now(), []byte("1"))
	mut.Set("links", "golang.org", bigtable.Now(), []byte("1"))
	err = table.Apply(ctx, "com.google.cloud", mut)
}
