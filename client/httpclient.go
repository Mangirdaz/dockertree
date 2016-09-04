package client

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	h "github.com/mangirdaz/dockertree/hellpers"
)

func (docker Docker) dorequest(method string, url string, payload []byte, retry int) (h.HTTPResponse, error) {
	log.Debugf("Do Request trigger method %s with retry %d", method, retry)

	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}

	counter := 1
	log.Debug("Start Requests Retry Cycle")
retry:
	resp, err := client.Do(req)

	if err != nil {
		h.Check(err)
		if counter <= retry {
			log.Debugf("Error detected in payload. Retry %d", counter)
			counter++

			time.Sleep(5000 * time.Millisecond)
			goto retry
		} else {
			h.DeleteTempFolder(docker.Config.TempDir)
			log.Fatalf("Error while reaching dockertree server API %s", url)
		}
	}
	defer resp.Body.Close()

	log.Debugf("response Status: [%s]", resp.Status)
	log.Debugf("response Headers: [%s]", resp.Header)
	body, err := ioutil.ReadAll(resp.Body)
	h.Check(err)
	log.Debugf("response Body: [%s]", string(body))

	//parse json response
	var res h.HTTPResponse
	log.Debug("Start unmarshar iteration data")
	err = json.Unmarshal(body, &res)
	h.Check(err)

	return res, err

}
