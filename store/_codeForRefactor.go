package store

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/coreos/etcd/mvcc/backend"
	"github.com/corestoreio/csfw/config"
	"github.com/corestoreio/csfw/config/cfgmock"
	"github.com/corestoreio/csfw/config/cfgmodel"
	"github.com/corestoreio/csfw/config/cfgpath"
	"github.com/corestoreio/csfw/store/scope"
	"github.com/corestoreio/csfw/util/errors"
	"github.com/gocraft/dbr"
	"github.com/stretchr/testify/assert"
)

// todo: refactor this all

/*
	Store Currency
*/

// CurrentCurrency TODO(cs)
// @see app/code/Magento/Store/Model/Store.php::getCurrentCurrency
func (s Store) CurrentCurrency() string {
	/*
		this returns just a string or string slice and no further
		involvement of the directory package.

		or those functions move directly into the directory package
	*/
	return ""
}

func (s Store) DefaultCurrency() string {
	return ""
}

func (s Store) AvailableCurrencyCodes() []string {
	return nil
}

// Path returns the sub path from the URL where CoreStore is installed
func (s Store) Path() string {
	url, err := s.BaseURL(config.URLTypeWeb, false)
	if err != nil {
		return "/"
	}
	return url.Path
}

//BaseURL returns a parsed and maybe cached URL from config.ScopedReader.
//It returns a copy of url.URL or an error. Possible URLTypes are:
//    - config.URLTypeWeb
//    - config.URLTypeStatic
//    - config.URLTypeMedia
func (s Store) BaseURL(ut config.URLType, isSecure bool) (url.URL, error) {

	switch isSecure {
	case true:
		if pu := s.urlcache.secure.Get(ut); pu != nil {
			return *pu, nil
		}
	case false:
		if pu := s.urlcache.unsecure.Get(ut); pu != nil {
			return *pu, nil
		}
	}

	var p cfgmodel.BaseURL
	switch ut {
	case config.URLTypeWeb:
		p = backend.Backend.WebUnsecureBaseURL
		if isSecure {
			p = backend.Backend.WebSecureBaseURL
		}
		break
	case config.URLTypeStatic:
		p = backend.Backend.WebUnsecureBaseStaticURL
		if isSecure {
			p = backend.Backend.WebSecureBaseStaticURL
		}
		break
	case config.URLTypeMedia:
		p = backend.Backend.WebUnsecureBaseMediaURL
		if isSecure {
			p = backend.Backend.WebSecureBaseMediaURL
		}
		break
	case config.URLTypeAbsent: // hack to clear the cache :-( refactor that
		_ = s.urlcache.unsecure.Clear()
		return url.URL{}, s.urlcache.secure.Clear()
	// TODO(cs) rethink that here and maybe add the other paths if needed.
	default:
		return url.URL{}, fmt.Errorf("Unsupported UrlType: %d", ut)
	}

	rawURL, _, err := p.Get(s.Config)
	if err != nil {
		return url.URL{}, err
	}

	if strings.Contains(rawURL, cfgmodel.PlaceholderBaseURL) {
		// TODO(cs) replace placeholder with \Magento\Framework\App\Request\Http::getDistroBaseUrl()
		// getDistroBaseUrl will be generated from the $_SERVER variable,
		base, err := s.baseConfig.String(cfgpath.MustNewByParts(config.PathCSBaseURL))
		if err != nil && !errors.IsNotFound(err) {
			base = config.CSBaseURL
		}
		rawURL = strings.Replace(rawURL, cfgmodel.PlaceholderBaseURL, base, 1)
	}
	rawURL = strings.TrimRight(rawURL, "/") + "/"

	if isSecure {
		retURL, retErr := s.urlcache.secure.Set(ut, rawURL)
		return *retURL, retErr
	}
	retURL, retErr := s.urlcache.unsecure.Set(ut, rawURL)
	return *retURL, retErr
}

//IsFrontURLSecure returns true from the config if the frontend must be secure.
func (s Store) IsFrontURLSecure() bool {
	return false // backend.Backend.WebSecureUseInFrontend.Get(s.Config)
}

//IsCurrentlySecure checks if a request for a give store aka. scope is secure. Checks
//include if base URL has been set and if front URL is secure
//This function might gets executed on every request.
func (s Store) IsCurrentlySecure(r *http.Request) bool {
	return false
	if httputil.IsSecure(s.cr, r) {
		return true
	}

	secureBaseURL, err := s.BaseURL(config.URLTypeWeb, true)
	if err != nil || false == s.IsFrontURLSecure() {
		PkgLog.Debug("store.Store.IsCurrentlySecure.BaseURL", "err", err, "secureBaseURL", secureBaseURL)
		return false
	}
	return secureBaseURL.Scheme == "https" && r.URL.Scheme == "https" // todo(cs) check for ports !? other schemes?
}

