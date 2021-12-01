import grpc from 'k6/net/grpc';
import ***REMOVED***check***REMOVED*** from "k6";

let client = new grpc.Client();

export default () => ***REMOVED***
	client.connect("127.0.0.1:10000", ***REMOVED***plaintext: true, reflect: true***REMOVED***)
	const response = client.invoke("main.FeatureExplorer/GetFeature", ***REMOVED***
		latitude: 410248224,
		longitude: -747127767
	***REMOVED***)

	check(response, ***REMOVED***"status is OK": (r) => r && r.status === grpc.StatusOK***REMOVED***);
	console.log(JSON.stringify(response.message))

	client.close()
***REMOVED***

