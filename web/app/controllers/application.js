import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  running: Ember.computed.alias('model.running'),
  actions: ***REMOVED***
    abort() ***REMOVED***
      var model = this.get('model');
      model.set('running', false);
      return model.save();
    ***REMOVED***,
  ***REMOVED***,
***REMOVED***);
