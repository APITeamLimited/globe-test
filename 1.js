import grpc from 'k6/net/grpc';

const client = new grpc.Client();
// Download addsvc.proto for https://grpcb.in/, located at:
// https://raw.githubusercontent.com/moul/pb/master/addsvc/addsvc.proto
// and put it in the same folder as this script.
client.load(null, 'any_test.proto');

export default function () ***REMOVED***
	client.connect('localhost:50051', ***REMOVED*** timeout: '5s',plaintext:true ***REMOVED***);

	const response = client.invoke('grpc.any.testing.AnyTestService/Sum', ***REMOVED***
		data: ***REMOVED***
			"@type": "type.googleapis.com/grpc.any.testing.SumRequestData",
			"a": 1,
			"b": 2,
		***REMOVED***,
	***REMOVED***);
	console.log(JSON.stringify(response))
	console.log(response.message.data.v); // should print 3

	client.close();
***REMOVED***;
