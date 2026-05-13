'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

class SettleLCWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
        this.lcIds = [];
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;

        for (let i = 0; i < 100; i++) {
            const lcId = `SETTLE_LC_${workerIndex}_${roundIndex}_${i}_${Date.now()}`;

            await this.sutAdapter.sendRequests({
                contractId: this.roundArguments.contractId,
                contractFunction: 'IssueLC',
                invokerMspId: 'BuyerMSP',
                invokerIdentity: 'Admin',
                targetPeers: ['peer0.buyer.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.farmer.agritrace.com'],
                contractArguments: [lcId, `BUYER_${workerIndex}`, `FARMER_${i}`, `BATCH_${i}`, '5000000', 'IDR'],
                readOnly: false
            });

            await this.waitUntilLCExists(lcId);
            this.lcIds.push(lcId);
        }

        await sleep(5000);
    }

    async waitUntilLCExists(lcId) {
        const queryRequest = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GetLC',
            invokerMspId: 'BuyerMSP',
            invokerIdentity: 'Admin',
            contractArguments: [lcId],
            readOnly: true
        };

        for (let attempt = 1; attempt <= 20; attempt++) {
            try {
                await this.sutAdapter.sendRequests(queryRequest);
                return;
            } catch (err) {
                if (attempt === 20) {
                    throw new Error(`LC ${lcId} was issued but not visible after ${attempt} checks: ${err.message}`);
                }
                await sleep(500);
            }
        }
    }

    async submitTransaction() {
        this.txIndex++;
        const lcId = this.lcIds.shift();

        if (!lcId) {
            throw new Error('No pre-created LC is available to settle. Increase setup LC count or reduce round txNumber.');
        }

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'SettleLC',
            invokerMspId: 'BuyerMSP',
            invokerIdentity: 'Admin',
            targetPeers: ['peer0.buyer.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.farmer.agritrace.com'],
            contractArguments: [lcId],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new SettleLCWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
