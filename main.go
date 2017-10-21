package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
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
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	DeviceID  string `json:"device_id"`
	SysName   string `json:"sysname"`
	SysVer    string `json:"sysver"`
	Timestamp string `json:"timestamp"`
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
		log.Println("エラー", err)
	}
	table = client.Open(tableName)
	return table, err
}

func write(table *bigtable.Table, rowKey string, lat string, lon string) (err error) {
	mut := bigtable.NewMutation()
	mut.Set(family, "lat", bigtable.Now(), []byte(lat))
	mut.Set(family, "lon", bigtable.Now(), []byte(lon))
	err = table.Apply(ctx, rowKey, mut)
	if err != nil {
		log.Println(err)
	}
	return
}

func decode(r io.ReadCloser) (*Info, error) {
	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	info := new(Info)
	if err := json.Unmarshal(bytes, &info); err != nil {
		return nil, err
	}
	return info, nil
}

func collect(w http.ResponseWriter, r *http.Request) {
	info, err := decode(r.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(info)
	rowKey := info.Timestamp + "#" + info.DeviceID
	table, err := openBigtable("latlon-table")
	if err != nil {
		log.Println(err)
	}

	err = write(table, rowKey, info.Latitude, info.Longitude)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	http.HandleFunc("/collect", collect)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatalln(err)
	}
}
