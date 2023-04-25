package main_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/gostorz/mgen"
)

var _ = Describe("gostorz", func() {
	It("mgen can generate", func() {
		err := mgen.Generate("test/model")
		Expect(err).To(BeNil())
	})
})
