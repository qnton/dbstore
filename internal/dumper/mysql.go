package dumper

import (
	"database/sql"
	"errors"
	"os"
)

type Dumper struct {
	db   *sql.DB
	name string
	dir  string
}

func Register(db *sql.DB, dir, name string) (*Dumper, error) {
	if !isDir(dir) {
		return nil, errors.New("invalid directory")
	}

	return &Dumper{
		db:   db,
		name: name,
		dir:  dir,
	}, nil
}

func (d *Dumper) Close() error {
	defer func() {
		d.db = nil
	}()
	return d.db.Close()
}

func exists(p string) (bool, os.FileInfo) {
	f, err := os.Open(p)
	if err != nil {
		return false, nil
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return false, nil
	}
	return true, fi
}

func isDir(p string) bool {
	if e, fi := exists(p); e {
		return fi.Mode().IsDir()
	}
	return false
}
