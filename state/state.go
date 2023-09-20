package state

import (
	"fmt"
	statecurrency "github.com/ProtoconNet/mitum-currency/v3/state/currency"
	currencytypes "github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"strconv"
	"strings"

	mitumbase "github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/pkg/errors"
)

var (
	StateKeyTimeStampPrefix     = "timestamp:"
	ServiceDesignStateValueHint = hint.MustNewHint("mitum-timestamp-service-design-state-value-v0.0.1")
	StateKeyServiceDesignSuffix = ":service"
)

func StateKeyTimeStampService(addr mitumbase.Address) string {
	return fmt.Sprintf("%s%s", StateKeyTimeStampPrefix, addr.String())
}

type ServiceDesignStateValue struct {
	hint.BaseHinter
	Design types.Design
}

func NewServiceDesignStateValue(design types.Design) ServiceDesignStateValue {
	return ServiceDesignStateValue{
		BaseHinter: hint.NewBaseHinter(ServiceDesignStateValueHint),
		Design:     design,
	}
}

func (sd ServiceDesignStateValue) Hint() hint.Hint {
	return sd.BaseHinter.Hint()
}

func (sd ServiceDesignStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid ServiceDesignStateValue")

	if err := sd.BaseHinter.IsValid(ServiceDesignStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := sd.Design.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (sd ServiceDesignStateValue) HashBytes() []byte {
	return sd.Design.Bytes()
}

func StateServiceDesignValue(st mitumbase.State) (types.Design, error) {
	v := st.Value()
	if v == nil {
		return types.Design{}, util.ErrNotFound.Errorf("service design not found in State")
	}

	d, ok := v.(ServiceDesignStateValue)
	if !ok {
		return types.Design{}, errors.Errorf("invalid service design value found, %T", v)
	}

	return d.Design, nil
}

func IsStateServiceDesignKey(key string) bool {
	return strings.HasSuffix(key, StateKeyServiceDesignSuffix)
}

func StateKeyServiceDesign(addr mitumbase.Address) string {
	return fmt.Sprintf("%s%s", StateKeyTimeStampService(addr), StateKeyServiceDesignSuffix)
}

type StateValueMerger struct {
	*mitumbase.BaseStateValueMerger
}

func NewStateValueMerger(height mitumbase.Height, key string, st mitumbase.State) *StateValueMerger {
	s := &StateValueMerger{
		BaseStateValueMerger: mitumbase.NewBaseStateValueMerger(height, key, st),
	}

	return s
}

func NewStateMergeValue(key string, stv mitumbase.StateValue) mitumbase.StateMergeValue {
	StateValueMergerFunc := func(height mitumbase.Height, st mitumbase.State) mitumbase.StateValueMerger {
		return NewStateValueMerger(height, key, st)
	}

	return mitumbase.NewBaseStateMergeValue(
		key,
		stv,
		StateValueMergerFunc,
	)
}

var (
	TimeStampLastIndexStateValueHint = hint.MustNewHint("mitum-timestamp-last-index-state-value-v0.0.1")
	StateKeyProjectLastIndexSuffix   = ":timestampidx"
)

type TimeStampLastIndexStateValue struct {
	hint.BaseHinter
	ProjectID string
	Index     uint64
}

func NewTimeStampLastIndexStateValue(pid string, index uint64) TimeStampLastIndexStateValue {
	return TimeStampLastIndexStateValue{
		BaseHinter: hint.NewBaseHinter(TimeStampLastIndexStateValueHint),
		ProjectID:  pid,
		Index:      index,
	}
}

func (ti TimeStampLastIndexStateValue) Hint() hint.Hint {
	return ti.BaseHinter.Hint()
}

func (ti TimeStampLastIndexStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TimeStampLastIndexStateValue")

	if err := ti.BaseHinter.IsValid(TimeStampLastIndexStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if len(ti.ProjectID) < 1 || len(ti.ProjectID) > types.MaxProjectIDLen {
		return errors.Errorf("invalid projectID length %v < 1 or > %v", len(ti.ProjectID), types.MaxProjectIDLen)
	}

	return nil
}

func (ti TimeStampLastIndexStateValue) HashBytes() []byte {
	return util.ConcatBytesSlice([]byte(ti.ProjectID), util.Uint64ToBytes(ti.Index))
}

func StateTimeStampLastIndexValue(st mitumbase.State) (uint64, error) {
	v := st.Value()
	if v == nil {
		return 0, util.ErrNotFound.Errorf("collection last nft index not found in State")
	}

	isv, ok := v.(TimeStampLastIndexStateValue)
	if !ok {
		return 0, errors.Errorf("invalid collection last nft index value found, %T", v)
	}

	return isv.Index, nil
}

func IsStateTimeStampLastIndexKey(key string) bool {
	return strings.HasSuffix(key, StateKeyProjectLastIndexSuffix)
}

func StateKeyTimeStampLastIndex(addr mitumbase.Address, pid string) string {
	return fmt.Sprintf("%s:%s%s", StateKeyTimeStampService(addr), pid, StateKeyProjectLastIndexSuffix)
}

var (
	TimeStampItemStateValueHint = hint.MustNewHint("mitum-timestamp-item-state-value-v0.0.1")
	StateKeyTimeStampItemSuffix = ":timestampitem"
)

type TimeStampItemStateValue struct {
	hint.BaseHinter
	TimeStampItem types.TimeStampItem
}

func NewTimeStampItemStateValue(item types.TimeStampItem) TimeStampItemStateValue {
	return TimeStampItemStateValue{
		BaseHinter:    hint.NewBaseHinter(TimeStampItemStateValueHint),
		TimeStampItem: item,
	}
}

func (ts TimeStampItemStateValue) Hint() hint.Hint {
	return ts.BaseHinter.Hint()
}

func (ts TimeStampItemStateValue) IsValid([]byte) error {
	e := util.ErrInvalid.Errorf("invalid TimeStampItemStateValue")

	if err := ts.BaseHinter.IsValid(TimeStampItemStateValueHint.Type().Bytes()); err != nil {
		return e.Wrap(err)
	}

	if err := ts.TimeStampItem.IsValid(nil); err != nil {
		return e.Wrap(err)
	}

	return nil
}

func (ts TimeStampItemStateValue) HashBytes() []byte {
	return ts.TimeStampItem.Bytes()
}

func StateTimeStampItemValue(st mitumbase.State) (types.TimeStampItem, error) {
	v := st.Value()
	if v == nil {
		return types.TimeStampItem{}, util.ErrNotFound.Errorf("TimeStampItem not found in State")
	}

	ts, ok := v.(TimeStampItemStateValue)
	if !ok {
		return types.TimeStampItem{}, errors.Errorf("invalid TimeStampItem value found, %T", v)
	}

	return ts.TimeStampItem, nil
}

func IsStateTimeStampItemKey(key string) bool {
	return strings.HasSuffix(key, StateKeyTimeStampItemSuffix)
}

func StateKeyTimeStampItem(addr mitumbase.Address, pid string, index uint64) string {
	return fmt.Sprintf("%s:%s:s%s%s", StateKeyTimeStampService(addr), pid, strconv.FormatUint(index, 10), StateKeyTimeStampItemSuffix)
}

func ParseStateKey(key string) ([]string, error) {
	parsedKey := strings.Split(key, ":")
	if parsedKey[0] != StateKeyTimeStampPrefix[:len(StateKeyTimeStampPrefix)-1] {
		return nil, errors.Errorf("State Key not include TimeStampPrefix, %s", parsedKey)
	}
	if len(parsedKey) < 3 {
		return nil, errors.Errorf("parsing State Key string failed, %s", parsedKey)
	} else {
		return parsedKey, nil
	}
}

func checkExistsState(
	key string,
	getState mitumbase.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case !found:
		return mitumbase.NewBaseOperationProcessReasonError("state, %q does not exist", key)
	default:
		return nil
	}
}

func checkNotExistsState(
	key string,
	getState mitumbase.GetStateFunc,
) error {
	switch _, found, err := getState(key); {
	case err != nil:
		return err
	case found:
		return mitumbase.NewBaseOperationProcessReasonError("state, %q already exists", key)
	default:
		return nil
	}
}

func existsState(
	k,
	name string,
	getState mitumbase.GetStateFunc,
) (mitumbase.State, error) {
	switch st, found, err := getState(k); {
	case err != nil:
		return nil, err
	case !found:
		return nil, mitumbase.NewBaseOperationProcessReasonError("%s does not exist", name)
	default:
		return st, nil
	}
}

func notExistsState(
	k,
	name string,
	getState mitumbase.GetStateFunc,
) (mitumbase.State, error) {
	var st mitumbase.State
	switch _, found, err := getState(k); {
	case err != nil:
		return nil, err
	case found:
		return nil, mitumbase.NewBaseOperationProcessReasonError("%s already exists", name)
	case !found:
		st = mitumbase.NewBaseState(mitumbase.NilHeight, k, nil, nil, nil)
	}
	return st, nil
}

func existsCurrencyPolicy(cid currencytypes.CurrencyID, getStateFunc mitumbase.GetStateFunc) (currencytypes.CurrencyPolicy, error) {
	var policy currencytypes.CurrencyPolicy

	switch st, found, err := getStateFunc(statecurrency.StateKeyCurrencyDesign(cid)); {
	case err != nil:
		return currencytypes.CurrencyPolicy{}, err
	case !found:
		return currencytypes.CurrencyPolicy{}, errors.Errorf("currency not found, %v", cid)
	default:
		design, ok := st.Value().(statecurrency.CurrencyDesignStateValue) //nolint:forcetypeassert //...
		if !ok {
			return currencytypes.CurrencyPolicy{}, errors.Errorf("expected CurrencyDesignStateValue, not %T", st.Value())
		}
		policy = design.CurrencyDesign.Policy()
	}

	return policy, nil
}
