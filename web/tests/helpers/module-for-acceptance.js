import ***REMOVED*** module ***REMOVED*** from 'qunit';
import Ember from 'ember';
import startApp from '../helpers/start-app';
import destroyApp from '../helpers/destroy-app';

const ***REMOVED*** RSVP: ***REMOVED*** Promise ***REMOVED*** ***REMOVED*** = Ember;

export default function(name, options = ***REMOVED******REMOVED***) ***REMOVED***
  module(name, ***REMOVED***
    beforeEach() ***REMOVED***
      this.application = startApp();

      if (options.beforeEach) ***REMOVED***
        return options.beforeEach.apply(this, arguments);
      ***REMOVED***
    ***REMOVED***,

    afterEach() ***REMOVED***
      let afterEach = options.afterEach && options.afterEach.apply(this, arguments);
      return Promise.resolve(afterEach).then(() => destroyApp(this.application));
    ***REMOVED***
  ***REMOVED***);
***REMOVED***
