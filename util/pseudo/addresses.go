package pseudo

import "strconv"

// Continent generates random continent
func (s *Service) Continent() string {
	return s.lookup(s.o.Lang, "continents", true)
}

// Country generates random country
func (s *Service) Country() string {
	return s.lookup(s.o.Lang, "countries", true)
}

// City generates random city
func (s *Service) City() string {
	city := s.lookup(s.o.Lang, "cities", true)
	switch s.r.Intn(5) {
	case 0:
		return join(s.cityPrefix(), city)
	case 1:
		return join(city, s.citySuffix())
	default:
		return city
	}
}

func (s *Service) cityPrefix() string {
	return s.lookup(s.o.Lang, "city_prefixes", false)
}

func (s *Service) citySuffix() string {
	return s.lookup(s.o.Lang, "city_suffixes", false)
}

// State generates random state
func (s *Service) State() string {
	return s.lookup(s.o.Lang, "states", false)
}

// StateAbbrev generates random state abbreviation
func (s *Service) StateAbbrev() string {
	return s.lookup(s.o.Lang, "state_abbrevs", false)
}

// Street generates random street name
func (s *Service) Street() string {
	street := s.lookup(s.o.Lang, "streets", true)
	return join(street, s.streetSuffix())
}

// StreetAddress generates random street name along with building number
func (s *Service) StreetAddress() string {
	return join(s.Street(), strconv.Itoa(s.r.Intn(100)))
}

func (s *Service) streetSuffix() string {
	return s.lookup(s.o.Lang, "street_suffixes", true)
}

// Zip generates random zip code using one of the formats specifies in zip_format file
func (s *Service) Zip() string {
	return s.generate(s.o.Lang, "zips", true)
}

// Phone generates random phone number using one of the formats format specified in phone_format file
func (s *Service) Phone() string {
	return s.generate(s.o.Lang, "phones", true)
}

// Address returns an american style address
func (s *Service) Address() string {
	panic("TODO implement")
	// use the same formatting as in M2 ;-)
	// return fmt.Sprintf("%d %s,\n%s, %s, %s", Number(100), Street(), City(), State(Small), PostalCode("US"))
}
