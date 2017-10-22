package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/bigtable"
	"cloud.google.com/go/logging"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

const (
	project       = "ca-intern-201710-team02"
	instance      = "teamb-bigtable1"
	pathToKeyFile = "ca-intern-201710-team02-4d5815ebcb43.json"
	family        = "Log"
	logName       = "collect-log"
)

var (
	ctx       = context.Background()
	table     *bigtable.Table
	logClient *logging.Client
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

func openBigtable(tableName string) (tbl *bigtable.Table, err error) {
	var client *bigtable.Client
	if isDevelop() {
		client, err = authenticate()
	} else {
		client, err = bigtable.NewClient(ctx, project, instance)
	}
	if err != nil {
		log.Fatal(err)
	}
	tbl = client.Open(tableName)
	return
}

func write(tbl *bigtable.Table, rowKey string, lat string, lon string) (err error) {
	mut := bigtable.NewMutation()
	mut.Set(family, "lat", bigtable.Now(), []byte(lat))
	mut.Set(family, "lon", bigtable.Now(), []byte(lon))
	err = tbl.Apply(ctx, rowKey, mut)
	if err != nil {
		log.Fatal(err)
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
	writeLog(*info)
	return info, nil
}

func collect(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	info, err := decode(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	rowKey := info.Timestamp + "#" + info.DeviceID
	if table == nil {
		table, err = openBigtable("latlon-table")
	}
	if err != nil {
		log.Fatal(err)
	}

	err = write(table, rowKey, info.Latitude, info.Longitude)
	if err != nil {
		log.Fatal(err)
	}

}

func writeLog(info Info) {
	logger := logClient.Logger(logName)
	// mapJSON := map[string]interface{}{
	// 	"lat":       info.Latitude,
	// 	"lon":       info.Longitude,
	// 	"timestamp": info.Timestamp,
	// 	"idfa":      info.DeviceID,
	// 	"sysname":   info.SysName,
	// 	"sysver":    info.SysVer,
	// }
	//logJSON, err := json.Marshal(info)
	// if err != nil {
	// 	log.Println(err)
	// }
	logger.Log(logging.Entry{Payload: info})
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("I'm Healthy"))
}

func init() {
	var err error
	table, err = openBigtable("latlon-table")
	if err != nil {
		log.Fatal(err)
	}

	logClient, err = logging.NewClient(ctx, project)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	http.HandleFunc("/collect", collect)
	http.HandleFunc("/hc", healthCheck)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
