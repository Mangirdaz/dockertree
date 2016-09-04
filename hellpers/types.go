package hellpers

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

//Config strcutre to read from configuration file and is used to pass data betwean server and client
type Config struct {
	Tag          string       `yaml:"tag"` // Supporting both JSON and YAML.
	Name         string       `yaml:"name"`
	Dependencies []string     `yaml:"dependencies"`
	Latest       bool         `yaml:"latest"`
	Force        bool         `yaml:"force"`
	Trigger      string       `yaml:"trigger"`
	Alias        []string     `yaml:"alias"`
	Mode         string       `yaml:"mode"`
	Server       string       `yaml:"server"`
	Docker       DockerConfig `yaml:"docker"`
	ModuleDir    string       `yaml:"moduledir"`
	TempDir      string       `yaml:"tempdir"`
	BuildTag     string       `yaml:"buildtag"`
	BaseImage    string       `yaml:"baseimage"`
	APIPrefix    string       `yaml:"apiprefix"`
	Iteration    string       `yaml:"iteration"`
	Version      string
}

//DockerConfig for setting docker client - optional
type DockerConfig struct {
	DefaultHeaders     []DefaultHeader `yaml:"defaultHeaders"`
	Socket             string          `yaml:"socket"`
	APIVersion         string          `yaml:"apiVersion"`
	FinalDefaultHeader map[string]string
}

type DefaultHeader struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Result struct {
	Key   string
	Value string
}

type HTTPResponse struct {
	Code    int    `json:"Code"`
	Message string `json:"Message"`
}
type Results struct {
	Result []Result
}

//DataStore used to define key value storage structure and is used to import/export data
type DataStore struct {
	Images struct {
		Image []struct {
			Name string `json:"name"`
			Tags []struct {
				Tag          string   `json:"tag"`
				Dependencies []string `json:"dependencies"`
				Iteration    int      `json:"iteration"`
				Metadata     struct {
					Lastcallip        string `json:"lastcallip"`
					Lastcallstarttime string `json:"lastcallstarttime"`
					Lastcallstoptime  string `json:"lastcallstoptime"`
				} `json:"metadata"`
				Lastonfig Config `json:"lastonfig"`
			} `json:"tag"`
		} `json:"image"`
	} `json:"images"`
}

type ImageList struct {
	Results []string `json:"images"`
}

type NamespaceList struct {
	Results []string `json:"namespaces"`
}

type TagList struct {
	Results []string `json:"tags"`
}

type TagMetadaList struct {
	Results []string `json:"metadata"`
}

type MetadaResults struct {
	Results struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	} `json:"Metadata"`
}

//ServerConfig defines server configuration
type ServerConfig struct {
	KVStorageIp   string
	KVStoragePort string
	ServerIP      string
	ServerPort    string
	DeVDummyData  bool
	Namespace     string
	APIPrefix     string
	KVStoreSubdir string
	Version       string
}

type ErrorMessage struct {
	Code    int
	Message string
}

//SetDefault for Client config
func (c *Config) SetDefault() {
	log.Debug("Set defaults")

	if len(c.Docker.DefaultHeaders) == 0 {
		log.Debug("Default Header Map is empty, add engine-api-cli version")
		c.Docker.DefaultHeaders = append(c.Docker.DefaultHeaders, DefaultHeader{Name: "User-Agent", Value: "engine-api-cli-1.0"})
	}
	if c.Docker.Socket == "" {
		log.Debug("Set Docker Socket Default")
		c.Docker.Socket = "unix:///var/run/docker.sock"
	}
	if c.Docker.APIVersion == "" {
		log.Debug("Set Docker APIVersion Default")
		c.Docker.APIVersion = "v1.22"
	}
	if len(c.BuildTag) == 0 {
		log.Debug("Set build flag to [build]")
		c.BuildTag = "build"
	}

	if len(c.APIPrefix) == 0 {
		log.Debugf("Set API prefix with [/api/v1]")
		c.APIPrefix = "/api/v1"
	}
}

//SetDefault for LibKV storage
func (c *ServerConfig) SetDefault() {

	log.Debug("Set LibKV Defaults")

	c.Version = "v.0.1"

	if len(c.KVStoreSubdir) == 0 {
		log.Debug("Set KVStoreSubdir to [dockertree/]")
		c.KVStoreSubdir = "dockertree/"
	}

	if os.Getenv("KEYVAL_STORAGE_IP") == "" {
		c.KVStorageIp = "0.0.0.0"
	} else {
		c.KVStorageIp = os.Getenv("KEYVAL_STORAGE_IP")
	}

	if os.Getenv("KEYVAL_STORAGE_PORT") == "" {
		c.KVStoragePort = "8500"
	} else {
		c.KVStoragePort = os.Getenv("KEYVAL_STORAGE_PORT")
	}

	if os.Getenv("DEV_DUMMYDATA") == "" {
		c.DeVDummyData = false
	} else {
		c.DeVDummyData = true
	}

	if os.Getenv("SERVER_IP") == "" {
		c.ServerIP = "0.0.0.0"
	} else {
		c.ServerIP = os.Getenv("SERVER_IP")
	}

	if os.Getenv("SERVER_PORT") == "" {
		c.ServerPort = "8080"
	} else {
		c.ServerPort = os.Getenv("SERVER_PORT")
	}

}
