package dbr

type StringWriter interface {
	WriteString(s string) (n int, err error)
	WriteRune(r rune) (n int, err error)
}

//type Buffer interface {
//	StringWriter
//
//	WriteValue(v ...interface{}) (err error)
//	Value() []interface{}
//}
//
//type buffer struct {
//	StringWriter
//	v []interface{}
//}

//func NewBuffer() StringWriter {
//	return &buffer{
//		StringWriter: bufferpool.Get(),
//	}
//}
//
//func (b *buffer) WriteValue(v ...interface{}) error {
//	b.v = append(b.v, v...)
//	return nil
//}
//
//func (b *buffer) Value() []interface{} {
//	return b.v
//}
//
//func PutBuffer(buf Buffer) {
//	if b, ok := buf.(*buffer); ok {
//		if bb, ok := b.StringWriter.(*bytes.Buffer); ok {
//			bufferpool.Put(bb)
//		} else {
//			panic(fmt.Sprintf("*bytes.Buffer not found in %#v", buf))
//		}
//	} else {
//		panic(fmt.Sprintf("*buffer not found in %#v", buf))
//	}
//}
