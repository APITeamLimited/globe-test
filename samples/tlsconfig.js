import http from "k6/http";
import ***REMOVED*** check ***REMOVED*** from "k6";

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
    let res = http.get("https://sha256.badssl.com");
    check(res, ***REMOVED***
        "is TLSv1.2": (r) => r.tls_version === http.TLS_1_2,
        "is sha256 cipher suite": (r) => r.tls_cipher_suite === "TLS_RSA_WITH_AES_128_GCM_SHA256"
    ***REMOVED***);
***REMOVED***;
