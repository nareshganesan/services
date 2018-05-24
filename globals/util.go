package globals

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// SetProjectHome returns services app project root path
func SetProjectHome() string {
	// Find home directory.
	projectHome, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(projectHome)

}

// CreateFile , creates file if it doesn't exist
func CreateFile(path string) {
	// detect if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("creating file: %s \n", path)
		// create file if not exists
		file, err := os.Create(path)
		if err != nil {
			fmt.Println("Error creating file")
			fmt.Println(err.Error())
			return
		}
		defer file.Close()
	}
}

// CreateFolder creates the path if it doesn't exist
func CreateFolder(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("creating folder: %s \n", path)
		os.Mkdir(path, os.FileMode(0755))
	}
}

// LoadJSON given filepath and interface , loads the interface with json data
func LoadJSON(filePath string, obj interface{}) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("file error")
		fmt.Println(err)
		os.Exit(1)
	}
	json.Unmarshal(file, &obj)
}
