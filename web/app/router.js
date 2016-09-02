import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend(***REMOVED***
  location: config.locationType,
  rootURL: config.rootURL
***REMOVED***);

Router.map(function() ***REMOVED***
***REMOVED***);

export default Router;
