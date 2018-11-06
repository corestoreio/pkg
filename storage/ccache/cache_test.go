package ccache

import (
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

func TestCache_DeletesAValue(t *testing.T) {
	cache := New(Configure())
	cache.Set("spice", "flow", time.Minute)
	cache.Set("worm", "sand", time.Minute)
	cache.Delete("spice")
	assert.Nil(t, cache.Get("spice"))
	assert.Exactly(t, "sand", cache.Get("worm").Value())
}

func TestCache_OnDeleteCallbackCalled(t *testing.T) {

	onDeleteFnCalled := false
	onDeleteFn := func(item *Item) {
		if item.key == "spice" {
			onDeleteFnCalled = true
		}
	}

	cache := New(Configure().OnDelete(onDeleteFn))
	cache.Set("spice", "flow", time.Minute)
	cache.Set("worm", "sand", time.Minute)

	time.Sleep(time.Millisecond * 10) // Run once to init
	cache.Delete("spice")
	time.Sleep(time.Millisecond * 10) // Wait for worker to pick up deleted items

	assert.Nil(t, cache.Get("spice"))
	assert.Exactly(t, "sand", cache.Get("worm").Value())
	assert.True(t, onDeleteFnCalled)
}

func TestCache_FetchesExpiredItems(t *testing.T) {
	cache := New(Configure())
	fn := func() (interface{}, error) { return "moo-moo", nil }

	cache.Set("beef", "moo", time.Second*-1)
	assert.Exactly(t, "moo", cache.Get("beef").Value())

	out, _ := cache.Fetch("beef", time.Second, fn)
	assert.Exactly(t, "moo-moo", out.Value())
}

func TestCache_GCsTheOldestItems(t *testing.T) {
	cache := New(Configure().ItemsToPrune(10))
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	// let the items get promoted (and added to our list)
	time.Sleep(time.Millisecond * 10)
	gcCache(cache)
	assert.Nil(t, cache.Get("9"))
	assert.Exactly(t, 10, cache.Get("10").Value())
}

func TestCache_PromotedItemsDontGetPruned(t *testing.T) {
	cache := New(Configure().ItemsToPrune(10).GetsPerPromote(1))
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	time.Sleep(time.Millisecond * 10) // run the worker once to init the list
	cache.Get("9")
	time.Sleep(time.Millisecond * 10)
	gcCache(cache)
	assert.Exactly(t, 9, cache.Get("9").Value())
	assert.Nil(t, cache.Get("10"))
	assert.Exactly(t, 11, cache.Get("11").Value())
}

func TestCache_TrackerDoesNotCleanupHeldInstance(t *testing.T) {
	cache := New(Configure().ItemsToPrune(10).Track())
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	item := cache.TrackingGet("0")
	time.Sleep(time.Millisecond * 10)
	gcCache(cache)
	assert.Exactly(t, 0, cache.Get("0").Value())
	assert.Nil(t, cache.Get("1"))
	item.Release()
	gcCache(cache)
	assert.Nil(t, cache.Get("0"))
}

func TestCache_RemovesOldestItemWhenFull(t *testing.T) {
	cache := New(Configure().MaxSize(5).ItemsToPrune(1))
	for i := 0; i < 7; i++ {
		cache.Set(strconv.Itoa(i), i, time.Minute)
	}
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, cache.Get("0"))
	assert.Nil(t, cache.Get("1"))
	assert.Exactly(t, 2, cache.Get("2").Value())
}

func TestCache_RemovesOldestItemWhenFullBySizer(t *testing.T) {
	cache := New(Configure().MaxSize(9).ItemsToPrune(2))
	for i := 0; i < 7; i++ {
		cache.Set(strconv.Itoa(i), &SizedItem{i, 2}, time.Minute)
	}
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, cache.Get("0"))
	assert.Nil(t, cache.Get("1"))
	assert.Nil(t, cache.Get("2"))
	assert.Nil(t, cache.Get("3"))
	assert.Exactly(t, 4, cache.Get("4").Value().(*SizedItem).id)
}

func TestCache_SetUpdatesSizeOnDelta(t *testing.T) {
	cache := New(Configure())
	cache.Set("a", &SizedItem{0, 2}, time.Minute)
	cache.Set("b", &SizedItem{0, 3}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 5)
	cache.Set("b", &SizedItem{0, 3}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 5)
	cache.Set("b", &SizedItem{0, 4}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 6)
	cache.Set("b", &SizedItem{0, 2}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 4)
	cache.Delete("b")
	time.Sleep(time.Millisecond * 100)
	checkSize(t, cache, 2)
}

func TestCache_ReplaceDoesNotchangeSizeIfNotSet(t *testing.T) {
	cache := New(Configure())
	cache.Set("1", &SizedItem{1, 2}, time.Minute)
	cache.Set("2", &SizedItem{1, 2}, time.Minute)
	cache.Set("3", &SizedItem{1, 2}, time.Minute)
	cache.Replace("4", &SizedItem{1, 2})
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 6)
}

func TestCache_ReplaceChangesSize(t *testing.T) {
	cache := New(Configure())
	cache.Set("1", &SizedItem{1, 2}, time.Minute)
	cache.Set("2", &SizedItem{1, 2}, time.Minute)

	cache.Replace("2", &SizedItem{1, 2})
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 4)

	cache.Replace("2", &SizedItem{1, 1})
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 3)

	cache.Replace("2", &SizedItem{1, 3})
	time.Sleep(time.Millisecond * 5)
	checkSize(t, cache, 5)
}

type SizedItem struct {
	id int
	s  int64
}

func (s *SizedItem) Size() int64 {
	return s.s
}

func checkSize(t *testing.T, cache *Cache, sz int64) {
	cache.Stop()
	assert.Exactly(t, sz, cache.size)
	cache.restart()
}

func gcCache(cache *Cache) {
	cache.Stop()
	cache.gc()
	cache.restart()
}
