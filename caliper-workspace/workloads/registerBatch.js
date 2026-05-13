'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegisterBatchWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.workerIndex = workerIndex;
    }

    async submitTransaction() {
        this.txIndex++;
        const batchId = `BATCH_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const ownerId = `ACTOR_${this.workerIndex}_1_${Date.now()}`; // Reference to an existing actor
        const crops = ['Padi', 'Jagung', 'Kedelai', 'Kopi', 'Kelapa Sawit'];
        const cropType = crops[this.txIndex % crops.length];
        const quantity = (Math.floor(Math.random() * 1000) + 100);

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'RegisterBatch',
            invokerMspId: 'FarmerMSP',
            contractArguments: [batchId, ownerId, cropType, quantity.toString()],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new RegisterBatchWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
