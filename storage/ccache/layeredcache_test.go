package ccache

import (
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

func TestLayered_GetsANonExistantValue(t *testing.T) {
	cache := newLayered()
	assert.Nil(t, cache.Get("spice", "flow"))
}

func TestLayered_SetANewValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "a value", time.Minute)
	assert.Exactly(t, "a value", cache.Get("spice", "flow").Value())
	assert.Nil(t, cache.Get("spice", "stop"))
}

func TestLayered_SetsMultipleValueWithinTheSameLayer(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	cache.Set("spice", "must", "value-b", time.Minute)
	cache.Set("leto", "sister", "ghanima", time.Minute)
	assert.Exactly(t, "value-a", cache.Get("spice", "flow").Value())
	assert.Exactly(t, "value-b", cache.Get("spice", "must").Value())
	assert.Nil(t, cache.Get("spice", "worm"))

	assert.Exactly(t, "ghanima", cache.Get("leto", "sister").Value())
	assert.Nil(t, cache.Get("leto", "brother"))
	assert.Nil(t, cache.Get("baron", "friend"))
}

func TestLayered_ReplaceDoesNothingIfKeyDoesNotExist(t *testing.T) {
	cache := newLayered()
	assert.False(t, cache.Replace("spice", "flow", "value-a"))
	assert.Nil(t, cache.Get("spice", "flow"))
}

func TestLayered_ReplaceUpdatesTheValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	assert.True(t, cache.Replace("spice", "flow", "value-b"))
	assert.Exactly(t, "value-b", cache.Get("spice", "flow").Value().(string))
	// not sure how to test that the TTL hasn't changed sort of a sleep..
}

func TestLayered_DeletesAValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	cache.Set("spice", "must", "value-b", time.Minute)
	cache.Set("leto", "sister", "ghanima", time.Minute)
	cache.Delete("spice", "flow")
	assert.Nil(t, cache.Get("spice", "flow"))
	assert.Exactly(t, "value-b", cache.Get("spice", "must").Value())
	assert.Nil(t, cache.Get("spice", "worm"))
	assert.Exactly(t, "ghanima", cache.Get("leto", "sister").Value())
}

func TestLayered_OnDeleteCallbackCalled(t *testing.T) {

	onDeleteFnCalled := false
	onDeleteFn := func(item *Item) {

		if item.group == "spice" && item.key == "flow" {
			onDeleteFnCalled = true
		}
	}

	cache := Layered(Configure().OnDelete(onDeleteFn))
	cache.Set("spice", "flow", "value-a", time.Minute)
	cache.Set("spice", "must", "value-b", time.Minute)
	cache.Set("leto", "sister", "ghanima", time.Minute)

	time.Sleep(time.Millisecond * 10) // Run once to init
	cache.Delete("spice", "flow")
	time.Sleep(time.Millisecond * 10) // Wait for worker to pick up deleted items

	assert.Nil(t, cache.Get("spice", "flow"))
	assert.Exactly(t, "value-b", cache.Get("spice", "must").Value())
	assert.Nil(t, cache.Get("spice", "worm"))
	assert.Exactly(t, "ghanima", cache.Get("leto", "sister").Value())

	assert.True(t, onDeleteFnCalled)
}

func TestLayered_DeletesALayer(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	cache.Set("spice", "must", "value-b", time.Minute)
	cache.Set("leto", "sister", "ghanima", time.Minute)
	cache.DeleteAll("spice")
	assert.Nil(t, cache.Get("spice", "flow"))
	assert.Nil(t, cache.Get("spice", "must"))
	assert.Nil(t, cache.Get("spice", "worm"))
	assert.Exactly(t, "ghanima", cache.Get("leto", "sister").Value())
}

func TestLayered_GCsTheOldestItems(t *testing.T) {
	cache := Layered(Configure().ItemsToPrune(10))
	cache.Set("xx", "a", 23, time.Minute)
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), "a", i, time.Minute)
	}
	cache.Set("xx", "b", 9001, time.Minute)
	// let the items get promoted (and added to our list)
	time.Sleep(time.Millisecond * 10)
	gcLayeredCache(cache)
	assert.Nil(t, cache.Get("xx", "a"))
	assert.Exactly(t, 9001, cache.Get("xx", "b").Value())
	assert.Nil(t, cache.Get("8", "a"))
	assert.Exactly(t, 9, cache.Get("9", "a").Value())
	assert.Exactly(t, 10, cache.Get("10", "a").Value())
}

