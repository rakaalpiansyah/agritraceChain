// caliper-workspace/workloads/workload.js
'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

class RegisterProductWorkload extends WorkloadModuleBase {
    constructor() {
        super();
        this.txIndex = 0;
    }

    async submitTransaction() {
        this.txIndex++;
        const productID = `PROD-${this.workerIndex}-${this.txIndex}`;
        const farmerID = `FARMER-ID-${Math.floor(Math.random() * 100)}`;
        
        const args = {
            contractId: 'sc01',
            contractFunction: 'RegisterProduct',
            contractArguments: [productID, farmerID, "-6.9147, 107.6098", "Arabica", "Organic", "A"],
            readOnly: false
        };

        await this.sutAdapter.sendRequests(args);
    }
}

function createWorkloadModule() {
    return new RegisterProductWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;