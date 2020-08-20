package pseudo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type creditCard struct {
	vendor   string
	length   int
	prefixes []int
}

const ccVendorCount = 4

var (
	ccVendorKey  = [...]string{"visa", "mastercard", "amex", "discover"}
	ccVendorName = [...]string{"VISA", "MasterCard", "American Express", "Discover"}
)

var creditCards = map[string]creditCard{
	ccVendorKey[0]: {ccVendorName[0], 16, []int{4539, 4556, 4916, 4532, 4929, 40240071, 4485, 4716, 4}},
	ccVendorKey[1]: {ccVendorName[1], 16, []int{51, 52, 53, 54, 55}},
	ccVendorKey[2]: {ccVendorName[2], 15, []int{34, 37}},
	ccVendorKey[3]: {ccVendorName[3], 16, []int{6011}},
}

// CreditCardType returns one of the following credit values:
// VISA, MasterCard, American Express and Discover
func (s *Service) CreditCardType() string {
	return ccVendorName[s.r.Intn(ccVendorCount)]
}

// CreditCardNum generated credit card number according to the card number rules
func (s *Service) CreditCardNum(vendor string) string {
	if vendor != "" {
		vendor = strings.ToLower(vendor)
	} else {
		vendor = ccVendorKey[s.r.Intn(ccVendorCount)]
	}
	card, ok := creditCards[vendor]
	if !ok {
		return fmt.Sprintf("CC Vendor %q not found, available: %v", vendor, ccVendorKey)
	}
	prefix := strconv.Itoa(card.prefixes[s.r.Intn(len(card.prefixes))])
	var buf strings.Builder
	buf.WriteString(prefix)
	for i := 1; i < card.length-len(prefix); i++ { // start 1 because last digit is check number
		fmt.Fprintf(&buf, "%d", s.r.Intn(10))
	}
	fmt.Fprintf(&buf, "%d", getCheckDigit(buf.String()))
	return buf.String()
}

func getCheckDigit(number string) int {
	var sum int
	for i := len(number) - 1; i >= 0; i-- { // reversed iteration
		digit, _ := strconv.Atoi(number[i : i+1])
		if (i % 2) == 0 {
			digit = digit * 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return ((int(math.Floor(float64(sum)/10))+1)*10 - sum) % 10
}
