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

func main() {
	if os.Getenv("DEV") == "1" {
		client := authenticate()
		table := client.Open("latlon-table")
		fmt.Println(table)
	} else {
		adminClient, err := bigtable.NewAdminClient(ctx, project, instance)
		if err != nil {
			log.Fatalln("エラー", err)
		}
		tables, err := adminClient.Tables(ctx)
		if err != nil {
			log.Fatalln("エラー", err)
		}
		log.Println(tables)
	}
}
