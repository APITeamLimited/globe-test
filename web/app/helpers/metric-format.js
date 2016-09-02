import Ember from 'ember';
import moment from 'moment';

export function metricFormat(params/*, hash*/) ***REMOVED***
  let [value, type] = params;
  if (type == "time") ***REMOVED***
    return moment.duration(value / 1000000, 'milliseconds').format('h[h]m[m]s[s]S[ms]');
  ***REMOVED***
  return value;
***REMOVED***

export default Ember.Helper.helper(metricFormat);
