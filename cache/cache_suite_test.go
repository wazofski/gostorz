package cache_test

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/gostorz/cache"
	"github.com/wazofski/gostorz/generated"
	"github.com/wazofski/gostorz/memory"
	"github.com/wazofski/gostorz/store"
)

func TestCache(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "cache suite")
}

var mainst store.Store
var cached store.Store

var _ = BeforeSuite(func() {
	sch := generated.Schema()

	mainst = store.New(sch, memory.Factory())
	cached = store.New(sch, cache.Factory(mainst, 1*time.Second))
})
