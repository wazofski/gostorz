package react

import (
	"context"
	"fmt"
	"log"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
)

type Action int
type Callback func(store.Object, store.Store) error

const (
	ActionCreate Action = 1
	ActionUpdate Action = 2
	ActionDelete Action = 3
)

type reactStore struct {
	Schema           store.SchemaHolder
	Store            store.Store
	Log              logger.Logger
	CallbackRegistry map[string]map[Action]Callback
}

type _Register struct {
	Kind     string
	Action   Action
	Callback Callback
}

func Subscribe(typ string, action Action, callback Callback) _Register {
	if action < 1 || action > 3 {
		log.Fatalf("invalid action %d", action)
	}

	return _Register{
		Kind:     typ,
		Action:   action,
		Callback: callback,
	}
}

func ReactFactory(data store.Store, callbacks ..._Register) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &reactStore{
			Schema:           schema,
			Store:            data,
			Log:              logger.Factory("react"),
			CallbackRegistry: make(map[string]map[Action]Callback),
		}

		for _, c := range callbacks {
			proto := schema.ObjectForKind(c.Kind)
			if proto == nil {
				continue
			}

			_, ok := client.CallbackRegistry[proto.Metadata().Kind()]
			if !ok {
				client.CallbackRegistry[proto.Metadata().Kind()] = make(map[Action]Callback)
			}

			_, ok = client.CallbackRegistry[proto.Metadata().Kind()][c.Action]
			if ok {
				return nil, fmt.Errorf("callback for %s %d already set", c.Kind, c.Action)
			}

			client.CallbackRegistry[proto.Metadata().Kind()][c.Action] = c.Callback
		}

		return client, nil
	}
}

func (d *reactStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("create %s", obj.PrimaryKey())
	err := d.runCallback(obj, ActionCreate)
	if err != nil {
		return nil, err
	}

	return d.Store.Create(ctx, obj, opt...)
}

func (d *reactStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	d.Log.Printf("update %s", identity.Path())
	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return nil, constants.ErrNoSuchObject
	}

	err := d.runCallback(existing, ActionUpdate)
	if err != nil {
		return nil, err
	}

	return d.Store.Update(ctx, identity, obj, opt...)
}

func (d *reactStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	d.Log.Printf("delete %s", identity.Path())
	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return constants.ErrNoSuchObject
	}

	err := d.runCallback(existing, ActionDelete)
	if err != nil {
		return err
	}

	return d.Store.Delete(ctx, identity, opt...)
}

func (d *reactStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	d.Log.Printf("get %s", identity.Path())
	return d.Store.Get(ctx, identity, opt...)
}

func (d *reactStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	d.Log.Printf("list %s", identity.Type())
	return d.Store.List(ctx, identity, opt...)
}

func (d *reactStore) runCallback(obj store.Object, action Action) error {
	_, ok := d.CallbackRegistry[obj.Metadata().Kind()]
	if !ok {
		return nil
	}

	_, ok = d.CallbackRegistry[obj.Metadata().Kind()][action]
	if !ok {
		return nil
	}

	return d.CallbackRegistry[obj.Metadata().Kind()][action](obj, d)
}
