package mgen_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/gostorz/generated"
)

var _ = Describe("mgen", func() {

	It("can factory", func() {
		world := generated.WorldFactory()
		Expect(world).ToNot(BeNil())
		Expect(world.External()).ToNot(BeNil())
		Expect(world.Internal()).ToNot(BeNil())
	})

	It("can call setters and getters ", func() {
		world := generated.WorldFactory()
		world.External().SetName("abc")
		Expect(world.External().Name()).To(Equal("abc"))

		world.External().Nested().SetAlive(true)
		Expect(world.External().Nested().Alive()).To(BeTrue())

		world.External().Nested().SetCounter(10)
		Expect(world.External().Nested().Counter()).To(Equal(10))

		world.Internal().SetDescription("qwe")
		Expect(world.Internal().Description()).To(Equal("qwe"))
	})

	It("has metadata", func() {
		world := generated.WorldFactory()
		Expect(world.Metadata().Kind()).To(Equal(generated.WorldKind()))
	})

	It("can deserialize", func() {
		world := generated.WorldFactory()
		world.External().Nested().SetCounter(10)
		world.External().Nested().SetAlive(true)
		world.External().Nested().SetAnotherDescription("qwe")
		world.External().SetName("abc")
		world.Internal().SetDescription("qwe")
		world.Internal().SetList([]generated.NestedWorld{
			generated.NestedWorldFactory(),
			generated.NestedWorldFactory(),
		})

		world.Internal().SetMap(map[string]generated.NestedWorld{
			"a": generated.NestedWorldFactory(),
			"b": generated.NestedWorldFactory(),
		})

		world.Internal().Map()["a"].SetL1([]bool{false, false, true})

		data, err := json.MarshalIndent(world, "", "  ")
		Expect(err).To(BeNil())

		// log.Println(string(data))

		newWorld := generated.WorldFactory()
		err = json.Unmarshal(data, &newWorld)

		Expect(err).To(BeNil())
		Expect(newWorld.External().Nested().Alive()).To(BeTrue())
		Expect(newWorld.External().Nested().Counter()).To(Equal(10))
		Expect(newWorld.External().Name()).To(Equal("abc"))
		Expect(newWorld.Internal().Description()).To(Equal("qwe"))
		Expect(len(newWorld.Internal().List())).To(Equal(2))

		data2, err := json.MarshalIndent(newWorld, "", "  ")
		Expect(err).To(BeNil())
		Expect(data).To(Equal(data2))
	})

	It("has working schema", func() {
		world := generated.WorldFactory()
		schema := generated.Schema()
		obj := schema.ObjectForKind(string(world.Metadata().Kind()))
		Expect(obj).ToNot(BeNil())
		anotherWorld := obj.(generated.World)
		Expect(anotherWorld).ToNot(BeNil())
	})

	It("can clone objects", func() {
		world := generated.WorldFactory()
		world.External().Nested().SetCounter(10)
		world.External().Nested().SetAlive(true)
		world.External().Nested().SetAnotherDescription("qwe")

		world.External().SetName("abc")
		world.Internal().SetDescription("qwe")
		world.Internal().SetList([]generated.NestedWorld{
			generated.NestedWorldFactory(),
			generated.NestedWorldFactory(),
		})

		world.Internal().SetMap(map[string]generated.NestedWorld{
			"a": generated.NestedWorldFactory(),
			"b": generated.NestedWorldFactory(),
		})

		world.Internal().Map()["a"].SetL1([]bool{false, false, true})

		newWorld := world.Clone().(generated.World)
		Expect(newWorld.External().Nested().Alive()).To(BeTrue())
		Expect(newWorld.External().Nested().Counter()).To(Equal(10))
		Expect(newWorld.External().Nested().AnotherDescription()).To(Equal("qwe"))
		Expect(newWorld.External().Name()).To(Equal("abc"))
		Expect(newWorld.Internal().Description()).To(Equal("qwe"))
		Expect(len(newWorld.Internal().List())).To(Equal(2))
	})

})
