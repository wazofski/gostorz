package react_test

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/memory"
	"github.com/wazofski/storz/react"
	"github.com/wazofski/storz/store"
)

func TestReact(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "React Suite")
}

var str store.Store
var ctx context.Context

var _ = BeforeSuite(func() {
	sch := generated.Schema()

	mem := store.New(
		sch,
		memory.Factory())

	str = store.New(
		sch,
		react.ReactFactory(mem,
			react.Subscribe(generated.WorldKind(), react.ActionDelete, WorldDeleteCb),
			react.Subscribe(generated.WorldKind(), react.ActionUpdate, WorldUpdateCb),
			react.Subscribe(generated.WorldKind(), react.ActionCreate, WorldCreateCb)))
})
