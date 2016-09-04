# DockerTree Decentralize Image build Utility
Utility purpose is to managea Docker image builds, maintain image tree dependencies, and enable different teams to use same base images, managed by different teams, with no need to have common CI&CD tool or central build managment chain.

## Usage
### Client
Client is default mode of the binary (`execution ./dockertree`). It will check for docker socket presence, and reads .dockertree configuration file for build configuration.

##### Namespaces
As much as we like `registry/image:tag` format, we do not support it. All images should be in the `namespaces`. example: `regsitry/namespace1/image:tag`. This applies to Base images used in the dockerfile and images you are building.

##### CLIENT Env Variables:
        DOCKERTREE_CONFIG_NAME - set alternative .dockertree config file name

### Server
Server is ran using ./dockertree --server flag on the same binary. 
Server expects Consul Key Value storage to be present. Its configurable via env Variables

##### SERVER Env Variables:
        KEYVAL_STORAGE_IP - ip address of keyvalue storage. defualt 0.0.0.0
        KEYVAL_STORAGE_PORT - key value storage port. default 8500
        SERVER_IP - server ip address to bind to. default 0.0.0.0
        SERVER_PORT - server port to bind to. default 8080.
        DEV_DUMMYDATA - set to true will load dummy data from dummydata.json on start. Default false

## Common Env Variables:
        LOG_LEVEL - options [debug, fatal, info] default - info

##### Important
If you have dependency:
`rhel7-base --> jboss` and you change it to `rhe7-base-->java-base-->jboss`, you need to make sure you removed first relationship from rhel7-base image in KVStorage, because it will do double trigger.
This feature we are working on. We will read previous config (which include base image) and if it changes, we remove relationship.

### .dockertree config file
```
tag: 7.2 #Image to to be built. Iteration tag-x will be added from server. Aliases do not have iteration
name: localhost:5000/mangirdaz/my-test-image #Image full name. If you use private regsitry, add it. If you use docker deamon config you can use shorter version. Registry will be used for push, BUT in the server we store data without regsitry
latest: true #if this image have to be tagged as latest too
alias: [] #Alias for other image names 
#  - test/my-test-image-new   #no tag will create latest tag
#  - test/test:latest
#  - myimage
server: localhost:3001 #server address with port
moduledir: dockerSrc/docker-debug-container #where Dockerfile sits relativly to the current dir
#buildtag: build  #image build tag before tagging.
#docker deamon/client configuration - optional
#docker:
#  defaultHeaders:
#    - name: User-Agent
#      value: engine-api-cli-1.0
# socket: unix:///var/run/docker.sock 
#  apiVersion: v1.22
```

## Build:
### GO:
### build for scratch image:
```
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .
```
### docker build:
```
docker build -t mangirdaz/dockertree .
```
### run:
```
For server:
docker run -d -p 8080:8080 mangirdaz/dockertree
For client:
./dockertree
```

TODO:
Add variable:
LOG_LEVEL rename to DOCKERTREE_LOG_LEVEL - set logging level

Client: 
* Cover code in tests
* Add image removal after push (housekeeping)
* Add help for CLI
* Add flags for CLI
* add Tesing framework - Build with certain files -> run test --> rebuild without files--> push
* add ability to create temp folder in custom locatiom
* detect if user has commented base image in dockerfile #FROM and uses just FROM in the line 2, etcd
* Enable Authentication
* Convert to docker plugin

Server:

* add trigger sequences for builds /webooks, api calls to go, jenkins, etc
* add server build capability with routines. If our build agents does not have docker we could build on central server
* add proxy mode. So we deploy client as docker plugin as docker build with certain tag will trigger remote api
* add context and tide up structs 
* if base image FROM was changed, remove dependency from data structure
* Enable authentication 
* 
 