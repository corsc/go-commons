package cache_test

import (
	"context"
	"encoding/json"

	"net/http"

	"github.com/corsc/go-commons/cache"
)

func ExampleClient_normalUsage() {
	// define a cache item builder
	buildFn := func(ctx context.Context, key string, dest cache.Binary) {
		// logic that builds/marshals the cacheable value
	}

	// init - called once; perhaps a global variable or member variable
	cacheClient := &cache.Client{
		Storage: &cache.RedigoStorage{},
		Builder: cache.BuilderFunc(buildFn),
	}

	// general usage
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	cacheKey := "cache.key"
	dest := &myDTO{}

	err := cacheClient.Get(ctx, cacheKey, dest)
	if err != nil {
		panic(err.Error())
	}
}

func ExampleClient_httpHandler() {
	// define a cache item builder
	getUser := func(ctx context.Context, key string, dest cache.Binary) {
		// logic that builds/marshals the cacheable value
	}

	// init - called once; perhaps a global variable or member variable
	userCache := &cache.Client{
		Storage: &cache.RedigoStorage{},
		Builder: cache.BuilderFunc(getUser),
	}

	handler := func(resp http.ResponseWriter, req *http.Request) {
		key := buildCacheKey(req)
		outputDTO := &myDTO{}

		err := userCache.Get(req.Context(), key, outputDTO)
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		data, err := outputDTO.MarshalBinary()
		if err != nil {
			http.Error(resp, err.Error(), http.StatusInternalServerError)
			return
		}

		resp.WriteHeader(http.StatusOK)
		resp.Write(data)
	}

	_ = http.ListenAndServe("/", http.HandlerFunc(handler))
}

func buildCacheKey(req *http.Request) string {
	// do something clever here that uses the inputs to generated a predictable key
	return ""
}

type myDTO struct {
	Name  string
	Email string
}

// MarshalBinary implements cache.Binary
func (m *myDTO) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary implements cache.Binary
func (m *myDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
