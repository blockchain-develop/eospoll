/* jshint esversion: 6 */ 
var ENV = 'dev';
var network;
var options;

if(ENV === 'dev'){
    // local testnet
    network = {
        blockchain: 'eos',
        host: '103.45.157.202',
        port: 8888,
        protocol: 'http',
        chainId: "f66c8e8800529cddebce531c1954bde660eb2d8824bebf519c3dccf2bf74b5ba",
        verbose: true,
        debug: true,
    };
    options = {
        broadcast: true,
        sign: true,
        chainId: "f66c8e8800529cddebce531c1954bde660eb2d8824bebf519c3dccf2bf74b5ba",
        httpEndpoint: "http://103.45.157.202:8888"
    };
} else if(ENV === 'testnet1'){
    // remote testnet
    network = {
        blockchain: 'eos',
        host: 'jungle.cryptolions.io',
        port: 18888,
        chainId: "038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca",
        protocol: "http"
    };
    options = {
        broadcast: true,
        sign: true,
        chainId: "038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca",
        httpEndpoint: "http://jungle.cryptolions.io"
    };
} else if(ENV === 'testnet'){
    // remote testnet
    network = {
        blockchain: 'eos',
        host: 'jungle.eosio.cr',
        port: 443,
        chainId: "038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca",
        protocol: "https"
    };
    options = {
        broadcast: true,
        sign: true,
        chainId: "038f4b0fc8ff18a4f0842a8f0564611f6e96e8535901dd45e43ac8691a1c4dca",
        httpEndpoint: "https://jungle.eosio.cr:443"
    };
} else if( ENV === 'mainnet') {
    // mainnet
    network = {
        blockchain: 'eos',
        host: 'nodes.get-scatter.com',
        port: 443,
        chainId: "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906",
        protocol: "https"
    };

    options = {
        broadcast: true,
        sign: true,
        chainId: "aca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906",
        httpEndpoint: "https://nodes.get-scatter.com:443"         
//      httpEndpoint: "http://mainnet.genereos.io:443"
    };
} else {
    throw("network config error");
}
