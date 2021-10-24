package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

func main() {
	templates := []string{
		"template/main.go",
		"template/go.mod",
		"template/cmd/example.go",
		"template/cmd/root.go",
	}

	cliPath := "./new_cli"
	if len(os.Args) > 1 {
		cliPath = os.Args[1]
	}

	headPath, cliName := path.Split(cliPath)
	templateDir := "template"

	data := struct {
		Name string
	}{cliName}

	os.Mkdir(filepath.Join(cliPath, data.Name), 0775)
	err := filepath.Walk(templateDir,
		func(p string, info os.FileInfo, err error) error {
			newp := filepath.Join(headPath, strings.ReplaceAll(p, templateDir, data.Name))
			if info.IsDir() {
				os.MkdirAll(newp, 0775)
			}

			if !inArray(templates, p) && !info.IsDir() {
				input, err := ioutil.ReadFile(p)
				if err != nil {
					fmt.Println(err)
					return nil
				}

				err = ioutil.WriteFile(newp, input, 0775)
				if err != nil {
					fmt.Println("Error creating", newp)
					fmt.Println(err)
					return nil
				}
			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	t, err := template.ParseFiles(templates...)
	if err != nil {
		panic(err)
	}

	for _, v := range templates {
		_, file := filepath.Split(v)
		destination, err := os.Create(filepath.Join(headPath, strings.ReplaceAll(v, templateDir, data.Name)))
		defer destination.Close()
		err = t.ExecuteTemplate(destination, file, data)
		if err != nil {
			panic(err)
		}
	}
}

func inArray(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
