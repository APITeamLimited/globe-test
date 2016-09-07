import DS from 'ember-data';

export default DS.JSONAPIAdapter.extend(***REMOVED***
  namespace: "v1",
  urlForFindRecord(id, modelName, snapshot) ***REMOVED***
    if (id === "default") ***REMOVED***
      return this.urlForFindAll(modelName, snapshot);
    ***REMOVED***
    return this._super(id, modelName, snapshot);
  ***REMOVED***,
***REMOVED***);
