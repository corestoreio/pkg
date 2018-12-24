package pseudo

// Brand generates brand name
func (s *Service) Brand() string {
	return s.Company()
}

// ProductName generates product name
func (s *Service) ProductName() string {
	productName := s.lookup(s.o.Lang, "adjectives", true) + " " + s.lookup(s.o.Lang, "nouns", true)
	if s.r.Intn(2) == 1 {
		productName = s.lookup(s.o.Lang, "adjectives", true) + " " + productName
	}
	return productName
}

// Product generates product title as brand + product name
func (s *Service) Product() string {
	return s.Brand() + " " + s.ProductName()
}

// Model generates model name that consists of letters and digits, optionally with a hyphen between them
func (s *Service) Model() string {
	seps := []string{"", " ", "-"}
	return s.CharactersN(s.r.Intn(3)+1) + seps[s.r.Intn(len(seps))] + s.Digits()
}
