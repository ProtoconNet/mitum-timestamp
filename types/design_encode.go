package types

import (
	"github.com/ProtoconNet/mitum-currency/v3/types"
	"github.com/ProtoconNet/mitum2/util/encoder"
	"github.com/ProtoconNet/mitum2/util/hint"
)

func (de *Design) unmarshal(
	_ encoder.Encoder,
	ht hint.Hint,
	svc string,
	prjs []string,
) error {
	de.BaseHinter = hint.NewBaseHinter(ht)
	de.service = types.ContractID(svc)
	de.projects = prjs

	return nil
}
