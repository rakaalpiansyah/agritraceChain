'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const sleep = ms => new Promise(resolve => setTimeout(resolve, ms));

class QueryActorWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
        this.actorIds = [];
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;

        // Pre-create some actors to query later
        for (let i = 0; i < 10; i++) {
            const actorId = `QUERY_ACTOR_${workerIndex}_${i}_${Date.now()}`;
            this.actorIds.push(actorId);

            const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'RegisterActor',
                invokerMspId: 'FarmerMSP',
                invokerIdentity: 'Admin',
                targetPeers: ['peer0.farmer.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.processor.agritrace.com'],
                contractArguments: [actorId, `Petani-Query-${i}`, 'Farmer', 'Jawa Barat'],
                readOnly: false
            };
            await this.sutAdapter.sendRequests(request);
            await this.waitUntilActorExists(actorId);
        }

        await sleep(5000);
    }

    async waitUntilActorExists(actorId) {
        const queryRequest = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GetActor',
            invokerMspId: 'FarmerMSP',
            invokerIdentity: 'Admin',
            targetPeers: ['peer0.farmer.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.processor.agritrace.com'],
            contractArguments: [actorId],
            readOnly: false
        };

        for (let attempt = 1; attempt <= 20; attempt++) {
            try {
                await this.sutAdapter.sendRequests(queryRequest);
                return;
            } catch (err) {
                if (attempt === 20) {
                    throw new Error(`Query actor ${actorId} was not visible after ${attempt} checks: ${err.message}`);
                }
                await sleep(500);
            }
        }
    }

    async submitTransaction() {
        this.txIndex++;
        const actorId = this.actorIds[this.txIndex % this.actorIds.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GetActor',
            invokerMspId: 'FarmerMSP',
            invokerIdentity: 'Admin',
            contractArguments: [actorId],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new QueryActorWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
