package types

import (
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	MaxProjectIDLen = 10
	MaxDataLen      = 1024
)

var ItemHint = hint.MustNewHint("mitum-timestamp-item-v0.0.1")

type Item struct {
	hint.BaseHinter
	projectID         string
	requestTimeStamp  uint64
	responseTimeStamp uint64
	timestampIdx      uint64
	data              string
}

func NewItem(
	pid string,
	reqTS,
	resTS,
	tidx uint64,
	data string,
) Item {
	return Item{
		BaseHinter:        hint.NewBaseHinter(ItemHint),
		projectID:         pid,
		requestTimeStamp:  reqTS,
		responseTimeStamp: resTS,
		timestampIdx:      tidx,
		data:              data,
	}
}

func (t Item) IsValid([]byte) error {
	if len(t.projectID) < 1 || len(t.projectID) > MaxProjectIDLen {
		return errors.Errorf("invalid projectID length %v < 1 or > %v", len(t.projectID), MaxProjectIDLen)
	}

	if !currencytypes.ReSpcecialChar.Match([]byte(t.projectID)) {
		return util.ErrInvalid.Errorf("invalid projectID due to the inclusion of special characters")
	}

	if len(t.data) < 1 || len(t.data) > MaxDataLen {
		return errors.Errorf("invalid data length %v < 1 or > %v", len(t.data), MaxDataLen)
	}

	return nil
}

func (t Item) Bytes() []byte {
	return util.ConcatBytesSlice(
		[]byte(t.projectID),
		util.Uint64ToBytes(t.requestTimeStamp),
		util.Uint64ToBytes(t.responseTimeStamp),
		util.Uint64ToBytes(t.timestampIdx),
		[]byte(t.data),
	)
}

func (t Item) ProjectID() string {
	return t.projectID
}

func (t Item) RequestTimeStamp() uint64 {
	return t.requestTimeStamp
}

func (t Item) ResponseTimeStamp() uint64 {
	return t.responseTimeStamp
}

func (t Item) TimestampID() uint64 {
	return t.timestampIdx
}

func (t Item) Data() string {
	return t.data
}

func (t Item) Equal(ct Item) bool {
	if t.projectID != ct.projectID {
		return false
	}

	if t.requestTimeStamp != ct.requestTimeStamp {
		return false
	}

	if t.responseTimeStamp != ct.responseTimeStamp {
		return false
	}

	if t.timestampIdx != ct.timestampIdx {
		return false
	}

	if t.data != ct.data {
		return false
	}

	return true
}
