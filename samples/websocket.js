import websocket from "k6/websocket";

export default function () ***REMOVED***
    var result = websocket.connect("wss://echo.websocket.org", function(socket) ***REMOVED***
        socket.on('open', function open() ***REMOVED***
            console.log('connected');
            socket.send(Date.now());

            socket.setInterval(function timeout() ***REMOVED***
                socket.ping();
                console.log("Pinging every 1sec (setInterval test)");
            ***REMOVED***, 1000);
        ***REMOVED***);

        socket.on('pong', function () ***REMOVED***
            console.log("PONG!");
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
            console.log('An error occured: ', e.error());
        ***REMOVED***);

        socket.setTimeout(function() ***REMOVED***
            console.log('5 seconds passed, closing the socket');
            socket.close();
        ***REMOVED***, 5000);
    ***REMOVED***);
***REMOVED***;
