import Ember from 'ember';

export default Ember.Controller.extend(***REMOVED***
  vus: Ember.computed.alias('model.metrics.vus.data.value'),
  vus_pooled: Ember.computed.alias('model.metrics.vus_pooled.data.value'),
  vus_max: Ember.computed('vus', 'vus_pooled', function() ***REMOVED***
    return this.get('vus') + this.get('vus_pooled');
  ***REMOVED***),
  vus_pct: Ember.computed('vus', 'vus_max', function() ***REMOVED***
    return (this.get('vus') / this.get('vus_max')) * 100;
  ***REMOVED***),
  metrics: Ember.computed('model.metrics', function() ***REMOVED***
    var metrics = this.get('model.metrics');
    var ret = [];
    for (var key in metrics) ***REMOVED***
      if (key !== 'vus' && key !== 'vus_pooled') ***REMOVED***
        ret.push(metrics[key]);
      ***REMOVED***
    ***REMOVED***
    return ret;
  ***REMOVED***),
  sortedMetrics: Ember.computed.sort('metrics', 'sortedMetricsBy'),
  sortedMetricsBy: ['name'],
***REMOVED***);
