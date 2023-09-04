package timestamp

import (
	"github.com/ProtoconNet/mitum-currency/v3/common"
	"go.mongodb.org/mongo-driver/bson"

	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"github.com/ProtoconNet/mitum2/util/valuehash"
)

func (fact CreateServiceFact) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint":    fact.Hint().String(),
			"hash":     fact.BaseFact.Hash().String(),
			"token":    fact.BaseFact.Token(),
			"sender":   fact.sender,
			"target":   fact.target,
			"service":  fact.service,
			"currency": fact.currency,
		},
	)
}

type CreateServiceFactBSONUnmarshaler struct {
	Hint     string `bson:"_hint"`
	Sender   string `bson:"sender"`
	Target   string `bson:"target"`
	Service  string `bson:"service"`
	Currency string `bson:"currency"`
}

func (fact *CreateServiceFact) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of CreateServiceFact")

	var u common.BaseFactBSONUnmarshaler

	err := enc.Unmarshal(b, &u)
	if err != nil {
		return e.Wrap(err)
	}

	fact.BaseFact.SetHash(valuehash.NewBytesFromString(u.Hash))
	fact.BaseFact.SetToken(u.Token)

	var uf CreateServiceFactBSONUnmarshaler
	if err := bson.Unmarshal(b, &uf); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(uf.Hint)
	if err != nil {
		return e.Wrap(err)
	}
	fact.BaseHinter = hint.NewBaseHinter(ht)

	return fact.unmarshal(enc, uf.Sender, uf.Target, uf.Service, uf.Currency)
}

func (op CreateService) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(
		bson.M{
			"_hint": op.Hint().String(),
			"hash":  op.Hash().String(),
			"fact":  op.Fact(),
			"signs": op.Signs(),
		},
	)
}

func (op *CreateService) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("failed to decode bson of CreateService")

	var ubo common.BaseOperation
	if err := ubo.DecodeBSON(b, enc); err != nil {
		return e.Wrap(err)
	}

	op.BaseOperation = ubo

	return nil
}
