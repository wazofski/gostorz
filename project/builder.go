package project

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"regexp"

	"github.com/wazofski/storz/utils"
)

func Generate(name string) string {
	if !validName(name) {
		return fmt.Sprintf("cannot create %s. invalid name", name)
	}

	if pathExists(name) {
		return fmt.Sprintf("cannot create %s. directory already exists", name)
	}

	err := os.Mkdir(name, 0755)
	if err != nil {
		return err.Error()
	}

	err = os.Mkdir(fmt.Sprintf("%s/cmd/", name), 0755)
	if err != nil {
		return err.Error()
	}

	err = os.Mkdir(fmt.Sprintf("%s/model/", name), 0755)
	if err != nil {
		return err.Error()
	}

	err = render("templates/go.modtext", name, "go.mod")
	if err != nil {
		return err.Error()
	}

	err = render("templates/go.sumtext", name, "go.sum")
	if err != nil {
		return err.Error()
	}

	err = render("templates/model.gotext", name, "model.go")
	if err != nil {
		return err.Error()
	}

	err = render("templates/main.gotext", name, "cmd/main.go")
	if err != nil {
		return err.Error()
	}

	err = render("templates/object.yamltext", name, "model/objects.yaml")
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("project %s generated", name)
}

func pathExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true

	} else if os.IsNotExist(err) {
		return false
	}

	return true
}

func validName(name string) bool {
	match, _ := regexp.Match(`^[a-z -']+$`, []byte(name))
	// log.Println(err)
	return match
}

func render(rtpath string, path string, name string) error {
	tpath := fmt.Sprintf("%s/%s", utils.RuntimeDir(), rtpath)

	t, err := template.ParseFiles(tpath)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, path)

	if err != nil {
		log.Fatal(err)
	}

	return utils.ExportFile(path, name, buf.String())
}
