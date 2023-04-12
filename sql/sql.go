package sql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
)

var log = logger.Factory("sql")

type _ConnectionMaker func(*sqlStore) (*sql.DB, error)

type sqlStore struct {
	Schema         store.SchemaHolder
	DB             *sql.DB
	MakeConnection _ConnectionMaker
}

func SqliteConnection(path string) _ConnectionMaker {
	return func(d *sqlStore) (*sql.DB, error) {
		return sql.Open("sqlite3", path)
	}
}

func MySqlConnection(path string) _ConnectionMaker {
	return func(d *sqlStore) (*sql.DB, error) {
		log.Printf("mysql connection %s", path)
		// username:password@tcp(127.0.0.1:3306)/test

		return sql.Open("mysql", path)
	}
}

func (d *sqlStore) TestConnection() error {
	if d.DB != nil {
		if d.DB.Ping() == nil {
			return nil
		}
	}

	var err error
	d.DB, err = d.MakeConnection(d)
	if err != nil {
		return err
	}

	err = d.DB.Ping()
	if err != nil {
		return err
	}

	return d.prepareTables()
}

func Factory(connector _ConnectionMaker) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		client := &sqlStore{
			Schema:         schema,
			MakeConnection: connector,
			DB:             nil,
		}

		return client, nil
	}
}

func (d *sqlStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	log.Printf("create %s", obj.PrimaryKey())

	var err error
	copt := options.CommonOptionHolderFactory()
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	lk := strings.ToLower(obj.Metadata().Kind())
	path := fmt.Sprintf("%s/%s", lk, obj.PrimaryKey())
	existing, _ := d.Get(ctx, store.ObjectIdentity(path))
	if existing != nil {
		return nil, constants.ErrObjectExists
	}

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	err = d.setIdentity(
		obj.Metadata().Identity().Path(),
		obj.PrimaryKey(),
		obj.Metadata().Kind())
	if err != nil {
		return nil, err
	}

	err = d.setObject(obj.PrimaryKey(), obj.Metadata().Kind(), obj)
	if err != nil {
		return nil, err
	}

	return obj.Clone(), nil
}

func (d *sqlStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	log.Printf("update %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return nil, constants.ErrNoSuchObject
	}

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	// log.Object("existing", existing)

	err = d.removeIdentity(existing.Metadata().Identity().Path())
	if err != nil {
		log.Printf("%s", err)
	}

	err = d.setIdentity(obj.Metadata().Identity().Path(),
		obj.PrimaryKey(), obj.Metadata().Kind())

	if err != nil {
		return nil, err
	}

	err = d.removeObject(existing.PrimaryKey(), existing.Metadata().Kind())
	if err != nil {
		return nil, err
	}

	err = d.setObject(obj.PrimaryKey(), obj.Metadata().Kind(), obj)
	if err != nil {
		return nil, err
	}

	return obj.Clone(), nil
}

func (d *sqlStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	log.Printf("delete %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return err
		}
	}

	existing, _ := d.Get(ctx, identity)
	if existing == nil {
		return constants.ErrNoSuchObject
	}

	err = d.TestConnection()
	if err != nil {
		return err
	}

	err = d.removeIdentity(existing.Metadata().Identity().Path())
	if err != nil {
		return err
	}

	return d.removeObject(existing.PrimaryKey(), existing.Metadata().Kind())
}

func (d *sqlStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	log.Printf("get %s", identity.Path())

	var err error
	copt := options.CommonOptionHolderFactory()
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	pkey, typ, err := d.getIdentity(identity.Path())
	if err == nil {
		return d.getObject(pkey, typ)
	}

	tokens := strings.Split(identity.Path(), "/")
	if len(tokens) == 2 {
		return d.getObject(tokens[1], tokens[0])
	}

	return nil, constants.ErrNoSuchObject
}

