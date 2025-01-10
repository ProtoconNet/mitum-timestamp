package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var timestampIndexModels = []mongo.IndexModel{
	{
		Keys: bson.D{
			bson.E{Key: "contract", Value: 1},
			bson.E{Key: "height", Value: -1},
			bson.E{Key: "isItem", Value: 1},
			bson.E{Key: "project_id", Value: 1},
			bson.E{Key: "timestamp_idx", Value: 1}},
		Options: options.Index().
			SetName(cdigest.IndexPrefix + "timestamp_idx_contract_height_isItem_projectID"),
	},
}

var DefaultIndexes = cdigest.DefaultIndexes

func init() {
	DefaultIndexes[DefaultColNameTimeStamp] = timestampIndexModels
}
