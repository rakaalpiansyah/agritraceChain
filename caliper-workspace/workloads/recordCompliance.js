'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RecordComplianceWorkload extends WorkloadModuleBase {
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
        const recordId = `COMPLIANCE_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const stages = ['COLLECTION', 'PROCESSING', 'PACKAGING', 'SHIPPING'];
        const statuses = ['COMPLIANT', 'COMPLIANT', 'COMPLIANT', 'NON_COMPLIANT']; // 75% compliant
        const stage = stages[this.txIndex % stages.length];
        const status = statuses[this.txIndex % statuses.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'RecordCompliance',
            invokerMspId: 'ProcessorMSP',
            contractArguments: [recordId, `BATCH_${this.txIndex}`, 'PROCESSOR_01', stage, status, `Check at stage ${stage}`],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new RecordComplianceWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