func TestLayered_PromotedItemsDontGetPruned(t *testing.T) {
	cache := Layered(Configure().ItemsToPrune(10).GetsPerPromote(1))
	for i := 0; i < 500; i++ {
		cache.Set(strconv.Itoa(i), "a", i, time.Minute)
	}
	time.Sleep(time.Millisecond * 10) //run the worker once to init the list
	cache.Get("9", "a")
	time.Sleep(time.Millisecond * 10)
	gcLayeredCache(cache)
	assert.Exactly(t, 9, cache.Get("9", "a").Value())
	assert.Nil(t, cache.Get("10", "a"))
	assert.Exactly(t, 11, cache.Get("11", "a").Value())
}

func TestLayered_TrackerDoesNotCleanupHeldInstance(t *testing.T) {
	cache := Layered(Configure().ItemsToPrune(10).Track())
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), "a", i, time.Minute)
	}
	item := cache.TrackingGet("0", "a")
	time.Sleep(time.Millisecond * 10)
	gcLayeredCache(cache)
	assert.Exactly(t, 0, cache.Get("0", "a").Value())
	assert.Nil(t, cache.Get("1", "a"))
	item.Release()
	gcLayeredCache(cache)
	assert.Nil(t, cache.Get("0", "a"))
}

func TestLayered_RemovesOldestItemWhenFull(t *testing.T) {
	cache := Layered(Configure().MaxSize(5).ItemsToPrune(1))
	cache.Set("xx", "a", 23, time.Minute)
	for i := 0; i < 7; i++ {
		cache.Set(strconv.Itoa(i), "a", i, time.Minute)
	}
	cache.Set("xx", "b", 9001, time.Minute)
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, cache.Get("xx", "a"))
	assert.Nil(t, cache.Get("0", "a"))
	assert.Nil(t, cache.Get("1", "a"))
	assert.Nil(t, cache.Get("2", "a"))
	assert.Exactly(t, 3, cache.Get("3", "a").Value())
	assert.Exactly(t, 9001, cache.Get("xx", "b").Value())
}

func newLayered() *LayeredCache {
	return Layered(Configure())
}

func TestLayered_RemovesOldestItemWhenFullBySizer(t *testing.T) {
	cache := Layered(Configure().MaxSize(9).ItemsToPrune(2))
	for i := 0; i < 7; i++ {
		cache.Set("pri", strconv.Itoa(i), &SizedItem{i, 2}, time.Minute)
	}
	time.Sleep(time.Millisecond * 10)
	assert.Nil(t, cache.Get("pri", "0"))
	assert.Nil(t, cache.Get("pri", "1"))
	assert.Nil(t, cache.Get("pri", "2"))
	assert.Nil(t, cache.Get("pri", "3"))
	assert.Exactly(t, 4, cache.Get("pri", "4").Value().(*SizedItem).id)
}

func TestLayered_SetUpdatesSizeOnDelta(t *testing.T) {
	cache := Layered(Configure())
	cache.Set("pri", "a", &SizedItem{0, 2}, time.Minute)
	cache.Set("pri", "b", &SizedItem{0, 3}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 5)
	cache.Set("pri", "b", &SizedItem{0, 3}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 5)
	cache.Set("pri", "b", &SizedItem{0, 4}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 6)
	cache.Set("pri", "b", &SizedItem{0, 2}, time.Minute)
	cache.Set("sec", "b", &SizedItem{0, 3}, time.Minute)
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 7)
	cache.Delete("pri", "b")
	time.Sleep(time.Millisecond * 10)
	checkLayeredSize(t, cache, 5)
}

func TestLayered_ReplaceDoesNotchangeSizeIfNotSet(t *testing.T) {
	cache := Layered(Configure())
	cache.Set("pri", "1", &SizedItem{1, 2}, time.Minute)
	cache.Set("pri", "2", &SizedItem{1, 2}, time.Minute)
	cache.Set("pri", "3", &SizedItem{1, 2}, time.Minute)
	cache.Replace("sec", "3", &SizedItem{1, 2})
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 6)
}

func TestLayered_ReplaceChangesSize(t *testing.T) {
	cache := Layered(Configure())
	cache.Set("pri", "1", &SizedItem{1, 2}, time.Minute)
	cache.Set("pri", "2", &SizedItem{1, 2}, time.Minute)

	cache.Replace("pri", "2", &SizedItem{1, 2})
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 4)

	cache.Replace("pri", "2", &SizedItem{1, 1})
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 3)

	cache.Replace("pri", "2", &SizedItem{1, 3})
	time.Sleep(time.Millisecond * 5)
	checkLayeredSize(t, cache, 5)
}

func checkLayeredSize(t *testing.T, cache *LayeredCache, sz int64) {
	cache.Stop()
	assert.Exactly(t, sz, cache.size)
	cache.restart()
}

func gcLayeredCache(cache *LayeredCache) {
	cache.Stop()
	cache.gc()
	cache.restart()
}
