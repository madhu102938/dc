package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gabriel-vasile/mimetype"
	"github.com/madhu102938/dc/constants"
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

func (s *ShellScriptDTO) addFileToScript(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("Error reading the file ", path)
	}

	s.shellScript += "cat >> " +
		"'" +
		path +
		"'" +
		"<< " +
		constants.EOF_constant +
		"\n" +
		string(content) +
		"\n" +
		constants.EOF_constant +
		"\n"
}

func isFileBinary(detectedMIME *mimetype.MIME) bool {
	isBinary := true
	for mtype := detectedMIME; mtype != nil; mtype = mtype.Parent() {
		if mtype.Is("text/plain") {
			isBinary = false
		}
	}

	return isBinary
}

func walkDirectories(
	basePath string,
	shellScriptDTO *ShellScriptDTO,
	exclusionPattern *regexp.Regexp,
	includeHidden bool) error {
	return filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Println("ERROR: error accessing path ", basePath)
		}

		if path == "." {
			return nil
		}

		if includeHidden == false &&
			strings.HasPrefix(path, ".") {
			return nil
		}

		if exclusionPattern.Match([]byte(path)) {
			fmt.Println("WARN: excluding path", path)
			return nil
		}

		if d.IsDir() {
			shellScriptDTO.addFolderToScript(path)
		} else {

			mType, err := mimetype.DetectFile(path)
			if err != nil {
				fmt.Println("ERROR: error finding MIMETYPE of file", path)
			}

			if isFileBinary(mType) {
				fmt.Println("WARN:", path, "not a text file")
				return nil
			}

			shellScriptDTO.addFileToScript(path)
		}
		return nil
	})
}

func traverseDirectories(paths []string, exclusionPattern *regexp.Regexp, copyToClipboard bool, includeHidden bool) {
	var shellScriptDTO ShellScriptDTO
	var err error
	for _, path := range paths {
		if exclusionPattern.Match([]byte(path)) {
			fmt.Println("WARN: excluding path", path)
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			fmt.Println("ERROR: error determining", path, "file or directory", err)
		}

		switch mode := fi.Mode(); {
		case mode.IsDir():
			err = walkDirectories(path, &shellScriptDTO, exclusionPattern, includeHidden)
			if err != nil {
				fmt.Println("ERROR: error going into the path", path, err)
			}
		case mode.IsRegular():
			shellScriptDTO.addFileToScript(path)
		}

	}

	if copyToClipboard {
		err = clipboard.WriteAll(shellScriptDTO.shellScript)
		if err != nil {
			fmt.Println(shellScriptDTO.shellScript)
			log.Fatalf("Failed to write to clipboard due to %v", err)
		} else {
			fmt.Println("Copied the script!")
		}
	} else {
		fileName := "script"+time.Now().String()+".sh"
		err = os.WriteFile(fileName, []byte(shellScriptDTO.shellScript), 0666)
		if err != nil {
			log.Fatalf("Failed to write to file due to %v", err)
		}
		fmt.Println("Script created with file name: ", fileName)
	}
}
