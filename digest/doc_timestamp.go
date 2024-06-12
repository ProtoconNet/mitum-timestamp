package digest

import (
	mongodb "github.com/ProtoconNet/mitum-currency/v3/digest/mongodb"
	bson "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	cstate "github.com/ProtoconNet/mitum-currency/v3/state"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/ProtoconNet/mitum2/util/encoder"
)

type DesignDoc struct {
	mongodb.BaseDoc
	st     base.State
	design types.Design
}

// NewDesignDoc get the State of TimeStamp Design
func NewDesignDoc(st base.State, enc encoder.Encoder) (DesignDoc, error) {
	design, err := state.GetDesignFromState(st)

	if err != nil {
		return DesignDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return DesignDoc{}, err
	}

	return DesignDoc{
		BaseDoc: b,
		st:      st,
		design:  design,
	}, nil
}

func (doc DesignDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.TimeStampStateKeyPrefix, 3)

	m["contract"] = parsedKey[1]
	m["height"] = doc.st.Height()
	m["isItem"] = false

	return bson.Marshal(m)
}

type ItemDoc struct {
	mongodb.BaseDoc
	st   base.State
	item types.Item
}

func NewItemDoc(st base.State, enc encoder.Encoder) (ItemDoc, error) {
	item, err := state.GetItemFromState(st)
	if err != nil {
		return ItemDoc{}, err
	}

	b, err := mongodb.NewBaseDoc(nil, st, enc)
	if err != nil {
		return ItemDoc{}, err
	}

	return ItemDoc{
		BaseDoc: b,
		st:      st,
		item:    item,
	}, nil
}

func (doc ItemDoc) MarshalBSON() ([]byte, error) {
	m, err := doc.BaseDoc.M()
	if err != nil {
		return nil, err
	}

	parsedKey, err := cstate.ParseStateKey(doc.st.Key(), state.TimeStampStateKeyPrefix, 5)
	if err != nil {
		return nil, err
	}

	m["contract"] = parsedKey[1]
	m["project_id"] = doc.item.ProjectID()
	m["timestamp_idx"] = doc.item.TimestampID()
	m["height"] = doc.st.Height()
	m["isItem"] = true

	return bson.Marshal(m)
}
