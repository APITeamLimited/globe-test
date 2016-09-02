import Resolver from '../../resolver';
import config from '../../config/environment';

const resolver = Resolver.create();

resolver.namespace = ***REMOVED***
  modulePrefix: config.modulePrefix,
  podModulePrefix: config.podModulePrefix
***REMOVED***;

export default resolver;
