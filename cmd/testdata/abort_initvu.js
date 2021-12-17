import exec from 'k6/execution';

// This won't fail on initial parsing of the script, but on VU initialization.
if (__VU == 1) ***REMOVED***
  exec.test.abort();
***REMOVED***

export default function() ***REMOVED******REMOVED***
