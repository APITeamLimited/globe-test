var forEach = function (obj, callback) ***REMOVED***
  for (var key in obj) ***REMOVED***
    if (obj.hasOwnProperty(key)) ***REMOVED***
      if (callback(key, obj[key])) ***REMOVED***
        break
      ***REMOVED***
    ***REMOVED***
  ***REMOVED***
***REMOVED***

var palette = ***REMOVED***
  bold: 1,
  faint: 2,
  red: 31,
  green: 32,
  cyan: 36,
  //TODO: add others?
***REMOVED***

var groupPrefix = '█'
var detailsPrefix = '↳'
var succMark = '✓'
var failMark = '✗'
var defaultOptions = ***REMOVED***
  indent: ' ',
  enableColors: true,
  summaryTimeUnit: null,
  summaryTrendStats: null,
***REMOVED***

// strWidth tries to return the actual width the string will take up on the
// screen, without any terminal formatting, unicode ligatures, etc.
function strWidth(s) ***REMOVED***
  // TODO: determine if NFC or NFKD are not more appropriate? or just give up? https://hsivonen.fi/string-length/
  var data = s.normalize('NFKC') // This used to be NFKD in Go, but this should be better
  var inEscSeq = false
  var inLongEscSeq = false
  var width = 0
  for (var char of data) ***REMOVED***
    if (char.done) ***REMOVED***
      break
    ***REMOVED***

    // Skip over ANSI escape codes.
    if (char == '\x1b') ***REMOVED***
      inEscSeq = true
      continue
    ***REMOVED***
    if (inEscSeq && char == '[') ***REMOVED***
      inLongEscSeq = true
      continue
    ***REMOVED***
    if (inEscSeq && inLongEscSeq && char.charCodeAt(0) >= 0x40 && char.charCodeAt(0) <= 0x7e) ***REMOVED***
      inEscSeq = false
      inLongEscSeq = false
      continue
    ***REMOVED***
    if (inEscSeq && !inLongEscSeq && char.charCodeAt(0) >= 0x40 && char.charCodeAt(0) <= 0x5f) ***REMOVED***
      inEscSeq = false
      continue
    ***REMOVED***

    if (!inEscSeq && !inLongEscSeq) ***REMOVED***
      width++
    ***REMOVED***
  ***REMOVED***
  return width
***REMOVED***

function summarizeCheck(indent, check, decorate) ***REMOVED***
  if (check.fails == 0) ***REMOVED***
    return decorate(indent + succMark + ' ' + check.name, palette.green)
  ***REMOVED***

  var succPercent = Math.floor((100 * check.passes) / (check.passes + check.fails))
  return decorate(
    indent +
    failMark +
    ' ' +
    check.name +
    '\n' +
    indent +
    ' ' +
    detailsPrefix +
    '  ' +
    succPercent +
    '% — ' +
    succMark +
    ' ' +
    check.passes +
    ' / ' +
    failMark +
    ' ' +
    check.fails,
    palette.red
  )
***REMOVED***

function summarizeGroup(indent, group, decorate) ***REMOVED***
  var result = []
  if (group.name != '') ***REMOVED***
    result.push(indent + groupPrefix + ' ' + group.name + '\n')
    indent = indent + '  '
  ***REMOVED***

  for (var i = 0; i < group.checks.length; i++) ***REMOVED***
    result.push(summarizeCheck(indent, group.checks[i], decorate))
  ***REMOVED***
  if (group.checks.length > 0) ***REMOVED***
    result.push('')
  ***REMOVED***
  for (var i = 0; i < group.groups.length; i++) ***REMOVED***
    Array.prototype.push.apply(result, summarizeGroup(indent, group.groups[i], decorate))
  ***REMOVED***

  return result
***REMOVED***

function displayNameForMetric(name) ***REMOVED***
  var subMetricPos = name.indexOf('***REMOVED***')
  if (subMetricPos >= 0) ***REMOVED***
    return '***REMOVED*** ' + name.substring(subMetricPos + 1, name.length - 1) + ' ***REMOVED***'
  ***REMOVED***
  return name
***REMOVED***

function indentForMetric(name) ***REMOVED***
  if (name.indexOf('***REMOVED***') >= 0) ***REMOVED***
    return '  '
  ***REMOVED***
  return ''
***REMOVED***

