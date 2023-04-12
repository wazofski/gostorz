package client

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/wazofski/storz/internal/constants"
	"github.com/wazofski/storz/internal/logger"
	"github.com/wazofski/storz/rest"
	"github.com/wazofski/storz/store"
	"github.com/wazofski/storz/store/options"
	"github.com/wazofski/storz/utils"
)

var log = logger.Factory("rest client")

type restStore struct {
	BaseURL     *url.URL
	Schema      store.SchemaHolder
	MakeRequest requestMaker
	Headers     []headerOption
}

type requestMaker func(path *url.URL, content []byte, method string, headers map[string]string) ([]byte, error)

type restOptions struct {
	options.CommonOptionHolder
	Headers map[string]string
}

func newRestOptions(d *restStore) restOptions {
	res := restOptions{
		CommonOptionHolder: options.CommonOptionHolderFactory(),
		Headers:            make(map[string]string),
	}

	for _, h := range d.Headers {
		h.ApplyFunction()(&res)
	}

	return res
}

func (d *restOptions) CommonOptions() *options.CommonOptionHolder {
	return &d.CommonOptionHolder
}

func Factory(serviceUrl string, headers ...headerOption) store.Factory {
	return func(schema store.SchemaHolder) (store.Store, error) {
		URL, err := url.Parse(serviceUrl)
		if err != nil {
			return nil, fmt.Errorf("invalid URL: %s", err)
		}

		client := &restStore{
			BaseURL:     URL,
			Schema:      schema,
			MakeRequest: makeHttpRequest,
			Headers:     headers,
		}

		log.Printf("initialized: %s", serviceUrl)
		return client, nil
	}
}

func makeHttpRequest(path *url.URL, content []byte, requestType string, headers map[string]string) ([]byte, error) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	req, err := http.NewRequest(requestType, path.String(), strings.NewReader(string(content)))
	if err != nil {
		return nil, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Close = true

	// req.ContentLength = contentLength
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	rd, err := utils.ReadStream(resp.Body)

	defer resp.Body.Close()

	if err != nil {
		return rd, err
	}

	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return rd, fmt.Errorf("http %d", resp.StatusCode)
	}

	return rd, nil
}

func processRequest(
	client *restStore,
	requestUrl *url.URL,
	content []byte,
	method string,
	headers map[string]string) ([]byte, error) {

	reqId := uuid.New().String()
	requestUrl.Path = strings.ReplaceAll(requestUrl.Path, "//", "/")
	origin := strings.ReplaceAll(requestUrl.String(), requestUrl.Path, "")
	headers["Origin"] = strings.ReplaceAll(origin, requestUrl.RawQuery, "")
	headers["X-Request-ID"] = reqId
	headers["Content-Type"] = "application/json"
	headers["X-Requested-With"] = "XMLHttpRequest"

	log.Printf("%s %s", strings.ToLower(method), requestUrl)
	// log.Printf("X-Request-ID %s", reqId)

	data, err := client.MakeRequest(requestUrl, content, method, headers)
	cerr := errorCheck(data)
	if err == nil {
		err = cerr
	} else if cerr != nil {
		err = fmt.Errorf("%s %s", err, cerr)
	}

	if err != nil {
		// log.Println(err)
		if len(content) > 0 {
			var js interface{}
			if json.Unmarshal([]byte(content), &js) == nil {
				log.Object("request content", js)
			} else {
				log.Printf("request content: %s", content)
			}
		}
		if len(data) > 0 {
			log.Printf("response content: %s", string(data))
		}
		return nil, err
	}

	return data, err
}

func errorCheck(response []byte) error {
	str := string(response)
	if len(str) == 0 {
		return nil
	}

	um := make(map[string]interface{})

	err := json.Unmarshal(response, &um)
	if err == nil {
		if v, found := um["errors"]; found {
			return errors.New(v.([]interface{})[0].(string))
		}
		if v, found := um["error"]; found {
			m := v.(map[string]interface{})

			return fmt.Errorf("%v %s",
				m["internal_code"],
				m["internal"])
		}
	}

	return nil
}

func makePathForType(baseUrl *url.URL, obj store.Object) *url.URL {
	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseUrl, strings.ToLower(obj.Metadata().Kind())))
	return u
}

func removeTrailingSlash(val string) string {
	if strings.HasSuffix(val, "/") {
		return val[:len(val)-1]
	}
	return val
}

func makePathForIdentity(baseUrl *url.URL, identity store.ObjectIdentity, params string) *url.URL {
	if len(params) > 0 {
		path := fmt.Sprintf("%s/%s?%s",
			baseUrl,
			removeTrailingSlash(identity.Path()),
			params)

		// log.Printf(`made path %s # %s # %s`, path, identity.Path(), string(identity))

		u, _ := url.ParseRequestURI(path)
		return u
	}

	u, _ := url.Parse(fmt.Sprintf("%s/%s", baseUrl, identity.Path()))
	return u
}

