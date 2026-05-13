'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class IssueLCWorkload extends WorkloadModuleBase {
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
        const lcId = `LC_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const amounts = [1000000, 2500000, 5000000, 7500000, 10000000];
        const amount = amounts[this.txIndex % amounts.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'IssueLC',
            invokerMspId: 'BuyerMSP',
            contractArguments: [lcId, `BUYER_${this.workerIndex}`, `FARMER_${this.txIndex}`, `BATCH_${this.txIndex}`, amount.toString(), 'IDR'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new IssueLCWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
