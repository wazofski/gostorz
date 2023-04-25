package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/wazofski/gostorz/internal/constants"
	"github.com/wazofski/gostorz/store"
)

type _MetaHolder struct {
	Metadata interface{} `json:"metadata,omitempty"`
	// External     interface{} `json:"external,omitempty"`
	// Internal   interface{} `json:"internal,omitempty"`
}

func CloneObject(obj store.Object, schema store.SchemaHolder) store.Object {
	ret := schema.ObjectForKind(obj.Metadata().Kind())
	jsn, err := json.Marshal(obj)
	if err != nil {
		log.Panic(err)
	}
	err = json.Unmarshal(jsn, &ret)
	if err != nil {
		log.Panic(err)
	}
	return ret
}

func ReadStream(r io.ReadCloser) ([]byte, error) {
	b := make([]byte, 0, 512)
	for {
		if len(b) == cap(b) {
			// Add more capacity (let append pick how much).
			b = append(b, 0)[:len(b)]
		}
		n, err := r.Read(b[len(b):cap(b)])
		b = b[:len(b)+n]
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return b, err
		}
	}
}

func UnmarshalObject(body []byte, schema store.SchemaHolder, kind string) (store.Object, error) {
	resource := schema.ObjectForKind(kind)
	err := json.Unmarshal(body, &resource)

	return resource, err
}

func ObjeectKind(response []byte) string {
	obj := _MetaHolder{}
	err := json.Unmarshal(response, &obj)
	if err != nil {
		return ""
	}

	if obj.Metadata == nil {
		return ""
	}

	return obj.Metadata.(map[string]interface{})["kind"].(string)
}

func PP(obj store.Object) string {
	jsn, err := json.MarshalIndent(obj, "", "    ")
	if err != nil {
		log.Panic(err)
	}

	// log.Println(obj)

	return string(jsn)
}

func Timestamp() string {
	return time.Now().Format(time.RFC3339)
}

func Serialize(mo store.Object) ([]byte, error) {
	if mo == nil {
		return nil, constants.ErrObjectNil
	}

	return json.Marshal(mo)
}

func ObjectPath(obj store.Object, path string) *string {
	data, _ := json.Marshal(obj)
	jsn, err := gabs.ParseJSON(data)
	if err != nil {
		log.Fatal(err)
	}
	if !jsn.Exists(strings.Split(path, ".")...) {
		return nil
	}
	ret := strings.ReplaceAll(jsn.Path(path).String(), "\"", "")
	return &ret
}

func ExportFile(targetDir string, name string, content string) error {
	os.Mkdir(targetDir, 0755)

	targetFile := fmt.Sprintf("%s/%s", targetDir, name)

	// log.Printf("exporting file %s", targetFile)

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}

	defer f.Close()

	f.WriteString(content)

	return nil
}

func RuntimeDir() string {
	_, file, _, ok := runtime.Caller(1)
	if ok {
		return filepath.Dir(file)
		// fmt.Printf("Called from %s, line #%d, func: %v\n",
		// 	file, line, runtime.FuncForPC(pc).Name())
	}
	return "./"
}
