import grpc from "k6/net/grpc";

var client = new grpc.Client();
// client.load([], "../../../../../grpc-go/examples/helloworld/helloworld/helloworld.proto");

export default function () ***REMOVED***
        if (!client) ***REMOVED***
                throw new Error("no client created");
        ***REMOVED***
        client.connect('localhost:50051', ***REMOVED*** plaintext: true, timeout: '3s', reflect: true ***REMOVED***);
        var resp = client.invoke('/helloworld.Greeter/SayHello', ***REMOVED******REMOVED***)
        if (!resp.message || resp.error ) ***REMOVED***
                throw new Error('unexpected response message: ' + JSON.stringify(resp.message))
        ***REMOVED***
        console.log(JSON.stringify(resp), resp.error)
***REMOVED***