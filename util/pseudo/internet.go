package pseudo

import (
	"fmt"
	"net"
	"strings"
)

// UserName generates user name in one of the following forms
// first name + last name, letter + last names or concatenation of from 1 to 3 lowercased words
func (s *Service) UserName() string {
	gender := s.randGender()
	switch s.r.Intn(3) {
	case 0:
		return s.lookup("en", gender+"_first_names", false) + s.lookup(s.o.Lang, gender+"_last_names", false)
	case 1:
		return s.Character() + s.lookup(s.o.Lang, gender+"_last_names", false)
	default:
		return strings.Replace(s.WordsN(s.r.Intn(3)+1, 20), " ", "_", -1)
	}
}

// TopLevelDomain generates random top level domain
func (s *Service) TopLevelDomain() string {
	return s.lookup(s.o.Lang, "top_level_domains", true)
}

// DomainName generates random domain name
func (s *Service) DomainName() string {
	c := strings.Replace(s.Company(), " ", "-", -1)
	c = strings.ToLower(c + "." + s.TopLevelDomain())
	var buf strings.Builder
	for _, r := range c {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9', r == '-', r == '.':
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// EmailAddress generates email address
func (s *Service) EmailAddress() string {
	return strings.ToLower(s.UserName()) + "@" + s.DomainName()
}

// EmailSubject generates random email subject
func (s *Service) EmailSubject(maxLen int) string {
	return s.Sentence(maxLen)
}

// EmailBody generates random email body
func (s *Service) EmailBody() string {
	return s.Paragraphs(1100)
}

// DomainZone generates random domain zone
func (s *Service) DomainZone() string {
	return s.lookup(s.o.Lang, "domain_zones", true)
}

// IPv4 generates IPv4 address
func (s *Service) IPv4() string {
	size := 4
	ip := make([]byte, size)
	for i := 0; i < size; i++ {
		ip[i] = byte(s.r.Intn(256))
	}
	return net.IP(ip).To4().String()
}

// IPv6 generates IPv6 address
func (s *Service) IPv6() string {
	size := 16
	ip := make([]byte, size)
	for i := 0; i < size; i++ {
		ip[i] = byte(s.r.Intn(256))
	}
	return net.IP(ip).To16().String()
}

// MacAddress generates random MacAddress
func (s *Service) MacAddress() string {
	var ip [6]byte
	for i := 0; i < 6; i++ {
		ip[i] = byte(s.r.Intn(256))
	}
	return net.HardwareAddr(ip[:]).String()
}

var urlFormats = []string{
	"http://www.%s/",
	"https://www.%s/",
	"http://%s/",
	"https://%s/",
	"http://www.%s/%s",
	"https://www.%s/%s",
	"http://%s/%s",
	"https://%s/%s",
	"http://%s/%s.html",
	"https://%s/%s.html",
	"http://%s/%s.php",
	"https://%s/%s.php",
}

// URL generates random URL standardised in urlFormats const
func (s *Service) URL() string {
	format := s.randomElementFromSliceString(urlFormats)
	countVerbs := strings.Count(format, "%s")
	if countVerbs == 1 {
		return fmt.Sprintf(format, s.DomainName())
	}
	return fmt.Sprintf(format, s.DomainName(), s.UserName())
}
