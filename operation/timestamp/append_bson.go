package timestamp

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ProtoconNet/mitum-currency/v3/common"
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact AppendFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":             fact.Hint().String(),
			"hash":              fact.BaseFact.Hash().String(),
			"token":             fact.BaseFact.Token(),
			"sender":            fact.sender,
			"target":            fact.target,
			"projectid":         fact.projectID,
			"request_timestamp": fact.requestTimeStamp,
			"data":              fact.data,
			"currency":          fact.currency,
		},
	)
}

type AppendFactBSONUnmarshaler struct {
	Hint             string `bson:"_hint"`
	Sender           string `bson:"sender"`
	Target           string `bson:"target"`
	ProjectID        string `bson:"projectid"`
	RequestTimeStamp uint64 `bson:"request_timestamp"`
	Data             string `bson:"data"`
	Currency         string `bson:"currency"`
}

func (fact *AppendFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of AppendFact")

	var u common.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf AppendFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unmarshal(enc, uf.Sender, uf.Target, uf.ProjectID, uf.RequestTimeStamp, uf.Data, uf.Currency)
}

func (op Append) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		})
}

func (op *Append) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of Mint")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
