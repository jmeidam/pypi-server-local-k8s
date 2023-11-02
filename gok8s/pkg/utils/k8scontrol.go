package utils

import (
	"log"
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

func WriteResourceToYAML(resources map[string]interface{}, outputDir string) {
	// resources: resourcename and resource object pairs
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		log.Fatalf("Error creating directory: %s", err)
	}
	for resourcename, obj := range resources {
		yamlData, err := yaml.Marshal(obj)
		if err != nil {
			log.Fatalf("Error converting object %s to yaml: %s", resourcename, err)
		}

		filename := resourcename + ".yaml"
		filePath := filepath.Join(outputDir, filename)
		if err := os.WriteFile(filePath, yamlData, 0644); err != nil {
			log.Fatalf("Error writing to %s: %s", filePath, err)
		}
		// these do not need to be log-entries, hence fmt instead of log
		fmt.Printf("Succesfully generated %s\n", filePath)
	}
}
