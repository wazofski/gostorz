package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/gorilla/mux"

	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

var log = logger.Factory("browser")

type _HandlerFunc func(http.ResponseWriter, *http.Request)

type _Server struct {
	Schema  store.SchemaHolder
	Store   store.Store
	Context context.Context
	Router  *mux.Router
}

func (d *_Server) Listen(port int) context.CancelFunc {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: d.Router,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("server closed\n")
		} else if err != nil {
			log.Printf("error listening: %s\n", err)
		} else {
			log.Printf("habidi dubidi")
		}
	}()

	log.Printf("listening on port %d", port)

	return func() { srv.Shutdown(context.Background()) }
}

func Server(schema store.SchemaHolder, store store.Store) store.Endpoint {
	server := &_Server{
		Schema:  schema,
		Store:   store,
		Context: context.Background(),
		Router:  mux.NewRouter(),
	}

	addHandler(server.Router, "/", makeIndexHandler(server))
	addHandler(server.Router, "/id/{id}", makeIdHandler(server))
	for _, k := range schema.Types() {
		addHandler(server.Router,
			fmt.Sprintf("/%s/{pkey}", strings.ToLower(k)),
			makeObjectHandler(server, k))
		addHandler(server.Router,
			fmt.Sprintf("/%s", strings.ToLower(k)),
			makeTypeHandler(server, k))
		addHandler(server.Router,
			fmt.Sprintf("/%s/", strings.ToLower(k)),
			makeTypeHandler(server, k))
	}

	return server
}

func addHandler(router *mux.Router, pattern string, handler _HandlerFunc) {
	// log.Printf("serving %s", pattern)
	router.HandleFunc(pattern, handler)
}

func makeIdHandler(server *_Server) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
		id := store.ObjectIdentity(mux.Vars(r)["id"])
		var robject store.Object = nil
		data, err := utils.ReadStream(r.Body)
		if err == nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema, utils.ObjeectKind(data))
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeIndexHandler(server *_Server) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)

		m := make(map[string]string)
		for _, t := range server.Schema.Types() {
			m[t] = strings.ToLower(t)
		}

		w.Write(render("templates/index.html", m))
	}
}

func makeObjectHandler(server *_Server, t string) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
		var robject store.Object = nil
		id := store.ObjectIdentity(strings.ToLower(t) + "/" + mux.Vars(r)["pkey"])
		data, err := utils.ReadStream(r.Body)
		if err == nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema, t)
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeTypeHandler(server *_Server, t string) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)

		opts := []options.ListOption{}

		ret, err := server.Store.List(
			server.Context,
			store.ObjectIdentity(
				fmt.Sprintf("%s/", strings.ToLower(t))),
			opts...)

		if err != nil {
			reportError(w, err, http.StatusBadRequest)
			return
		} else if ret != nil {
			resp, _ := json.Marshal(ret)
			writeResponse(w, t+" objects", string(resp))
		}
	}
}

func (d *_Server) handlePath(
	w http.ResponseWriter,
	r *http.Request,
	identity store.ObjectIdentity,
	object store.Object) {

	var ret store.Object = nil
	var err error = nil
	switch r.Method {
	case http.MethodGet:
		ret, err = d.Store.Get(d.Context, identity)
		if err != nil {
			reportError(w, err, http.StatusNotFound)
			return
		}
	}

	if err == nil && ret != nil {
		resp, _ := json.Marshal(ret)
		writeResponse(w, identity.Path(), string(resp))
	}
}

func reportError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

type _Page struct {
	Title string
	Json  string
}

func writeResponse(w http.ResponseWriter, title, data string) {
	w.Write(render("templates/base.html",
		_Page{
			Title: title,
			Json:  data,
		}))
}

func prepResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", strings.ToLower(r.Method), r.URL)
	w.Header().Add("Content-Type", "text/html")
}

func render(rpath string, data interface{}) []byte {
	path := fmt.Sprintf("%s/%s", utils.RuntimeDir(), rpath)

	t, err := template.ParseFiles(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBufferString("")
	err = t.Execute(buf, data)

	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
