package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	h "github.com/mangirdaz/dockertree/hellpers"
	"strings"
)

func (l *LibKVBackend) AddDependency(config h.Config) {

	log.Infof("Add Dependency for image [%s] to image [%s]", config.Name, config.BaseImage)

	config.Name, _ = h.ParseImageName(config.Name)
	config.BaseImage, _ = h.ParseImageName(config.BaseImage)
	log.Debugf("Parsed image name [%s] and base image [%s]", config.Name, config.BaseImage)

	//check if base image project exist, if not create empty key
	baseNamespace := strings.Split(config.BaseImage, "/")[0]
	baseName := strings.Split(config.BaseImage, "/")[1]
	var baseTag string
	if strings.Contains(baseName, ":") {
		temp := strings.Split(baseName, ":")
		baseName = temp[0]
		baseTag = temp[1]
	} else {
		baseTag = "latest"
	}
	config.BaseImage = fmt.Sprintf("%s/%s:%s", baseNamespace, baseName, baseTag)
	config.Name = fmt.Sprintf("%s:%s", config.Name, config.Tag)

	log.Infof("Base image namespace [%s], name [%s], tag [%s]", baseNamespace, baseName, baseTag)

	//pathLastconfig := fmt.Sprintf("images/%s/%s/tags/%s/lastconfig", baseNamespace, baseName, baseTag)
	pathDependencies := fmt.Sprintf("images/%s/%s/tags/%s/dependencies", baseNamespace, baseName, baseTag)
	_, err := l.Get(pathDependencies)
	if err != nil {
		log.Debugf("Path %s does not exist, base image was not registred with us yet", pathDependencies)
		//if image is not registred we just create dependecies section and add this new image
		var dep []string

		dep = append(dep, config.Name)

		json, err := json.Marshal(dep)
		h.Check(err)
		log.Debugf("Payload for dependencie %s", string(json))
		l.storeKey(pathDependencies, json)
	} else {
		log.Debugf("Base Image under path %s exist", pathDependencies)
		pair, err := l.Get(pathDependencies)
		h.Check(err)

		//form array from data
		var dep []string
		json.Unmarshal([]byte(pair.Value), &dep)

		log.Debugf("Dependencies [%s] [%s]", pair.Key, pair.Value)
		dep = append(dep, config.Name)
		dep = h.RemoveDuplicates(dep)
		log.Debugf("New dep Array [%s]", dep)

		//store new dep array
		json, _ := json.Marshal(dep)
		l.storeKey(pathDependencies, json)
	}

}
