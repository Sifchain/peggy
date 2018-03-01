package withdraw

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
)

type WithdrawTxMapper struct {
	cdc *wire.Codec
	key sdk.StoreKey
}

func NewWithdrawTxMapper(key sdk.StoreKey) WitnessTxMapper {
	cdc := wire.NewCodec()
	cdc.RegisterConcrete(WithdrawData{}, "com.cosmos.peggy.WithdrawData", nil)
	cdc.RegisterConcrete(WithdrawTx{}, "com.cosmos.peggy.WithdrawTx", nil)

	return WithdrawTxMapper{
		cdc: cdc,
		key: key,
	}
}

type WithdrawObject struct {
	tx WithdrawTx
	WithdrawData
}

func (w WithdrawTxMapper) NewWithdrawObject(ctx sdk.Context, wdtx WithdrawTx) *WithdrawObject {
	return &WithdrawObject{
		tx: wdtx,
		// sets and mapping
	}
}

func (w WithdrawTxMapper) GetWithdrawData(ctx sdk.Context, wdtx WithdrawTx) *WithdrawObject {
	key := HashWithdrawTx(wdtx)
	kv := ctx.KVStore(w.key)
	bz := kv.Get(key)
	if bz == nil {
		return nil
	}
	var wd WithdrawData
	err := w.cdc.UnmarshalBinary(bz, &wd)
	if err != nil {
		panic(err)
	}
	return &WithdrawObject{
		tx:           wdtx,
		WithdrawData: wd,
	}
}

func (w WithdrawTxMapper) SetWithdrawData(ctx sdk.Context, wo *WithdrawObject) {
	kv := ctx.KVStore(w.key)
	bz, err := w.cdc.MarshalBinary(wo.WithdrawData)
	if err != nil {
		panic(err)
	}
	key := HashWithdrawTx(wo.tx)
	kv.Set(key, bz)
}

func HashWithdrawTx(wdtx WithdrawTx) []byte {
	WdTxHash := crypto.SHA256([]byte(fmt.Sprintf("%v", wdtx)))
	return WdTxHash
}
