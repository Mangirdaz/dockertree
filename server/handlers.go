package server

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	h "github.com/mangirdaz/dockertree/hellpers"
	"net/http"
)

//Index index method for API
func Index(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "DockerTree Server API")
}

func GetNamespaces(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	path := "images"

	log.Info("Get Images")
	images, err := storage.GetSubdir(path)
	log.Debugf("Image got [%s]", images)
	if err != nil {
		log.Error(fmt.Sprintf("Directory [%s] not found", path))
		printStatus(w, http.StatusNoContent, "Directory not found")
	} else {
		var namespacelist h.NamespaceList
		namespacelist.Results = images
		json, err := json.Marshal(namespacelist)
		if err != nil {
			log.Error(fmt.Sprintf("Result not recieved with error [%s]", err))
			printStatus(w, http.StatusNoContent, "Error while getting values")
		}
		fmt.Fprintln(w, string(json))
	}
}

func GetImages(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	log.Infof("Requested path images/%s", vars["namespace"])
	path := fmt.Sprintf("images/%s", vars["namespace"])
	log.Info("Get Images")
	images, err := storage.GetSubdir(path)
	log.Debugf("Image got [%s]", images)
	if err != nil {
		log.Error(fmt.Sprintf("Directory [%s] not found", path))
		printStatus(w, http.StatusNoContent, "Directory not found")
	} else {
		var imagelist h.ImageList
		imagelist.Results = images
		json, err := json.Marshal(imagelist)
		if err != nil {
			log.Error(fmt.Sprintf("Result not recieved with error [%s]", err))
			printStatus(w, http.StatusNoContent, "Error while getting values")
		}
		fmt.Fprintln(w, string(json))
	}
}

func AddImage(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	log.Info(fmt.Sprintf("Add Images"))
	decoder := json.NewDecoder(r.Body)
	var config h.Config
	err := decoder.Decode(&config)
	log.Debug(config.Name)
	if err != nil {
		log.Error(fmt.Sprintf("Error while parsing payload [%s]", err))
		printStatus(w, http.StatusExpectationFailed, "Error while parsing payload ")
	} else {
		//var storage LibKVBackend
		log.Infof("Add images with payload [%s]", config)
		pair, err := storage.ParseAndLoadImage(config)
		if err == nil {
			var results h.Result
			results.Key = pair.Key
			results.Value = string(pair.Value)
			json, err := json.Marshal(results)
			h.Check(err)
			log.Debugf("Payload for response [%s]", string(json))
			printStatus(w, http.StatusOK, string(json))
		} else {
			log.Error(fmt.Sprintf("Image not stored with error [%s]", err))
			printStatus(w, http.StatusNotImplemented, "Image Not Stored")
		}
		//add dependency for the image
		go storage.AddDependency(config)
	}

}

func GetTag(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	log.Infof("Requested path %s/%s/tags/%s", vars["namespace"], vars["name"], vars["tag"])
	path := fmt.Sprintf("images/%s/%s/tags/%s", vars["namespace"], vars["name"], vars["tag"])
	log.Info("Get tags")
	images, err := storage.GetSubdir(path)
	if err != nil {
		log.Error(fmt.Sprintf("Directory [%s] not found", path))
		printStatus(w, http.StatusNoContent, "Directory not found")
	} else {
		var list h.TagMetadaList
		list.Results = images
		json, err := json.Marshal(list)
		if err != nil {
			log.Error(fmt.Sprintf("Result not recieved with error [%s]", err))
			printStatus(w, http.StatusNoContent, "Error while getting values")
		}
		fmt.Fprintln(w, string(json))
	}
}

func GetMetadata(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	log.Infof("Requested path %s/%s/tags/%s/%s", vars["namespace"], vars["name"], vars["tag"], vars["metadata"])
	path := fmt.Sprintf("images/%s/%s/tags/%s/%s", vars["namespace"], vars["name"], vars["tag"], vars["metadata"])
	log.Info("GetMetadata")
	result, err := storage.Get(path)
	if err != nil {
		log.Error(fmt.Sprintf("Directory [%s] not found", path))
		printStatus(w, http.StatusNoContent, "Directory not found")
	} else {
		var list h.MetadaResults
		list.Results.Key = result.Key
		list.Results.Value = string(result.Value)
		json, err := json.Marshal(list)
		if err != nil {
			log.Error(fmt.Sprintf("Result not recieved with error [%s]", err))
			printStatus(w, http.StatusNoContent, "Error while getting values")
		}
		fmt.Fprintln(w, string(json))
	}
}

func UpdateIteration(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	log.Info("Update Metadata")
	vars := mux.Vars(r)
	path := fmt.Sprintf("images/%s/%s/tags/%s/iteration", vars["namespace"], vars["name"], vars["tag"])

	decoder := json.NewDecoder(r.Body)
	var config string
	err := decoder.Decode(&config)
	h.Check(err)
	storage.storeKey(path, []byte(config))
}

func GetImage(w http.ResponseWriter, r *http.Request, storage LibKVBackend) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	vars := mux.Vars(r)
	log.Infof("Requested image [%s/%s]", vars["namespace"], vars["name"])

	path := fmt.Sprintf("images/%s/%s/tags", vars["namespace"], vars["name"])
	log.Infof("Get Image Tags with path [%s]", path)
	images, err := storage.GetSubdir(path)
	log.Debugf("Image got [%s]", images)
	if err != nil {
		log.Error(fmt.Sprintf("Directory [%s] not found", path))
		printStatus(w, http.StatusNoContent, "Directory not found")
	} else {
		var tagList h.TagList
		tagList.Results = images
		json, err := json.Marshal(tagList)
		if err != nil {
			log.Error(fmt.Sprintf("Note not recieved with error [%s]", err))
			printStatus(w, http.StatusNoContent, "Error while getting values")
		}
		fmt.Fprintln(w, string(json))
	}
}

//extend standart handler with our required storage backend details
type backendHandler func(w http.ResponseWriter, r *http.Request, storage LibKVBackend)

type Handler interface {
	ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

//retunr what mux expects
func mybackendHandler(handler backendHandler, storage LibKVBackend) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, storage)
	}
}
func CheckAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	user, pass, _ := r.BasicAuth()
	if !checkPass(user, pass) {
		printStatus(w, http.StatusUnauthorized, "Unauthorized")
		return
	}
	next(w, r)
}

func checkPass(user, pass string) bool {
	log.Info(fmt.Sprintf("User [%s] and Pass [%s]", user, pass))
	//if user == "mj" && pass == "test" {
	log.Info("Pass OK")
	return true
	//} else {
	//	log.Info("Pass Error")
	//	return false
	//	}
	//return false
}

func printStatus(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	response := h.ErrorMessage{status, message}
	json, err := json.Marshal(response)
	log.Info("Message to return: " + string(json))
	log.Info(err)
	if err == nil {
		fmt.Fprintln(w, string(json))
	}

}
