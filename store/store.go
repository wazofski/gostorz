package store

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
	"github.com/wazofski/storz/store/options"
)

type Endpoint interface {
	Listen(int) context.CancelFunc
}

type Object interface {
	MetaHolder
	Clone() Object
	UnmarshalJSON(data []byte) error
	PrimaryKey() string
}

type ExternalHolder interface {
	ExternalInternalSet(interface{})
	ExternalInternal() interface{}
}

type ObjectList []Object
type ObjectIdentity string

func ObjectIdentityFactory() ObjectIdentity {
	id := uuid.New().String()
	id = strings.ReplaceAll(id, "-", "")
	id = id[5:25]

	return ObjectIdentity(id)
}

func (o ObjectIdentity) Path() string {
	if strings.Index(string(o), "/") > 0 {
		tok := strings.Split(string(o), "/")
		return fmt.Sprintf("%s/%s", strings.ToLower(tok[0]), tok[1])
	}

	return fmt.Sprintf("id/%s", o)
}

func (o ObjectIdentity) Type() string {
	tokens := strings.Split(o.Path(), "/")
	return strings.ToLower(tokens[0])
}

func (o ObjectIdentity) Key() string {
	tokens := strings.Split(o.Path(), "/")
	if len(tokens) > 1 {
		return tokens[1]
	}
	return ""
}

type Store interface {
	Get(context.Context, ObjectIdentity, ...options.GetOption) (Object, error)
	List(context.Context, ObjectIdentity, ...options.ListOption) (ObjectList, error)
	Create(context.Context, Object, ...options.CreateOption) (Object, error)
	Delete(context.Context, ObjectIdentity, ...options.DeleteOption) error
	Update(context.Context, ObjectIdentity, Object, ...options.UpdateOption) (Object, error)
}

type SchemaHolder interface {
	ObjectForKind(kind string) Object
	// ObjectMethods() map[string][]string
	Types() []string
}

type Factory func(schema SchemaHolder) (Store, error)

func New(schema SchemaHolder, factory Factory) Store {
	store, err := factory(schema)
	if err != nil {
		log.Fatal(err)
	}
	return store
}
