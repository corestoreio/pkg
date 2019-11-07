package pseudo

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	"github.com/corestoreio/errors"
	"github.com/corestoreio/pkg/storage/null"
	"github.com/corestoreio/pkg/sync/bgwork"
	"github.com/corestoreio/pkg/util/assert"
)

type SomeStruct struct {
	Inta    int
	Int8    int8
	Int16   int16
	Int32   int32
	Int64   int64
	Float32 float32
	Float64 float64

	UInta  uint
	UInt8  uint8
	UInt16 uint16
	UInt32 uint32
	UInt64 uint64

	Latitude           float32 `faker:"lat"`
	LATITUDE           float64 `faker:"lat"`
	Long               float32 `faker:"long"`
	LONG               float64 `faker:"long"`
	String             string
	CreditCardType     string `faker:"cc_type"`
	CreditCardNumber   string `faker:"cc_number"`
	Email              string `faker:"email"`
	IPV4               string `faker:"ipv4"`
	IPV6               string `faker:"ipv6"`
	Bool               bool
	SString            []string
	SInt               []int
	SInt8              []int8
	SInt16             []int16
	SInt32             []int32
	SInt64             []int64
	SFloat32           []float32
	SFloat64           []float64
	SBool              []bool
	Struct             AStruct
	Time               time.Time
	Stime              []time.Time
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"price"`
	AmountWithCurrency string  `faker:"price_currency"`
	ID                 int64   `faker:"id"`
	UUID               string  `faker:"uuid"`
	HyphenatedID       string  `faker:"uuid_string"`

	MapStringString        map[string]string
	MapStringStruct        map[string]AStruct
	MapStringStructPointer map[string]*AStruct
}

type AStruct struct {
	Number        int64
	Height        int64
	AnotherStruct CStruct
}

type BStruct struct {
	Image string
}

type CStruct struct {
	BStruct
	Name string
}

type TaggedStruct struct {
	Latitude           float32 `faker:"lat"`
	Longitude          float32 `faker:"long"`
	CreditCardNumber   string  `faker:"cc_number"`
	CreditCardType     string  `faker:"cc_type"`
	Email              string  `faker:"email"`
	IPV4               string  `faker:"ipv4"`
	IPV6               string  `faker:"ipv6"`
	Password           string  `faker:"password"`
	PhoneNumber        string  `faker:"phone_number"`
	MacAddress         string  `faker:"mac_address"`
	URL                string  `faker:"url"`
	UserName           string  `faker:"username"`
	FirstName          string  `faker:"first_name"`
	FirstNameMale      string  `faker:"male_first_name"`
	FirstNameFemale    string  `faker:"female_first_name"`
	LastName           string  `faker:"last_name"`
	Name               string  `faker:"name"`
	UnixTime           int64   `faker:"unix_time"`
	Date               string  `faker:"date"`
	Time               string  `faker:"time"`
	MonthName          string  `faker:"month_name"`
	Year               string  `faker:"year"`
	Month              string  `faker:"month"`
	DayOfWeek          string  `faker:"week_day"`
	Timestamp          string  `faker:"timestamp"`
	TimeZone           string  `faker:"timezone"`
	Word               string  `faker:"word"`
	Sentence           string  `faker:"sentence"`
	Paragraph          string  `faker:"paragraph"`
	Currency           string  `faker:"currency"`
	Amount             float64 `faker:"price"`
	AmountWithCurrency string  `faker:"price_currency"`
	ID                 []byte  `faker:"uuid"`
	HyphenatedID       string  `faker:"uuid_string"`
}

type NotTaggedStruct struct {
	Latitude         float32
	Long             float32
	CreditCardType   string
	CreditCardNumber string
	Email            string
	IPV4             string
	IPV6             string
}

func TestFakerData(t *testing.T) {
	s := MustNewService(0, nil)
	var a SomeStruct
	err := s.FakeData(&a)
	assert.NoError(t, err, "\n%+v", err)

	data, err := json.Marshal(&a)
	assert.NoError(t, err, "\n%+v", err)
	assert.LenBetween(t, data, 100, 140000)

	// t.Logf("SomeStruct: %+v\n", a)

	var b TaggedStruct
	err = s.FakeData(&b)
	assert.NoError(t, err, "%+v", err)

	// repr.Println(b)

	// t.Logf("TaggedStruct: %+v\n", b)
	data, err = json.Marshal(&a)
	assert.NoError(t, err, "\n%+v", err)
	assert.LenBetween(t, data, 100, 139000)

	// Example Result :
	// {Int:8906957488773767119 Int8:6 Int16:14 Int32:391219825 Int64:2374447092794071106 String:poraKzAxVbWVkMkpcZCcWlYMd Bool:false SString:[MehdV aVotHsi] SInt:[528955241289647236 7620047312653801973 2774096449863851732] SInt8:[122 -92 -92] SInt16:[15679 -19444 -30246] SInt32:[1146660378 946021799 852909987] SInt64:[6079203475736033758 6913211867841842836 3269201978513619428] SFloat32:[0.019562425 0.12729558 0.36450312] SFloat64:[0.7825838989890364 0.9732903338838912 0.8316541489234004] SBool:[true false true] Struct:{Number:7693944638490551161 Height:6513508020379591917}}
}

func TestUnsuportedMapStringInterface(t *testing.T) {
	s := MustNewService(0, nil)

	type Sample struct {
		Map map[string]interface{}
	}
	sample := new(Sample)
	err := s.FakeData(sample)
	assert.NoError(t, err, "%+v", err)
	assert.Empty(t, sample.Map, "sample.Map should be empty")
}

func TestSetDataIfArgumentNotPtr(t *testing.T) {
	s := MustNewService(0, nil)
	temp := struct{}{}
	err := s.FakeData(temp)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataIfArgumentNotHaveReflect(t *testing.T) {
	temp := func() {}
	s := MustNewService(0, nil)
	err := s.FakeData(temp)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataErrorDataParseTagStringType(t *testing.T) {
	temp := &struct {
		Test string `faker:"test"`
	}{}
	t.Logf("%+v ", temp)
	s := MustNewService(0, nil)
	err := s.FakeData(temp)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

func TestSetDataErrorDataParseTagIntType(t *testing.T) {
	temp := &struct {
		Test int `faker:"test"`
	}{}
	s := MustNewService(0, nil)
	if err := s.FakeData(temp); err == nil {
		t.Error("Exptected error Unsupported tag, but got nil")
	}
}

func TestSetDataWithTagIfFirstArgumentNotPtr(t *testing.T) {
	temp := struct{}{}
	s := MustNewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(temp), "", 0, false)
	assert.True(t, errors.NotSupported.Match(err), "%+v", err)
}

func TestSetDataWithTagIfFirstArgumentSlice(t *testing.T) {
	temp := []int{}
	s := MustNewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(&temp), "", 0, false)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

func TestSetDataWithTagIfFirstArgumentNotFound(t *testing.T) {
	temp := struct{}{}
	s := MustNewService(0, nil)
	err := s.setDataWithTag(reflect.ValueOf(&temp), "", 0, false)
	assert.True(t, errors.NotFound.Match(err), "%+v", err)
}

type PointerStructA struct {
	SomeStruct *SomeStruct
}

type PointerStructB struct {
	PointA PointerStructA
}

type PointerC struct {
	TaggedStruct *TaggedStruct
}

func TestStructPointer(t *testing.T) {
	s := MustNewService(0, nil)
	a := new(PointerStructB)
	err := s.FakeData(a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	// t.Logf("A value: %+v , Somestruct Value: %+v  ", a, a)

	tagged := new(PointerC)
	err = s.FakeData(tagged)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	// t.Logf(" tagged value: %+v , TaggedStruct Value: %+v  ", a, a.PointA.SomeStruct)
}

type CustomString string
type CustomInt int
type CustomMap map[string]string
type CustomPointerStruct PointerStructB
type CustomTypeStruct struct {
	CustomString        CustomString
	CustomInt           CustomInt
	CustomMap           CustomMap
	CustomPointerStruct CustomPointerStruct
}

func TestCustomType(t *testing.T) {
	a := new(CustomTypeStruct)
	s := MustNewService(0, nil)
	err := s.FakeData(a)
	assert.NoError(t, err)
	// t.Logf("A value: %+v , Somestruct Value: %+v  ", a, a)
}

type SampleStruct struct {
	name string
	Age  int
}

func TestUnexportedFieldStruct(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has unexported field
	a := new(SampleStruct)
	s := MustNewService(0, nil)
	err := s.FakeData(a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , SampleStruct Value: %+v  ", a, a)
}

func TestPointerToCustomScalar(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has field
	a := new(CustomInt)
	s := MustNewService(0, nil)
	err := s.FakeData(a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , Custom scalar Value: %+v  ", a, a)
}

type PointerCustomIntStruct struct {
	V *CustomInt
}

func TestPointerToCustomIntStruct(t *testing.T) {
	// This test is to ensure that the faker won't panic if trying to fake data on struct that has field
	a := new(PointerCustomIntStruct)
	s := MustNewService(0, nil)
	err := s.FakeData(a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}
	t.Logf("A value: %+v , PointerCustomIntStruct scalar Value: %+v  ", a, a)
}

func TestSkipField(t *testing.T) {
	// This test is to ensure that the faker won't fill field with tag skip

	a := struct {
		ID              int
		ShouldBeSkipped int `faker:"-"`
	}{}
	s := MustNewService(0, nil)
	err := s.FakeData(&a)
	if err != nil {
		t.Error("Expected Not Error, But Got: ", err)
	}

	if a.ShouldBeSkipped != 0 {
		t.Error("Expected that field will be skipped")
	}
}

func TestExtend(t *testing.T) {
	// This test is to ensure that faker can be extended new providers

	a := struct {
		ID string `faker:"test"`
	}{}
	s := MustNewService(0, nil, WithTagFakeFunc("test", func(maxLen int) (interface{}, error) {
		return "test", nil
	}))

	assert.NoError(t, s.FakeData(&a))
	assert.Exactly(t, "test", a.ID)
}

func TestTagAlreadyExists(t *testing.T) {
	s := MustNewService(0, nil, WithTagFakeFunc("email", func(maxLen int) (interface{}, error) { return "", nil }))
	assert.NotNil(t, s)
}

func TestSetLang(t *testing.T) {
	s := MustNewService(0, nil)
	err := s.SetLang("ru")
	if err != nil {
		t.Error("SetLang should successfully set lang")
	}

	err = s.SetLang("sd")
	if err == nil {
		t.Error("SetLang with nonexistent lang should return error")
	}
}

func TestFakerRuWithoutCallback(t *testing.T) {
	s := MustNewService(0, &Options{
		EnFallback: false,
	})
	assert.NoError(t, s.SetLang("ru"))
	brand := s.Brand()
	if brand != "" {
		t.Error("Fake call with no samples should return blank string")
	}
}

func TestFakerRuWithCallback(t *testing.T) {
	s := MustNewService(0, &Options{
		EnFallback: true,
	})
	assert.NoError(t, s.SetLang("ru"))
	brand := s.Brand()
	if brand == "" {
		t.Error("Fake call for name with no samples with callback should not return blank string")
	}
}

// TestConcurrentSafety runs fake methods in multiple go routines concurrently.
// This test should be run with the race detector enabled.
func TestConcurrentSafety(t *testing.T) {
	s := MustNewService(0, &Options{
		EnFallback: true,
	})

	funcs := []func() string{
		s.FirstName,
		s.LastName,
		s.Gender,
		s.FullName,
		s.WeekDayShort,
		s.Country,
		s.Company,
		s.Industry,
		s.Street,
	}

	bgwork.Wait(len(funcs), func(idx int) {
		_ = funcs[idx]()
		for i := 0; i < 20; i++ {
			for _, fn := range funcs {
				_ = fn()
			}
		}
	})
}

type CoreConfigData struct {
	ConfigID      uint32       `json:"config_id,omitempty" max_len:"10"`
	Scope         string       `json:"scope,omitempty" max_len:"8"`
	ScopeID       int32        `json:"scope_id" xml:"scope_id"`
	Expires       null.Time    `json:"expires,omitempty" `
	Path          string       `json:"x_path" xml:"y_path" max_len:"255"`
	Value         null.String  `json:"value,omitempty" max_len:"65535"`
	ColDecimal100 null.Decimal `json:"col_decimal_10_0,omitempty"  max_len:"10"`
	ColBlob       []byte       `json:"col_blob,omitempty"  max_len:"65535"`
}

func TestMaxLen(t *testing.T) {
	t.Run("CoreConfigData", func(t *testing.T) {
		s := MustNewService(0, nil)

		for i := 0; i < 100; i++ {
			ccd := new(CoreConfigData)
			err := s.FakeData(ccd)
			assert.NoError(t, err, "%+v", err)
			assert.LenBetween(t, ccd.ConfigID, 0, 10)
			assert.LenBetween(t, ccd.Scope, 1, 8)
			assert.LenBetween(t, ccd.ScopeID, 0, math.MaxInt32)
			assert.LenBetween(t, ccd.Path, 1, 255)
			if ccd.Value.Valid {
				assert.LenBetween(t, ccd.Value.Data, 1, 65535)
			} else {
				assert.Exactly(t, null.String{}, ccd.Value)
			}
			if len(ccd.ColBlob) > 0 {
				assert.LenBetween(t, ccd.ColBlob, 1, 65535)
			} else {
				assert.Empty(t, ccd.ColBlob)
			}
			// t.Logf("%#v", ccd.ColDecimal100)
		}
	})
}

type MyNullString struct {
	String string `max_len:"10"`
	Valid  bool
}

func TestRespectValidField(t *testing.T) {
	s := MustNewService(0, &Options{
		RespectValidField: true,
	})

	for i := 0; i < 1000; i++ {
		ns := new(MyNullString)
		assert.NoError(t, s.FakeData(ns))
		if ns.Valid {
			assert.LenBetween(t, ns.String, 1, 10)
		} else {
			assert.Empty(t, ns.String)
		}
	}
}

func TestToSnakeCase(t *testing.T) {
	tests := []struct {
		have, want string
	}{
		{"Hello", "hello"},
		{"RegionID", "region_id"},
		{"EntityID", "entity_id"},
		{"PasswordToken", "password_token"},
	}

	bgwork.Wait(len(tests), func(idx int) {
		for _, test := range tests {
			h := toSnakeCase(test.have)
			assert.Exactly(t, test.want, h, "Want: %q Have: %q", test.want, h)
		}
	})
}

type CustomerEntity struct {
	EntityID     uint32      `max_len:"10"`
	WebsiteID    null.Uint32 `max_len:"5"`
	Email        null.String `max_len:"255"`
	GroupID      uint32      `max_len:"5"`
	Prefix       null.String `max_len:"40"`
	Firstname    null.String `max_len:"255"`
	Middlename   null.String `max_len:"255"`
	Lastname     null.String `max_len:"255"`
	Dob          null.Time
	PasswordHash null.String `max_len:"128"`
	RandomWeird  null.String `max_len:"13"`
}

func TestCustomerEntity_Fieldnames(t *testing.T) {
	s := MustNewService(0, nil)

	for i := 0; i < 10; i++ {
		var a CustomerEntity
		err := s.FakeData(&a)
		assert.NoError(t, err, "\n%+v", err)
		// t.Logf("%#v", a)
		assert.LenBetween(t, fmt.Sprintf("%#v", a), 260, 545)
		if a.Email.Valid {
			assert.Regexp(t, "^[a-z0-9\\-_]+@.+\\.[a-z0-9\\-]+$", a.Email.Data, "Email address")
		}
	}
}

type MaxStringLen struct {
	ColBlob      []byte      `json:"col_blob,omitempty"  max_len:"65535"`
	ColLongtext1 null.String `json:"col_longtext_1,omitempty"  max_len:"4294967295"`
	ColLongtext2 string      `json:"col_longtext_2,omitempty"  max_len:"4294967295"`
}

// TestMaxStringLen is flaky
func TestMaxStringLen(t *testing.T) {
	s := MustNewService(0, nil) // defaults to maxLenStringLimit

	for i := 0; i < 10; i++ {
		var a MaxStringLen
		err := s.FakeData(&a)
		assert.NoError(t, err, "\n%+v", err)

		assert.LenBetween(t, a.ColBlob, 1, maxLenStringLimit, "Field a.ColBlob")
		if a.ColLongtext1.Valid {
			assert.LenBetween(t, a.ColLongtext1.Data, 1, maxLenStringLimit, "Field ColLongtext1.String")
		} else {
			assert.Empty(t, a.ColLongtext1.Data, "Field ColLongtext1.String")
		}
		assert.LenBetween(t, a.ColLongtext2, 1, maxLenStringLimit, "Field ColLongtext2")
	}
}

type testServiceDateTime struct {
	ID           int32
	ColDate2     time.Time
	ColDatetime1 null.Time
}

func TestService_CustomFakeFunc_DateTime(t *testing.T) {
	// When i have for a field ColDate2 a custom fake func ... then why gets it
	// iterated to the next time.Time field? becuase ColDatetime1 has null.Time
	// which is later just Time because the pkg path is not prefiexed
	s := MustNewService(0, nil,
		WithTagFakeFunc("col_date2", func(maxLen int) (interface{}, error) {
			x := "2018-04-30"
			return x, nil
		}),
	)

	a := new(testServiceDateTime)
	err := s.FakeData(&a)
	assert.NoError(t, err, "\n%+v", err)

	// t.Logf("%#v", a)
	assert.Exactly(t, "2018-04-30 00:00:00", a.ColDate2.Format("2006-01-02 15:04:05"))
	assert.NotEmpty(t, a.ColDatetime1.String())
}

type testServiceDecimal struct {
	ID           int32
	ColDatetime2 time.Time
	ColPrice1    null.Decimal
	ColPrice2    null.Decimal
}

func TestService_CustomFakeFunc_Decimal(t *testing.T) {
	s := MustNewService(0, nil,
		WithTagFakeFunc("col_price1", func(maxLen int) (interface{}, error) {
			x := "123.4567"
			return x, nil
		}),
		WithTagFakeFuncAlias("col_price2", "col_price1"),
		WithTagFakeFunc("pseudo.testServiceDecimal.ID", func(maxLen int) (interface{}, error) {
			return 999, nil
		}),
	)
	for i := 0; i < 20; i++ {

		a := new(testServiceDecimal)
		err := s.FakeData(&a)
		assert.NoError(t, err, "\nIDX:%d %+v", i, err)

		// t.Logf("%#v", a)
		assert.Exactly(t, int32(999), a.ID, "IDX %d for field ID", i)
		assert.Exactly(t, "123.4567", a.ColPrice1.String(), "IDX %d for field ColPrice1", i)
		assert.Exactly(t, "123.4567", a.ColPrice2.String(), "IDX %d for field ColPrice2", i)
	}
}
