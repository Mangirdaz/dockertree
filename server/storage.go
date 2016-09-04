package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libkv"
	"github.com/docker/libkv/store"
	consul "github.com/docker/libkv/store/consul"
	h "github.com/mangirdaz/dockertree/hellpers"
	"strconv"
	"strings"
	"time"
)

// LibKVBackend - libkv container
type LibKVBackend struct {
	Store  store.Store
	Config h.ServerConfig
}

//NewLibKVBackend creates and sets key value storage backend
func NewLibKVBackend() (storage LibKVBackend) {
	log.Info("Create New LibKV backend")
	var s LibKVBackend
	s.Config.SetDefault()
	config := s.Config
	consul.Register()

	client := config.KVStorageIp + ":" + config.KVStoragePort

	log.Info("Init New store")
	kv, err := libkv.NewStore(
		store.CONSUL, // or "consul"
		[]string{client},
		&store.Config{
			ConnectionTimeout: 10 * time.Second,
		},
	)
	if err != nil {
		log.Fatal("Cannot create store consul")

	}

	err = kv.Put(s.Config.KVStoreSubdir, []byte(s.Config.Version), nil)
	for err != nil {
		err = kv.Put(s.Config.KVStoreSubdir, []byte(s.Config.Version), nil)
		log.Debugf("Connection To storage Failed... [%s]", err)
		time.Sleep(5 * time.Second)
	}

	s.Store = kv
	log.Info("Store init")

	if config.DeVDummyData {
		log.Info("Dev Mode detected, loading dummy data")
		LoadDummyData(s)
	}
	return s

}

//GetSubdir gives all subdirs in the kv storage. Used for storing data in Keys
func (l *LibKVBackend) GetSubdir(path string) (results []string, err error) {
	log.Infof("Try get subdirs of [%s]", path)
	path = l.Config.KVStoreSubdir + path
	pair, err := l.Store.List(path)
	var subdirs []string
	if err != nil {
		log.Error("Directory not found: ", err)
	} else {
		//if we split by work dept will be 2 always (0 and 1)
		dept := 1
		//found where our subdir is located relativly to path

		//append only subdir to separaet subdirs array
		for _, v := range pair {
			temp := strings.Split(v.Key, path)
			log.Debug(temp)
			subdirs = append(subdirs, strings.Split(temp[dept], "/")[1])
		}
		subdirs = h.RemoveDuplicates(subdirs)
		return subdirs, nil
	}
	return subdirs, err
}

func (l *LibKVBackend) Get(namespace string) (results *store.KVPair, err error) {

	pair, err := l.Store.Get(l.Config.KVStoreSubdir + namespace)

	if err != nil {
		log.Error("Directory not found: ", err)
	} else {
		pair.Key = strings.Replace(pair.Key, l.Config.KVStoreSubdir, "", 1)
		return pair, nil
	}
	return pair, err
}

func (l *LibKVBackend) ParseAndLoadImage(config h.Config) (results *store.KVPair, err error) {

	config.Name, err = h.ParseImageName(config.Name)
	h.Check(err)
	pathIteration := fmt.Sprintf("images/%s/tags/%s/%s", config.Name, config.Tag, "iteration")
	if err == nil {
		b, err := json.Marshal(config)
		h.Check(err)
		//image name is already with namespace from ParseImageName
		l.storeKey(fmt.Sprintf("images/%s/tags/%s/%s", config.Name, config.Tag, "lastconfig"), b)
		log.Debug("Get [%s]", pathIteration)
		_, err = l.Get(pathIteration)
		if err == nil {
			log.Debug("Existing image, getting iteration")
			value, err := l.Get(pathIteration)
			return value, err
		} else {
			log.Debug("New image, set iteration to 0")
			l.storeKey(pathIteration, []byte(strconv.Itoa(0)))
			pair, err := l.Get(pathIteration)
			return pair, err
		}
	}
	return results, err
}

func (l *LibKVBackend) ParseAndLoadDummyData(datastore h.DataStore) (err error) {
	for _, image := range datastore.Images.Image {
		for _, tag := range image.Tags {
			//convert tag element to json to same as a value
			//TODO: add check if data is empty, not override.
			b, err := json.Marshal(tag.Dependencies)
			h.Check(err)

			image.Name, err = h.ParseImageName(image.Name)
			h.Check(err)
			if err == nil {
				l.storeKey(fmt.Sprintf("images/%s/tags/%s/%s", image.Name, tag.Tag, "dependencies"), b)
				l.storeKey(fmt.Sprintf("images/%s/tags/%s/%s", image.Name, tag.Tag, "iteration"), []byte(strconv.Itoa(tag.Iteration)))
				b, err = json.Marshal(tag.Metadata)
				h.Check(err)
				l.storeKey(fmt.Sprintf("images/%s/tags/%s/%s", image.Name, tag.Tag, "metadata"), b)
				b, err = json.Marshal(tag.Lastonfig)
				h.Check(err)
				l.storeKey(fmt.Sprintf("images/%s/tags/%s/%s", image.Name, tag.Tag, "lastconfig"), b)
			}
		}
	}
	return err
}

func (l *LibKVBackend) storeKey(key string, value []byte) {
	if string(value) != "null" {
		log.Debugf("Key [%s] with value [%s] not empty - store", key, string(value))
		err := l.Store.Put(l.Config.KVStoreSubdir+key, value, nil)
		h.Check(err)
	} else {
		log.Debugf("Key [%s] empty - ignore", key)
	}
}
