package digest

import (
	"github.com/ProtoconNet/mitum-timestamp/state"
	mitumbase "github.com/ProtoconNet/mitum2/base"
	"go.mongodb.org/mongo-driver/mongo"
)

func (bs *BlockSession) prepareTimeStamps() error {
	if len(bs.sts) < 1 {
		return nil
	}

	var timestampModels []mongo.WriteModel
	for i := range bs.sts {
		st := bs.sts[i]
		switch {
		case state.IsDesignStateKey(st.Key()):
			j, err := bs.handleTimeStampDesignState(st)
			if err != nil {
				return err
			}
			timestampModels = append(timestampModels, j...)
		case state.IsItemStateKey(st.Key()):
			j, err := bs.handleTimeStampItemState(st)
			if err != nil {
				return err
			}
			timestampModels = append(timestampModels, j...)
		default:
			continue
		}
	}

	bs.timestampModels = timestampModels

	return nil
}

func (bs *BlockSession) handleTimeStampDesignState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if serviceDesignDoc, err := NewDesignDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(serviceDesignDoc),
		}, nil
	}
}

func (bs *BlockSession) handleTimeStampItemState(st mitumbase.State) ([]mongo.WriteModel, error) {
	if TimeStampItemDoc, err := NewItemDoc(st, bs.st.Encoder()); err != nil {
		return nil, err
	} else {
		return []mongo.WriteModel{
			mongo.NewInsertOneModel().SetDocument(TimeStampItemDoc),
		}, nil
	}
}
