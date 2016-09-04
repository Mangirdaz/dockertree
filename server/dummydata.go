package server

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	h "github.com/mangirdaz/dockertree/hellpers"
	"io"
	"io/ioutil"
	"os"
)

func LoadDummyData(storage LibKVBackend) {
	log.Info("Load dummy data")
	var configFile io.Reader
	if h.FileExists("dummydata.json") {
		configFile, _ = os.Open("dummydata.json")
	} else if h.FileExists("server/dummydata.json") {
		configFile, _ = os.Open("server/dummydata.json")
	} else {
		log.Error("No dummydata.json file file found in current dir or server/ ")
	}

	var datastore h.DataStore
	hah, err := ioutil.ReadAll(configFile)
	h.Check(err)
	log.Debug("Start unmarshar dummy data")
	error := json.Unmarshal(hah, &datastore)
	h.Check(error)

	log.Debug("Start Loading Dummy Data")

	storage.ParseAndLoadDummyData(datastore)

}
