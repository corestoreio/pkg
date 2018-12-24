package pseudo

// Latitude generates latitude (from -90.0 to 90.0)
func (s *Service) Latitude() float64 {
	return float64(s.r.Float32()*180 - 90)
}

// LatitudeDegrees generates latitude degrees (from -90 to 90)
func (s *Service) LatitudeDegrees() int {
	return s.r.Intn(180) - 90
}

// LatitudeMinutes generates latitude minutes (from 0 to 60)
func (s *Service) LatitudeMinutes() int {
	return s.r.Intn(60)
}

// LatitudeSeconds generates latitude seconds (from 0 to 60)
func (s *Service) LatitudeSeconds() int {
	return s.r.Intn(60)
}

// LatitudeDirection generates latitude direction (N(orth) o S(outh))
func (s *Service) LatitudeDirection() string {
	if s.r.Intn(2) == 0 {
		return "N"
	}
	return "S"
}

// Longitude generates longitude (from -180 to 180)
func (s *Service) Longitude() float64 {
	return float64(s.r.Float32()*360 - 180)
}

// LongitudeDegrees generates longitude degrees (from -180 to 180)
func (s *Service) LongitudeDegrees() int {
	return s.r.Intn(360) - 180
}

// LongitudeMinutes generates (from 0 to 60)
func (s *Service) LongitudeMinutes() int {
	return s.r.Intn(60)
}

// LongitudeSeconds generates (from 0 to 60)
func (s *Service) LongitudeSeconds() int {
	return s.r.Intn(60)
}

// LongitudeDirection generates (W(est) or E(ast))
func (s *Service) LongitudeDirection() string {
	if s.r.Intn(2) == 0 {
		return "W"
	}
	return "E"
}
