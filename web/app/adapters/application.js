import DS from 'ember-data';

export default DS.JSONAPIAdapter.extend(***REMOVED***
  namespace: "v1",
  buildURL(modelName, id, snapshot, requestType, query) ***REMOVED***
    if (id == "default") ***REMOVED***
      return this.urlForFindAll(modelName, snapshot);
    ***REMOVED***
    return this._super(modelName, id, snapshot, requestType, query);
  ***REMOVED***
***REMOVED***);
