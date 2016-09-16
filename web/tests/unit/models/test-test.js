import ***REMOVED*** moduleForModel, test ***REMOVED*** from 'ember-qunit';

moduleForModel('test', 'Unit | Model | test', ***REMOVED***
  // Specify the other units that are required for this test.
  needs: ['model:group']
***REMOVED***);

test('it exists', function(assert) ***REMOVED***
  let model = this.subject();
  // let store = this.store();
  assert.ok(!!model);
***REMOVED***);
