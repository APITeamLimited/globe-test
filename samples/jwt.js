import crypto from "k6/crypto";
import encoding from "k6/encoding";
import ***REMOVED***sleep***REMOVED*** from "k6";

const algToHash = ***REMOVED***
    HS256: "sha256",
    HS384: "sha384",
    HS512: "sha512"
***REMOVED***;

function sign(data, hashAlg, secret) ***REMOVED***
    let hasher = crypto.createHMAC(hashAlg, secret);
    hasher.update(data);

    // Some manual base64 rawurl encoding as `Hasher.digest(encodingType)`
    // doesn't support that encoding type yet.
    return hasher.digest("base64").replace(/\//g, "_").replace(/\+/g, "-").replace(/=/g, "");
***REMOVED***

function encode(payload, secret, algorithm) ***REMOVED***
    algorithm = algorithm || "HS256";
    let header = encoding.b64encode(JSON.stringify(***REMOVED*** typ: "JWT", alg: algorithm ***REMOVED***), "rawurl");
    payload = encoding.b64encode(JSON.stringify(payload), "rawurl");
    let sig = sign(header + "." + payload, algToHash[algorithm], secret);
    return [header, payload, sig].join(".");
***REMOVED***

function decode(token, secret, algorithm) ***REMOVED***
    let parts = token.split('.');
    let header = JSON.parse(encoding.b64decode(parts[0], "rawurl"));
    let payload = JSON.parse(encoding.b64decode(parts[1], "rawurl"));
    algorithm = algorithm || algToHash[header.alg];
    if (sign(parts[0] + "." + parts[1], algorithm, secret) != parts[2]) ***REMOVED***
        throw Error("JWT signature verification failed");
    ***REMOVED***
    return payload;
***REMOVED***

export default function() ***REMOVED***
    let message = ***REMOVED*** key2: "value2" ***REMOVED***;
    let token = encode(message, "secret");
    console.log("encoded", token);
    let payload = decode(token, "secret");
    console.log("decoded", JSON.stringify(payload));
    sleep(1)
***REMOVED***
