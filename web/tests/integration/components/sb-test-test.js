import ***REMOVED*** moduleForComponent, test ***REMOVED*** from 'ember-qunit';
import hbs from 'htmlbars-inline-precompile';

moduleForComponent('sb-test', 'Integration | Component | sb test', ***REMOVED***
  integration: true
***REMOVED***);

test('it renders', function(assert) ***REMOVED***
  // Set any properties with this.set('myProperty', 'value');
  // Handle any actions with this.on('myAction', function(val) ***REMOVED*** ... ***REMOVED***);

  this.render(hbs`***REMOVED******REMOVED***sb-test***REMOVED******REMOVED***`);

  assert.equal(this.$().text().trim(), '');

  // Template block usage:
  this.render(hbs`
    ***REMOVED******REMOVED***#sb-test***REMOVED******REMOVED***
      template block text
    ***REMOVED******REMOVED***/sb-test***REMOVED******REMOVED***
  `);

  assert.equal(this.$().text().trim(), 'template block text');
***REMOVED***);
