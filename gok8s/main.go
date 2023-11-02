package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	"github.com/jmeidam/pypi-server-local-k8s/gok8s/pkg/pypik8s"
	"github.com/jmeidam/pypi-server-local-k8s/gok8s/pkg/utils"
)


func getSecrets(secretsString string) map[string]string {
	// If no login details provided via commandline argument, try the environment variable
	if secretsString == "{}" {
		if pypiloginsEnv := os.Getenv("PYPILOGINS"); pypiloginsEnv != "" {
			secretsString = pypiloginsEnv
		}
	}
	if secretsString == "{}" {
		log.Fatalln("No username and password provided for the service")
	}

	var secretMap map[string]string
	err := json.Unmarshal([]byte(secretsString), &secretMap)
	if err != nil {
		log.Fatalln("Error parsing secrets JSON:", err)
	}


	// Validate if 'username' and 'password' keys are present
	if username, usernameExists := secretMap["username"]; !usernameExists || username == "" {
		log.Fatalln("In the secrets JSON, key 'username' is missing or empty")
	}

	if password, passwordExists := secretMap["password"]; !passwordExists || password == "" {
		log.Fatalln("In the secrets JSON, key 'password' is missing or empty")
	}

	return secretMap
}

func main(){
	var outputFolder *string
	var secrets *string
	var image *string

	outputFolder = flag.String("outpath", "generated_yaml", "path to folder to store the generated yaml files")
	secrets = flag.String("pypilogins", "{}", "JSON string with login name-secret pairs. If empty, will try PYPILOGINS")
	image = flag.String("image", "jmeidam/pypiserver", "Location of the image")

	flag.Parse()

	secretMap := getSecrets(*secrets)
	appName := "pypi-server"
	pvName := "pypi-pv"
	pvcName := "pypi-pvc"
	secretsName := "pypisecret"

	resources := map[string]interface{}{
		"secrets": pypik8s.Secrets(secretsName, secretMap),
		"pv": pypik8s.PV(pvName),
		"pvc": pypik8s.PVClaim(pvcName),
		"deployment": pypik8s.Deployment(appName, secretsName, pvcName, *image),
	}
	utils.WriteResourceToYAML(resources, *outputFolder)
}
