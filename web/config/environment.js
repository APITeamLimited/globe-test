/* jshint node: true */

module.exports = function(environment) ***REMOVED***
  var ENV = ***REMOVED***
    modulePrefix: 'speedboat',
    environment: environment,
    rootURL: '/',
    locationType: 'auto',
    EmberENV: ***REMOVED***
      FEATURES: ***REMOVED***
        // Here you can enable experimental features on an ember canary build
        // e.g. 'with-controller': true
      ***REMOVED***
    ***REMOVED***,

    APP: ***REMOVED***
      // Here you can pass flags/options to your application instance
      // when it is created
    ***REMOVED***
  ***REMOVED***;

  if (environment === 'development') ***REMOVED***
    // ENV.APP.LOG_RESOLVER = true;
    // ENV.APP.LOG_ACTIVE_GENERATION = true;
    // ENV.APP.LOG_TRANSITIONS = true;
    // ENV.APP.LOG_TRANSITIONS_INTERNAL = true;
    // ENV.APP.LOG_VIEW_LOOKUPS = true;
  ***REMOVED***

  if (environment === 'test') ***REMOVED***
    // Testem prefers this...
    ENV.locationType = 'none';

    // keep test console output quieter
    ENV.APP.LOG_ACTIVE_GENERATION = false;
    ENV.APP.LOG_VIEW_LOOKUPS = false;

    ENV.APP.rootElement = '#ember-testing';
  ***REMOVED***

  if (environment === 'production') ***REMOVED***

  ***REMOVED***

  return ENV;
***REMOVED***;
