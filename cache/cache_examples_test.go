package cache_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/corsc/go-commons/cache"
)

func ExampleClient_normalUsage() {
	// init - called once; perhaps a global variable or member variable
	cacheClient := &cache.Client{
		Storage: &cache.RedisStorage{},
	}

	// general usage
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	cacheKey := "cache.key"
	dest := &myDTO{}

	err := cacheClient.Get(ctx, cacheKey, dest, cache.BuilderFunc(func(ctx context.Context, key string, dest cache.BinaryEncoder) error {
		// logic that builds/marshals the cacheable value
		return errors.New("not implemented")
	}))

	if err != nil {
		panic(err.Error())
	}
}

func ExampleClient_httpHandler() {
	// init - called once; perhaps a global variable or member variable
	userCache := &cache.Client{
		Storage: &cache.RedisStorage{},
	}

	// the HTTP Handler
	handler := func(resp http.ResponseWriter, req *http.Request) {
		key := buildCacheKey(req)
		outputDTO := &myDTO{}

		err := userCache.Get(req.Context(), key, outputDTO, cache.BuilderFunc(func(ctx context.Context, key string, dest cache.BinaryEncoder) error {
			// logic that builds/marshals the cacheable value
			return errors.New("not implemented")
		}))

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
		_, _ = resp.Write(data)
	}

	_ = http.ListenAndServe("", http.HandlerFunc(handler))
}

func buildCacheKey(_ *http.Request) string {
	// do something clever here that uses the inputs to generated a predictable key
	return ""
}

type myDTO struct {
	Name  string
	Email string
}

// MarshalBinary implements cache.BinaryEncoder
func (m *myDTO) MarshalBinary() (data []byte, err error) {
	return json.Marshal(m)
}

// UnmarshalBinary implements cache.BinaryEncoder
func (m *myDTO) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, m)
}
