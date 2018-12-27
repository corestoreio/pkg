package pseudo

func (s *Service) randGender() string {
	g := "male"
	if s.r.Intn(2) == 0 {
		g = "female"
	}
	return g
}

func (s *Service) firstName(gender string) string {
	return s.lookup(s.o.Lang, gender+"_first_names", true)
}

// MaleFirstName generates male first name
func (s *Service) MaleFirstName() string {
	return s.firstName("male")
}

// FemaleFirstName generates female first name
func (s *Service) FemaleFirstName() string {
	return s.firstName("female")
}

// FirstName generates first name
func (s *Service) FirstName() string {
	return s.firstName(s.randGender())
}

func (s *Service) lastName(gender string) string {
	return s.lookup(s.o.Lang, gender+"_last_names", true)
}

// MaleLastName generates male last name
func (s *Service) MaleLastName() string {
	return s.lastName("male")
}

// FemaleLastName generates female last name
func (s *Service) FemaleLastName() string {
	return s.lastName("female")
}

// LastName generates last name
func (s *Service) LastName() string {
	return s.lastName(s.randGender())
}

func (s *Service) patronymic(gender string) string {
	return s.lookup(s.o.Lang, gender+"_patronymics", false)
}

// MalePatronymic generates male patronymic
func (s *Service) MalePatronymic() string {
	return s.patronymic("male")
}

// FemalePatronymic generates female patronymic
func (s *Service) FemalePatronymic() string {
	return s.patronymic("female")
}

// Patronymic generates patronymic
func (s *Service) Patronymic() string {
	return s.patronymic(s.randGender())
}

// Prefix returns a random prefix for either a female or male, sometimes an
// empty prefix.
func (s *Service) Suffix() string {
	if s.r.Intn(101)%2 == 0 {
		return ""
	}
	return s.suffix(s.randGender())
}

// Prefix returns a random prefix for either a female or male, sometimes an
// empty prefix.
func (s *Service) Prefix() string {
	if s.r.Intn(101)%5 == 0 {
		return ""
	}
	return s.prefix(s.randGender())
}
func (s *Service) prefix(gender string) string {
	return s.lookup(s.o.Lang, gender+"_name_prefixes", false)
}

func (s *Service) suffix(gender string) string {
	return s.lookup(s.o.Lang, gender+"_name_suffixes", false)
}

func (s *Service) fullNameWithPrefix(gender string) string {
	return join(s.prefix(gender), s.firstName(gender), s.lastName(gender))
}

// MaleFullNameWithPrefix generates prefixed male full name
// if prefixes for the given language are available
func (s *Service) MaleFullNameWithPrefix() string {
	return s.fullNameWithPrefix("male")
}

// FemaleFullNameWithPrefix generates prefixed female full name
// if prefixes for the given language are available
func (s *Service) FemaleFullNameWithPrefix() string {
	return s.fullNameWithPrefix("female")
}

// FullNameWithPrefix generates prefixed full name
// if prefixes for the given language are available
func (s *Service) FullNameWithPrefix() string {
	return s.fullNameWithPrefix(s.randGender())
}

func (s *Service) fullNameWithSuffix(gender string) string {
	return join(s.firstName(gender), s.lastName(gender), s.suffix(gender))
}

// MaleFullNameWithSuffix generates suffixed male full name
// if suffixes for the given language are available
func (s *Service) MaleFullNameWithSuffix() string {
	return s.fullNameWithPrefix("male")
}

// FemaleFullNameWithSuffix generates suffixed female full name
// if suffixes for the given language are available
func (s *Service) FemaleFullNameWithSuffix() string {
	return s.fullNameWithPrefix("female")
}

// FullNameWithSuffix generates suffixed full name
// if suffixes for the given language are available
func (s *Service) FullNameWithSuffix() string {
	return s.fullNameWithPrefix(s.randGender())
}

func (s *Service) fullName(gender string) string {
	switch s.r.Intn(10) {
	case 0:
		return s.fullNameWithPrefix(gender)
	case 1:
		return s.fullNameWithSuffix(gender)
	default:
		return join(s.firstName(gender), s.lastName(gender))
	}
}

// MaleFullName generates male full name
// it can occasionally include prefix or suffix
func (s *Service) MaleFullName() string {
	return s.fullName("male")
}

// FemaleFullName generates female full name
// it can occasionally include prefix or suffix
func (s *Service) FemaleFullName() string {
	return s.fullName("female")
}

// FullName generates full name
// it can occasionally include prefix or suffix
func (s *Service) FullName() string {
	return s.fullName(s.randGender())
}
