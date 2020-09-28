module.exports = async () => {
    /*******************************************
     *** Set up
     ******************************************/
    const Web3 = require("web3");
    const HDWalletProvider = require("@truffle/hdwallet-provider");

    // Contract abstraction
    const truffleContract = require("truffle-contract");

    const oracleContract = truffleContract(
        require("../build/contracts/Oracle.json")
    );

    /*******************************************
     *** Constants
     ******************************************/
    // Config values
    const NETWORK_ROPSTEN =
        process.argv[4] === "--network" && process.argv[5] === "ropsten";

    let prophecyID = 2;
    let message = "0x3fe12969726048afa44e45d718c30a1245feadef820bc911f560690afcb9524f"
    let signature = "0xe1987eddf47013562d1f03ddf212ade198aa9fa6c16d2de5f7149a590b2000097355f42d122873deeb99ac541e201c8c69dec1ae9572c474f7c5869400a1aeb51b";

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

    console.log(`Attempting to send newOracleClaim() tx...`);

    try {
        var {
            logs
        } = await oracleContract.deployed().then(function (instance) {
            return instance.newOracleClaim(prophecyID, message, signature, {
                from: accounts[1],
                // value: 0,
                // gas: 300000 // 300,000 Gwei
            });
        });
    } catch (error) {
        console.log(error.message)
        return
    }

    // Get event logs
    const claim_event = logs.find(e => e.event === "LogNewOracleClaim");

    if (claim_event) {
        console.log(`\n\tProphecy ${claim_event.args._prophecyID} claimed`);
        console.log("-------------------------------------------");
        console.log(`Prophecy ID is:\t ${claim_event.args._prophecyID}`);
        console.log(`Message is:\t ${claim_event.args._message}`);
        console.log(`Validator is:\t ${claim_event.args._validatorAddress}`);
        console.log(`Signature is:\t ${claim_event.args._signature}`);
        console.log("-------------------------------------------");
    } else {
        console.error("Error: no result from transaction!");
    }

    const process_event = logs.find(e => e.event === "LogProphecyProcessed");

    if (process_event) {
        console.log(`\n\tProphecy ${process_event.args._prophecyID} processed`);
        console.log("-------------------------------------------");
        console.log(`Submitter:\t ${process_event.args._submitter}`);
        console.log(`Current prophecy power:\t ${process_event.args._prophecyPowerCurrent}`);
        console.log(`Prophecy power threshold:\t ${process_event.args._prophecyPowerThreshold}`);
        console.log("-------------------------------------------");
    } else {
        console.error("Error: no result from transaction!");
    }

    return;
};