package client

import (
	"bytes"
	"os"
	str "strings"

	log "github.com/Sirupsen/logrus"
	client "github.com/fsouza/go-dockerclient"
	h "github.com/mangirdaz/dockertree/hellpers"
)

//Docker client together with cofiguration
type Docker struct {
	h.Config
	Client *client.Client
}

func (docker *Docker) DockerInit() {
	log.Debug("Docker Init")
	docker.Config = GetConfig()
	var err error
	docker.Client, err = client.NewClient(docker.Config.Docker.Socket)
	h.Check(err)
}

func (docker *Docker) Build(archivePath string) {

	dockerBuildContext, _ := os.Open(archivePath)
	defer dockerBuildContext.Close()

	var buf bytes.Buffer
	opts := client.BuildImageOptions{
		Name:                docker.Config.Name + ":" + docker.Config.BuildTag,
		Pull:                true,
		NoCache:             true,
		SuppressOutput:      false,
		RmTmpContainer:      true,
		ForceRmTmpContainer: true,
		InputStream:         dockerBuildContext,
		OutputStream:        &buf,
	}

	log.Infof("Build Image [%s]", docker.Config.Name+":"+docker.Config.BuildTag)
	err := docker.Client.BuildImage(opts)
	h.Check(err)
}

func (docker *Docker) TagImage() {

	//tag for main tag
	log.Infof("Tag from [%s] to [%s]", docker.Config.Name+":"+docker.Config.BuildTag, docker.Config.Name+":"+docker.Config.Tag+"-"+docker.Config.Iteration)
	opts := client.TagImageOptions{
		Repo:  docker.Config.Name,
		Tag:   docker.Config.Tag + "-" + docker.Config.Iteration,
		Force: true,
	}
	//todo: add error checking for 500
	err := docker.Client.TagImage(docker.Config.Name+":"+docker.Config.BuildTag, opts)
	h.Check(err)

	//tag for aliases
	for _, value := range docker.Config.Alias {
		log.Infof("Tag to [%s]", value)
		image := str.Split(value, ":")
		if len(image) == 2 {
			log.Info("[" + image[0] + "][" + image[1] + "]")
			opts := client.TagImageOptions{
				Repo:  image[0],
				Tag:   image[1],
				Force: true,
			}
			err := docker.Client.TagImage(docker.Config.Name+":"+docker.Config.BuildTag, opts)
			h.Check(err)
		} else {
			opts := client.TagImageOptions{
				Repo:  image[0],
				Force: true,
			}
			err := docker.Client.TagImage(docker.Config.Name+":"+docker.Config.BuildTag, opts)
			h.Check(err)
		}

	}

	//tag for latest if latest is set
	if docker.Config.Latest {
		log.Infof("Tag to [%s]", "latest")
		opts := client.TagImageOptions{
			Tag:   "latest",
			Repo:  docker.Config.Name,
			Force: true,
		}
		err := docker.Client.TagImage(docker.Config.Name+":"+docker.Config.BuildTag, opts)
		h.Check(err)
	}

}

func (docker *Docker) PushImage() {

	log.Infof("Push Images to Registry")
	opts := client.PushImageOptions{
		Name: docker.Config.Name,
		Tag:  docker.Config.Tag + "-" + docker.Config.Iteration,
	}
	var auth client.AuthConfiguration
	log.Debugf("Push [%s:%s] Image", docker.Config.Name, docker.Config.Tag+"-"+docker.Config.Iteration)
	err := docker.Client.PushImage(opts, auth)
	h.Check(err)

	for _, value := range docker.Config.Alias {
		image := str.Split(value, ":")
		if len(image) == 2 {
			log.Info("[" + image[0] + "][" + image[1] + "]")
			opts := client.PushImageOptions{
				Name: image[0],
				Tag:  image[1],
			}
			log.Debugf("Push [%s:%s] Image", image[0], image[1])
			err := docker.Client.PushImage(opts, auth)
			if err != nil {
				log.Debugf("You may see Unauthorized: authentication required error. This means Docker was not able to push to alias/image. Check docker config ")
			}
			h.Check(err)
		} else {
			opts := client.PushImageOptions{
				Name: image[0],
				Tag:  "latest",
			}
			log.Debugf("Push [%s:%s] Image", image[0], "latest")
			err := docker.Client.PushImage(opts, auth)
			h.Check(err)
		}
	}

}
