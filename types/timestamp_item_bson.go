package types

import (
	bsonenc "github.com/ProtoconNet/mitum-currency/v3/digest/util/bson"
	"github.com/ProtoconNet/mitum2/util"
	"github.com/ProtoconNet/mitum2/util/hint"
	"go.mongodb.org/mongo-driver/bson"
)

func (t Item) MarshalBSON() ([]byte, error) {
	return bsonenc.Marshal(bson.M{
		"_hint":              t.Hint().String(),
		"project_id":         t.projectID,
		"request_timestamp":  t.requestTimeStamp,
		"response_timestamp": t.responseTimeStamp,
		"timestamp_idx":      t.timestampIdx,
		"data":               t.data,
	})
}

type TimeStampItemBSONUnmarshaler struct {
	Hint              string `bson:"_hint"`
	ProjectID         string `bson:"project_id"`
	RequestTimeStamp  uint64 `bson:"request_timestamp"`
	ResponseTimeStamp uint64 `bson:"response_timestamp"`
	TimeStampIdx      uint64 `bson:"timestamp_idx"`
	Data              string `bson:"data"`
}

func (t *Item) DecodeBSON(b []byte, enc *bsonenc.Encoder) error {
	e := util.StringError("decode bson of TimeStampItem")

	var u TimeStampItemBSONUnmarshaler
	if err := enc.Unmarshal(b, &u); err != nil {
		return e.Wrap(err)
	}

	ht, err := hint.ParseHint(u.Hint)
	if err != nil {
		return e.Wrap(err)
	}

	return t.unmarshal(ht, u.ProjectID, u.RequestTimeStamp, u.ResponseTimeStamp, u.TimeStampIdx, u.Data)
}