func toBytes(obj interface{}) []byte {
	if obj == nil {
		return []byte{}
	}

	jsn, _ := json.Marshal(obj)

	return []byte(string(jsn))
}

func listParameters(ropt restOptions) string {
	opt := ropt.CommonOptions()

	q := url.Values{}
	if len(opt.OrderBy) > 0 {
		q.Add(rest.OrderByArg, opt.OrderBy)
		q.Add(rest.IncrementalArg, strconv.FormatBool(opt.OrderIncremental))
	}

	if opt.PageOffset > 0 {
		q.Add(rest.PageOffsetArg, fmt.Sprintf("%d", opt.PageOffset))
	}

	if opt.PageSize > 0 {
		q.Add(rest.PageSizeArg, fmt.Sprintf("%d", opt.PageSize))
	}

	if opt.PropFilter != nil {
		content, err := json.Marshal(opt.PropFilter)
		if err != nil {
			log.Fatal(err)
		}

		if len(content) > 0 {
			q.Add(rest.PropFilterArg, string(content))
		}
	}

	if opt.KeyFilter != nil {
		content, err := json.Marshal(opt.KeyFilter)
		if err != nil {
			log.Fatal(err)
		}

		if len(content) > 0 {
			q.Add(rest.KeyFilterArg, string(content))
		}
	}

	return q.Encode()
}

func (d *restStore) Create(
	ctx context.Context,
	obj store.Object,
	opt ...options.CreateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	log.Printf("get %s", obj.Metadata().Identity().Path())

	copt := newRestOptions(d)
	var err error
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	data, err := stripSerialize(obj)
	if err != nil {
		return nil, err
	}

	data, err = processRequest(d,
		makePathForType(d.BaseURL, obj),
		data,
		http.MethodPost,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	clone := obj.Clone()
	err = json.Unmarshal(data, &clone)
	if err != nil {
		log.Printf(string(data))
		clone = nil
	}

	return clone, err
}

func (d *restStore) Update(
	ctx context.Context,
	identity store.ObjectIdentity,
	obj store.Object,
	opt ...options.UpdateOption) (store.Object, error) {

	if obj == nil {
		return nil, constants.ErrObjectNil
	}

	log.Printf("update %s", identity.Path())

	copt := newRestOptions(d)
	var err error
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	data, err := stripSerialize(obj)
	if err != nil {
		return nil, err
	}

	data, err = processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		data,
		http.MethodPut,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	clone := obj.Clone()
	err = json.Unmarshal(data, &clone)

	return clone, err
}

func (d *restStore) Delete(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.DeleteOption) error {

	log.Printf("delete %s", identity.Path())

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return err
		}
	}

	_, err = processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		[]byte{},
		http.MethodDelete,
		copt.Headers)

	return err
}

func (d *restStore) Get(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.GetOption) (store.Object, error) {

	log.Printf("get %s", identity.Path())

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	resp, err := processRequest(d,
		makePathForIdentity(d.BaseURL, identity, ""),
		[]byte{},
		http.MethodGet,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	tp := identity.Type()
	if tp == "id" {
		tp = utils.ObjeectKind(resp)
	}

	return utils.UnmarshalObject(resp, d.Schema, tp)
}

func (d *restStore) List(
	ctx context.Context,
	identity store.ObjectIdentity,
	opt ...options.ListOption) (store.ObjectList, error) {

	log.Printf("list %s", identity)

	var err error
	copt := newRestOptions(d)
	for _, o := range opt {
		err = o.ApplyFunction()(&copt)
		if err != nil {
			return nil, err
		}
	}

	params := listParameters(copt)
	path := makePathForIdentity(d.BaseURL, identity, params)
	res, err := processRequest(
		d,
		path,
		[]byte{},
		http.MethodGet,
		copt.Headers)

	if err != nil {
		return nil, err
	}

	parsed := []*json.RawMessage{}
	err = json.Unmarshal(res, &parsed)
	if err != nil {
		return nil, err
	}

	marshalledResult := store.ObjectList{}
	if len(parsed) == 0 {
		return marshalledResult, nil
	}

	resource := d.Schema.ObjectForKind(utils.ObjeectKind(*parsed[0]))

	for _, r := range parsed {
		clone := resource.Clone()
		clone.UnmarshalJSON(toBytes(r))

		marshalledResult = append(marshalledResult, clone)
	}

	return marshalledResult, nil
}

type strippedObject struct {
	External map[string]*json.RawMessage `json:"external"`
}

func stripSerialize(object store.Object) ([]byte, error) {
	data, err := utils.Serialize(object)
	if err != nil {
		return nil, err
	}
	obj := strippedObject{}
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return json.Marshal(obj)
}
