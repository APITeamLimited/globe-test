import ***REMOVED*** moduleForModel, test ***REMOVED*** from 'ember-qunit';

moduleForModel('group', 'Unit | Model | group', ***REMOVED***
  // Specify the other units that are required for this test.
  needs: ['model:group', 'model:test']
***REMOVED***);

test('it exists', function(assert) ***REMOVED***
  let model = this.subject();
  // let store = this.store();
  assert.ok(!!model);
***REMOVED***);
