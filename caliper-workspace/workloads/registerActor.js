'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegisterActorWorkload extends WorkloadModuleBase {
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
        const actorId = `ACTOR_${this.workerIndex}_${this.txIndex}_${Date.now()}`;
        const roles = ['Farmer', 'Aggregator', 'Processor', 'Buyer'];
        const locations = ['Jawa Barat', 'Jawa Tengah', 'Jawa Timur', 'Sulawesi Selatan', 'Sumatera Utara'];
        const role = roles[this.txIndex % roles.length];
        const location = locations[this.txIndex % locations.length];

        const request = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'RegisterActor',
            invokerMspId: 'FarmerMSP',
            invokerIdentity: 'Admin',
            targetPeers: ['peer0.farmer.agritrace.com', 'peer0.aggregator.agritrace.com', 'peer0.processor.agritrace.com'],
            contractArguments: [actorId, `Petani-${actorId}`, role, location],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(request);
    }
}

function createWorkloadModule() {
    return new RegisterActorWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
