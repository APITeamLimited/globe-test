(function () ***REMOVED***
    var jslib = ***REMOVED******REMOVED***;
    (function (module, exports) ***REMOVED***
        /*JSLIB_SUMMARY_CODE*/;
    ***REMOVED***)(***REMOVED*** exports: jslib ***REMOVED***, jslib);

    var forEach = function (obj, callback) ***REMOVED***
        for (var key in obj) ***REMOVED***
            if (obj.hasOwnProperty(key)) ***REMOVED***
                if (callback(key, obj[key])) ***REMOVED***
                    break;
                ***REMOVED***
            ***REMOVED***
        ***REMOVED***
    ***REMOVED***

    var transformGroup = function (group) ***REMOVED***
        if (Array.isArray(group.groups)) ***REMOVED***
            var newFormatGroups = group.groups;
            group.groups = ***REMOVED******REMOVED***;
            for (var i = 0; i < newFormatGroups.length; i++) ***REMOVED***
                group.groups[newFormatGroups[i].name] = transformGroup(newFormatGroups[i]);
            ***REMOVED***
        ***REMOVED***
        if (Array.isArray(group.checks)) ***REMOVED***
            var newFormatChecks = group.checks;
            group.checks = ***REMOVED******REMOVED***;
            for (var i = 0; i < newFormatChecks.length; i++) ***REMOVED***
                group.checks[newFormatChecks[i].name] = newFormatChecks[i];
            ***REMOVED***
        ***REMOVED***
        return group;
    ***REMOVED***;

    var oldJSONSummary = function (data) ***REMOVED***
        // Quick copy of the data, since it's easiest to modify it in place.
        var results = JSON.parse(JSON.stringify(data));
        delete results.options;
        delete results.state;

        forEach(results.metrics, function (metricName, metric) ***REMOVED***
            var oldFormatMetric = metric.values;
            if (metric.thresholds && Object.keys(metric.thresholds).length > 0) ***REMOVED***
                var newFormatThresholds = metric.thresholds;
                oldFormatMetric.thresholds = ***REMOVED******REMOVED***;
                forEach(newFormatThresholds, function (thresholdName, threshold) ***REMOVED***
                    oldFormatMetric.thresholds[thresholdName] = !threshold.ok;
                ***REMOVED***);
            ***REMOVED***
            if (metric.type == 'rate' && oldFormatMetric.hasOwnProperty('rate')) ***REMOVED***
                oldFormatMetric.value = oldFormatMetric.rate; // sigh...
                delete oldFormatMetric.rate;
            ***REMOVED***
            results.metrics[metricName] = oldFormatMetric;
        ***REMOVED***);

        results.root_group = transformGroup(results.root_group);

        return JSON.stringify(results, null, 4);
    ***REMOVED***;

    return function (exportedSummaryCallback, jsonSummaryPath, data) ***REMOVED***
        var getDefaultSummary = function () ***REMOVED***
            var enableColors = (!data.options.noColor && data.state.isStdOutTTY);
            return ***REMOVED***
                'stdout': '\n' + jslibWorker.textSummary(data, ***REMOVED*** indent: ' ', enableColors: enableColors ***REMOVED***) + '\n\n',
            ***REMOVED***;
        ***REMOVED***;

        var result = ***REMOVED******REMOVED***;
        if (exportedSummaryCallback) ***REMOVED***
            try ***REMOVED***
                result = exportedSummaryCallback(data);
            ***REMOVED*** catch (e) ***REMOVED***
                console.error('handleSummary() failed with error "' + e + '", falling back to the default summary');
                result = getDefaultSummary();
            ***REMOVED***
        ***REMOVED*** else ***REMOVED***
            result = getDefaultSummary();
        ***REMOVED***

        // TODO: ensure we're returning a map of strings or null/undefined...
        // and if not, log an error and generate the default summary?

        if (jsonSummaryPath != '') ***REMOVED***
            result[jsonSummaryPath] = oldJSONSummary(data);
        ***REMOVED***

        return result;
    ***REMOVED***;
***REMOVED***)();