
import ***REMOVED*** check ***REMOVED*** from 'k6';
import jsexttest from 'k6/x/jsexttest';

export let options = ***REMOVED***
    iterations: 5,
    thresholds: ***REMOVED***
        checks: ['rate===1'],
    ***REMOVED***
***REMOVED***;

export function handleSummary(data) ***REMOVED***
    return ***REMOVED***
        'summary-results.txt': data.metrics.foos.values.count.toString(),
    ***REMOVED***;
***REMOVED***


export default function () ***REMOVED***
    check(null, ***REMOVED***
        "foo is true": () => jsexttest.foo(__ITER),
    ***REMOVED***);
***REMOVED***