package logger

import (
	"context"

	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

type loggerStore struct {
	Schema store.SchemaHolder
	Store  store.Store
	Logger Logger
}

func StoreFactory(module string, st store.Store) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &loggerStore{
			Schema: schema,
			Store:  st,
			Logger: Factory(module),
		}

		return client, nil
	}
}

func (d *loggerStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	ret, err := d.Store.Create(ctx, obj, opt...)
	d.Logger.Object("ret", ret)
	if err != nil {
		d.Logger.Printf(err.Error())
	}
	return ret, err
}

func (d *loggerStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	ret, err := d.Store.Update(ctx, identity, obj, opt...)
	d.Logger.Object("ret", ret)
	if err != nil {
		d.Logger.Printf(err.Error())
	}
	return ret, err
}

func (d *loggerStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	err := d.Store.Delete(ctx, identity, opt...)
	if err != nil {
		d.Logger.Printf(err.Error())
	}
	return err
}

func (d *loggerStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	ret, err := d.Store.Get(ctx, identity, opt...)
	d.Logger.Object("ret", ret)
	if err != nil {
		d.Logger.Printf(err.Error())
	}
	return ret, err
}

func (d *loggerStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	ret, err := d.Store.List(ctx, identity, opt...)
	d.Logger.Object("ret", ret)
	if err != nil {
		d.Logger.Printf(err.Error())
	}
	return ret, err
}
