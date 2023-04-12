package client_test

import (
	"log"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/client"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

var _ = Describe("client", func() {
	worldName := "c137"
	worldDescription := "is the main world"

	It("can specify client.Headers", func() {
		_, err := stc.Get(
			ctx, "",
			client.Header("setting", "a client.Header"),
			client.Header("setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)

		_, err = stc.List(
			ctx, "",
			options.PropFilter("metadata.ID", "value"),
			client.Header("setting", "a client.Header"),
			client.Header("another setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)

		_, err = stc.Update(
			ctx, "", nil,
			// Options for other APIs are not accepted
			// options.PropFilter("metadata.ID", "value"),
			client.Header("setting", "another client.Header"),
		)
		Expect(err).ToNot(BeNil())
		log.Printf("expected error: %s", err)
	})

	It("cannot GET non-allowed", func() {
		ret, err := stc.Get(
			ctx, generated.ThirdWorldIdentity(worldName))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 404"))

		ret, err = stc.Get(
			ctx,
			store.ObjectIdentity("id/aliksjdlsakjdaslkjdaslkj"))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot CREATE non-allowed", func() {
		w := generated.ThirdWorldFactory()
		ret, err := stc.Create(ctx, w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot UPDATE non-allowed", func() {
		w := generated.ThirdWorldFactory()
		ret, err := stc.Update(ctx,
			generated.ThirdWorldIdentity(worldName), w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 405"))

		ret, err = stc.Update(ctx,
			store.ObjectIdentity("id/aliksjdlsakjdaslkjdaslkj"), w)

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot DELETE non-allowed", func() {
		w := generated.SecondWorldFactory()
		w.External().SetName(worldName)
		ret, err := stc.Create(ctx, w)
		Expect(err).To(BeNil())

		err = stc.Delete(
			ctx, generated.SecondWorldIdentity(worldName))

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))

		sw := ret.(generated.SecondWorld)

		err = stc.Delete(
			ctx,
			sw.Metadata().Identity())

		Expect(err).ToNot(BeNil())
		Expect(err.Error()).To(Equal("http 405"))
	})

	It("cannot LIST non-allowed", func() {
		ret, err := stc.List(
			ctx, generated.ThirdWorldKindIdentity())

		Expect(err).ToNot(BeNil())
		// Expect(err.Error()).To(Equal("http 405"))

		Expect(len(ret)).To(Equal(0))
	})

	It("can initialize metadata", func() {
		world := generated.WorldFactory()
		world.External().SetName(worldName)

		obj, err := stc.Create(ctx, world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)
		Expect(newWorld).ToNot(BeNil())

		Expect(newWorld.Metadata().Identity()).ToNot(Equal(world.Metadata().Identity()))

		Expect(len(newWorld.Metadata().Created())).ToNot(Equal(0))
		Expect(len(newWorld.Metadata().Updated())).To(Equal(0))

		time.Sleep(1 * time.Second)
	})

	It("can update metadata", func() {
		world := generated.WorldFactory()
		world.External().SetName(worldName)

		obj, err := stc.Update(ctx, generated.WorldIdentity(worldName), world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)
		Expect(len(newWorld.Metadata().Updated())).ToNot(Equal(0))
		Expect(newWorld.Metadata().Identity()).ToNot(
			Equal(world.Metadata().Identity()))
		Expect(newWorld.Metadata().Updated()).ToNot(
			Equal(world.Metadata().Updated()))
		Expect(newWorld.Metadata().Created()).ToNot(
			Equal(world.Metadata().Created()))

		time.Sleep(1 * time.Second)
	})

	It("can reset internal", func() {
		obj, err := stc.Get(ctx, generated.WorldIdentity(worldName))
		Expect(err).To(BeNil())

		world := obj.(generated.World)
		world.External().SetName(worldName)
		world.Internal().SetDescription(worldDescription)

		// log.Println(utils.PP(world))

		obj, err = stc.Update(ctx, generated.WorldIdentity(worldName), world)
		Expect(err).To(BeNil())

		newWorld := obj.(generated.World)

		// log.Println(utils.PP(world))
		// log.Println(utils.PP(newWorld))

		Expect(newWorld.Metadata().Identity()).To(
			Equal(world.Metadata().Identity()))
		Expect(newWorld.Metadata().Created()).To(
			Equal(world.Metadata().Created()))
		Expect(newWorld.Metadata().Updated()).ToNot(
			Equal(world.Metadata().Updated()))

		Expect(newWorld.Internal().Description).ToNot(Equal(
			world.Internal().Description()))
	})

})
