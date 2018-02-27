package withdraw

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
)

// ------------------------------
// WithdrawTx

type WithdrawTx struct {
	destination []byte // Ethereum address of destinatary
	coin Coin	// contains: {denom: type of token to be send, amount: amount of token to send}
  signature sdk.StdSignature // a signature that authorises this transaction
}

func NewWithdrawTx(destination []byte, denom string, amount int64, signature sdk.StdSignature) WithdrawTx {
	coin := sdk.Coin{Denom: denom, Amount: amount}
	return WithdrawTx{
		destination: destination,
		coin: Coin,
		signature: signature
	 }
}

var _ sdk.Msg = (*WithdrawTx)(nil)

func (wtx WithdrawTx) ValidateBasic() sdk.Error {
	//TODO validate all fields from Withdraw Tx
	return nil
}

func (wtx WithdrawTx) Type() string {
	return "WithdrawTx"
}

func (wtx WithdrawTx) String() string {
	return fmt.Sprintf("WithdrawTx{%s, %v, %v}", wtx.destination, wtx.coin, wtx.signature)
}

func (wtx WithdrawTx) GetMsg() sdk.Msg {
	return wtx
}

func (wtx WithdrawTx) ValidateBasic() sdk.Error {

	// check if the fields are present
	// TODO errors
	// if len(wdata.SignedWithdraw) == 0 {
	// 	return sdk.NewError(code, msg)
	// }
	return nil
}

// ------------------------------
// WithdrawData

type WithdrawData struct {
	SignedWithdraw []SignTx	// Accumulates SignTxs until it reaches +2/3 of total power
	AccumulatedPower int // sum of each validator power
}

func (wdata WithdrawData) Type() string {
	return "WithdrawData"
}

func (wdata WithdrawData) GetAccumulatedPower() int64 {
	return wdata.AccumulatedPower
}

func (wdata WithdrawData) ValidateBasic() sdk.Error {
	totalPower = 0
	if len(wdata.SignedWithdraw) == 0 {
		return sdk.NewError(code, msg)
	}
	for _, stx := range wdata.SignedWithdraw {
		if err := stx.ValidateBasic(); err != nil {
			return err.Trace("")
		}
		// TODO: get power from validator using his signature
		// pubKey:= stx.signature.PubKey

		// ––––––– XXX functions not implemented ––––––––

		// if !isCurrentValidator(pubKey) {
		// 	errorMsg:= "signature {signature} doesn't come from a registered validator"
		// 	return sdk.newError("pubKey is not a validator")
		// }
		// validator:= getValidator(pubKey)

		// ––––––––––––––––––––––––––––––––––––––––––––––––

		// totalPower = totalPower + validator.Power
	}
	// if totalPower != wdata.AccumulatedPower {
	// 	"AcummulatedPower {number} doesn't match with the total signers power"
	// 	return error
	// }
	return nil
}
// TODO not implemented
// func (wdata WithdrawData) HasReachedSupermajority() (string, bool) {
// 	percentage:= wdata.AccumulatedPower*100 /chain.totalPower // round to decimals
// 	if 3*wdata.AccumulatedPower < 2*chain.totalPower {
// 		resBool:= false
// 	}  else {
// 		resBool:= true
// 	}
// 	resMsg:= ""
// 	return resMsg, resBool
// }



func NewWithdrawData() WithdrawData {
	return WithdrawData{
		SignedWithdraw: []SignTx,
		AccumulatedPower: 0
	}
}

// --------------------------------
// SignTx

type SignTx struct {
	signatureBytes []bytes   // signature bytes over the concatenation of the destination and Coin fields of WithdrawTx
	signature      sdk.StdSignature // signature that authorises this transaction is coming from a validator
	// XXX how are we assuring that the tx is coming from a validator ??
}

func (stx SignTx) Type() string {
	return "SignTx"
}

func (stx SignTx) GetMsg() sdk.Msg {
	return stx
}

func (stx SignTx) ValidateBasic() sdk.Error {

	// TODO errors
	// if len(wdata.SignedWithdraw) == 0 {
	// 	return sdk.NewError(code, msg)
	// }
	return nil
}
