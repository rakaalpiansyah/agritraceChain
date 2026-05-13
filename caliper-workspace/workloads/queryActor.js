'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

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
                contractArguments: [actorId, `Petani-Query-${i}`, 'Farmer', 'Jawa Barat'],
                readOnly: false
            };
            await this.sutAdapter.sendRequests(request);
        }
    }

    async submitTransaction() {
        this.txIndex++;
        const actorId = this.actorIds[this.txIndex % this.actorIds.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'GetActor',
            invokerMspId: 'FarmerMSP',
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
