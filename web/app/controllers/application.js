import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  flashMessages: Ember.inject.service(),

  running: Ember.computed.alias('model.running'),

  actions: ***REMOVED***
    abort() ***REMOVED***
      var model = this.get('model');
      model.set('running', false);
      return model.save().then(() => ***REMOVED***
        this.get('flashMessages').success("Test stopped");
      ***REMOVED***, (err) => ***REMOVED***
        for (var e of err.errors) ***REMOVED***
          this.get('flashMessages').danger(e.title);
        ***REMOVED***
      ***REMOVED***);
    ***REMOVED***,
  ***REMOVED***,
***REMOVED***);
