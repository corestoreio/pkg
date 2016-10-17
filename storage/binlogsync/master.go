package binlogsync

import (
	"sync"
	"time"

	"github.com/corestoreio/csfw/util/errors"
	"github.com/siddontang/go-mysql/mysql"
)

type masterInfo struct {
	FileName string
	Position uint32

	// mu protects the underlying storage engine to save the position and the filename.
	mu sync.Mutex

	lastSaveTime time.Time
}

func loadMasterInfo(exec mysql.Executer) (*masterInfo, error) {
	// todo(CyS) refactor
	res, err := exec.Execute("SHOW MASTER STATUS")
	if err != nil {
		return nil, errors.Wrap(err, "[binlogsync] Failed to execute SHOW MASTER STATUS")
	}

	name, err := res.GetString(0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "[binlogsync] Failed to fetch first row with 1st column")
	}

	pos, err := res.GetUint(0, 1)
	if err != nil {
		return nil, errors.Wrap(err, "[binlogsync] Failed to fetch first row with 2nd column")
	}

	m := &masterInfo{
		FileName: name,
		Position: uint32(pos),
	}
	return m, nil
}

// Save todo: implement saving
func (m *masterInfo) Save(force bool) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	n := time.Now()
	if !force && n.Sub(m.lastSaveTime) < time.Second {
		return nil
	}

	//var buf bytes.Buffer
	//e := toml.NewEncoder(&buf)
	//
	//e.Encode(m)
	//
	//var err error
	//if err = ioutil2.WriteFileAtomic(m.name, buf.Bytes(), 0644); err != nil {
	//	log.Errorf("canal save master info to file %s err %v", m.name, err)
	//}

	m.lastSaveTime = n

	return nil
}

func (m *masterInfo) Update(name string, pos uint32) {
	m.mu.Lock()
	m.FileName = name
	m.Position = pos
	m.mu.Unlock()
}

func (m *masterInfo) Pos() mysql.Position {
	var pos mysql.Position
	m.mu.Lock()
	pos.Name = m.FileName
	pos.Pos = m.Position
	m.mu.Unlock()

	return pos
}

func (m *masterInfo) Close() {
	m.Save(true)
}
