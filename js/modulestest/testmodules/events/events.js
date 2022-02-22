import exec from 'k6/execution';
import ***REMOVED*** setTimeout, clearTimeout, setInterval, clearInterval ***REMOVED*** from 'k6/events'


export let options = ***REMOVED***
    scenarios: ***REMOVED***
        'foo': ***REMOVED***
            executor: 'constant-vus',
            vus: 1,
            duration: '3.8s',
            gracefulStop: '0s',
        ***REMOVED***
    ***REMOVED***
***REMOVED***;

function debug(arg) ***REMOVED***
    let t = String((new Date()) - exec.scenario.startTime).padStart(6, ' ')
    console.log(`[$***REMOVED***t***REMOVED***ms, iter=$***REMOVED***exec.scenario.iterationInTest***REMOVED***] $***REMOVED***arg***REMOVED***`);
***REMOVED***

export default function () ***REMOVED***
    debug('default start');

    let tickCount = 1;
    let f0 = (arg) => ***REMOVED***
        debug(`$***REMOVED***arg***REMOVED*** $***REMOVED***tickCount++***REMOVED***`);
    ***REMOVED***
    let t0 = setInterval(f0, 500, 'tick')

    let f1 = (arg) => ***REMOVED***
        debug(arg);
        clearInterval(t0);
    ***REMOVED***
    let t1 = setTimeout(f1, 2000, 'third');

    let t2 = setTimeout(debug, 1500, 'never happening');

    let f3 = (arg) => ***REMOVED***
        debug(arg);
        clearTimeout(t2);
        setTimeout(debug, 600, 'second');
    ***REMOVED***
    let t3 = setTimeout(f3, 1000, 'first');

    debug('default end');
    if (exec.scenario.iterationInTest == 1) ***REMOVED***
        debug(`expected last iter, the interval ID is $***REMOVED***t0***REMOVED***, we also expect timer $***REMOVED***t1***REMOVED*** to be interrupted`)
    ***REMOVED***
***REMOVED***
