package client

import (
	log "github.com/Sirupsen/logrus"
	h "github.com/mangirdaz/dockertree/hellpers"
)

func Start() {

	var docker Docker
	//init docker client
	docker.DockerInit()
	//create temp folder
	docker.Config.TempDir = h.CreateTempFolder()
	//create archive from Module dir (project root)
	archivePath := h.CreateTarArchive(docker.Config.TempDir, docker.Config.ModuleDir)
	//build image
	log.Debug("Start Dockertree build")
	//before start build we send our info to dockertree server
	// *image from for dependencies update
	// *our config
	from, err := h.GetBaseImageFromDockerfile(docker.Config.ModuleDir + "/Dockerfile")
	h.Check(err)
	log.Debugf("Base image [%s]", from)
	//store base image for dependency managment
	docker.Config.BaseImage = from
	iteration := docker.PostConfigToServer()
	//+1 for next build
	docker.Config.Iteration = h.IntToStr(h.StrToInt(iteration) + 1)
	log.Debugf("Build iteration [%s]", docker.Config.Iteration)
	docker.Build(archivePath)
	docker.TagImage()
	docker.PushImage()
	defer h.DeleteTempFolder(docker.Config.TempDir)
	docker.PostIterationToServer()

}
