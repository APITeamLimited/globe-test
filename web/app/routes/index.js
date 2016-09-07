import Ember from 'ember';

export default Ember.Route.extend(***REMOVED***
  model() ***REMOVED***
    return Ember.RSVP.hash(***REMOVED***
      "metrics": this.get('store').findAll('metric'),
    ***REMOVED***);
  ***REMOVED***,
***REMOVED***);
