package types

import (
	"bytes"
	"sort"

	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

var DesignHint = hint.MustNewHint("mitum-timestamp-design-v0.0.1")

type Design struct {
	hint.BaseHinter
	projects []string
}

func NewDesign(projects ...string) Design {
	return Design{
		BaseHinter: hint.NewBaseHinter(DesignHint),
		projects:   projects,
	}
}

func (de Design) IsValid([]byte) error {
	if err := util.CheckIsValiders(nil, false,
		de.BaseHinter,
	); err != nil {
		return err
	}

	return nil
}

func (de Design) Bytes() []byte {
	bytesArray := make([][]byte, len(de.projects))

	sort.Slice(de.projects, func(i, j int) bool {
		return bytes.Compare([]byte(de.projects[j]), []byte(de.projects[i])) < 0
	})

	for i := range de.projects {
		bytesArray[i] = []byte(de.projects[i])
	}

	return util.ConcatBytesSlice(bytesArray...)
}

func (de Design) Hash() util.Hash {
	return de.GenerateHash()
}

func (de Design) GenerateHash() util.Hash {
	return valuehash.NewSHA256(de.Bytes())
}

func (de Design) Projects() []string {
	return de.projects
}

func (de *Design) AddProject(project string) {
	for i := range de.projects {
		if de.projects[i] == project {
			return
		}
	}
	projects := append(de.projects, project)
	de.projects = projects
}

func (de Design) Equal(cd Design) bool {
	if len(de.projects) != len(cd.projects) {
		return false
	}

	sort.Slice(de.projects, func(i, j int) bool {
		return bytes.Compare([]byte(de.projects[i]), []byte(de.projects[j])) < 0
	})

	for i := range de.projects {
		if de.projects[i] != cd.projects[i] {
			return false
		}
	}

	return true
}
