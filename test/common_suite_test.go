package common_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/cache"
	"github.com/wazofski/storz/client"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/mongo"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/rest"
	"github.com/wazofski/storz/sql"
	"github.com/wazofski/storz/store"
)

type initializer func()

var clt store.Store
var ctx context.Context = context.Background()
var cancel context.CancelFunc

func suites(suite string) {
	sch := generated.Schema()

	stores := make(map[string]initializer)

	stores["memory"] = func() {
		clt = store.New(
			sch,
			memory.Factory())
	}

	stores["react"] = func() {
		clt = store.New(
			sch,
			react.ReactFactory(
				store.New(
					generated.Schema(),
					memory.Factory())))
	}

	stores["client"] = func() {
		mem := store.New(sch, memory.Factory())

		srv := rest.Server(sch, mem,
			rest.TypeMethods(generated.WorldKind(),
				rest.ActionGet, rest.ActionCreate,
				rest.ActionDelete, rest.ActionUpdate),
			rest.TypeMethods(generated.SecondWorldKind(),
				rest.ActionGet, rest.ActionCreate, rest.ActionDelete))

		cancel = srv.Listen(8000)

		clt = store.New(
			sch,
			client.Factory("http://localhost:8000/"))
	}

	stores["sqlite"] = func() {
		clt = store.New(
			sch,
			logger.StoreFactory("SQLite",
				store.New(
					generated.Schema(),
					sql.Factory(sql.SqliteConnection("test.sqlite")))))
	}

	stores["mysql"] = func() {
		clt = store.New(
			sch,
			logger.StoreFactory("mySQL",
				store.New(
					generated.Schema(),
					sql.Factory(sql.MySqlConnection(
						"root:qwerasdf@tcp(127.0.0.1:3306)/test")))))
	}

	stores["mongo"] = func() {
		clt = store.New(
			sch,
			mongo.Factory("mongodb://localhost:27017/", "tests"))
	}

	stores["cache"] = func() {
		s1 := store.New(
			sch,
			memory.Factory())

		clt = store.New(
			sch,
			cache.Factory(s1))
	}

	stores[suite]()
}

func TestCommon(t *testing.T) {
	RegisterFailHandler(Fail)

	argKey := "store="
	for _, arg := range os.Args {
		if strings.HasPrefix(arg, argKey) {
			sarg := strings.Split(arg, "=")[1]

			suites(sarg)
			RunSpecs(t, fmt.Sprintf("Common Suite %s", sarg))

			break
		}
	}
}

var _ = AfterSuite(func() {
	if cancel != nil {
		cancel()
	}
})
