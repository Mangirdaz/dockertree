package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() {

	//create new router
	router := mux.NewRouter().StrictSlash(true)
	storage := NewLibKVBackend()
	config := storage.Config
	config.SetDefault()

	//api backend init
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	apiV1.Methods("GET").Path("/").Name("Index").Handler(http.HandlerFunc(mybackendHandler(Index, storage)))

	//Images methods
	apiV1.Methods("GET").Path("/images").Name("NameSpaces").Handler(http.HandlerFunc(mybackendHandler(GetNamespaces, storage)))
	apiV1.Methods("GET").Path("/images/{namespace}").Name("Images").Handler(http.HandlerFunc(mybackendHandler(GetImages, storage)))
	apiV1.Methods("GET").Path("/images/{namespace}/{name}").Name("ImageName").Handler(http.HandlerFunc(mybackendHandler(GetImages, storage)))

	//add image
	apiV1.Methods("POST").Path("/images").Name("AddImage").Handler(http.HandlerFunc(mybackendHandler(AddImage, storage)))
	apiV1.Methods("POST").Path("/images/{namespace}/{name}").Name("AddImage").Handler(http.HandlerFunc(mybackendHandler(AddImage, storage)))

	//tags methods
	apiV1.Methods("GET").Path("/images/{namespace}/{name}/tags").Name("Tag List").Handler(http.HandlerFunc(mybackendHandler(GetImage, storage)))

	//get tag specific metadata options available (list of subdirs)
	apiV1.Methods("GET").Path("/images/{namespace}/{name}/tags/{tag}").Name("Tag data List").Handler(http.HandlerFunc(mybackendHandler(GetTag, storage)))
	apiV1.Methods("GET").Path("/images/{namespace}/{name}/tags/{tag}/{metadata}").Name("Tag metadata List").Handler(http.HandlerFunc(mybackendHandler(GetMetadata, storage)))

	//update metadata
	apiV1.Methods("POST").Path("/images/{namespace}/{name}/tags/{tag}/iteration").Name("Update iteration").Handler(http.HandlerFunc(mybackendHandler(UpdateIteration, storage)))

	//middleware intercept
	midd := http.NewServeMux()
	midd.Handle("/", router)
	midd.Handle("/api/v1/images", negroni.New(
		negroni.HandlerFunc(CheckAuth),
		negroni.Wrap(apiV1),
	))
	n := negroni.Classic()
	n.UseHandler(midd)
	log.Fatal(http.ListenAndServe(config.ServerIP+":"+config.ServerPort, n))

	//return router

}
