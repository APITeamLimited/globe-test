import Ember from 'ember';
import InflectorInitializer from 'k6/initializers/inflector';
import ***REMOVED*** module, test ***REMOVED*** from 'qunit';

let application;

module('Unit | Initializer | inflector', ***REMOVED***
  beforeEach() ***REMOVED***
    Ember.run(function() ***REMOVED***
      application = Ember.Application.create();
      application.deferReadiness();
    ***REMOVED***);
  ***REMOVED***
***REMOVED***);

// Replace this with your real tests.
test('it works', function(assert) ***REMOVED***
  InflectorInitializer.initialize(application);

  // you would normally confirm the results of the initializer here
  assert.ok(true);
***REMOVED***);
