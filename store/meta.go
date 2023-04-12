package store

type Meta interface {
	Kind() string
	Identity() ObjectIdentity
	Created() string
	Updated() string
}

type MetaSetter interface {
	SetKind(string)
	SetIdentity(ObjectIdentity)
	SetCreated(string)
	SetUpdated(string)
}

type MetaHolder interface {
	Metadata() Meta
}

type metaWrapper struct {
	Kind_     *string         `json:"kind"`
	Identity_ *ObjectIdentity `json:"identity"`
	Created_  *string         `json:"created"`
	Updated_  *string         `json:"updated"`
}

func (m *metaWrapper) Kind() string {
	return *m.Kind_
}

func (m *metaWrapper) Created() string {
	return *m.Created_
}

func (m *metaWrapper) Updated() string {
	return *m.Updated_
}

func (m *metaWrapper) Identity() ObjectIdentity {
	return *m.Identity_
}

func (m *metaWrapper) SetKind(kind string) {
	m.Kind_ = &kind
}

func (m *metaWrapper) SetIdentity(identity ObjectIdentity) {
	m.Identity_ = &identity
}

func (m *metaWrapper) SetCreated(created string) {
	m.Created_ = &created
}

func (m *metaWrapper) SetUpdated(updated string) {
	m.Updated_ = &updated
}

func MetaFactory(kind string) Meta {
	emptyIdentity := ObjectIdentityFactory()
	emptyString1 := ""
	emptyString2 := ""
	mw := metaWrapper{
		Kind_:     &kind,
		Identity_: &emptyIdentity,
		Created_:  &emptyString1,
		Updated_:  &emptyString2,
	}

	return &mw
}
