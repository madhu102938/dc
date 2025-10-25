package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/madhu102938/directory-cloner/constants"
)

type ShellScriptDTO struct {
	shellScript string
}

func (s *ShellScriptDTO) addFolderToScript(path string) {
	s.shellScript += "mkdir -p " + 
		"'" +
		path +
		"'" +
		"\n"
}

func (s *ShellScriptDTO) addFileToScript(path, content string) {
	s.shellScript += "cat >> " + 
		"'" +
		path +
		"'" +
		"<< " +
 		constants.EOF_constant +
		"\n" + 
		content + 
		"\n" + 
		constants.EOF_constant +
		"\n"
}

func main() {
	basePath := "."
	var shellScriptDTO ShellScriptDTO
	
	// shellScriptDTO.shellScript = ""
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if (err != nil) {
			fmt.Println("Error accessing path ", basePath)
		}
		
		if (path == ".") {
			return nil
		}
		
		if d.IsDir() {
			shellScriptDTO.addFolderToScript(path)
		} else {
			content, err := os.ReadFile(path)
			
			if err != nil {
				fmt.Println("Error reading the file ", path)
			}
			
			shellScriptDTO.addFileToScript(path, string(content))
		}
		return nil
	})
	
	if err != nil {
		fmt.Println("Error walking the directory ", basePath)
	}
	
	err = clipboard.WriteAll(shellScriptDTO.shellScript)
	if err != nil {
		fmt.Println(shellScriptDTO.shellScript)
		log.Fatalf("Failed to write to clipboard due to %v", err)
	} else {
		fmt.Println("Copied the script!")
	}
}