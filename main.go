package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	c "github.com/mangirdaz/dockertree/client"
	h "github.com/mangirdaz/dockertree/hellpers"
	s "github.com/mangirdaz/dockertree/server"
	"os"
)

func init() {
	log.SetOutput(os.Stderr)
	h.SetLoggingLevel()
}

func main() {

	log.Info("Start DockerTree")
	//check if server of client
	wordPtr := flag.String("runmode", "client", "Run mode [client, server]")
	flag.Parse()

	if *wordPtr == "server" || os.Getenv("RUNMODE") == "server" {
		log.Info("Starting in Server mode")
		s.NewRouter()

	} else {
		log.Info("Starting in Client mode")
		c.Start()
	}

}