func (d *sqlStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	log.Printf("list %s", identity)

	if len(identity.Key()) > 0 {
		return nil, constants.ErrInvalidPath
	}

	var err error
	copt := options.CommonOptionHolderFactory()
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	err = d.TestConnection()
	if err != nil {
		return nil, err
	}

	query := `SELECT Object FROM Objects
		WHERE Type = ?`

	// pkey filter
	if copt.KeyFilter != nil {
		query = query + fmt.Sprintf(
			" AND Pkey IN ('%s')",
			strings.Join(*copt.KeyFilter, "', '"))
	}

	// prop filter
	if copt.PropFilter != nil {
		obj := d.Schema.ObjectForKind(identity.Type())
		if obj == nil {
			return nil, constants.ErrNoSuchObject
		}
		if utils.ObjectPath(obj, copt.PropFilter.Key) == nil {
			return nil, constants.ErrInvalidFilter
		}

		query = query + fmt.Sprintf(
			" AND json_extract(Object, '$.%s') = '%s'",
			copt.PropFilter.Key, copt.PropFilter.Value)
	}

	if len(copt.OrderBy) > 0 {
		query = fmt.Sprintf(`SELECT Object
			FROM Objects
			WHERE Type = ?
			ORDER BY json_extract(Object, '$.%s')`, copt.OrderBy)

		if copt.OrderIncremental {
			query = query + " ASC"
		} else {
			query = query + " DESC"
		}
	}

	if copt.PageSize > 0 {
		query = query + fmt.Sprintf(" LIMIT %d", copt.PageSize)
	}

	if copt.PageOffset > 0 {
		query = query + fmt.Sprintf(" OFFSET %d", copt.PageOffset)
	}

	log.Printf(query)

	rows, err := d.DB.Query(query, identity.Type())
	if err != nil {
		return nil, err
	}

	res := d.parseObjectRows(rows, identity.Type())
	rows.Close()

	return res, nil
}

func (d *sqlStore) prepareTables() error {
	// log.Printf("preparing tables")

	create := `
		CREATE TABLE IF NOT EXISTS IdIndex (
		Path VARCHAR(25) NOT NULL PRIMARY KEY,
		Pkey NVARCHAR(50) NOT NULL,
		Type VARCHAR(25) NOT NULL);`

	_, err := d.DB.Exec(create)
	if err != nil {
		return err
	}

	create = `
		CREATE TABLE IF NOT EXISTS Objects (
		Pkey NVARCHAR(50) NOT NULL,
		Type VARCHAR(25) NOT NULL,
		Object JSON,
		PRIMARY KEY (Pkey,Type));`

	_, err = d.DB.Exec(create)
	if err != nil {
		return err
	}

	return nil
}

func (d *sqlStore) getIdentity(path string) (string, string, error) {
	row := d.DB.QueryRow("SELECT Pkey, Type FROM IdIndex WHERE Path=?", path)

	var pkey string = ""
	var typ string = ""

	err := row.Scan(&pkey, &typ)
	return pkey, typ, err
}

func (d *sqlStore) setIdentity(path string, pkey string, typ string) error {
	// log.Printf("setting identity %s %s %s", path, pkey, typ)

	query := ""
	_, _, err := d.getIdentity(path)

	if err == nil {
		query = `update IdIndex set Pkey=?, Type=? where Path = ?`
	} else {
		query = `insert into IdIndex (Pkey, Type, Path) values (?, ?, ?)`
	}

	_, err = d.DB.Exec(query, pkey, strings.ToLower(typ), path)

	return err
}

func (d *sqlStore) removeIdentity(path string) error {
	query := "DELETE FROM IdIndex WHERE Path = ?"

	_, err := d.DB.Exec(query, path)
	return err
}

func (d *sqlStore) getObject(pkey string, typ string) (store.Object, error) {
	// log.Printf("getting %s %s", pkey, typ)

	return d.parseObjectRow(
		d.DB.QueryRow("SELECT Object FROM Objects WHERE Pkey=? AND Type=?",
			pkey, strings.ToLower(typ)), typ)
}

func (d *sqlStore) setObject(pkey string, typ string, obj store.Object) error {
	query := ""
	_, err := d.getObject(pkey, typ)
	if err == nil {
		query = `update Objects set Object=@obj where Pkey = @pkey AND Type = @typ`
	} else {
		query = `insert into Objects (Object, Pkey, Type) values (?, ?, ?)`
	}

	data, err := utils.Serialize(obj)

	if err != nil {
		return err
	}

	_, err = d.DB.Exec(query, string(data), pkey, strings.ToLower(typ))
	return err
}

func (d *sqlStore) removeObject(pkey string, typ string) error {
	query := "DELETE FROM Objects WHERE Pkey = ? AND Type = ?"

	_, err := d.DB.Exec(query, pkey, strings.ToLower(typ))

	return err
}

func (d *sqlStore) parseObjectRow(row *sql.Row, typ string) (store.Object, error) {
	var data string = ""

	err := row.Scan(&data)

	if err != nil {
		// log.Fatal(err)
		return nil, err
	}

	return utils.UnmarshalObject([]byte(data), d.Schema, typ)
}

func (d *sqlStore) parseObjectRows(rows *sql.Rows, typ string) store.ObjectList {
	res := store.ObjectList{}
	for rows.Next() {
		var data string = ""
		err := rows.Scan(&data)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		ret, err := utils.UnmarshalObject([]byte(data), d.Schema, typ)
		if err != nil {
			log.Fatal(err)
			return nil
		}

		res = append(res, ret)
	}

	return res
}
