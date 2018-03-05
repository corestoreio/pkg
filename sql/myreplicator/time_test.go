package myreplicator

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFracTime_String(t *testing.T) {
	t.Parallel()

	t.Run("toString", func(t *testing.T) {

		tbls := []struct {
			year     int
			month    int
			day      int
			hour     int
			min      int
			sec      int
			microSec int
			frac     int
			expected string
		}{
			{2000, 1, 1, 1, 1, 1, 1, 0, "2000-01-01 01:01:01"},
			{2000, 1, 1, 1, 1, 1, 1, 1, "2000-01-01 01:01:01.0"},
			{2000, 1, 1, 1, 1, 1, 1, 6, "2000-01-01 01:01:01.000001"},
		}
		for _, test := range tbls {
			t1 := fracTime{time.Date(test.year, time.Month(test.month), test.day, test.hour, test.min, test.sec, test.microSec*1000, time.UTC), test.frac}
			assert.Exactly(t, test.expected, t1.String())
		}

	})

	t.Run("zero", func(t *testing.T) {

		zeroTbls := []struct {
			frac     int
			dec      int
			expected string
		}{
			{0, 1, "0000-00-00 00:00:00.0"},
			{1, 1, "0000-00-00 00:00:00.0"},
			{123, 3, "0000-00-00 00:00:00.000"},
			{123000, 3, "0000-00-00 00:00:00.123"},
			{123, 6, "0000-00-00 00:00:00.000123"},
			{123000, 6, "0000-00-00 00:00:00.123000"},
		}

		for _, test := range zeroTbls {
			assert.Exactly(t, test.expected, formatZeroTime(test.frac, test.dec))
		}
	})
}
