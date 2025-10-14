const { connect, Identity, Signer, signers } = require('@hyperledger/fabric-gateway');
const grpc = require('@grpc/grpc-js');
const { promises: fs } = require('fs');
const path = require('path');
const crypto = require('crypto');

const channelName = 'didchannel';
const chaincodeName = 'didregistry';
const mspId = 'Org1MSP';

const cryptoPath = path.resolve(__dirname, '..', '..', 'fabric-samples', 'test-network',
    'organizations', 'peerOrganizations', 'org1.example.com');


const keyDirectoryPath = path.resolve(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'keystore');
const certPath = path.resolve(cryptoPath, 'users', 'User1@org1.example.com', 'msp', 'signcerts', 'cert.pem');
const tlsCertPath = path.resolve(cryptoPath, 'peers', 'peer0.org1.example.com', 'tls', 'ca.crt');

const peerEndpoint = 'localhost:7051';

async function newGrpcConnection() {
    const tlsRootCert = await fs.readFile(tlsCertPath);
    const tlsCredentials = grpc.credentials.createSsl(tlsRootCert);
    return new grpc.Client(peerEndpoint, tlsCredentials, {
        'grpc.ssl_target_name_override': 'peer0.org1.example.com',
    });
}

async function newIdentity() {
    const credentials = await fs.readFile(certPath);
    return { mspId, credentials };
}

async function newSigner() {
    const files = await fs.readdir(keyDirectoryPath);
    const keyPath = path.resolve(keyDirectoryPath, files[0]);
    const privateKeyPem = await fs.readFile(keyPath);
    const privateKey = crypto.createPrivateKey(privateKeyPem);
    return signers.newPrivateKeySigner(privateKey);
}

async function getContract() {
    const client = await newGrpcConnection();
    const gateway = connect({
        client,
        identity: await newIdentity(),
        signer: await newSigner(),
        evaluateOptions: () => {
            return { deadline: Date.now() + 5000 }; 
        },
        endorseOptions: () => {
            return { deadline: Date.now() + 15000 };
        },
        submitOptions: () => {
            return { deadline: Date.now() + 5000 }; 
        },
        commitStatusOptions: () => {
            return { deadline: Date.now() + 60000 }; 
        },
    });

    try {
        const network = gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);
        return contract;
    } catch (error) {
        console.error('Failed to connect to the gateway:', error);
        throw error;
    }
}

module.exports = {
    getContract
};
