package digest

import (
	"net/http"
	"strconv"
	"time"

	cdigest "github.com/ProtoconNet/mitum-currency/v3/digest"
	"github.com/ProtoconNet/mitum-timestamp/types"
	"github.com/ProtoconNet/mitum2/base"
	"github.com/gorilla/mux"
)

func (hd *Handlers) handleTimeStampDesign(w http.ResponseWriter, r *http.Request) {
	cacheKey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cacheKey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	if v, err, shared := hd.rg.Do(cacheKey, func() (interface{}, error) {
		return hd.handleTimeStampDesignInGroup(contract)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cacheKey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampDesignInGroup(contract string) ([]byte, error) {
	var de types.Design
	var st base.State

	de, st, err := TimestampDesign(hd.database, contract)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStampDesign(contract, de, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildTimeStampDesign(contract string, de types.Design, st base.State) (cdigest.Hal, error) {
	h, err := hd.combineURL(HandlerPathTimeStampDesign, "contract", contract)
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(de, cdigest.NewHalLink(h, nil))

	h, err = hd.combineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}

func (hd *Handlers) handleTimeStampItem(w http.ResponseWriter, r *http.Request) {
	cachekey := cdigest.CacheKeyPath(r)
	if err := cdigest.LoadFromCache(hd.cache, cachekey, w); err == nil {
		return
	}

	contract, err, status := cdigest.ParseRequest(w, r, "contract")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	project, err, status := cdigest.ParseRequest(w, r, "project_id")
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, status)

		return
	}

	s, found := mux.Vars(r)["timestamp_idx"]
	if !found {
		cdigest.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}
	idx, err := parseIdxFromPath(s)
	if err != nil {
		cdigest.HTTP2ProblemWithError(w, err, http.StatusBadRequest)

		return
	}

	if v, err, shared := hd.rg.Do(cachekey, func() (interface{}, error) {
		return hd.handleTimeStampItemInGroup(contract, project, idx)
	}); err != nil {
		cdigest.HTTP2HandleError(w, err)
	} else {
		cdigest.HTTP2WriteHalBytes(hd.encoder, w, v.([]byte), http.StatusOK)

		if !shared {
			cdigest.HTTP2WriteCache(w, cachekey, time.Second*3)
		}
	}
}

func (hd *Handlers) handleTimeStampItemInGroup(contract, project string, idx uint64) ([]byte, error) {
	var it types.Item
	var st base.State

	it, st, err := TimestampItem(hd.database, contract, project, idx)
	if err != nil {
		return nil, err
	}

	i, err := hd.buildTimeStampItem(contract, it, st)
	if err != nil {
		return nil, err
	}
	return hd.encoder.Marshal(i)
}

func (hd *Handlers) buildTimeStampItem(contract string, it types.Item, st base.State) (cdigest.Hal, error) {
	h, err := hd.combineURL(
		HandlerPathTimeStampItem,
		"contract", contract, "project_id", it.ProjectID(), "timestamp_idx",
		strconv.FormatUint(it.TimestampID(), 10))
	if err != nil {
		return nil, err
	}

	var hal cdigest.Hal
	hal = cdigest.NewBaseHal(it, cdigest.NewHalLink(h, nil))

	h, err = hd.combineURL(cdigest.HandlerPathBlockByHeight, "height", st.Height().String())
	if err != nil {
		return nil, err
	}
	hal = hal.AddLink("block", cdigest.NewHalLink(h, nil))

	for i := range st.Operations() {
		h, err := hd.combineURL(cdigest.HandlerPathOperation, "hash", st.Operations()[i].String())
		if err != nil {
			return nil, err
		}
		hal = hal.AddLink("operations", cdigest.NewHalLink(h, nil))
	}

	return hal, nil
}
