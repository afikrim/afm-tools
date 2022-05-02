package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/afikrim/afm-tools/config"
	"github.com/afikrim/afm-tools/lib"
)

func main() {
	httpClient := http.DefaultClient
	options := parseOptions()
	cfg, err := config.LoadConfig(*options["initPostman"].(*bool))
	if err != nil {
		panic(err)
	}

	if !*options["isSyncPostman"].(*bool) {
		return
	}
	if *options["isSyncPostman"].(*bool) {
		syncPostman := lib.NewSyncPostman(cfg, httpClient)
		err := syncPostman.SyncPostmanCollection(options["collectionName"].(string))
		if err != nil {
			panic(err)
		}
	}
}

func parseOptions() map[string]interface{} {
	var options = make(map[string]interface{})

	flag.Usage = func() {
		fmt.Print("Usage: afm-tools [options] <...args>\nOptions:\n")
		flag.PrintDefaults()
	}

	options["initPostman"] = flag.Bool("init-postman", false, "Initialize postman credential")
	options["isSyncPostman"] = flag.Bool("sync-postman", false, "Sync postman collection")
	if *options["isSyncPostman"].(*bool) {
		options["isListPostmanCollections"] = flag.Bool("list", false, "List all available postman collections")
	}

	flag.Parse()

	if *options["isSyncPostman"].(*bool) {
		options["collectionName"] = flag.Arg(0)
	}

	return options
}
