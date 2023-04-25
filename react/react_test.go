package react_test

import (
	"context"
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/gostorz/generated"
	"github.com/wazofski/gostorz/store"
)

func WorldCreateCb(obj store.Object, str store.Store) error {
	world := obj.(generated.World)
	world.Internal().SetDescription("abc")

	return nil
}

func WorldUpdateCb(obj store.Object, str store.Store) error {
	anotherWorld := generated.SecondWorldFactory()
	anotherWorld.External().SetName("def")

	_, err := str.Create(context.Background(), anotherWorld)
	return err
}

func WorldDeleteCb(obj store.Object, str store.Store) error {
	return fmt.Errorf("cannot delete")
}

var _ = Describe("react", func() {

	It("can set internal on CREATE", func() {
		world := generated.WorldFactory()
		world.External().SetName("abc")

		ret, err := str.Create(ctx, world)

		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world = ret.(generated.World)

		ret, err = str.Get(ctx, world.Metadata().Identity())
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world = ret.(generated.World)
		Expect(world.Internal().Description()).To(Equal("abc"))
	})

	It("can creat objects on UPDATE", func() {
		ret, err := str.Get(ctx, generated.WorldIdentity("abc"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		world := ret.(generated.World)
		world.External().SetDescription("qwe")
		ret, err = str.Update(ctx, generated.WorldIdentity("abc"), world)
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())

		ret, err = str.Get(ctx, generated.SecondWorldIdentity("def"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
	})

	It("can reject DELETE", func() {
		err := str.Delete(ctx, generated.WorldIdentity("abc"))
		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("cannot delete"))

		ret, err := str.Get(ctx, generated.WorldIdentity("abc"))
		Expect(ret).ToNot(BeNil())
		Expect(err).To(BeNil())
	})
})
