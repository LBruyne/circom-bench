pragma circom 2.0.0;

include "../poseidon.circom";

template TestPoseidon() {
    var nInputs = 16;
    var nHashes = 600;
    signal input inputs[nInputs];
    signal output out[nHashes];

    component h[nHashes];
    for (var i = 0; i < nHashes; i++) {
        h[i] = Poseidon(nInputs);
    }
    
    for (var i = 0; i < nHashes; i++) {
        for(var j = 0; j < nInputs; j++) {
            h[i].inputs[j] <== inputs[j];
        }
    }

    for (var i = 0; i < nHashes; i++) {
        out[i] <== h[i].out;
    }
}

component main = TestPoseidon();