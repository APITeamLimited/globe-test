import Ember from 'ember';

export default Ember.Route.extend(***REMOVED***
  model() ***REMOVED***
    return Ember.RSVP.hash(***REMOVED***
      "metrics": this.get('store').findAll('metric'),
      "groups": this.get('store').findAll('group'),
    ***REMOVED***);
  ***REMOVED***,
  afterModel(model) ***REMOVED***
    model["defaultGroup"] = this.get('store').peekRecord('group', 0);
  ***REMOVED***,
***REMOVED***);