//TODO move net related functions into the storenet package
func TestStoreBaseURLandPath(t *testing.T) {

	t.Skip("@todo refactor and move these functions into another package")

	s, err := store.NewStore(
		cfgmock.NewService(),
		&store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 1, Name: "Germany", SortOrder: 10, IsActive: true},
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("admin"), Name: dbr.NewNullString("Admin"), SortOrder: 0, DefaultGroupID: 0, IsDefault: dbr.NewNullBool(false)},
		&store.TableGroup{GroupID: 1, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 1},
	)
	assert.NoError(t, err)
	if s == nil {
		t.Fail()
	}

	tests := []struct {
		haveR        config.Getter
		haveUT       config.URLType
		haveIsSecure bool
		wantBaseUrl  string
		wantPath     string
	}{
		{
			cfgmock.NewService(cfgmock.WithString(
				func(path string) (string, error) {

					switch path {
					// scope is here store but config.ScopedGetter must fall back to default
					case backend.Backend.WebSecureBaseURL.String():
						return "https://corestore.io", nil
					case backend.Backend.WebUnsecureBaseURL.String():
						return "http://corestore.io", nil
					}
					return "", errors.NewNotFoundf("Invalid path: %s", path)
				},
			)),
			config.URLTypeWeb, true, "https://corestore.io/", "/",
		},
		{
			cfgmock.NewService(cfgmock.WithString(
				func(path string) (string, error) {
					switch path {
					case backend.Backend.WebSecureBaseURL.String():
						return "https://myplatform.io/customer1", nil
					case backend.Backend.WebUnsecureBaseURL.String():
						return "http://myplatform.io/customer1", nil
					}
					return "", errors.NewNotFoundf("Invalid path: %s", path)
				},
			)),
			config.URLTypeWeb, false, "http://myplatform.io/customer1/", "/customer1/",
		},
		{
			cfgmock.NewService(cfgmock.WithString(
				func(p string) (string, error) {
					switch p {
					case backend.Backend.WebSecureBaseURL.String():
						return cfgmodel.PlaceholderBaseURL, nil
					case backend.Backend.WebUnsecureBaseURL.String():
						return cfgmodel.PlaceholderBaseURL, nil
					case cfgpath.MustNewByParts(config.PathCSBaseURL).String():
						return config.CSBaseURL, nil
					}
					return "", errors.NewNotFoundf("Invalid path: %s", p)
				},
			)),
			config.URLTypeWeb, false, config.CSBaseURL, "/",
		},
	}

	for i, test := range tests {
		s.Options(store.WithStoreConfig(test.haveR))
		assert.NotNil(t, s.Config, "Index %d", i)
		baseURL, err := s.BaseURL(test.haveUT, test.haveIsSecure)
		assert.NoError(t, err)
		assert.EqualValues(t, test.wantBaseUrl, baseURL.String())
		assert.EqualValues(t, test.wantPath, s.Path())

		_, err = s.BaseURL(config.URLTypeAbsent, false)
		assert.NoError(t, err)
	}
}

//TODO
func getWebsiteBaseCurrency(priceScope int, curGlobal, curWebsite string) (*store.Website, error) {
	return store.NewWebsite(
		&store.TableWebsite{WebsiteID: 1, Code: dbr.NewNullString("euro"), Name: dbr.NewNullString("Europe"), SortOrder: 0, DefaultGroupID: 1, IsDefault: dbr.NewNullBool(true)},
		store.SetWebsiteGroupsStores(
			store.TableGroupSlice{
				0: &store.TableGroup{GroupID: 0, WebsiteID: 1, Name: "Default", RootCategoryID: 0, DefaultStoreID: 1},
			},
			store.TableStoreSlice{
				0: &store.TableStore{StoreID: 0, Code: dbr.NewNullString("Admin"), WebsiteID: 1, GroupID: 0, Name: "Admin", SortOrder: 0, IsActive: true},
				1: &store.TableStore{StoreID: 1, Code: dbr.NewNullString("de"), WebsiteID: 1, GroupID: 0, Name: "Germany", SortOrder: 10, IsActive: true},
			},
		),
		store.SetWebsiteConfig(
			cfgmock.NewService(cfgmock.PathValue{
				catconfig.Backend.CatalogPriceScope.FQPathInt64(scope.StrDefault, 0):    priceScope,
				directory.Backend.CurrencyOptionsBase.FQPathInt64(scope.StrDefault, 0):  curGlobal,
				directory.Backend.CurrencyOptionsBase.FQPathInt64(scope.StrWebsites, 1): curWebsite,
			}),
		),
	)
}

func TestWebsiteBaseCurrency(t *testing.T) {

	tests := []struct {
		priceScope int
		curGlobal  string
		curWebsite string
		curWant    string
		wantErr    error
	}{
		{catconfig.PriceScopeGlobal, "USD", "EUR", "USD", nil},
		{catconfig.PriceScopeGlobal, "ZZ", "EUR", "XXX", errors.New("currency: tag is not well-formed")},
		{catconfig.PriceScopeWebsite, "USD", "EUR", "EUR", nil},
		{catconfig.PriceScopeWebsite, "USD", "YYY", "XXX", errors.New("currency: tag is not a recognized currency")},
	}

	for _, test := range tests {
		w, err := getWebsiteBaseCurrency(test.priceScope, test.curGlobal, test.curWebsite)
		assert.NoError(t, err)
		if false == assert.NotNil(t, w) {
			t.Fatal("website is nil")
		}

		haveCur, haveErr := w.BaseCurrency()

		if test.wantErr != nil {
			assert.EqualError(t, haveErr, test.wantErr.Error())
			assert.Exactly(t, test.curWant, haveCur.Unit.String())
			continue
		}

		assert.NoError(t, haveErr)

		wantCur, err := directory.NewCurrencyISO(test.curWant)
		assert.NoError(t, err)
		assert.Exactly(t, wantCur, haveCur)
	}
}
