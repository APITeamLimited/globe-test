import http from 'k6/http';

export let options = ***REMOVED***
    tlsCipherSuites: [
        "TLS_RSA_WITH_RC4_128_SHA",
        "TLS_RSA_WITH_AES_128_GCM_SHA256",
    ],
    tlsVersion: ***REMOVED***
        min: "ssl3.0",
        max: "tls1.2"
    ***REMOVED***
***REMOVED***;

export default function() ***REMOVED***
  const response = http.get("https://sha256.badssl.com");
***REMOVED***;

