import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  application: Ember.inject.controller(),

  vus: Ember.computed.alias('application.model.vus'),
  vusMax: Ember.computed.alias('application.model.vusMax'),
  vusPercent: Ember.computed('vus', 'vusMax', function() ***REMOVED***
    return (this.get('vus') / this.get('vusMax')) * 100;
  ***REMOVED***),
  vusLabel: Ember.computed('vus', 'vusMax', function() ***REMOVED***
    return this.get('vus') + ' / ' + this.get('vusMax');
  ***REMOVED***),
***REMOVED***);
