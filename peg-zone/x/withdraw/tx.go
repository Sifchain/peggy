package withdraw

import (
	"fmt"
	"regexp"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ------------------------------
// WithdrawTx

type WithdrawTx struct {
	destination []byte           // Ethereum address of destinatary
	coin        Coin             // contains: {denom: type of token to be send, amount: amount of token to send}
	signature   sdk.StdSignature // a signature that authorises this transaction
}

func NewWithdrawTx(destination []byte, denom string, amount int64, signature sdk.StdSignature) WithdrawTx {
	coin := sdk.Coin{Denom: denom, Amount: amount}
	return WithdrawTx{
		destination: destination,
		coin:        Coin,
		signature:   signature,
	}
}

var _ sdk.Msg = (*WithdrawTx)(nil)

func (wtx WithdrawTx) Type() string {
	return "WithdrawTx"
}

func (wtx WithdrawTx) String() string {
	return fmt.Sprintf("WithdrawTx{\n\t%v,\n\t %v,\n\t %v\n}", wtx.destination, wtx.coin, wtx.signature)
}

func (wtx WithdrawTx) GetMsg() sdk.Msg {
	return wtx
}

var (
	reDnm  = `[[:alpha:]][[:alnum:]]{2,15}`
	reCoin = regexp.MustCompile(fmt.Sprintf(`^(%s)$`, reDnm))
)

func (wtx WithdrawTx) ValidateBasic() sdk.Error {

	if len(wtx.destination) == 0 {
		return ErrEmptyByteArray(wtx.destination)
	}

	matches := reCoin.FindStringSubmatch(wtx.coin.Denom)
	if matches == nil {
		return ErrInvalidCoin(wtx.coin.Denom)
	}
	if wtx.coin.IsZero() {
		return ErrZeroValueWithdrawTx()
	}
	if wtx.signature.IsZero() {

	}
	return nil
}

// ------------------------------
// WithdrawData

type WithdrawData struct {
	SignedWithdraw   []SignTx // Accumulates SignTxs until it reaches +2/3 of total power
	AccumulatedPower int64    // sum of each validator power
}

func IsDoubleSigning(wdata WithdrawData, newStx SignTx) bool {
	for _, stx := range wdata.SignedWithdraw {
		if stx.signature.PubKey == newStx.signature.PubKey {
			return true
		}
	}
	return false
}

func (wdata *WithdrawData) AddSignTx(stx SignTx) sdk.Error {
	err := stx.ValidateBasic()
	if err != nil {
		return err
	}
	if len(wdata.SignedWithdraw) != 0 && wdata.SignedWithdraw[len(wdata.SignedWithdraw)-1].signatureBytes != stx.signatureBytes {
		return ErrSignedDataMismatch()
	}
	if IsDoubleSigning(wdata, stx) {
		return ErrDoubleSign()
	}
	wdata.SignedWithdraw = append(wdata.SignedWithdraw, stx)
	// TODO add validator power AccumulatedPower
	// validator := stx.signature.PublicKey //get crypto.Address from signature PublicKey
	// wdata.AccumulatedPower = wdata.AccumulatedPower.Plus(validator.power) // make a secure int64 operation wo under/overflow
	return nil
}

func (wdata WithdrawData) Type() string {
	return "WithdrawData"
}

func (wdata WithdrawData) GetAccumulatedPower() int64 {
	return wdata.AccumulatedPower
}

func (wdata WithdrawData) String() string {
	return fmt.Sprintf("WithdrawData {\n\t%v,\n\t %d\n}", wdata.SignedWithdraw, wdata.AccumulatedPower)
}

func (wdata WithdrawData) ValidateBasic() sdk.Error {
	totalPower = 0
	if len(wdata.SignedWithdraw) == 0 {
		return ErrEmptySignedWithdraw()
	}
	for _, stx := range wdata.SignedWithdraw {
		if err := stx.ValidateBasic(); err != nil {
			return err.Trace("")
		}
		// TODO: get power from validator using his signature

		// validator, _ := getValidator(stx.signature.PubKey)

		// totalPower = totalPower + validator.Power
	}
	// if totalPower != wdata.AccumulatedPower {
	// 	"AcummulatedPower {number} doesn't match with the total signers power"
	// 	return error
	// }
	return nil
}

// XXX not implemented
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
		SignedWithdraw:   make([]SignTx),
		AccumulatedPower: 0,
	}
}

// --------------------------------
// SignTx

type SignTx struct {
	signatureBytes []bytes          // signature bytes over the concatenation of the destination and Coin fields of WithdrawTx
	signature      sdk.StdSignature // signature that authorises this transaction is coming from a validator
	// XXX how are we assuring that the tx is coming from a validator ??
}

func (stx SignTx) Type() string {
	return "SignTx"
}

func (stx SignTx) GetMsg() sdk.Msg {
	return stx
}

func (stx SignTx) String() sdk.Msg {
	return fmt.Sprintf("SignTx {\n\t%v,\n\t %v\n}", stx.signatureBytes, stx.signature)
}

func (stx SignTx) ValidateBasic() sdk.Error {

	zeroSignature := stx.signature.IsZero() // what does it mean to be a zero signature ?
	if zeroSignature {
		// ErrZeroSignature
	}
	if len(signatureBytes) == 0 {
		return ErrEmptyByteSlice("signatureBytes")
	}
	// TODO valiate that signature comes from a current validator

	// validator, err := getValidator(pubKey)
	// if error != nil {
	// 	errorMsg:= "signature {signature} doesn't come from a registered validator"
	// 	return sdk.newError("pubKey is not a validator")
	// }
	return nil
}

func NewSignTx(signatureBytes []bytes, signature sdk.StdSignature) SignTx {
	return SignTx{
		signatureBytes,
		signature,
	}
}
