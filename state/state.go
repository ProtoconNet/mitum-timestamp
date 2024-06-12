package state

import (
	"fmt"
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"strconv"
	"strings"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	DesignStateValueHint    = hint.MustNewHint("mitum-timestamp-design-state-value-v0.0.1")
	TimeStampStateKeyPrefix = "timestamp"
	DesignStateKeySuffix    = "design"
)

func TimeStampStateKey(addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s", TimeStampStateKeyPrefix, addr.String())
}

type DesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewDesignStateValue(design types.Design) DesignStateValue {
	return DesignStateValue{
		BaseHinter: hint.NewBaseHinter(DesignStateValueHint),
		Design:     design,
	}
}

func (sv DesignStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv DesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid DesignStateValue")

	if err := sv.BaseHinter.IsValid(DesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv DesignStateValue) HashBytes() []byte {
	return sv.Design.Bytes()
}

func GetDesignFromState(st mitumbase.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, errors.Errorf("state value is nil")
	}

	d, ok := v.(DesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("expected DesignStateValue but %T", v)
	}

	return d.Design, nil
}

func IsDesignStateKey(key string) bool {
	return strings.HasSuffix(key, DesignStateKeySuffix)
}

func DesignStateKey(addr mitumbase.Address) string {
	return fmt.Sprintf("%s:%s", TimeStampStateKey(addr), DesignStateKeySuffix)
}

var (
	LastIdxStateValueHint = hint.MustNewHint("mitum-timestamp-last-idx-state-value-v0.0.1")
	LastIdxStateKeySuffix = "timestampIdx"
)

type LastIdxStateValue struct {
	hint.BaseHinter
	ProjectID string
	Index     uint64
}

func NewLastIdxStateValue(pid string, idx uint64) LastIdxStateValue {
	return LastIdxStateValue{
		BaseHinter: hint.NewBaseHinter(LastIdxStateValueHint),
		ProjectID:  pid,
		Index:      idx,
	}
}

func (sv LastIdxStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv LastIdxStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TimeStampLastIdxStateValue")

	if err := sv.BaseHinter.IsValid(LastIdxStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(sv.ProjectID) < 1 || len(sv.ProjectID) > types.MaxProjectIDLen {
		return common.ErrValOOR.Wrap(
			errors.Errorf("invalid projectID length %v < 1 or > %v", len(sv.ProjectID), types.MaxProjectIDLen))
	}

	return nil
}

func (sv LastIdxStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice([]byte(sv.ProjectID), util.Uint64ToBytes(sv.Index))
}

func GetLastIdxFromState(st mitumbase.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, errors.Errorf("state value is nil")
	}

	isv, ok := v.(LastIdxStateValue)
	if !ok {
		return 0, errors.Errorf("expected LastIdxStateValue but, %T", v)
	}

	return isv.Index, nil
}

func IsLastIdxStateKey(key string) bool {
	return strings.HasSuffix(key, LastIdxStateKeySuffix)
}

func LastIdxStateKey(addr mitumbase.Address, pid string) string {
	return fmt.Sprintf("%s:%s:%s", TimeStampStateKey(addr), pid, LastIdxStateKeySuffix)
}

var (
	ItemStateValueHint = hint.MustNewHint("mitum-timestamp-item-state-value-v0.0.1")
	ItemStateKeySuffix = "item"
)

type ItemStateValue struct {
	hint.BaseHinter
	Item types.Item
}

func NewItemStateValue(item types.Item) ItemStateValue {
	return ItemStateValue{
		BaseHinter: hint.NewBaseHinter(ItemStateValueHint),
		Item:       item,
	}
}

func (sv ItemStateValue) Hint() hint.Hint {
	return sv.BaseHinter.Hint()
}

func (sv ItemStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ItemStateValue")

	if err := sv.BaseHinter.IsValid(ItemStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sv.Item.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sv ItemStateValue) HashBytes() []byte {
	return sv.Item.Bytes()
}

func GetItemFromState(st mitumbase.State) (types.Item, error) {
	v := st.Value()
	if v == nil {
		return types.Item{}, errors.Errorf("State value is nil")
	}

	ts, ok := v.(ItemStateValue)
	if !ok {
		return types.Item{}, common.ErrTypeMismatch.Wrap(errors.Errorf("expected ItemStateValue found, %T", v))
	}

	return ts.Item, nil
}

func IsItemStateKey(key string) bool {
	return strings.HasSuffix(key, ItemStateKeySuffix)
}

func ItemStateKey(addr mitumbase.Address, pid string, index uint64) string {
	return fmt.Sprintf("%s:%s:%s:%s", TimeStampStateKey(addr), pid, strconv.FormatUint(index, 10), ItemStateKeySuffix)
}
