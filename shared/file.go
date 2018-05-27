package shared

import (
	"fmt"
	"io/ioutil"
	"strings"
)

// GetFiles returns list of filename given directory absolute path
func GetFiles(folderPath, fType string) []string {
	var files []string
	filesArr, err := ioutil.ReadDir(folderPath)
	if err != nil {
		fmt.Println(err)
	}

	for _, f := range filesArr {
		if !f.IsDir() && strings.Contains(f.Name(), fType) {
			files = append(files, f.Name())
		}

	}
	return files
}
