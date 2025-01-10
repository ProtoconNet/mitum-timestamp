package digest

import (
	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	utilc "github.com/ProtoconNet/mitum-currency/v3/digest/util"
	"github.com/ProtoconNet/mitum-timestamp/state"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	utilm "github.com/ProtoconNet/mitum2/util"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	DefaultColNameTimeStamp = "digest_ts"
)

func TimestampDesign(st *cdigest.Database, contract string) (types.Design, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("isItem", false)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var sta base.State
	if err := st.MongoClient().GetByFilter(
		DefaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, st.Encoders())
			if err != nil {
				return err
			}
			sta = i
			return nil
		},
		opt,
	); err != nil {
		return types.Design{}, nil, utilm.ErrNotFound.WithMessage(err, "timestamp design by contract account %v", contract)
	}

	if sta != nil {
		de, err := state.GetDesignFromState(sta)
		if err != nil {
			return types.Design{}, nil, err
		}
		return de, sta, nil
	} else {
		return types.Design{}, nil, errors.Errorf("state is nil")
	}
}

func TimestampItem(db *cdigest.Database, contract, project string, idx uint64) (types.Item, base.State, error) {
	filter := utilc.NewBSONFilter("contract", contract)
	filter = filter.Add("project_id", project)
	filter = filter.Add("timestamp_idx", idx)
	filter = filter.Add("isItem", true)
	q := filter.D()

	opt := options.FindOne().SetSort(
		utilc.NewBSONFilter("height", -1).D(),
	)
	var st base.State
	if err := db.MongoClient().GetByFilter(
		DefaultColNameTimeStamp,
		q,
		func(res *mongo.SingleResult) error {
			i, err := cdigest.LoadState(res.Decode, db.Encoders())
			if err != nil {
				return err
			}
			st = i
			return nil
		},
		opt,
	); err != nil {
		return types.Item{}, nil, utilm.ErrNotFound.WithMessage(err, "timestamp item by contract account %s, project %s, timestamp idx %s", contract, project, idx)
	}

	if st != nil {
		it, err := state.GetItemFromState(st)
		if err != nil {
			return types.Item{}, nil, err
		}
		return it, st, nil
	} else {
		return types.Item{}, nil, errors.Errorf("state is nil")
	}
}
