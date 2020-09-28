module.exports = async () => {
    /*******************************************
     *** Set up
     ******************************************/
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    // Contract abstraction
    const truffleContract = require("truffle-contract");

    const oracleContract = truffleContract(
        require("../build/contracts/CosmosBridge.json")
    );

    /*******************************************
     *** Constants
     ******************************************/
    // Config values
    const NETWORK_ROPSTEN =
        process.argv[4] === "--network" && process.argv[5] === "ropsten";

    /*******************************************
     *** processBridgeProphecy transaction parameters
     ******************************************/
    let prophecyID;

    if (NETWORK_ROPSTEN) {
        prophecyID = Number(process.argv[6]);
    } else {
        prophecyID = Number(process.argv[4]);
    }

    /*******************************************
     *** Web3 provider
     *** Set contract provider based on --network flag
     ******************************************/
    let provider;
    if (NETWORK_ROPSTEN) {
        provider = new HDWalletProvider(
            process.env.MNEMONIC,
            "https://ropsten.infura.io/v3/".concat(process.env.INFURA_PROJECT_ID)
        );
    } else {
        provider = new Web3.providers.HttpProvider(process.env.LOCAL_PROVIDER);
    }

    const web3 = new Web3(provider);

    console.log("Fetching Oracle contract...");
    oracleContract.setProvider(web3.currentProvider);

    /*******************************************
     *** Contract interaction
     ******************************************/
    // Get current accounts
    const accounts = await web3.eth.getAccounts();

    cosmosSender = "0x0000000000000000000000000000000000000000";

    ethereumReceiver = accounts[1];
    // ethTokenAddress = "0x0000000000000000000000000000000000000000";
    // symbol = "ETH";
    nativeCosmosAssetDenom = "ATOM";
    // prefixedNativeCosmosAssetDenom = "PEGGYATOM";
    // amountWei = "1000000000000000000";
    amountNativeCosmos = "1000000000000000000";
    CLAIM_TYPE_LOCK = 2;

    console.log("Attempting to send processBridgeProphecy() tx...");
    try {
        var {
            logs
        } = await oracleContract.deployed().then(function (instance) {
            return instance.newProphecyClaim(
                CLAIM_TYPE_LOCK,
                cosmosSender,
                ethereumReceiver,
                nativeCosmosAssetDenom,
                amountNativeCosmos, {
                from: accounts[1],
            });
        });
    } catch (error) {
        console.log(error.message)
        return
    }

    // Get event logs
    const event = logs.find(e => e.event === "LogNewProphecyClaim");

    if (event) {
        console.log(`\n\tProphecy ${event.args._prophecyID} processed`);
        console.log("-------------------------------------------");
        console.log(`Submitter:\t ${event.args._claimType}`);
        console.log(`Current prophecy power:\t ${event.args._cosmosSender}`);
        console.log(`Prophecy power threshold:\t ${event.args._ethereumReceiver}`);
        console.log(`Token address:\t ${event.args._tokenAddress}`);
        console.log(`Validator address:\t ${event.args._validatorAddress}`);
        console.log(`Symbol:\t ${event.args._symbol}`);
        console.log(`Amount:\t ${event.args._amount}`);
        console.log("-------------------------------------------");
    } else {
        console.error("Error: no result from transaction!");
    }

    return;
};