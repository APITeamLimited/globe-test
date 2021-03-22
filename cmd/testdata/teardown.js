import ***REMOVED*** abortTest ***REMOVED*** from 'k6';

export default function () ***REMOVED***
    abortTest();
***REMOVED***

export function teardown() ***REMOVED***
    console.log('Calling teardown function after abortTest()');
***REMOVED***