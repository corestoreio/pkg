package storenet

// requestIDService default prefix generator. Creates a prefix once the middleware
// is set up.
type RequestIDService struct {
	prefix string
}

// Prefix returns a unique prefix string for the current (micro) service.
// This id gets reset once you restart the service.
func (rp *RequestIDService) Init() {
	// algorithm taken from https://github.com/zenazn/goji/blob/master/web/middleware/request_id.go#L40-L52
	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}
	var buf [12]byte
	var b64 string
	for len(b64) < 10 {
		rand.Read(buf[:])
		b64 = base64.StdEncoding.EncodeToString(buf[:])
		b64 = strings.NewReplacer("+", "", "/", "").Replace(b64)
	}
	rp.prefix = fmt.Sprintf("%s/%s-", hostname, b64[0:10])
}

// NewID returns a new ID unique for the current compilation.
func (rp *RequestIDService) NewID(_ *http.Request) string {
	return rp.prefix + strconv.FormatInt(atomic.AddInt64(&reqID, 1), 10)
}
