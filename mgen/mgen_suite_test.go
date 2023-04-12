package mgen_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestMgen(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "mGen Suite")
}
