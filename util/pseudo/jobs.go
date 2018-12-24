package pseudo

import (
	"strings"

	"golang.org/x/exp/errors/fmt"
)

// Company generates company name
func (s *Service) Company() string {
	return s.lookup(s.o.Lang, "companies", true)
}

// Company generates company name
func (s *Service) CompanyLegal() string {
	return fmt.Sprintf("%s %s",
		s.lookup(s.o.Lang, "companies", true),
		s.lookup(s.o.Lang, "company_entity", true),
	)
}

// JobTitle generates job title
func (s *Service) JobTitle() string {
	job := s.lookup(s.o.Lang, "jobs", true)
	return strings.Replace(job, "#{N}", s.jobTitleSuffix(), 1)
}

func (s *Service) jobTitleSuffix() string {
	return s.lookup(s.o.Lang, "jobs_suffixes", false)
}

// Industry generates industry name
func (s *Service) Industry() string {
	return s.lookup(s.o.Lang, "industries", true)
}
