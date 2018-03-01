package withdraw

import (
	"fmt"
	"reflect"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

func NewHandler(with WithdrawTxMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case WithdrawTx:
			return handleWithdrawTx(ctx, with, msg)
		default:
			errMsg := "Unrecognized withdraw Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// TODO

func buildMsg(from crypto.Address, to crypto.Address, pegCoin string, amount int64) (sdk.Msg, error) {
	strAmount = strconv.Itoa(amount)
	strCoin := fmt.Sprintf("%s%s", strAmount, pegCoin)
	coin, err := sdk.ParseCoin(strCoin)
	if err != nil {
		return nil, err
	}

	Eth := sdk.Coin{"Eth", amount}

	return msg, err
}

//
func handleWithdrawTx(ctx sdk.Context, with WithdrawTxMapper) sdk.Result {

	//

	for _, in := range msg.Inputs {
		_, err := ck.SubtractCoins(ctx, in.Address, in.Coins)
		if err != nil {
			return err.Result()
		}
	}

	for _, out := range msg.Outputs {
		_, err := ck.AddCoins(ctx, out.Address, out.Coins)
		if err != nil {
			return err.Result()
		}
	}

	return sdk.Result{} // TODO
}
