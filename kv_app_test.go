package iavlapp

import (
	"fmt"
	"os"
	"testing"

	"cosmossdk.io/log"
	"github.com/cometbft/cometbft/abci/types"
)

func TestKvApp(t *testing.T) {
	testLogger := log.NewLogger(os.Stdout)
	app := NewMerkleApp(testLogger)
	app.DeliverTx(types.RequestDeliverTx{
		Tx: []byte("cnode=cool"),
	})
	fmt.Println("apphash")
	fmt.Println(app.appHash)
	app.Commit()
	fmt.Println("apphash")
	fmt.Println(app.appHash)
	response := app.Query(types.RequestQuery{
		Path:  "/key",
		Data:  []byte("cnode"),
		Prove: true,
	})
	fmt.Println(response)
}
