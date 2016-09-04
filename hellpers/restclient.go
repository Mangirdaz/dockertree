package hellpers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/jmcvetta/napping"
)

type ResponseUserAgent struct {
	Useragent string `json:"user-agent"`
}

func ApiGet(url string) interface{} {

	// Start Session
	s := napping.Session{}
	log.Infof("Get URL: %s", url)

	res := ResponseUserAgent{}
	resp, err := s.Get(url, nil, &res, nil)
	if err != nil {
		log.Fatal(err)
	}
	//
	// Process response
	//

	log.Debugf("response Status: %s", resp.Status())
	log.Debugf("Header [%s]", resp.HttpResponse().Header)
	log.Debugf("RawText [%s]", resp.RawText())

	return resp.RawText
}
