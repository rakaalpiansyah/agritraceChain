'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class IssueCertificateWorkload extends WorkloadModuleBase {
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
        const certId = `CERT_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const certTypes = ['Organic', 'FairTrade', 'GAP', 'GradeA', 'Halal'];
        const certType = certTypes[this.txIndex % certTypes.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'IssueCertificate',
            invokerMspId: 'RegulatorMSP',
            invokerIdentity: 'Admin',
            targetPeers: ['peer0.regulator.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.farmer.agritrace.com'],
            contractArguments: [certId, 'REGULATOR_01', `BATCH_${this.txIndex}`, certType, '2027-12-31'],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new IssueCertificateWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
