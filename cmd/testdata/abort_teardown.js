import exec from 'k6/execution';

export default function () ***REMOVED***
    exec.test.abort();
***REMOVED***

export function teardown() ***REMOVED***
    console.log('Calling teardown function after test.abort()');
***REMOVED***
