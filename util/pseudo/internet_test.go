package pseudo

import (
	"testing"

	"github.com/corestoreio/pkg/util/assert"
)

func TestInternet(t *testing.T) {
	s := MustNewService(0, nil)
	for _, lang := range s.GetLangs() {
		assert.NoError(t, s.SetLang(lang))

		v := s.UserName()
		t.Logf("UserName %q", v)
		assert.LenBetween(t, v, 1, 35, "UserName failed with lang %s", lang)

		v = s.TopLevelDomain()
		t.Logf("TopLevelDomain %q", v)
		assert.LenBetween(t, v, 1, 30, "TopLevelDomain failed with lang %s", lang)

		v = s.DomainName()
		t.Logf("DomainName %q", v)
		assert.LenBetween(t, v, 1, 70, "DomainName failed with lang %s", lang)

		v = s.EmailAddress()
		t.Logf("EmailAddress %q", v)
		assert.LenBetween(t, v, 1, 90, "EmailAddress failed with lang %s", lang)

		v = s.EmailSubject(30)
		t.Logf("EmailSubject %q", v)
		assert.LenBetween(t, v, 1, 100, "EmailSubject failed with lang %s", lang)

		v = s.EmailBody()
		t.Logf("EmailBody %q", v)
		assert.LenBetween(t, v, 1, 1230, "EmailBody failed with lang %s", lang)

		v = s.DomainZone()
		t.Logf("DomainZone %q", v)
		assert.LenBetween(t, v, 1, 30, "DomainZone failed with lang %s", lang)

		v = s.IPv4()
		t.Logf("IPv4 %q", v)
		assert.LenBetween(t, v, 1, 30, "IPv4 failed with lang %s", lang)

		v = s.IPv6()
		t.Logf("IPv6 %q", v)
		assert.LenBetween(t, v, 1, 39, "IPv6 failed with lang %s", lang)

	}
}
