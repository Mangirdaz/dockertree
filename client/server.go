package client

import (
	"encoding/json"
	"fmt"

	log "github.com/Sirupsen/logrus"
	h "github.com/mangirdaz/dockertree/hellpers"
)

type ResponseUserAgent struct {
	Useragent string `json:"user-agent"`
}

func (docker Docker) PostConfigToServer() (iteration string) {

	url := docker.Config.Server + docker.Config.APIPrefix + "/images"
	log.Info("URL:>", url)
	output, _ := json.Marshal(docker.Config)
	res, err := docker.dorequest("POST", url, output, 10)
	h.Check(err)
	result := h.Result{}
	json.Unmarshal([]byte(res.Message), &result)
	return result.Value
}

func (docker Docker) PostIterationToServer() {

	docker.Config.Name, _ = h.ParseImageName(docker.Config.Name)
	url := fmt.Sprintf("%s%s/images/%s/tags/%s/%s", docker.Config.Server, docker.Config.APIPrefix, docker.Config.Name, docker.Config.Tag, "iteration")

	output, _ := json.Marshal(docker.Config.Iteration)
	_, err := docker.dorequest("POST", url, output, 5)
	h.Check(err)

}
