package withdraw

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	CodeSignatureMismatch      CodeType = 201
	CodeInvalidEthereumAddress CodeType = 202
	CodeEmptySignedWithdraw    CodeType = 203
	CodeZeroValueWithdrawTx    CodeType = 204
	CodeEmptyByteSlice         CodeType = 205
	CodeInvalidCoin            CodeType = 206
	CodeSignedMsgMismatch      CodeType = 207
	CodeUnknownRequest         CodeType = sdk.CodeUnknownRequest
)

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeSignatureMismatch:
		return "Signature doesn't match with a registered validator"
	case CodeInvalidEthereumAddress:
		return "Invalid Ethereum address"
	case CodeEmptySignedWithdraw:
		return "No SignTx stored in slice"
	case CodeZeroValueWithdrawTx:
		return "Amount of WithdrawTx is zero"
	case CodeEmptyByteSlice:
		return "Empty byte slice"
	case CodeInvalidCoin:
		return "Invalid coin denomination"
	case CodeSignedMsgMismatch:
		return "Invalid coin denomination"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

//----------------------------------------
// Error constructors

func ErrSignatureMismatch(signature sdk.StdSignature) sdk.Error {
	return newError(CodeSignatureMismatch, fmt.Sprintf("Signature %v doesn't match with a registered validator", signature)) // TODO
}

func ErrInvalidEthereumAddress(address []byte) sdk.Error {
	return newError(CodeInvalidEthereumAddress, fmt.Sprintf("Invalid Ethereum address: %v", address))
}

func ErrEmptySignedWithdraw() sdk.Error {
	return newError(CodeEmptySignedWithdraw, "")
}

func ErrZeroValueWithdrawTx() sdk.Error {
	return newError(CodeZeroValueWithdrawTx, "")
}

func ErrEmptyByteSlice(byteSlice string) sdk.Error {
	return newError(CodeEmptyByteArray, fmt.Sprintf("Empty byte slice: %s", byteSlice))
}

func ErrInvalidCoin(denom string) sdk.Error {
	return newError(CodeInvalidCoin, fmt.Sprintf("Invalid coin denom: %s", denom))
}

func ErrSignedMsgMismatch() sdk.Error {
	return newError(CodeSignedMsgMismatch, "")
}

func ErrUnknownRequest(msg string) sdk.Error {
	return newError(CodeUnknownRequest, msg)
}

//----------------------------------------

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)

}

func newError(code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(code, msg)
}
