# Demo for peggy 

## what's peggy
Basically, peggy works as the asset transfer bridge between Etherum and Cosmos.

From developer, peggy include the depolyed contract in ethereum, module in cosmos,
and a relayer to listen on cross-chain asset transfer in both ethereum and cosmos.

## four user cases
From user interface point of view, just four types operations.
1. lock eth/erc20 in smart contract, then you will see your asset in cosmos
2. burn your peggied ethereum asset in cosmos, then your get locked asset back to your ethereum
3. lock cosmos asset and then you see new erc20 asset in ethereum
4. burn erc20 asset then get your cosmos asset back

## most important implementation in solidity
BridgeBank.sol
lock
unlock
mintBridgeTokens
burn

## most important method in cosmos
switch msg := msg.(type) {
		case MsgCreateEthBridgeClaim:  new eth asset come in
			return handleMsgCreateEthBridgeClaim(ctx, cdc, bridgeKeeper, msg)
            type ClaimType int
            const (
                LockText = ClaimType(iota) // case 1
                BurnText       // case 4
            )
		case MsgBurn: ethereum asset burn and then go back to ethereum // case 2
			return handleMsgBurn(ctx, cdc, accountKeeper, bridgeKeeper, msg)
		case MsgLock: cosmos asset locked and go to ethereum  // case 3
			return handleMsgLock(ctx, cdc, accountKeeper, bridgeKeeper, msg)

## common used account
bridge registry address: 0x30753E4A8aad7F8597332E813735Def5dD395028
bridge bank address: 0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4
operator address: 0x627306090abaB3A6e1400e9345bC60c78a8BEf57
most powerful validator address: 0xf17f52151EbEF6C7334FAD080c5704D77216b732 
bridge token addresss: 0x345cA3e014Aaf5dcA488057592ee47305D9B3e10

## Demo include three parts, (ethereum, cosmos, and end-to-end)

## demo part one, just ethereum
1. start the ganache and deploy smart contracts.
yarn develop
yarn migrate in other console 
yarn peggy:setup

### case 1: lock eth of operator to smart contract
1. lock some eth to cosmos account
check the bridge bank balance
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 eth
check the sender, aka operator balance
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth

yarn peggy:lock cosmos1wpqd7x8yyylmzvpwtlw74n0a5d3588fk9385c0  0x0000000000000000000000000000000000000000 1000000000000000000

check contract eth balance after transfer
   yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 eth
check contract eth balance after transfer
   yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth


### case 2: unlock eth in cosmos, then asset back to ethereum
check contract eth balance before claim sent
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 eth
check contract eth balance before claim send
yarn peggy:getTokenBalance  0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth

1. claim the an account burn 1 eth in cosmos chain, ethereum receiver is validator
truffle exec scripts/sendProphecyClaimBurn.js

check contract eth balance after claim sent
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 eth
check contract eth balance after claim send
yarn peggy:getTokenBalance  0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth

2. validator verify it, finalize after threshold reached.
truffle exec scripts/sendOracleClaimBurn.js

check contract eth balance after claim finalized
yarn peggy:getTokenBalance  0x2C2B9C9a4a25e24B174f26114e8926a9f2128FE4 eth
check contract eth balance after claim finalized
yarn peggy:getTokenBalance  0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth


### case 3: issue 1 atom in ethereum send to 0xf17f52151EbEF6C7334FAD080c5704D77216b732
1. claim of cosmos asset lock happend in ethereum
truffle  exec scripts/sendProphecyClaimLock.js

check the balance after claim sent
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA

2. validator verify claim and finalize it after threshold reached.
truffle exec scripts/sendOracleClaimLock.js

check the balance after claim verified.
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA

### case 4: burn cosmos asset in ethereum
truffle exec scripts/sendBurnTx.js 0x0000000000000000000000000000000000000000 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA 123456789

check the balance after peggyatom burn.
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA

## demo part two, just cosmos side
setup and start cosmos chain

ebd init local --chain-id=peggy
ebcli config keyring-backend test

ebcli config chain-id peggy
ebcli config trust-node true
ebcli config indent true
ebcli config output json

ebcli keys add validator
ebcli keys add testuser

ebd add-genesis-account $(ebcli keys show validator -a) 1000000000stake,1000000000atom
ebd gentx --name validator --keyring-backend test
ebd collect-gentxs

ebd start

check the running of cosmos chain.
ebcli tx send validator $(ebcli keys show testuser -a) 10atom --yes

ebcli query account $(ebcli keys show validator -a)
ebcli query account $(ebcli keys show testuser -a)

# Confirm your validator was created correctly, and has become Bonded
ebcli query staking validators

## case 1: eth locked in ethereum, and validator send claim to cosmos
ebcli query account $(ebcli keys show testuser -a)
ebcli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 0 eth 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 5 lock --token-contract-address=0x0000000000000000000000000000000000000000 --ethereum-chain-id=3 --from=validator --yes
ebcli q tx TXHASH
ebcli query account $(ebcli keys show testuser -a)


## case 2: eth burned in cosmos then back to ethereum
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1 peggyeth --ethereum-chain-id=3 --from=testuser --yes
ebcli q tx TXHASH

## case 3: lock cosmos asset, atom will locked in module's address, then transfer to ethereum
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 1 atom  --ethereum-chain-id=3 --from=testuser --yes
ebcli q tx TXHASH

## case 4: burn cosmos asset in ethereum, then token back to cosmos
ebcli tx ethbridge create-claim 0x30753E4A8aad7F8597332E813735Def5dD395028 1 atom 0x11111111262b236c9ac9a9a8c8e4276b5cf6b2c9 $(ebcli keys show testuser -a) $(ebcli keys show validator -a --bech val) 1 burn --ethereum-chain-id=3 --token-contract-address=0x345cA3e014Aaf5dcA488057592ee47305D9B3e10 --from=validator --yes
ebcli q tx TXHASH

## demo part 3, end to end
1. start truffle ethereum and deploy contract
yarn develop
yarn migrate
truffle exec scripts/setOracleAndBridgeBank.js

2. start cosmos chain
ebcli query account $(ebcli keys show validator -a)
ebcli tx send validator $(ebcli keys show testuser -a) 10atom --yes
ebcli query account $(ebcli keys show testuser -a)

3. start ebrelayer
ebrelayer generate
// the adderss should be from yarn peggy:address
ebrelayer init tcp://localhost:26657 ws://localhost:7545/ 0x30753E4A8aad7F8597332E813735Def5dD395028 validator --chain-id=peggy

### case 1: lock eth and send to cosmos test user from eth operator account
1. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
2. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x30753E4A8aad7F8597332E813735Def5dD395028  eth
3. check the testuser balance before lock
ebcli query account $(ebcli keys show testuser -a)

yarn peggy:lock cosmos14qx47vr5kh9xr47ht67r8rw446jgc408z3h54m  0x0000000000000000000000000000000000000000 1000000000000000000

4. check the balance of operator before lock
yarn peggy:getTokenBalance  0x627306090abaB3A6e1400e9345bC60c78a8BEf57 eth
5. check the ballance of contract before lock
yarn peggy:getTokenBalance  0x30753E4A8aad7F8597332E813735Def5dD395028  eth
6. check the testuser balance before lock
ebcli query account $(ebcli keys show testuser -a)

### case 2: burn testuser's eth in cosmos then asset to back to ethereum's validator account
1. check the validator's balance before burn
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
ebcli query account $(ebcli keys show testuser -a)

2. send burn tx in cosmos
ebcli tx ethbridge burn $(ebcli keys show testuser -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1000000000000000000 peggyeth --ethereum-chain-id=5777 --from=testuser --yes

3. check testuser's account 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732 eth
ebcli query account $(ebcli keys show testuser -a)

### case 3: lock atom in cosmos then issue the token in ethereum
ebcli tx ethbridge lock $(ebcli keys show testuser -a) 0xf17f52151EbEF6C7334FAD080c5704D77216b732 1 atom  --ethereum-chain-id=5777 --from=testuser --yes

1. check the balance of validator peggyatom in ethereum
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
ebcli query account $(ebcli keys show testuser -a)

### case 4: burn atom in ethereum and atom will be back to cosmos
truffle exec scripts/sendBurnTx.js cosmos14qx47vr5kh9xr47ht67r8rw446jgc408z3h54m 0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA 1
1. check balance after burn 
yarn peggy:getTokenBalance 0xf17f52151EbEF6C7334FAD080c5704D77216b732  0x409Ba3dd291bb5D48D5B4404F5EFa207441F6CbA
ebcli query account $(ebcli keys show testuser -a)


