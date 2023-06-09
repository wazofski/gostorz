package client_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/gostorz/client"
	"github.com/wazofski/gostorz/generated"
	"github.com/wazofski/gostorz/memory"
	"github.com/wazofski/gostorz/rest"
	"github.com/wazofski/gostorz/store"
)

var stc store.Store
var ctx context.Context
var cancel context.CancelFunc

var _ = BeforeSuite(func() {
	sch := generated.Schema()

	mem := store.New(sch, memory.Factory())
	// rct := store.New(sch, react.ReactFactory(mhr))

	srv := rest.Server(sch, mem,
		rest.TypeMethods(generated.WorldKind(),
			rest.ActionGet, rest.ActionCreate,
			rest.ActionDelete, rest.ActionUpdate),
		rest.TypeMethods(generated.SecondWorldKind(),
			rest.ActionGet, rest.ActionCreate))

	cancel = srv.Listen(8000)

	stc = store.New(
		generated.Schema(),
		client.Factory(
			"http://localhost:8000/",
			client.Header("test", "header")))
})

var _ = AfterSuite(func() {
	if cancel != nil {
		cancel()
	}
})

func TestClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Client Suite")
}
