package null

import "math/bits"

func (d Decimal) SizeVT() (n int) {
	if !d.Valid {
		return 0
	}
	l := len(d.PrecisionStr)
	if l > 0 {
		n += 1 + l + sov(uint64(l))
	}
	if d.Precision != 0 {
		n += 1 + sov(d.Precision)
	}
	if d.Scale != 0 {
		n += 1 + sov(uint64(d.Scale))
	}
	if d.Negative {
		n += 2
	}
	if d.Valid {
		n += 2
	}
	if d.Quote {
		n += 2
	}
	return n
}

func sov(x uint64) (n int) {
	return (bits.Len64(x|1) + 6) / 7
}
