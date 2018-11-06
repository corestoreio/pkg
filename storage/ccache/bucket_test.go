package ccache

import (
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

func TestBucket_GetMissFromBucket(t *testing.T) {
	bucket := testBucket()
	assert.Nil(t, bucket.get("invalid"))
}

func TestBucket_GetHitFromBucket(t *testing.T) {
	bucket := testBucket()
	item := bucket.get("power")
	assertValue(t, item, "9000")
}

func TestBucket_DeleteItemFromBucket(t *testing.T) {
	bucket := testBucket()
	bucket.delete("power")
	assert.Nil(t, bucket.get("power"))
}

func TestBucket_SetsANewBucketItem(t *testing.T) {
	bucket := testBucket()
	item, existing := bucket.set("spice", TestValue("flow"), time.Minute)
	assertValue(t, item, "flow")
	item = bucket.get("spice")
	assertValue(t, item, "flow")
	assert.Nil(t, existing)
}

func TestBucket_SetsAnExistingItem(t *testing.T) {
	bucket := testBucket()
	item, existing := bucket.set("power", TestValue("9001"), time.Minute)
	assertValue(t, item, "9001")
	item = bucket.get("power")
	assertValue(t, item, "9001")
	assertValue(t, existing, "9000")
}

func testBucket() *bucket {
	b := &bucket{lookup: make(map[string]*Item)}
	b.lookup["power"] = &Item{
		key:   "power",
		value: TestValue("9000"),
	}
	return b
}

func assertValue(t *testing.T, item *Item, expected string) {
	value := item.value.(TestValue)
	assert.Exactly(t, TestValue(expected), value)
}

type TestValue string

func (v TestValue) Expires() time.Time {
	return time.Now()
}
