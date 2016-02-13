package phpserialize

const (
	TOKEN_NULL              rune = 'N'
	TOKEN_BOOL              rune = 'b'
	TOKEN_INT               rune = 'i'
	TOKEN_FLOAT             rune = 'd'
	TOKEN_STRING            rune = 's'
	TOKEN_ARRAY             rune = 'a'
	TOKEN_OBJECT            rune = 'O'
	TOKEN_OBJECT_SERIALIZED rune = 'C'
	TOKEN_REFERENCE         rune = 'R'
	TOKEN_REFERENCE_OBJECT  rune = 'r'
	TOKEN_SPL_ARRAY         rune = 'x'
	TOKEN_SPL_ARRAY_MEMBERS rune = 'm'

	SEPARATOR_VALUE_TYPE rune = ':'
	SEPARATOR_VALUES     rune = ';'

	DELIMITER_STRING_LEFT  rune = '"'
	DELIMITER_STRING_RIGHT rune = '"'
	DELIMITER_OBJECT_LEFT  rune = '{'
	DELIMITER_OBJECT_RIGHT rune = '}'

	FORMATTER_FLOAT     byte = 'g'
	FORMATTER_PRECISION int  = 17
)

func NewPhpObject(className string) *PhpObject {
	return &PhpObject{
		className: className,
		members:   PhpArray{},
	}
}

type SerializedDecodeFunc func([]byte) (PhpValue, error)

type SerializedEncodeFunc func(PhpValue) (string, error)

type PhpValue interface{}

type PhpArray map[PhpValue]PhpValue

type PhpSlice []PhpValue

type PhpObject struct {
	className string
	members   PhpArray
}

func (p *PhpObject) GetClassName() string {
	return p.className
}

func (p *PhpObject) SetClassName(name string) *PhpObject {
	p.className = name
	return p
}

func (p *PhpObject) GetMembers() PhpArray {
	return p.members
}

func (p *PhpObject) SetMembers(members PhpArray) *PhpObject {
	p.members = members
	return p
}

func (p *PhpObject) GetPrivate(name string) (v PhpValue, ok bool) {
	v, ok = p.members["\x00"+p.className+"\x00"+name]
	return
}

func (p *PhpObject) SetPrivate(name string, value PhpValue) *PhpObject {
	p.members["\x00"+p.className+"\x00"+name] = value
	return p
}

func (p *PhpObject) GetProtected(name string) (v PhpValue, ok bool) {
	v, ok = p.members["\x00*\x00"+name]
	return
}

func (p *PhpObject) SetProtected(name string, value PhpValue) *PhpObject {
	p.members["\x00*\x00"+name] = value
	return p
}

func (p *PhpObject) GetPublic(name string) (v PhpValue, ok bool) {
	v, ok = p.members[name]
	return
}

func (p *PhpObject) SetPublic(name string, value PhpValue) *PhpObject {
	p.members[name] = value
	return p
}

func NewPhpObjectSerialized(className string) *PhpObjectSerialized {
	return &PhpObjectSerialized{
		className: className,
	}
}

type PhpObjectSerialized struct {
	className string
	data      string
	value     PhpValue
}

func (p *PhpObjectSerialized) GetClassName() string {
	return p.className
}

func (p *PhpObjectSerialized) SetClassName(name string) *PhpObjectSerialized {
	p.className = name
	return p
}

func (p *PhpObjectSerialized) GetData() string {
	return p.data
}

func (p *PhpObjectSerialized) SetData(data string) *PhpObjectSerialized {
	p.data = data
	return p
}

func (p *PhpObjectSerialized) GetValue() PhpValue {
	return p.value
}

func (p *PhpObjectSerialized) SetValue(value PhpValue) *PhpObjectSerialized {
	p.value = value
	return p
}

func NewPhpSplArray(array, properties PhpValue) *PhpSplArray {
	if array == nil {
		array = make(PhpArray)
	}

	if properties == nil {
		properties = make(PhpArray)
	}

	return &PhpSplArray{
		array:      array,
		properties: properties,
	}
}

type PhpSplArray struct {
	flags      int
	array      PhpValue
	properties PhpValue
}

func (p *PhpSplArray) GetFlags() int {
	return p.flags
}

func (p *PhpSplArray) SetFlags(value int) {
	p.flags = value
}

func (p *PhpSplArray) GetArray() PhpValue {
	return p.array
}

func (p *PhpSplArray) SetArray(value PhpValue) {
	p.array = value
}

func (p *PhpSplArray) GetProperties() PhpValue {
	return p.properties
}

func (p *PhpSplArray) SetProperties(value PhpValue) {
	p.properties = value
}
