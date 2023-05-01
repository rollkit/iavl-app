package iavlapp

import (
	"os"
	"testing"

	"cosmossdk.io/log"
	storetypes "cosmossdk.io/store/types"
	"github.com/cometbft/cometbft/abci/types"
	"github.com/cometbft/cometbft/crypto/merkle"
	"github.com/stretchr/testify/require"
)

func TestKvApp(t *testing.T) {
	require := require.New(t)
	testLogger := log.NewLogger(os.Stdout)
	app := NewMerkleApp(testLogger)
	app.DeliverTx(types.RequestDeliverTx{
		Tx: []byte("cnode=cool"),
	})
	app.Commit()
	response := app.Query(types.RequestQuery{
		Path:   "/key",
		Data:   []byte("cnode"),
		Height: 1,
		Prove:  true,
	})
	decoder := merkle.NewProofRuntime()
	decoder.RegisterOpDecoder("ics23:iavl", storetypes.CommitmentOpDecoder)
	err := decoder.VerifyValue(response.ProofOps, app.appHash, "/cnode", []byte("cool"))
	require.NoError(err)
	err = decoder.VerifyValue(response.ProofOps, app.appHash, "/cnode2", []byte("cool"))
	require.Error(err)
	err = decoder.VerifyValue(response.ProofOps, app.appHash, "/cnode", []byte("cringe"))
	require.Error(err)
	app.DeliverTx(types.RequestDeliverTx{
		Tx: []byte("cnode=based"),
	})
	app.Commit()
	response = app.Query(types.RequestQuery{
		Path:   "/key",
		Data:   []byte("cnode"),
		Prove:  true,
		Height: 2,
	})
	err = decoder.VerifyValue(response.ProofOps, app.appHash, "/cnode", []byte("based"))
	require.NoError(err)
	err = decoder.VerifyValue(response.ProofOps, app.appHash, "/cnode", []byte("cringee"))
	require.Error(err)
}
