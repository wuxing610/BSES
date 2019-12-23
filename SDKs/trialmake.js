'use strict';
var Fabric_Client = require('fabric-client');
var path = require('path');
var util = require('util');
var os = require('os');
var argument = process.argv.splice(2);
var conflict_name = argument[0];
var user = argument[1];
var service_provider = argument[2];
var name_1 = argument[3];
var faith_1 = argument[4];
var select_1 = argument[5];
var content_1 = argument[6];
var name_2 = argument[7];
var faith_2 = argument[8];
var select_2 = argument[9];
var content_2 = argument[10];
var name_3 = argument[11];
var faith_3 = argument[12];
var select_3 = argument[13];
var content_3 = argument[14];
var name_4 = argument[15];
var faith_4 = argument[16];
var select_4 = argument[17];
var content_4 = argument[18];
var name_5 = argument[19];
var faith_5 = argument[20];
var select_5 = argument[21];
var content_5 = argument[22];
var fabric_client = new Fabric_Client();
var channel = fabric_client.newChannel('conflictchannel');
var peer = fabric_client.newPeer('grpc://localhost:7051');
channel.addPeer(peer);
var order = fabric_client.newOrderer('grpc://localhost:7050')
channel.addOrderer(order);

//
var member_user = null;
var store_path = path.join(__dirname, 'user_key');
// console.log('Store path:'+store_path);
var tx_id = null;

Fabric_Client.newDefaultKeyValueStore({ path: store_path
}).then((state_store) => {
	// assign the store to the fabric client
	fabric_client.setStateStore(state_store);
	var crypto_suite = Fabric_Client.newCryptoSuite();
	var crypto_store = Fabric_Client.newCryptoKeyStore({path: store_path});
	crypto_suite.setCryptoKeyStore(crypto_store);
	fabric_client.setCryptoSuite(crypto_suite);

	// get the enrolled user from persistence, this user will sign all requests
	return fabric_client.getUserContext('user1', true);
}).then((user_from_store) => {
	if (user_from_store && user_from_store.isEnrolled()) {
		// console.log('Successfully loaded user1 from persistence');
		member_user = user_from_store;
	} else {
		throw new Error('Failed to get user1.... run registerUser.js');
	}

	tx_id = fabric_client.newTransactionID();
	// console.log("Assigning transaction_id: ", tx_id._transaction_id);
	var request = {
		//targets: let default to the peer assigned to the client
		chaincodeId: 'CCconflict',
		fcn: 'makeTrial',
		args: [conflict_name, user, service_provider, name_1, faith_1, select_1, content_1, name_2, faith_2, select_2, content_2, name_3, faith_3, select_3, content_3, name_4, faith_4, select_4, content_4, name_5, faith_5, select_5, content_5],
		chainId: 'conflictchannel',
		txId: tx_id
	};
	return channel.sendTransactionProposal(request);
}).then((results) => {
	var proposalResponses = results[0];
	var proposal = results[1];
	let isProposalGood = false;
	if (proposalResponses && proposalResponses[0].response &&
		proposalResponses[0].response.status === 200) {
			isProposalGood = true;
			// console.log('Transaction proposal was good');
		} else {
			console.error('Transaction proposal was bad');
		}
	if (isProposalGood) {
		console.log(util.format(
			'Successfully sent Proposal and received ProposalResponse: Status - %s, message - "%s"',
			proposalResponses[0].response.status, proposalResponses[0].response.message));
		var request = {
			proposalResponses: proposalResponses,
			proposal: proposal
		};

		var transaction_id_string = tx_id.getTransactionID(); //Get the transaction ID string to be used by the event processing
		var promises = [];

		var sendPromise = channel.sendTransaction(request);
		promises.push(sendPromise); //we want the send transaction first, so that we know where to check status

		let event_hub = channel.newChannelEventHub(peer);
		let txPromise = new Promise((resolve, reject) => {
			let handle = setTimeout(() => {
				event_hub.disconnect();
				resolve({event_status : 'TIMEOUT'}); //we could use reject(new Error('Trnasaction did not complete within 30 seconds'));
			}, 3000);
			event_hub.connect();
			event_hub.registerTxEvent(transaction_id_string, (tx, code) => {
				clearTimeout(handle);
				event_hub.unregisterTxEvent(transaction_id_string);
				event_hub.disconnect();
				var return_status = {event_status : code, tx_id : transaction_id_string};
				if (code !== 'VALID') {
					console.error('The transaction was invalid, code = ' + code);
					resolve(return_status); // we could use reject(new Error('Problem with the tranaction, event status ::'+code));
				} else {
					// console.log('The transaction has been committed on peer ' + event_hub.getPeerAddr);
					resolve(return_status);
				}
			}, (err) => {
				reject(new Error('There was a problem with the eventhub ::'+err));
			});
		});
		promises.push(txPromise);

		return Promise.all(promises);
	} else {
		console.error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
		throw new Error('Failed to send Proposal or receive valid response. Response null or status is not 200. exiting...');
	}
}).then((results) => {
	// console.log('Send transaction promise and event listener promise have completed');
	if (results && results[0] && results[0].status === 'SUCCESS') {
		// console.log('Successfully sent transaction to the orderer.');
	} else {
		console.error('Failed to order the transaction. Error code: ' + response.status);
	}

	if(results && results[1] && results[1].event_status === 'VALID') {
		// console.log('Successfully committed the change to the ledger by the peer');
	} else {
		console.log('Transaction failed to be committed to the ledger due to ::'+results[1].event_status);
	}
}).catch((err) => {
	console.error('Failed to invoke successfully :: ' + err);
});
