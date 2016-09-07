import DS from 'ember-data';

export default DS.JSONAPIAdapter.extend(***REMOVED***
  namespace: "v1",
  pathForType(modelName) ***REMOVED***
    switch (modelName) ***REMOVED***
    case 'status':
      return modelName
    default:
      return this._super(modelName);
    ***REMOVED***
  ***REMOVED***,
  urlForFindRecord(id, modelName, snapshot) ***REMOVED***
    if (id === "default") ***REMOVED***
      return this.urlForFindAll(modelName, snapshot);
    ***REMOVED***
    return this._super(id, modelName, snapshot);
  ***REMOVED***,
***REMOVED***);
