package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/mgen"
)

var _ = Describe("storz", func() {
	It("mgen can generate", func() {
		err := mgen.Generate("test/model")
		Expect(err).To(BeNil())
	})
})
