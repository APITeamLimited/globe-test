import DS from 'ember-data';

export default DS.Model.extend(***REMOVED***
  name: DS.attr('string'),
  group: DS.belongsTo('group'),
  passes: DS.attr('number'),
  fails: DS.attr('number')
***REMOVED***);
