import Ember from 'ember';
import DS from 'ember-data';

export default DS.Model.extend(***REMOVED***
  name: DS.attr('string'),
  parent: DS.belongsTo('group', ***REMOVED*** inverse: 'groups' ***REMOVED***),
  groups: DS.hasMany('group', ***REMOVED*** inverse: 'parent' ***REMOVED***),
  tests: DS.hasMany('test'),

  testsSortedBy: ['id'],
  testsSorted: Ember.computed.sort('tests', 'testsSortedBy')
***REMOVED***);
