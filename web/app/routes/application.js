import Ember from 'ember';

export default Ember.Route.extend(***REMOVED***
  _scheduleRefresh: Ember.on('init', function() ***REMOVED***
    Ember.run.later(()=> ***REMOVED***
      this.refresh();
      if (this.get('controller.model.running')) ***REMOVED***
        this._scheduleRefresh();
      ***REMOVED***
    ***REMOVED***, 5000);
  ***REMOVED***),
  model() ***REMOVED***
    return Ember.$.getJSON("/v1/status");
  ***REMOVED***,
  actions: ***REMOVED***
    abort() ***REMOVED***
      return Ember.$.post("/v1/abort").then(()=> ***REMOVED***
        this.refresh();
      ***REMOVED***);
    ***REMOVED***,
  ***REMOVED***,
***REMOVED***);
