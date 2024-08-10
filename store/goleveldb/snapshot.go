package goleveldb

import (
	"github.com/opentoys/ledisdb/pkg/leveldb"
	"github.com/opentoys/ledisdb/store/driver"
)

type Snapshot struct {
	db  *DB
	snp *leveldb.Snapshot
}

func (s *Snapshot) Get(key []byte) ([]byte, error) {
	return s.snp.Get(key, s.db.iteratorOpts)
}

func (s *Snapshot) NewIterator() driver.IIterator {
	it := &Iterator{
		s.snp.NewIterator(nil, s.db.iteratorOpts),
	}
	return it
}

func (s *Snapshot) Close() {
	s.snp.Release()
}
