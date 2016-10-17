package binlogsync

import (
	"strconv"
	"strings"

	"github.com/corestoreio/csfw/util/errors"
)

// Position contains the binlog filename and the position based replication. It
// can parse into itself from a string.
type Position struct {
	Name string
	Pos  uint32
}

func (p Position) Compare(other Position) int {
	// First compare binlog name
	if p.Name > other.Name {
		return 1
	} else if p.Name < other.Name {
		return -1
	} else {
		// Same binlog file, compare position
		if p.Pos > other.Pos {
			return 1
		} else if p.Pos < other.Pos {
			return -1
		} else {
			return 0
		}
	}
}

func (p Position) String() string {
	return p.Name + ";" + strconv.FormatUint(uint64(p.Pos), 10)
}

func (p *Position) FromString(str string) error {
	c := strings.IndexByte(str, ';')
	p.Name = str[:c]
	pos, err := strconv.ParseUint(str[c+1:], 10, 32)
	if err != nil {
		return errors.Wrap(err, "[binlogsync] FromString.ParseUint")
	}
	p.Pos = uint32(pos)
	return nil
}
