package ccache

import (
	"math"
	"testing"
	"time"

	"github.com/corestoreio/pkg/util/assert"
)

func TestPromotability(t *testing.T) {
	item := &Item{promotions: 4}
	assert.True(t, item.shouldPromote(5))
	assert.False(t, item.shouldPromote(5))
}

func TestExpired(t *testing.T) {
	now := time.Now().UnixNano()
	item1 := &Item{expires: now + (10 * int64(time.Millisecond))}
	item2 := &Item{expires: now - (10 * int64(time.Millisecond))}
	assert.False(t, item1.Expired())
	assert.True(t, item2.Expired())
}

func TestTTL(t *testing.T) {
	now := time.Now().UnixNano()
	item1 := &Item{expires: now + int64(time.Second)}
	item2 := &Item{expires: now - int64(time.Second)}
	assert.Exactly(t, 1, int(math.Ceil(item1.TTL().Seconds())))
	assert.Exactly(t, -1, int(math.Ceil(item2.TTL().Seconds())))
}

func TestExpires(t *testing.T) {
	now := time.Now().UnixNano()
	item := &Item{expires: now + (10)}
	assert.Exactly(t, now+10, item.Expires().UnixNano())
}

func TestExtend(t *testing.T) {
	item := &Item{expires: time.Now().UnixNano() + 10}
	item.Extend(time.Minute * 2)
	assert.Exactly(t, time.Now().Unix()+120, item.Expires().Unix())
}
