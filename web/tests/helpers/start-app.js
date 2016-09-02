import Ember from 'ember';
import Application from '../../app';
import config from '../../config/environment';

export default function startApp(attrs) ***REMOVED***
  let application;

  let attributes = Ember.merge(***REMOVED******REMOVED***, config.APP);
  attributes = Ember.merge(attributes, attrs); // use defaults, but you can override;

  Ember.run(() => ***REMOVED***
    application = Application.create(attributes);
    application.setupForTesting();
    application.injectTestHelpers();
  ***REMOVED***);

  return application;
***REMOVED***
