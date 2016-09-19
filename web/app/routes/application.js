import Ember from 'ember';

export default Ember.Route.extend(***REMOVED***
  _scheduleRefresh: Ember.on('init', function(delay = 5000) ***REMOVED***
    Ember.run.later(()=> ***REMOVED***
      this.refresh();
      this._scheduleRefresh(this.get('controller.model.running') ? 5000 : 15000);
    ***REMOVED***, delay);
  ***REMOVED***),
  model() ***REMOVED***
    return this.get('store').findRecord('status', 'default');
  ***REMOVED***
***REMOVED***);
