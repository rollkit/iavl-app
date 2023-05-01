package iavlapp

import (
	"bytes"

	"os"

	"cosmossdk.io/log"
	"cosmossdk.io/store/iavl"
	"cosmossdk.io/store/metrics"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/abci/example/code"
	"github.com/cometbft/cometbft/abci/types"
	cmdb "github.com/cosmos/cosmos-db"
)

var _ types.Application = (*Application)(nil)

type Application struct {
	types.BaseApplication
	store   storetypes.CommitKVStore
	appHash []byte
}

func NewMerkleApp() Application {
	log := log.NewLogger(os.Stdout)
	db := cmdb.NewMemDB()
	key := storetypes.NewKVStoreKey("data")
	commitID := storetypes.CommitID{
		Version: 0,
		Hash:    []byte(""),
	}
	metrics := metrics.NewNoOpMetrics()
	store, err := iavl.LoadStore(db, log, key, commitID, false, 10, true, metrics)
	if err != nil {
		panic("Unable to create IAVL store")
	}
	appHash := make([]byte, 8)
	return Application{
		store:   store,
		appHash: appHash,
	}
}

func (a *Application) DeliverTx(req types.RequestDeliverTx) types.ResponseDeliverTx {

	var key, value string
	parts := bytes.Split(req.Tx, []byte("="))
	if len(parts) == 2 {
		key, value = string(parts[0]), string(parts[1])
	} else {
		key, value = string(req.Tx), string(req.Tx)
	}
	a.store.Set([]byte(key), []byte(value))

	return types.ResponseDeliverTx{
		Code: code.CodeTypeOK,
	}
}

func (a *Application) Commit() types.ResponseCommit {
	resp := a.store.Commit()
	a.appHash = resp.Hash
	return types.ResponseCommit{
		Data: resp.Hash,
	}
}

func (a *Application) Query(req types.RequestQuery) types.ResponseQuery {
	return a.store.(*iavl.Store).Query(req)
}
