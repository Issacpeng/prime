package models

import (
	_"path/filepath"

	_"github.com/wrench/db"
)

type Repository struct {
	Storage *RAWObjectStorage
	URL     string
	Branches []Ref
}

type Ref struct {
	Name string
	Id   Hash
}

/*
func Open(dirname string) (*Repository, error) {
	headfiles, err := ioutil.ReadDir(filepath.Join(dirname, "refs/heads"))
	if err != nil {
		return nil, err
	}
	repo := new(Repository)
	repo.Path = dirname
	for _, h := range headfiles {
		s, err := ioutil.ReadFile(filepath.Join(dirname, "refs/heads", h.Name()))
		if err != nil {
			return nil, err
		}
		ref := Ref{Name: h.Name()}
		s = bytes.TrimSpace(s)
		n, err := hex.Decode(ref.Id[:], s)
	    fmt.Printf("####### repo ref:%v #######\r\n", ref)
		if err != nil {
			return nil, err
		}
		if n < 20 {
			return nil, errTruncatedHead
		}
		repo.Branches = append(repo.Branches, ref)
	}
	fmt.Printf("####### repo :%v #######\r\n", repo)
	return repo, nil
}
*/