package rest

import (
	"context"
	"fmt"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

type internalStore struct {
	Schema store.SchemaHolder
	Store  store.Store
	Log    logger.Logger
}

func internalFactory(data store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &internalStore{
			Schema: schema,
			Store:  data,
			Log:    logger.Factory("server internal"),
		}

		return client, nil
	}
}

func (d *internalStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())

	// initialize metadata
	original := d.Schema.ObjectForKind(obj.Metadata().Kind())
	if original == nil {
		return nil, fmt.Errorf("unknown kind %s", obj.Metadata().Kind())
	}

	// update external
	externalHolder := original.(store.ExternalHolder)
	if externalHolder != nil {
		externalHolder.ExternalInternalSet(obj.(store.ExternalHolder).ExternalInternal())
	}

	ms := original.Metadata().(store.MetaSetter)

	ms.SetIdentity(store.ObjectIdentityFactory())
	ms.SetCreated(utils.Timestamp())

	return d.Store.Create(ctx, original, opt...)
}

func (d *internalStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())
	// read the real object
	original, err := d.Store.Get(ctx, identity)

	// if doesn't exist return error
	if err != nil {
		return nil, err
	}
	if original == nil {
		return nil, constants.ErrNoSuchObject
	}

	// update external
	externalHolder := original.(store.ExternalHolder)
	if externalHolder != nil && obj != nil {
		objExternalHolder := obj.(store.ExternalHolder)
		if objExternalHolder != nil {
			externalHolder.ExternalInternalSet(
				objExternalHolder.ExternalInternal())
		}
	}

	ms := original.Metadata().(store.MetaSetter)
	ms.SetUpdated(utils.Timestamp())

	return d.Store.Update(ctx, identity, original, opt...)
}

func (d *internalStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *internalStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())

	return d.Store.Get(ctx, identity, opt...)
}

func (d *internalStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())

	return d.Store.List(ctx, identity, opt...)
}
