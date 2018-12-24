package pseudo

import "fmt"

// Currency generates currency name.
func (s *Service) Currency() string {
	return s.lookup(s.o.Lang, "currencies", true)
}

// CurrencyCode generates currency code.
func (s *Service) CurrencyCode() string {
	return s.lookup(s.o.Lang, "currency_codes", true)
}

// Amount returns a random floating price amount
// with a random precision of [1,2] up to (10**8 - 1)
func (s *Service) Price() float64 {
	return float64(s.r.Intn(10000)) + float64(s.r.Intn(99))/100
}

// PriceWithCurrency combines both price and currency together
func (s *Service) PriceWithCurrency() string {
	return fmt.Sprintf("%.02f %s", s.Price(), s.CurrencyCode())
}
