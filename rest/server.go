package rest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/exp/slices"

	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

var log = logger.Factory("rest server")

const (
	PropFilterArg  = "pf"
	KeyFilterArg   = "kf"
	IncrementalArg = "inc"
	PageSizeArg    = "pageSize"
	PageOffsetArg  = "pageOffset"
	OrderByArg     = "orderBy"
)

type _HandlerFunc func(http.ResponseWriter, *http.Request)

type _Server struct {
	Schema  store.SchemaHolder
	Store   store.Store
	Context context.Context
	Router  *mux.Router
	Exposed map[string][]Action
}

func (d *_Server) Listen(port int) context.CancelFunc {
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: d.Router,
	}

	go func() {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			log.Printf("server closed")
		} else if err != nil {
			log.Printf("error listening: %s", err)
		} else {
			log.Printf("habidi dubidi")
		}
	}()

	log.Printf("listening on port %d", port)

	return func() { srv.Shutdown(context.Background()) }
}

type Action string

const (
	ActionCreate Action = http.MethodPost
	ActionUpdate Action = http.MethodPut
	ActionDelete Action = http.MethodDelete
	ActionGet    Action = http.MethodGet
)

type _TypeMethods struct {
	Kind    string
	Actions []Action
}

func TypeMethods(kind string, actions ...Action) _TypeMethods {
	return _TypeMethods{
		Kind:    kind,
		Actions: actions,
	}
}

func Server(schema store.SchemaHolder, stor store.Store, exposed ..._TypeMethods) store.Endpoint {
	server := &_Server{
		Schema:  schema,
		Store:   store.New(schema, internalFactory(stor)),
		Context: context.Background(),
		Router:  mux.NewRouter(),
		Exposed: make(map[string][]Action),
	}

	addHandler(server.Router, "/id/{id}", makeIdHandler(server))
	for _, e := range exposed {
		server.Exposed[e.Kind] = e.Actions

		addHandler(server.Router,
			fmt.Sprintf("/%s/{pkey}", strings.ToLower(e.Kind)),
			makeObjectHandler(server, e.Kind, e.Actions))
		addHandler(server.Router,
			fmt.Sprintf("/%s", strings.ToLower(e.Kind)),
			makeTypeHandler(server, e.Kind, e.Actions))
		addHandler(server.Router,
			fmt.Sprintf("/%s/", strings.ToLower(e.Kind)),
			makeTypeHandler(server, e.Kind, e.Actions))
	}

	return server
}

func addHandler(router *mux.Router, pattern string, handler _HandlerFunc) {
	router.HandleFunc(pattern, handler)
}

func makeIdHandler(server *_Server) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
		id := store.ObjectIdentity(mux.Vars(r)["id"])
		existing, _ := server.Store.Get(server.Context, id)

		var robject store.Object = nil
		if existing != nil {
			kind := existing.Metadata().Kind()
			data, err := utils.ReadStream(r.Body)
			if err == nil {
				robject, _ = utils.UnmarshalObject(data, server.Schema, kind)
			}

			// method validation
			objMethods := server.Exposed[kind]
			if objMethods == nil || !slices.Contains(objMethods, Action(r.Method)) {
				reportError(w,
					constants.ErrInvalidMethod,
					http.StatusMethodNotAllowed)
				return
			}
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeObjectHandler(server *_Server, t string, methods []Action) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)
		var robject store.Object = nil
		id := store.ObjectIdentity(strings.ToLower(t) + "/" + mux.Vars(r)["pkey"])
		data, err := utils.ReadStream(r.Body)
		if err == nil {
			robject, _ = utils.UnmarshalObject(data, server.Schema, t)
		}

		// method validation
		if !slices.Contains(methods, Action(r.Method)) {
			reportError(w,
				constants.ErrInvalidMethod,
				http.StatusMethodNotAllowed)
			return
		}

		server.handlePath(w, r, id, robject)
	}
}

func makeTypeHandler(server *_Server, t string, methods []Action) _HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		prepResponse(w, r)

		// method validation
		if !slices.Contains(methods, Action(r.Method)) {
			reportError(w,
				constants.ErrInvalidMethod,
				http.StatusMethodNotAllowed)
			return
		}

		switch r.Method {
		case http.MethodGet:
			opts := []options.ListOption{}

			vals := r.URL.Query()
			filter, ok := vals[PropFilterArg]
			if ok {
				flt := options.PropFilterSetting{}
				err := json.Unmarshal([]byte(filter[0]), &flt)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
					return
				}
				opts = append(opts, options.PropFilter(flt.Key, flt.Value))
			}

			keyFilter, ok := vals[KeyFilterArg]
			if ok {
				flt := options.KeyFilterSetting{}
				err := json.Unmarshal([]byte(keyFilter[0]), &flt)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
					return
				}
				opts = append(opts, options.KeyFilter(flt...))
			}

			pageSize, ok := vals[PageSizeArg]
			if ok {
				ps, _ := strconv.Atoi(pageSize[0])
				opts = append(opts, options.PageSize(ps))
			}

			pageOffset, ok := vals[PageOffsetArg]
			if ok {
				ps, _ := strconv.Atoi(pageOffset[0])
				opts = append(opts, options.PageOffset(ps))
			}

			orderBy, ok := vals[OrderByArg]
			if ok {
				ob := orderBy[0]
				opts = append(opts, options.OrderBy(ob))
			}

			orderInc, ok := vals[IncrementalArg]
			if ok {
				ob := true
				err := json.Unmarshal([]byte(orderInc[0]), &ob)
				if err != nil {
					reportError(w, err, http.StatusBadRequest)
					return
				}
				if !ob {
					opts = append(opts, options.OrderDescending())
				}
			}

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
				writeResponse(w, resp)
			}
		case http.MethodPost:
			data, err := utils.ReadStream(r.Body)
			if err != nil {
				reportError(w,
					err,
					http.StatusBadRequest)
				return
			}

			robject, err := utils.UnmarshalObject(data, server.Schema, t)
			if err != nil {
				reportError(w,
					err,
					http.StatusBadRequest)
				return
			}

			server.handlePath(w, r, store.ObjectIdentity(t+"/"), robject)
		default:
			reportError(w,
				constants.ErrInvalidMethod,
				http.StatusMethodNotAllowed)
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
	case http.MethodPost:
		ret, err = d.Store.Create(d.Context, object)
		if err != nil {
			reportError(w, err, http.StatusNotAcceptable)
			return
		}
	case http.MethodPut:
		ret, err = d.Store.Update(d.Context, identity, object)
		if err != nil {
			reportError(w, err, http.StatusNotAcceptable)
			return
		}
	case http.MethodDelete:
		err = d.Store.Delete(d.Context, identity)
		if err != nil {
			reportError(w, err, http.StatusNotFound)
			return
		}
	}

	if err == nil && ret != nil {
		resp, _ := json.Marshal(ret)
		writeResponse(w, resp)
	}
}

func reportError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func writeResponse(w http.ResponseWriter, data []byte) {
	w.Write(data)
}

func prepResponse(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", strings.ToLower(r.Method), r.URL)
	w.Header().Add("Content-Type", "application/json")
}
