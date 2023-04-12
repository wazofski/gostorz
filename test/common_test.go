package common_test

import (
	"log"
	"sort"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wazofski/storz/generated"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

var _ = Describe("common", func() {

	worldName := "c137zxczx"
	anotherWorldName := "j19zeta7 qweqw"
	worldDescription := "zxkjhajkshdas world of argo"
	newWorldDescription := "is only beoaoqwiewioqu"

	It("can CLEAR everything", func() {
		ret, err := clt.List(ctx, generated.WorldKindIdentity())
		Expect(err).To(BeNil())
		for _, r := range ret {
			err = clt.Delete(ctx, r.Metadata().Identity())
			Expect(err).To(BeNil())
		}

		ret, err = clt.List(ctx, generated.SecondWorldKindIdentity())
		Expect(err).To(BeNil())
		for _, r := range ret {
			err = clt.Delete(ctx, r.Metadata().Identity())
			Expect(err).To(BeNil())
		}

		ret, _ = clt.List(ctx, generated.SecondWorldKindIdentity())
		Expect(len(ret)).To(Equal(0))
		ret, _ = clt.List(ctx, generated.WorldKindIdentity())
		Expect(len(ret)).To(Equal(0))

		// ret, err = clt.List(ctx, generated.ThirdWorldKindIdentity())
		// Expect(err).To(BeNil())
		// for _, r := range ret {
		// 	err = clt.Delete(ctx, r.Metadata().Identity())
		// 	Expect(err).To(BeNil())
		// }
	})

	It("can LIST empty lists", func() {
		ret, err := clt.List(
			ctx, generated.WorldKindIdentity())

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(0))
	})

	It("can POST objects", func() {
		w := generated.WorldFactory()

		w.External().SetName("abc")

		ret, err := clt.Create(ctx, w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))
	})

	It("can LIST single object", func() {
		ret, err := clt.List(
			ctx, generated.WorldKindIdentity())

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal("abc"))
	})

	It("can POST other objects", func() {
		w := generated.SecondWorldFactory()

		w.External().SetName("abc")

		ret, err := clt.Create(ctx, w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))

		ret, err = clt.Get(ctx, ret.Metadata().Identity())
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		w = ret.(generated.SecondWorld)
		Expect(w).ToNot(BeNil())

		ret, err = clt.Get(ctx, generated.SecondWorldIdentity("abc"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		w = ret.(generated.SecondWorld)
		Expect(w).ToNot(BeNil())
	})

	It("can GET objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("abc"))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret.Metadata().Identity())).ToNot(Equal(0))

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("cannot double POST objects", func() {
		w := generated.WorldFactory()

		w.External().SetName("abc")

		ret, err := clt.Create(ctx, w)

		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can PUT objects", func() {
		w := generated.WorldFactory()

		w.External().SetName("abc")
		w.External().SetDescription("def")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("abc"),
			w)

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.External().Description()).To(Equal("def"))
	})

	It("can PUT change naming props", func() {
		w := generated.WorldFactory()

		w.External().SetName("def")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("abc"), w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.External().Name()).To(Equal("def"))

		ret, err = clt.Get(ctx,
			generated.WorldIdentity("abc"))
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can PUT objects BY ID", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
		world.External().SetDescription("zxc")

		log.Println(utils.PP(world))

		ret, err = clt.Update(ctx,
			world.Metadata().Identity(), world)

		log.Println(utils.PP(ret))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world = ret.(generated.World)
		Expect(world).ToNot(BeNil())
		Expect(world.External().Description()).To(Equal("zxc"))
	})

	It("cannot PUT non-existent objects", func() {
		world := generated.WorldFactory()
		Expect(world).ToNot(BeNil())
		world.External().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("zcxzcxzc"), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot PUT non-existent objects BY ID", func() {
		world := generated.WorldFactory()
		world.External().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			world.Metadata().Identity(), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot PUT objects of wrong type", func() {
		world := generated.SecondWorldFactory()
		world.External().SetName("zxcxzcxz")

		ret, err := clt.Update(ctx,
			generated.WorldIdentity("qwe"), world)
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can GET objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("can GET objects BY ID", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("def"))
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world := ret.(generated.World)
		Expect(world).ToNot(BeNil())

		ret, err = clt.Get(ctx,
			world.Metadata().Identity())
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		world = ret.(generated.World)
		Expect(world).ToNot(BeNil())
	})

	It("cannot GET non-existent objects", func() {
		ret, err := clt.Get(ctx,
			generated.WorldIdentity("zxcxzczx"))
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot GET non-existent objects BY ID", func() {
		ret, err := clt.Get(ctx,
			store.ObjectIdentity("id/kjjakjjsadldkjalkdajs"))
		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("can DELETE objects", func() {
		w := generated.WorldFactory()
		w.External().SetName("tobedeleted")

		ret, err := clt.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())

		err = clt.Delete(ctx,
			generated.WorldIdentity(w.External().Name()))
		Expect(err).To(BeNil())

		_, err = clt.Get(ctx,
			generated.WorldIdentity(w.External().Name()))
		Expect(err).ToNot(BeNil())
	})

	It("can DELETE objects BT ID", func() {
		w := generated.WorldFactory()
		w.External().SetName("tobedeleted")

		ret, err := clt.Create(ctx, w)
		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		w = ret.(generated.World)

		err = clt.Delete(ctx, w.Metadata().Identity())
		Expect(err).To(BeNil())

		_, err = clt.Get(ctx, w.Metadata().Identity())
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE non-existent objects", func() {
		err := clt.Delete(ctx,
			generated.WorldIdentity("akjsdhsajkhdaskjh"))
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE non-existent objects BY ID", func() {
		err := clt.Delete(ctx,
			store.ObjectIdentity("id/kjjakjjsadldkjalkdajs"))
		Expect(err).ToNot(BeNil())
	})

	It("cannot GET nil identity", func() {
		_, err := clt.Get(ctx, "")
		Expect(err).ToNot(BeNil())
	})

	It("cannot CREATE nil object", func() {
		_, err := clt.Create(ctx, nil)
		Expect(err).ToNot(BeNil())
	})

	It("cannot PUT nil identity", func() {
		_, err := clt.Update(ctx,
			"", generated.WorldFactory())
		Expect(err).ToNot(BeNil())
	})

	It("cannot PUT nil object", func() {
		_, err := clt.Update(ctx,
			generated.WorldIdentity("qwe"), nil)
		Expect(err).ToNot(BeNil())
	})

	It("cannot DELETE nil identity", func() {
		err := clt.Delete(ctx, "")
		Expect(err).ToNot(BeNil())
	})

	It("can CREATE multiple objects", func() {
		ret, err := clt.List(ctx, generated.WorldKindIdentity())
		Expect(err).To(BeNil())
		for _, r := range ret {
			err = clt.Delete(ctx, r.Metadata().Identity())
			Expect(err).To(BeNil())
		}

		world := generated.WorldFactory()
		world.External().SetName(worldName)
		world.External().SetDescription(worldDescription)

		world2 := generated.WorldFactory()
		world2.External().SetName(anotherWorldName)
		world2.External().SetDescription(newWorldDescription)

		_, err = clt.Create(ctx, world)
		Expect(err).To(BeNil())
		_, err = clt.Create(ctx, world2)
		Expect(err).To(BeNil())

		world3 := generated.SecondWorldFactory()
		world3.External().SetName(anotherWorldName)
		world3.External().SetDescription(newWorldDescription)

		_, err = clt.Create(ctx, world3)
		Expect(err).To(BeNil())
	})

	It("can LIST multiple objects", func() {
		ret, err := clt.List(
			ctx, generated.WorldKindIdentity())

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		sort.Slice(ret, func(i, j int) bool {
			return ret[i].(generated.World).External().Name() < ret[j].(generated.World).External().Name()
		})

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world.External().Description()).To(Equal(worldDescription))

		world2 := ret[1].(generated.World)
		Expect(world2.External().Name()).To(Equal(anotherWorldName))
		Expect(world2.External().Description()).To(Equal(newWorldDescription))
	})

	It("can LIST and sort multiple objects", func() {
		ret, err := clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.OrderBy("external.name"))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world.External().Description()).To(Equal(worldDescription))

		world2 := ret[1].(generated.World)
		Expect(world2.External().Name()).To(Equal(anotherWorldName))
		Expect(world2.External().Description()).To(Equal(newWorldDescription))

		ret, err = clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.OrderBy("external.name"),
			options.OrderDescending())

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(2))

		world = ret[1].(generated.World)
		world2 = ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world2.External().Name()).To(Equal(anotherWorldName))
	})

	It("can LIST and paginate multiple objects", func() {
		ret, err := clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.OrderBy("external.name"),
			options.PageSize(1))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world.External().Description()).To(Equal(worldDescription))

		ret, err = clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.OrderBy("external.name"),
			options.PageSize(1),
			options.PageOffset(1))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world = ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(anotherWorldName))

		ret, err = clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.OrderBy("external.name"),
			options.PageOffset(1),
			options.PageSize(1000))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world = ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(anotherWorldName))
	})

	It("can LIST and filter by primary key", func() {
		ret, err := clt.List(
			ctx, generated.WorldKindIdentity())

		Expect(err).To(BeNil())

		keys := []string{}
		for _, o := range ret {
			keys = append(keys, o.PrimaryKey())
		}

		Expect(len(keys)).To(Equal(2))

		ret, err = clt.List(
			ctx, generated.WorldKindIdentity(),
			options.KeyFilter(keys[0], keys[1]))

		Expect(err).To(BeNil())
		Expect(len(ret)).To(Equal(2))

		for _, k := range keys {
			ret, err = clt.List(
				ctx, generated.WorldKindIdentity(),
				options.KeyFilter(k))

			Expect(err).To(BeNil())
			Expect(len(ret)).To(Equal(1))
			Expect(ret[0].PrimaryKey()).To(Equal(k))
		}
	})

	It("cannot LIST and FILTER BY nonexistent props", func() {
		ret, err := clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.PropFilter("metadata.askdjhasd", "asdsadas"))

		Expect(err).ToNot(BeNil())
		Expect(ret).To(BeNil())
	})

	It("cannot LIST specific object", func() {
		ret, err := clt.List(
			ctx, generated.WorldIdentity(worldName))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	It("cannot LIST specific nonexistent object", func() {
		ret, err := clt.List(
			ctx, generated.WorldIdentity("akjhdsjkhdaskjhdaskj"))

		Expect(ret).To(BeNil())
		Expect(err).ToNot(BeNil())
	})

	worldId := store.ObjectIdentity("")

	It("can LIST and FILTER", func() {
		ret, err := clt.List(
			ctx,
			generated.WorldKindIdentity(),
			options.PropFilter("external.name", worldName))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world.External().Description()).To(Equal(worldDescription))
		worldId = world.Metadata().Identity()
	})

	It("can LIST and FILTER BY ID", func() {
		ret, err := clt.List(
			ctx, generated.WorldKindIdentity(),
			options.PropFilter("metadata.identity", string(worldId)))

		Expect(err).To(BeNil())
		Expect(ret).ToNot(BeNil())
		Expect(len(ret)).To(Equal(1))

		world := ret[0].(generated.World)
		Expect(world.External().Name()).To(Equal(worldName))
		Expect(world.External().Description()).To(Equal(worldDescription))
	})

})
