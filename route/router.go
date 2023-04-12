package route

import (
	"context"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

type routeStore struct {
	Schema  store.SchemaHolder
	Log     logger.Logger
	Mapping map[string]store.Store
	Default store.Store
}

type Mapping struct {
	Kind  string
	Store store.Store
}

func Factory(deault store.Store, mappings ...Mapping) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &routeStore{
			Schema:  schema,
			Log:     logger.Factory("route"),
			Mapping: make(map[string]store.Store),
			Default: deault,
		}

		for _, m := range mappings {
			client.Mapping[m.Kind] = m.Store
		}

		return client, nil
	}
}

func (d *routeStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())

	return d.Default.Create(ctx, obj, opt...)
}

func (d *routeStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())

	return d.Default.Update(ctx, identity, obj, opt...)
}

func (d *routeStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())

	return d.Default.Delete(ctx, identity, opt...)
}

func (d *routeStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())

	return d.Default.Get(ctx, identity, opt...)
}

func (d *routeStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())

	return d.Default.List(ctx, identity, opt...)
}
