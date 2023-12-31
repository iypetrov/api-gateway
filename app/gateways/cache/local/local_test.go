package local_test

import (
	"context"
	"github.com/IliyaYavorovPetrov/api-gateway/app/common/models"
	"github.com/IliyaYavorovPetrov/api-gateway/test"
	"testing"
)

func TestAddAndGet(t *testing.T) {
	test.ClearLocalCache()
	cache := test.GetLocalCache()

	key := "key1"
	rri := models.ReqRoutingInfo{
		SourceURL:      "https://src/1",
		DestinationURL: "http://dest/1",
		MethodHTTP:     "GET",
		IsAuthNeeded:   false,
	}

	err := cache.Add(test.GetCtx(), key, rri)
	if err != nil {
		t.Errorf("Add failed: %v", err)
	}

	res, err := cache.Get(test.GetCtx(), key)
	if err != nil {
		t.Errorf("Get failed: %v", err)
	}
	if *res != rri {
		t.Errorf("expected %+v, but got %+v", rri, res)
	}
}

func TestAddAllItems(t *testing.T) {
	test.ClearLocalCache()
	cache := test.GetLocalCache()

	data := map[string]models.ReqRoutingInfo{
		"key1": {
			SourceURL:      "https://src/1",
			DestinationURL: "http://dest/1",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
		"key2": {
			SourceURL:      "https://src/2",
			DestinationURL: "http://dest/2",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
	}

	err := cache.AddAllItems(context.Background(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	for key, val := range data {
		res, err := cache.Get(test.GetCtx(), key)
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if *res != val {
			t.Errorf("expected %+v, but got %+v", val, res)
		}
	}
}

func TestGetAllKeysByPrefix(t *testing.T) {
	test.ClearLocalCache()
	cache := test.GetLocalCache()

	data := map[string]models.ReqRoutingInfo{
		"test:key:key1": {
			SourceURL:      "https://src/1",
			DestinationURL: "http://dest/1",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
		"test:key:key2": {
			SourceURL:      "https://src/2",
			DestinationURL: "http://dest/2",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
		"wrong:key:key2": models.ReqRoutingInfo{
			SourceURL:      "https://src/3",
			DestinationURL: "http://dest/3",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
	}

	err := cache.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	prefix := "test:key:"
	keys, err := cache.GetAllKeysByPrefix(test.GetCtx(), prefix)
	if err != nil {
		t.Errorf("GetAllKeysByPrefix failed: %v", err)
	}

	exp := []string{"test:key:key1", "test:key:key2"}
	for _, key := range exp {
		if !test.ContainsItem(key, keys) {
			t.Errorf("expected to contain %v, but it did not", key)
		}
	}
}

func TestGetAllItems(t *testing.T) {
	test.ClearLocalCache()
	cache := test.GetLocalCache()

	data := map[string]models.ReqRoutingInfo{
		"test:key:key1": {
			SourceURL:      "https://src/1",
			DestinationURL: "http://dest/1",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
		"test:key:key2": {
			SourceURL:      "https://src/2",
			DestinationURL: "http://dest/2",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
	}

	err := cache.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	items, err := cache.GetAllItems(test.GetCtx())
	if err != nil {
		t.Errorf("GetAllItems failed: %v", err)
	}

	if !test.MapEqual(items, data) {
		t.Errorf("expected %v, but got %v", data, items)
	}
}

func TestDelete(t *testing.T) {
	test.ClearLocalCache()
	cache := test.GetLocalCache()

	data := map[string]models.ReqRoutingInfo{
		"test:key:key1": {
			SourceURL:      "https://src/1",
			DestinationURL: "http://dest/1",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
		"test:key:key2": {
			SourceURL:      "https://src/2",
			DestinationURL: "http://dest/2",
			MethodHTTP:     "GET",
			IsAuthNeeded:   false,
		},
	}

	delKey := "test:key:key1"

	err := cache.AddAllItems(test.GetCtx(), data)
	if err != nil {
		t.Errorf("AddAllItems failed: %v", err)
	}

	err = cache.Delete(test.GetCtx(), delKey)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = cache.Get(test.GetCtx(), delKey)
	if err == nil {
		t.Errorf("expected error for Get after deletion, but got nil")
	}
}