function humanizeBytes(bytes) ***REMOVED***
  var units = ['B', 'kB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
  var base = 1000
  if (bytes < 10) ***REMOVED***
    return bytes + ' B'
  ***REMOVED***

  var e = Math.floor(Math.log(bytes) / Math.log(base))
  var suffix = units[e | 0]
  var val = Math.floor((bytes / Math.pow(base, e)) * 10 + 0.5) / 10
  return val.toFixed(val < 10 ? 1 : 0) + ' ' + suffix
***REMOVED***

var unitMap = ***REMOVED***
  s: ***REMOVED*** unit: 's', coef: 0.001 ***REMOVED***,
  ms: ***REMOVED*** unit: 'ms', coef: 1 ***REMOVED***,
  us: ***REMOVED*** unit: 'µs', coef: 1000 ***REMOVED***,
***REMOVED***

function toFixedNoTrailingZeros(val, prec) ***REMOVED***
  // TODO: figure out something better?
  return parseFloat(val.toFixed(prec)).toString()
***REMOVED***

function toFixedNoTrailingZerosTrunc(val, prec) ***REMOVED***
  var mult = Math.pow(10, prec)
  return toFixedNoTrailingZeros(Math.trunc(mult * val) / mult, prec)
***REMOVED***

function humanizeGenericDuration(dur) ***REMOVED***
  if (dur === 0) ***REMOVED***
    return '0s'
  ***REMOVED***

  if (dur < 0.001) ***REMOVED***
    // smaller than a microsecond, print nanoseconds
    return Math.trunc(dur * 1000000) + 'ns'
  ***REMOVED***
  if (dur < 1) ***REMOVED***
    // smaller than a millisecond, print microseconds
    return toFixedNoTrailingZerosTrunc(dur * 1000, 2) + 'µs'
  ***REMOVED***
  if (dur < 1000) ***REMOVED***
    // duration is smaller than a second
    return toFixedNoTrailingZerosTrunc(dur, 2) + 'ms'
  ***REMOVED***

  var result = toFixedNoTrailingZerosTrunc((dur % 60000) / 1000, dur > 60000 ? 0 : 2) + 's'
  var rem = Math.trunc(dur / 60000)
  if (rem < 1) ***REMOVED***
    // less than a minute
    return result
  ***REMOVED***
  result = (rem % 60) + 'm' + result
  rem = Math.trunc(rem / 60)
  if (rem < 1) ***REMOVED***
    // less than an hour
    return result
  ***REMOVED***
  return rem + 'h' + result
***REMOVED***

function humanizeDuration(dur, timeUnit) ***REMOVED***
  if (timeUnit !== '' && unitMap.hasOwnProperty(timeUnit)) ***REMOVED***
    return (dur * unitMap[timeUnit].coef).toFixed(2) + unitMap[timeUnit].unit
  ***REMOVED***

  return humanizeGenericDuration(dur)
***REMOVED***

function humanizeValue(val, metric, timeUnit) ***REMOVED***
  if (metric.type == 'rate') ***REMOVED***
    // Truncate instead of round when decreasing precision to 2 decimal places
    return (Math.trunc(val * 100 * 100) / 100).toFixed(2) + '%'
  ***REMOVED***

  switch (metric.contains) ***REMOVED***
    case 'data':
      return humanizeBytes(val)
    case 'time':
      return humanizeDuration(val, timeUnit)
    default:
      return toFixedNoTrailingZeros(val, 6)
  ***REMOVED***
***REMOVED***

function nonTrendMetricValueForSum(metric, timeUnit) ***REMOVED***
  switch (metric.type) ***REMOVED***
    case 'counter':
      return [
        humanizeValue(metric.values.count, metric, timeUnit),
        humanizeValue(metric.values.rate, metric, timeUnit) + '/s',
      ]
    case 'gauge':
      return [
        humanizeValue(metric.values.value, metric, timeUnit),
        'min=' + humanizeValue(metric.values.min, metric, timeUnit),
        'max=' + humanizeValue(metric.values.max, metric, timeUnit),
      ]
    case 'rate':
      return [
        humanizeValue(metric.values.rate, metric, timeUnit),
        succMark + ' ' + metric.values.passes,
        failMark + ' ' + metric.values.fails,
      ]
    default:
      return ['[no data]']
  ***REMOVED***
***REMOVED***

function summarizeMetrics(options, data, decorate) ***REMOVED***
  var indent = options.indent + '  '
  var result = []

  var names = []
  var nameLenMax = 0

  var nonTrendValues = ***REMOVED******REMOVED***
  var nonTrendValueMaxLen = 0
  var nonTrendExtras = ***REMOVED******REMOVED***
  var nonTrendExtraMaxLens = [0, 0]

  var trendCols = ***REMOVED******REMOVED***
  var numTrendColumns = options.summaryTrendStats.length
  var trendColMaxLens = new Array(numTrendColumns).fill(0)
  forEach(data.metrics, function (name, metric) ***REMOVED***
    names.push(name)
    // When calculating widths for metrics, account for the indentation on submetrics.
    var displayName = indentForMetric(name) + displayNameForMetric(name)
    var displayNameWidth = strWidth(displayName)
    if (displayNameWidth > nameLenMax) ***REMOVED***
      nameLenMax = displayNameWidth
    ***REMOVED***

    if (metric.type == 'trend') ***REMOVED***
      var cols = []
      for (var i = 0; i < numTrendColumns; i++) ***REMOVED***
        var tc = options.summaryTrendStats[i]
        var value = metric.values[tc]
        if (tc === 'count') ***REMOVED***
          value = value.toString()
        ***REMOVED*** else ***REMOVED***
          value = humanizeValue(value, metric, options.summaryTimeUnit)
        ***REMOVED***
        var valLen = strWidth(value)
        if (valLen > trendColMaxLens[i]) ***REMOVED***
          trendColMaxLens[i] = valLen
        ***REMOVED***
        cols[i] = value
      ***REMOVED***
      trendCols[name] = cols
      return
    ***REMOVED***
    var values = nonTrendMetricValueForSum(metric, options.summaryTimeUnit)
    nonTrendValues[name] = values[0]
    var valueLen = strWidth(values[0])
    if (valueLen > nonTrendValueMaxLen) ***REMOVED***
      nonTrendValueMaxLen = valueLen
    ***REMOVED***
    nonTrendExtras[name] = values.slice(1)
    for (var i = 1; i < values.length; i++) ***REMOVED***
      var extraLen = strWidth(values[i])
      if (extraLen > nonTrendExtraMaxLens[i - 1]) ***REMOVED***
        nonTrendExtraMaxLens[i - 1] = extraLen
      ***REMOVED***
    ***REMOVED***
  ***REMOVED***)

  names.sort()

  var getData = function (name) ***REMOVED***
    if (trendCols.hasOwnProperty(name)) ***REMOVED***
      var cols = trendCols[name]
      var tmpCols = new Array(numTrendColumns)
      for (var i = 0; i < cols.length; i++) ***REMOVED***
        tmpCols[i] =
          options.summaryTrendStats[i] +
          '=' +
          decorate(cols[i], palette.cyan) +
          ' '.repeat(trendColMaxLens[i] - strWidth(cols[i]))
      ***REMOVED***
      return tmpCols.join(' ')
    ***REMOVED***

    var value = nonTrendValues[name]
    var fmtData = decorate(value, palette.cyan) + ' '.repeat(nonTrendValueMaxLen - strWidth(value))

    var extras = nonTrendExtras[name]
    if (extras.length == 1) ***REMOVED***
      fmtData = fmtData + ' ' + decorate(extras[0], palette.cyan, palette.faint)
    ***REMOVED*** else if (extras.length > 1) ***REMOVED***
      var parts = new Array(extras.length)
      for (var i = 0; i < extras.length; i++) ***REMOVED***
        parts[i] =
          decorate(extras[i], palette.cyan, palette.faint) +
          ' '.repeat(nonTrendExtraMaxLens[i] - strWidth(extras[i]))
      ***REMOVED***
      fmtData = fmtData + ' ' + parts.join(' ')
    ***REMOVED***

    return fmtData
  ***REMOVED***

  for (var name of names) ***REMOVED***
    var metric = data.metrics[name]
    var mark = ' '
    var markColor = function (text) ***REMOVED***
      return text
    ***REMOVED*** // noop

    if (metric.thresholds) ***REMOVED***
      mark = succMark
      markColor = function (text) ***REMOVED***
        return decorate(text, palette.green)
      ***REMOVED***
      forEach(metric.thresholds, function (name, threshold) ***REMOVED***
        if (!threshold.ok) ***REMOVED***
          mark = failMark
          markColor = function (text) ***REMOVED***
            return decorate(text, palette.red)
          ***REMOVED***
          return true // break
        ***REMOVED***
      ***REMOVED***)
    ***REMOVED***
    var fmtIndent = indentForMetric(name)
    var fmtName = displayNameForMetric(name)
    fmtName =
      fmtName +
      decorate(
        '.'.repeat(nameLenMax - strWidth(fmtName) - strWidth(fmtIndent) + 3) + ':',
        palette.faint
      )

    result.push(indent + fmtIndent + markColor(mark) + ' ' + fmtName + ' ' + getData(name))
  ***REMOVED***

  return result
***REMOVED***

function generateTextSummary(data, options) ***REMOVED***
  var mergedOpts = Object.assign(***REMOVED******REMOVED***, defaultOptions, data.options, options)
  var lines = []

  // TODO: move all of these functions into an object with methods?
  var decorate = function (text) ***REMOVED***
    return text
  ***REMOVED***
  if (mergedOpts.enableColors) ***REMOVED***
    decorate = function (text, color /*, ...rest*/) ***REMOVED***
      var result = '\x1b[' + color
      for (var i = 2; i < arguments.length; i++) ***REMOVED***
        result += ';' + arguments[i]
      ***REMOVED***
      return result + 'm' + text + '\x1b[0m'
    ***REMOVED***
  ***REMOVED***

  Array.prototype.push.apply(
    lines,
    summarizeGroup(mergedOpts.indent + '    ', data.root_group, decorate)
  )

  Array.prototype.push.apply(lines, summarizeMetrics(mergedOpts, data, decorate))

  return lines.join('\n')
***REMOVED***

exports.humanizeValue = humanizeValue
exports.textSummary = generateTextSummary