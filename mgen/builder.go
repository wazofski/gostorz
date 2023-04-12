package mgen

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/wazofski/storz/utils"
)

func Generate(model string) error {
	structs, resources := loadModel(model)

	imports := []string{
		// "errors",
		// "log",
		// "strings",
		"fmt",
		"encoding/json",
		"github.com/wazofski/storz/utils",
		"github.com/wazofski/storz/store",
	}

	var b strings.Builder
	b.WriteString(render("templates/imports.gotext", imports))
	b.WriteString(compileResources(resources))
	b.WriteString(compileStructs(structs))

	str := strings.ReplaceAll(b.String(), "&#34;", "\"")
	res, err := format.Source([]byte(str))

	if err != nil {
		log.Println(err)
		res = []byte(str)
	}

	targetDir := "generated"
	os.RemoveAll(targetDir)

	return utils.ExportFile(targetDir, "objects.go", string(res))
}

type _Interface struct {
	Name       string
	Methods    []string
	Implements []string
}

func compileResources(resources []_Resource) string {
	var b strings.Builder

	for _, r := range resources {
		props := []_Prop{
			{
				Name:    "Meta",
				Type:    "store.Meta",
				Json:    "metadata",
				Default: fmt.Sprintf("store.MetaFactory(\"%s\")", r.Name),
			},
		}

		if len(r.External) > 0 {
			props = append(props,
				_Prop{
					Name: "External",
					Type: r.External,
					Json: "external",
				})
		}

		if len(r.Internal) > 0 {
			props = append(props,
				_Prop{
					Name: "Internal",
					Type: r.Internal,
					Json: "internal",
				})
		}

		s := _Struct{
			Name:   r.Name,
			Props:  props,
			Embeds: []string{},
			Implements: []string{
				"store.Object",
			},
		}

		b.WriteString(compileStruct(s))
		b.WriteString(render("templates/meta.gotext", r))
		b.WriteString(render("templates/clone.gotext", s))
	}

	b.WriteString(render("templates/schema.gotext", resources))

	return b.String()
}

func compileStructs(structs []_Struct) string {
	var b strings.Builder

	for _, s := range structs {
		b.WriteString(compileStruct(s))
	}

	return b.String()
}

type _Tuple struct {
	A string
	B string
}

func compileStruct(s _Struct) string {
	var b strings.Builder
	methods := []string{}

	s.Props = addDefaultPropValues(s.Props)

	for _, p := range s.Props {
		if p.Name != "Meta" {
			methods = append(methods,
				fmt.Sprintf("%s() %s", p.Name, p.Type))
		}

		if p.Name != "Meta" && p.Name != "External" && p.Name != "Internal" {
			methods = append(methods,
				fmt.Sprintf("Set%s(v %s)", p.Name, p.Type))
		}

		if p.Name == "External" {
			b.WriteString(
				render("templates/specinternal.gotext",
					_Tuple{A: s.Name, B: p.Type}))
		}
	}

	impl := append(s.Implements, "json.Unmarshaler")

	b.WriteString(render("templates/interface.gotext", _Interface{
		Name:       s.Name,
		Methods:    methods,
		Implements: impl,
	}))

	b.WriteString(render("templates/structure.gotext", s))
	b.WriteString(render("templates/unmarshall.gotext", s))

	return b.String()
}

func render(rpath string, data interface{}) string {
	path := fmt.Sprintf("%s/%s", utils.RuntimeDir(), rpath)

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)

	if err != nil {
		log.Fatal(err)
	}

	return buf.String()
}

func addDefaultPropValues(props []_Prop) []_Prop {
	res := []_Prop{}

	for _, p := range props {
		if len(p.Default) > 0 {
			res = append(res, p)
			continue
		}

		res = append(res, _Prop{
			Name:    p.Name,
			Json:    p.Json,
			Type:    p.Type,
			Default: typeDefault(p.Type),
		})
	}

	return res
}

func typeDefault(tp string) string {
	if strings.HasPrefix(tp, "[]") {
		return fmt.Sprintf("%s {}", tp)
	}
	if strings.HasPrefix(tp, "map") {
		return fmt.Sprintf("make(%s)", tp)
	}

	switch tp {
	case "string":
		return "fmt.Sprint()"
	case "bool":
		return "false"
	case "int":
		return "0"
	case "float":
		return "0"
	default:
		return fmt.Sprintf("%sFactory()", tp)
	}
}

func (u _Prop) IsMap() bool {
	if len(u.Type) < 3 {
		return false
	}

	return u.Type[:3] == "map"
}

func (u _Prop) IsArray() bool {
	if len(u.Type) < 2 {
		return false
	}

	return u.Type[:2] == "[]"
}

func (u _Prop) StrippedType() string {
	if u.IsMap() {
		return u.Type[strings.LastIndex(u.Type, "]")+1:]
	}
	if u.IsArray() {
		return u.Type[2:]
	}
	return u.Type
}

func (u _Prop) StrippedDefault() string {
	return typeDefault(u.StrippedType())
}
