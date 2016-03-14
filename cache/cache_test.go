package cache

import (
	"encoding/json"
	testconfig "gost/tests/config"
	"testing"
	"time"
)

type cacheTest struct {
	X int
	Y int
	Z int
}

func TestCache(t *testing.T) {
	const cacheExpireTime = 500 * time.Millisecond

	var cacheKeys = []string{
		MapKey("/testKey1"),
		MapKey("/testKey2"),
		MapKey("/testTheThirdKey"),
		MapKey("/thisNeedsToExpire"),
	}

	var items []cacheTest
	var cachedItems []*Cache
	var expiringItem *Cache

	testconfig.InitTestsDatabase()
	StartCachingSystem(cacheExpireTime)
	defer StopCachingSystem()

	items = testInitItems(t)

	testFetchInexistentCache(t, cacheKeys[0])
	cachedItems, expiringItem = testAddingToCache(t, items, cacheKeys)
	testFetchingFromCache(t, cachedItems)
	testRemovingFromCache(t, cachedItems)
	testFetchInexistentCache(t, cacheKeys[1])
	testExpiringItem(t, expiringItem, cacheExpireTime)
}

func testExpiringItem(t *testing.T, expiringItem *Cache, cacheExpireTime time.Duration) {
	t.Log("Testing the expired cache invalidation system")

	expiringKey := expiringItem.Key

	time.Sleep(1000 * time.Millisecond)

	_, err := QueryByKey(expiringKey)

	if err == nil || err != ErrKeyInvalidated {
		t.Fatal("The cache items did not properly expire")
	}
}

func testFetchInexistentCache(t *testing.T, mockQuery string) {
	t.Log("Testing the cache querying system with inexistent or invalid data")

	// Will never be added
	data, _ := QueryByKey("keySFAFSAGKAGHAJSKfhaskfhaskf")
	if data != nil {
		t.Fatal("Unexpected output from cache")
	}

	// Will be added later during the test
	data, _ = QueryByKey(mockQuery)
	if data != nil {
		t.Fatal("Unexpected output from cache")
	}
}

func testFetchingFromCache(t *testing.T, cachedItems []*Cache) {
	t.Log("Testing the cache querying system with valid data")

	var q1 *Cache
	var q2 *Cache
	var q3 *Cache
	i := 0

	for i < 2 {
		q1, _ = QueryByKey(cachedItems[0].Key)
		q2, _ = QueryByKey(cachedItems[1].Key)
		q3, _ = QueryByKey(cachedItems[2].Key)

		if q1 == nil || q2 == nil || q3 == nil {
			t.Fatal("Cache didn't properly return test items")
		}

		i++
	}

	if q1.Key != cachedItems[0].Key || q2.Key != cachedItems[1].Key || q3.Key != cachedItems[2].Key {
		t.Fatal("Wrong cache values were returned")
	}
}

func testAddingToCache(t *testing.T, items []cacheTest, cacheKeys []string) ([]*Cache, *Cache) {
	t.Log("Testing the data caching system")

	var cachedItems = make([]*Cache, 3)
	var expiringCacheItem *Cache

	var q1 []cacheTest
	var q2 []cacheTest
	var q3 []cacheTest

	// First type
	for i := 0; i < len(items); i++ {
		if items[i].X%2 == 0 {
			q1 = append(q1, items[i])
		}
	}
	j1, _ := json.MarshalIndent(q1, "", "  ")
	c1 := &Cache{
		Key:  cacheKeys[0],
		Data: j1,
	}
	cachedItems[0] = c1

	// Second type
	for i := 0; i < len(items); i++ {
		if ((items[i].X + items[i].Y) % items[i].Z) > 1 {
			q2 = append(q2, items[i])
		}
	}
	j2, _ := json.MarshalIndent(q2, "", "  ")
	c2 := &Cache{
		Key:  cacheKeys[1],
		Data: j2,
	}
	cachedItems[1] = c2

	// Third type
	for i := 0; i < len(items); i++ {
		if items[i].Z > 550 {
			q3 = append(q3, items[i])
		}
	}
	j3, _ := json.MarshalIndent(q3, "", "  ")
	c3 := &Cache{
		Key:  cacheKeys[2],
		Data: j3,
	}
	cachedItems[2] = c3

	// Expiring type
	expiringCacheItem = &Cache{
		Key:  cacheKeys[3],
		Data: j1,
	}

	expiringCacheItem.Cache()
	for _, cachedItem := range cachedItems {
		cachedItem.Cache()
	}

	time.Sleep(500 * time.Millisecond)

	return cachedItems, expiringCacheItem
}

func testRemovingFromCache(t *testing.T, cachedItems []*Cache) {
	t.Log("Testing the cache invalidation system")

	for _, it := range cachedItems {
		it.Invalidate()
	}
}

func testInitItems(t *testing.T) []cacheTest {
	var items []cacheTest

	for i := 1; i < 1000; i++ {
		items = append(items, cacheTest{
			X: i,
			Y: i * 11 / 3,
			Z: i * 3,
		})
	}

	return items
}
