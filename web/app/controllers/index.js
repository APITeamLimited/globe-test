import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  application: Ember.inject.controller(),

  vusActive: Ember.computed.alias('application.model.activeVUs'),
  vusInactive: Ember.computed.alias('application.model.inactiveVUs'),
  vusMax: Ember.computed('vusActive', 'vusInactive', function() ***REMOVED***
    return this.get('vusActive') + this.get('vusInactive');
  ***REMOVED***),
  vusPercent: Ember.computed('vusActive', 'maxVUs', function() ***REMOVED***
    return (this.get('vusActive') / this.get('maxVUs')) * 100;
  ***REMOVED***),
  vusLabel: Ember.computed('vusActive', 'maxVUs', function() ***REMOVED***
    return this.get('vusActive') + ' / ' + this.get('vusMax');
  ***REMOVED***),
***REMOVED***);
