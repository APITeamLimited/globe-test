import Ember from 'ember';

export default Ember.Route.extend(***REMOVED***
  model() ***REMOVED***
    return Ember.RSVP.hash(***REMOVED***
      "metrics": Ember.$.getJSON("/v1/metrics"),
    ***REMOVED***);
  ***REMOVED***,
***REMOVED***);
