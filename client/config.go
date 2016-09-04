package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/go-yaml/yaml"
	h "github.com/mangirdaz/dockertree/hellpers"
	"io/ioutil"
	"os"
	"path/filepath"
)

func GetConfig() (c h.Config) {
	log.Info("Read Configuration file")

	var filename string
	if os.Getenv("DOCKERTREE_CONFIG_NAME") == "" {
		filename, _ = filepath.Abs("./.dockertree")
	} else {
		filename, _ = filepath.Abs(os.Getenv("DOCKERTREE_CONFIG_NAME"))
	}

	yamlFile, err := ioutil.ReadFile(filename)
	h.Check(err)

	err = yaml.Unmarshal(yamlFile, &c)
	h.Check(err)

	//check if DockerFile project is preset and if not if moduedir is set and file is present there
	if len(c.ModuleDir) != 0 {
		log.Debug("Module dir is set. Checking for Dockerfile")
		if string(c.ModuleDir[0]) == "/" {
			log.Debug("Absolute Path detected")
			if h.FileExists(c.ModuleDir + "/Dockerfile") {
				log.Infof("Dockerfile found under [%s]", c.ModuleDir+"/Dockerfile")
				c.ModuleDir = filepath.Dir(c.ModuleDir)
			} else {
				log.Errorf("Dockerfile not found under [%s]", c.ModuleDir+"/Dockerfile")
			}
		} else {
			log.Debug("Relative Path detected")
			file, err := filepath.Abs(c.ModuleDir + "/Dockerfile")
			h.Check(err)
			if h.FileExists(file) {
				log.Infof("Dockerfile found under [%s]", file)
				c.ModuleDir = filepath.Dir(file)
			} else {
				log.Errorf("Dockerfile not found under [%s]", file)
				os.Exit(1)
			}
		}
	} else {
		log.Debug("Module dir not set, searching for Dockerfile in current folder")
		file, err := filepath.Abs("./Dockerfile")
		h.Check(err)
		if h.FileExists(file) {
			log.Infof("Dockerfile found under [%s]", file)
			c.ModuleDir = filepath.Dir(file)
		} else {
			log.Errorf("Dockerfile not found under [%s]", file)
			os.Exit(1)
		}
	}
	log.Debugf("Module dir set to [%s]", c.ModuleDir)
	c.SetDefault()

	//construct header value for docker client
	header := make(map[string]string)
	for _, value := range c.Docker.DefaultHeaders {
		header[value.Name] = value.Value
	}
	c.Docker.FinalDefaultHeader = header

	return c

}
