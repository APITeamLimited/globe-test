import ws from "k6/ws";
import ***REMOVED*** check ***REMOVED*** from "k6";

export default function () ***REMOVED***
    var url = "ws://echo.websocket.org";
    var params = ***REMOVED*** "tags": ***REMOVED*** "my_tag": "hello" ***REMOVED*** ***REMOVED***;

    var response = ws.connect(url, params, function (socket) ***REMOVED***
        socket.on('open', function open() ***REMOVED***
            console.log('connected');
            socket.send(Date.now());

            socket.setInterval(function timeout() ***REMOVED***
                socket.ping();
                console.log("Pinging every 1sec (setInterval test)");
            ***REMOVED***, 1000);
        ***REMOVED***);

        socket.on('ping', function () ***REMOVED***
            console.log("PING!");
        ***REMOVED***);

        socket.on('pong', function () ***REMOVED***
            console.log("PONG!");
        ***REMOVED***);

        socket.on('pong', function () ***REMOVED***
            // Multiple event handlers on the same event
            console.log("OTHER PONG!");
        ***REMOVED***);

        socket.on('message', function incoming(data) ***REMOVED***
            console.log(`Roundtrip time: $***REMOVED***Date.now() - data***REMOVED*** ms`);
            socket.setTimeout(function timeout() ***REMOVED***
                socket.send(Date.now());
            ***REMOVED***, 500);
        ***REMOVED***);

        socket.on('close', function close() ***REMOVED***
            console.log('disconnected');
        ***REMOVED***);

        socket.on('error', function (e) ***REMOVED***
            if (e.error() != "websocket: close sent") ***REMOVED***
                console.log('An unexpected error occurred: ', e.error());
            ***REMOVED***
        ***REMOVED***);

        socket.setTimeout(function () ***REMOVED***
            console.log('2 seconds passed, closing the socket');
            socket.close();
        ***REMOVED***, 2000);
    ***REMOVED***);

    check(response, ***REMOVED*** "status is 101": (r) => r && r.status === 101 ***REMOVED***);
***REMOVED***;
