package ccache

import (
	"strconv"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

func TestSecondary_GetsANonExistantValue(t *testing.T) {
	cache := newLayered().GetOrCreateSecondaryCache("foo")
	assert.NotNil(t, cache)
}

func TestSecondary_SetANewValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "a value", time.Minute)
	sCache := cache.GetOrCreateSecondaryCache("spice")
	assert.Exactly(t, "a value", sCache.Get("flow").Value())
	assert.Nil(t, sCache.Get("stop"))
}

func TestSecondary_ValueCanBeSeenInBothCaches1(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "a value", time.Minute)
	sCache := cache.GetOrCreateSecondaryCache("spice")
	sCache.Set("orinoco", "another value", time.Minute)
	assert.Exactly(t, "another value", sCache.Get("orinoco").Value())
	assert.Exactly(t, "another value", cache.Get("spice", "orinoco").Value())
}

func TestSecondary_ValueCanBeSeenInBothCaches2(t *testing.T) {
	cache := newLayered()
	sCache := cache.GetOrCreateSecondaryCache("spice")
	sCache.Set("flow", "a value", time.Minute)
	assert.Exactly(t, "a value", sCache.Get("flow").Value())
	assert.Exactly(t, "a value", cache.Get("spice", "flow").Value())
}

func TestSecondary_DeletesAreReflectedInBothCaches(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "a value", time.Minute)
	cache.Set("spice", "sister", "ghanima", time.Minute)
	sCache := cache.GetOrCreateSecondaryCache("spice")

	cache.Delete("spice", "flow")
	assert.Nil(t, cache.Get("spice", "flow"))
	assert.Nil(t, sCache.Get("flow"))

	sCache.Delete("sister")
	assert.Nil(t, cache.Get("spice", "sister"))
	assert.Nil(t, sCache.Get("sister"))
}

func TestSecondary_ReplaceDoesNothingIfKeyDoesNotExist(t *testing.T) {
	cache := newLayered()
	sCache := cache.GetOrCreateSecondaryCache("spice")
	assert.False(t, sCache.Replace("flow", "value-a"))
	assert.Nil(t, cache.Get("spice", "flow"))
}

func TestSecondary_ReplaceUpdatesTheValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	sCache := cache.GetOrCreateSecondaryCache("spice")
	assert.True(t, sCache.Replace("flow", "value-b"))
	assert.Exactly(t, "value-b", cache.Get("spice", "flow").Value().(string))
}

func TestSecondary_FetchReturnsAnExistingValue(t *testing.T) {
	cache := newLayered()
	cache.Set("spice", "flow", "value-a", time.Minute)
	sCache := cache.GetOrCreateSecondaryCache("spice")
	val, _ := sCache.Fetch("flow", time.Minute, func() (interface{}, error) { return "a fetched value", nil })
	assert.Exactly(t, "value-a", val.Value().(string))
}

func TestSecondary_FetchReturnsANewValue(t *testing.T) {
	cache := newLayered()
	sCache := cache.GetOrCreateSecondaryCache("spice")
	val, _ := sCache.Fetch("flow", time.Minute, func() (interface{}, error) { return "a fetched value", nil })
	assert.Exactly(t, "a fetched value", val.Value().(string))
}

func TestSecondary_TrackerDoesNotCleanupHeldInstance(t *testing.T) {
	cache := Layered(Configure().ItemsToPrune(10).Track())
	for i := 0; i < 10; i++ {
		cache.Set(strconv.Itoa(i), "a", i, time.Minute)
	}
	sCache := cache.GetOrCreateSecondaryCache("0")
	item := sCache.TrackingGet("a")
	time.Sleep(time.Millisecond * 10)
	gcLayeredCache(cache)
	assert.Exactly(t, 0, cache.Get("0", "a").Value())
	assert.Nil(t, cache.Get("1", "a"))
	item.Release()
	gcLayeredCache(cache)
	assert.Nil(t, cache.Get("0", "a"))
}
