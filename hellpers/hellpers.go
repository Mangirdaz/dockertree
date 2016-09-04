package hellpers

import (
	"bytes"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/jhoonb/archivex"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
)

func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func Check(e error) {
	if e != nil {
		log.Errorf("Error detected: [%s]", e)
	}
}

func SetLoggingLevel() {

	loglevel := strings.ToLower(os.Getenv("LOG_LEVEL"))

	switch {
	case loglevel == "debug":
		log.SetLevel(log.DebugLevel)
	case loglevel == "fatal":
		log.SetLevel(log.FatalLevel)
	case loglevel == "info":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func CreateTempFolder() (folder string) {
	log.Debug("Create temp Folder")
	folder, err := ioutil.TempDir(os.TempDir(), "dockertree")
	Check(err)
	return folder
}

func DeleteTempFolder(folder string) {
	log.Debug("Delete temp Folder")
	err := os.RemoveAll(folder)
	Check(err)
}

func CreateTarArchive(tempDir string, source string) (path string) {
	log.Infof("Create tar archive in [%s]", tempDir+"/Dockerfile.tar")
	tar := new(archivex.TarFile)
	tar.Create(tempDir + "/Dockerfile.tar")
	tar.AddAll(source, false)
	tar.Close()
	return tempDir + "/Dockerfile.tar"
}

//StreamToByte stream to bytes conversion
func StreamToByte(stream io.Reader) []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf.Bytes()
}

//Becase we call data from struct we are key sensitive. So we make this hack. TODO: migrate to annotations
func UpcaseInitial(str string) string {
	for i, v := range str {
		var end string
		end = str[i+1:]
		return string(unicode.ToUpper(v)) + strings.ToLower(end)
	}
	return ""
}

func RemoveDuplicates(xs []string) (result []string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range xs {
		if !found[x] {
			found[x] = true
			(xs)[j] = (xs)[i]
			j++
		}
	}
	xs = (xs)[:j]
	return xs
}

func GetBaseImageFromDockerfile(path string) (string, error) {
	log.Debugf("Get base image from [%s]", path)
	input, err := ioutil.ReadFile(path)
	if err != nil {
		log.Errorf("Failed to Open Dockerfile [%s] with error [%s]", path, err)
	}
	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, "FROM") {
			result := strings.Replace(strings.Replace(lines[i], "FROM", "", 1), " ", "", 1)
			return result, nil
		}
	}
	return "", err
}

func ParseImageName(name string) (string, error) {
	log.Debugf("Parsing Image [%s] name", name)
	var err error
	temp := strings.Split(name, "/")
	switch {
	case len(temp) == 3:
		log.Debug("Image name has registry/namespace/name in it")
		return fmt.Sprintf("%s/%s", temp[1], temp[2]), nil
	case len(temp) == 2:
		log.Debug("Image name has namespace/name in it")
		return fmt.Sprintf("%s/%s", temp[0], temp[1]), nil
	case len(temp) == 1:
		log.Error("Image does not have namespace")
		err = errors.New("Image does not have namespace")
	}
	return name, err
}

func (box *Results) AddItem(item Result) []Result {
	box.Result = append(box.Result, item)
	return box.Result
}

func StrToInt(str string) int {
	nonFractionalPart := strings.Split(str, ".")
	temp, err := strconv.Atoi(nonFractionalPart[0])
	Check(err)
	return temp
}

func IntToStr(integer int) string {
	temp := strconv.Itoa(integer)
	return temp
}
