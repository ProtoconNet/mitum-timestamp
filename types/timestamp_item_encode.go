package types

import (
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (t *TimeStampItem) unmarshal(
	ht hint.Hint,
	pid string,
	rqts,
	rsts,
	tsid uint64,
	data string,
) error {
	t.BaseHinter = hint.NewBaseHinter(ht)
	t.projectID = pid
	t.requestTimeStamp = rqts
	t.responseTimeStamp = rsts
	t.timestampID = tsid
	t.data = data

	return nil
}
