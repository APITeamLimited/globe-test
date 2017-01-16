import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  flashMessages: Ember.inject.service(),

  paused: Ember.computed.alias('model.paused'),

  actions: ***REMOVED***
    pause() ***REMOVED***
      var model = this.get('model');
      model.set('paused', true);
      return model.save().catch((err) => ***REMOVED***
        for (var e of err.errors) ***REMOVED***
          this.get('flashMessages').danger(e.title);
        ***REMOVED***
      ***REMOVED***);
    ***REMOVED***,
    resume() ***REMOVED***
      var model = this.get('model');
      model.set('paused', false);
      return model.save().catch((err) => ***REMOVED***
        for (var e of err.errors) ***REMOVED***
          this.get('flashMessages').danger(e.title);
        ***REMOVED***
      ***REMOVED***);
    ***REMOVED***,
  ***REMOVED***,
***REMOVED***);
