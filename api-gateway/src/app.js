const express = require('express');
const cors = require('cors');
const winston = require('winston');
const { getContract } = require('./fabric-connection');

// Configure logging
const logger = winston.createLogger({
    level: 'info',
    format: winston.format.combine(
        winston.format.timestamp(),
        winston.format.json()
    ),
    transports: [
        new winston.transports.Console(),
        new winston.transports.File({ filename: 'error.log', level: 'error' }),
        new winston.transports.File({ filename: 'combined.log' })
    ]
});

const app = express();
app.use(cors());
app.use(express.json());

// Fabric Gateway connection configuration
const channelName = 'didchannel';
const chaincodeName = 'did-registry';
const mspId = 'Org1MSP';

// API Endpoints
app.post('/api/did/create', async (req, res) => {
    try {
        const { didID, pubKey } = req.body;
        
        if (!didID || !pubKey) {
            return res.status(400).json({ error: 'Missing required parameters' });
        }

        // Convert the simple pubKey format to JWK format
        const publicKeyJwk = {
            kty: "EC",
            crv: "P-256",
            x: pubKey.substring(0, pubKey.length/2),
            y: pubKey.substring(pubKey.length/2)
        };

        const contract = await getContract();
        await contract.submitTransaction('CreateDID', didID, didID, JSON.stringify(publicKeyJwk));
        
        logger.info(`Created DID: ${didID}`);
        res.status(201).json({ message: 'DID created successfully', didID });
    } catch (error) {
        logger.error('Error creating DID:', error);
        res.status(500).json({ error: error.message });
    }
});

app.get('/api/did/resolve/:id', async (req, res) => {
    try {
        const didID = req.params.id;
        const contract = await getContract();
        const result = await contract.evaluateTransaction('ResolveDID', didID);
        
        const didDoc = JSON.parse(result.toString());
        
        // Check if DID is revoked
        if (didDoc.status === 'revoked') {
            return res.status(401).json({ error: 'DID has been revoked' });
        }
        
        logger.info(`Resolved DID: ${didID}`);
        res.json(didDoc);
    } catch (error) {
        logger.error('Error resolving DID:', error);
        res.status(500).json({ error: error.message });
    }
});

app.post('/api/did/revoke', async (req, res) => {
    try {
        const { didID } = req.body;
        
        if (!didID) {
            return res.status(400).json({ error: 'Missing DID ID' });
        }

        const contract = await getContract();
        await contract.submitTransaction('RevokeDID', didID);
        
        logger.info(`Revoked DID: ${didID}`);
        res.json({ message: 'DID revoked successfully', didID });
    } catch (error) {
        logger.error('Error revoking DID:', error);
        res.status(500).json({ error: error.message });
    }
});

// Helper function to get contract instance is imported from fabric-connection.js

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
    logger.info(`Server running on port ${PORT}`);
});
