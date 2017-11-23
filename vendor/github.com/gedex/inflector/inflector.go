// Copyright 2013 Akeda Bagus <admin@gedex.web.id>. All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package inflector pluralizes and singularizes English nouns.

There are only two exported functions: `Pluralize` and `Singularize`.

	s := "People"
	fmt.Println(inflector.Singularize(s)) // will print "Person"

	s2 := "octopus"
	fmt.Println(inflector.Pluralize(s2)) // will print "octopuses"

*/
package inflector

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Rule represents name of the inflector rule, can be
// Plural or Singular
type Rule int

const (
	Plural = iota
	Singular
)

// InflectorRule represents inflector rule
type InflectorRule struct ***REMOVED***
	Rules               []*ruleItem
	Irregular           []*irregularItem
	Uninflected         []string
	compiledIrregular   *regexp.Regexp
	compiledUninflected *regexp.Regexp
	compiledRules       []*compiledRule
***REMOVED***

type ruleItem struct ***REMOVED***
	pattern     string
	replacement string
***REMOVED***

type irregularItem struct ***REMOVED***
	word        string
	replacement string
***REMOVED***

// compiledRule represents compiled version of Inflector.Rules.
type compiledRule struct ***REMOVED***
	replacement string
	*regexp.Regexp
***REMOVED***

// threadsafe access to rules and caches
var mutex sync.Mutex
var rules = make(map[Rule]*InflectorRule)

// Words that should not be inflected
var uninflected = []string***REMOVED***
	`Amoyese`, `bison`, `Borghese`, `bream`, `breeches`, `britches`, `buffalo`,
	`cantus`, `carp`, `chassis`, `clippers`, `cod`, `coitus`, `Congoese`,
	`contretemps`, `corps`, `debris`, `diabetes`, `djinn`, `eland`, `elk`,
	`equipment`, `Faroese`, `flounder`, `Foochowese`, `gallows`, `Genevese`,
	`Genoese`, `Gilbertese`, `graffiti`, `headquarters`, `herpes`, `hijinks`,
	`Hottentotese`, `information`, `innings`, `jackanapes`, `Kiplingese`,
	`Kongoese`, `Lucchese`, `mackerel`, `Maltese`, `.*?media`, `mews`, `moose`,
	`mumps`, `Nankingese`, `news`, `nexus`, `Niasese`, `Pekingese`,
	`Piedmontese`, `pincers`, `Pistoiese`, `pliers`, `Portuguese`, `proceedings`,
	`rabies`, `rice`, `rhinoceros`, `salmon`, `Sarawakese`, `scissors`,
	`sea[- ]bass`, `series`, `Shavese`, `shears`, `siemens`, `species`, `swine`,
	`testes`, `trousers`, `trout`, `tuna`, `Vermontese`, `Wenchowese`, `whiting`,
	`wildebeest`, `Yengeese`,
***REMOVED***

// Plural words that should not be inflected
var uninflectedPlurals = []string***REMOVED***
	`.*[nrlm]ese`, `.*deer`, `.*fish`, `.*measles`, `.*ois`, `.*pox`, `.*sheep`,
	`people`,
***REMOVED***

// Singular words that should not be inflected
var uninflectedSingulars = []string***REMOVED***
	`.*[nrlm]ese`, `.*deer`, `.*fish`, `.*measles`, `.*ois`, `.*pox`, `.*sheep`,
	`.*ss`,
***REMOVED***

type cache map[string]string

// Inflected words that already cached for immediate retrieval from a given Rule
var caches = make(map[Rule]cache)

// map of irregular words where its key is a word and its value is the replacement
var irregularMaps = make(map[Rule]cache)

func init() ***REMOVED***

	rules[Plural] = &InflectorRule***REMOVED***
		Rules: []*ruleItem***REMOVED***
			***REMOVED***`(?i)(s)tatus$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***tatuses`***REMOVED***,
			***REMOVED***`(?i)(quiz)$`, `$***REMOVED***1***REMOVED***zes`***REMOVED***,
			***REMOVED***`(?i)^(ox)$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***en`***REMOVED***,
			***REMOVED***`(?i)([m|l])ouse$`, `$***REMOVED***1***REMOVED***ice`***REMOVED***,
			***REMOVED***`(?i)(matr|vert|ind)(ix|ex)$`, `$***REMOVED***1***REMOVED***ices`***REMOVED***,
			***REMOVED***`(?i)(x|ch|ss|sh)$`, `$***REMOVED***1***REMOVED***es`***REMOVED***,
			***REMOVED***`(?i)([^aeiouy]|qu)y$`, `$***REMOVED***1***REMOVED***ies`***REMOVED***,
			***REMOVED***`(?i)(hive)$`, `$1s`***REMOVED***,
			***REMOVED***`(?i)(?:([^f])fe|([lre])f)$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***ves`***REMOVED***,
			***REMOVED***`(?i)sis$`, `ses`***REMOVED***,
			***REMOVED***`(?i)([ti])um$`, `$***REMOVED***1***REMOVED***a`***REMOVED***,
			***REMOVED***`(?i)(p)erson$`, `$***REMOVED***1***REMOVED***eople`***REMOVED***,
			***REMOVED***`(?i)(m)an$`, `$***REMOVED***1***REMOVED***en`***REMOVED***,
			***REMOVED***`(?i)(c)hild$`, `$***REMOVED***1***REMOVED***hildren`***REMOVED***,
			***REMOVED***`(?i)(buffal|tomat)o$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***oes`***REMOVED***,
			***REMOVED***`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|vir)us$`, `$***REMOVED***1***REMOVED***i`***REMOVED***,
			***REMOVED***`(?i)us$`, `uses`***REMOVED***,
			***REMOVED***`(?i)(alias)$`, `$***REMOVED***1***REMOVED***es`***REMOVED***,
			***REMOVED***`(?i)(ax|cris|test)is$`, `$***REMOVED***1***REMOVED***es`***REMOVED***,
			***REMOVED***`s$`, `s`***REMOVED***,
			***REMOVED***`^$`, ``***REMOVED***,
			***REMOVED***`$`, `s`***REMOVED***,
		***REMOVED***,
		Irregular: []*irregularItem***REMOVED***
			***REMOVED***`atlas`, `atlases`***REMOVED***,
			***REMOVED***`beef`, `beefs`***REMOVED***,
			***REMOVED***`brother`, `brothers`***REMOVED***,
			***REMOVED***`cafe`, `cafes`***REMOVED***,
			***REMOVED***`child`, `children`***REMOVED***,
			***REMOVED***`cookie`, `cookies`***REMOVED***,
			***REMOVED***`corpus`, `corpuses`***REMOVED***,
			***REMOVED***`cow`, `cows`***REMOVED***,
			***REMOVED***`ganglion`, `ganglions`***REMOVED***,
			***REMOVED***`genie`, `genies`***REMOVED***,
			***REMOVED***`genus`, `genera`***REMOVED***,
			***REMOVED***`graffito`, `graffiti`***REMOVED***,
			***REMOVED***`hoof`, `hoofs`***REMOVED***,
			***REMOVED***`loaf`, `loaves`***REMOVED***,
			***REMOVED***`man`, `men`***REMOVED***,
			***REMOVED***`money`, `monies`***REMOVED***,
			***REMOVED***`mongoose`, `mongooses`***REMOVED***,
			***REMOVED***`move`, `moves`***REMOVED***,
			***REMOVED***`mythos`, `mythoi`***REMOVED***,
			***REMOVED***`niche`, `niches`***REMOVED***,
			***REMOVED***`numen`, `numina`***REMOVED***,
			***REMOVED***`occiput`, `occiputs`***REMOVED***,
			***REMOVED***`octopus`, `octopuses`***REMOVED***,
			***REMOVED***`opus`, `opuses`***REMOVED***,
			***REMOVED***`ox`, `oxen`***REMOVED***,
			***REMOVED***`penis`, `penises`***REMOVED***,
			***REMOVED***`person`, `people`***REMOVED***,
			***REMOVED***`sex`, `sexes`***REMOVED***,
			***REMOVED***`soliloquy`, `soliloquies`***REMOVED***,
			***REMOVED***`testis`, `testes`***REMOVED***,
			***REMOVED***`trilby`, `trilbys`***REMOVED***,
			***REMOVED***`turf`, `turfs`***REMOVED***,
			***REMOVED***`potato`, `potatoes`***REMOVED***,
			***REMOVED***`hero`, `heroes`***REMOVED***,
			***REMOVED***`tooth`, `teeth`***REMOVED***,
			***REMOVED***`goose`, `geese`***REMOVED***,
			***REMOVED***`foot`, `feet`***REMOVED***,
		***REMOVED***,
	***REMOVED***
	prepare(Plural)

	rules[Singular] = &InflectorRule***REMOVED***
		Rules: []*ruleItem***REMOVED***
			***REMOVED***`(?i)(s)tatuses$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***tatus`***REMOVED***,
			***REMOVED***`(?i)^(.*)(menu)s$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***`***REMOVED***,
			***REMOVED***`(?i)(quiz)zes$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(matr)ices$`, `$***REMOVED***1***REMOVED***ix`***REMOVED***,
			***REMOVED***`(?i)(vert|ind)ices$`, `$***REMOVED***1***REMOVED***ex`***REMOVED***,
			***REMOVED***`(?i)^(ox)en`, `$1`***REMOVED***,
			***REMOVED***`(?i)(alias)(es)*$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(alumn|bacill|cact|foc|fung|nucle|radi|stimul|syllab|termin|viri?)i$`, `$***REMOVED***1***REMOVED***us`***REMOVED***,
			***REMOVED***`(?i)([ftw]ax)es`, `$1`***REMOVED***,
			***REMOVED***`(?i)(cris|ax|test)es$`, `$***REMOVED***1***REMOVED***is`***REMOVED***,
			***REMOVED***`(?i)(shoe|slave)s$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(o)es$`, `$1`***REMOVED***,
			***REMOVED***`ouses$`, `ouse`***REMOVED***,
			***REMOVED***`([^a])uses$`, `$***REMOVED***1***REMOVED***us`***REMOVED***,
			***REMOVED***`(?i)([m|l])ice$`, `$***REMOVED***1***REMOVED***ouse`***REMOVED***,
			***REMOVED***`(?i)(x|ch|ss|sh)es$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(m)ovies$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***ovie`***REMOVED***,
			***REMOVED***`(?i)(s)eries$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***eries`***REMOVED***,
			***REMOVED***`(?i)([^aeiouy]|qu)ies$`, `$***REMOVED***1***REMOVED***y`***REMOVED***,
			***REMOVED***`(?i)(tive)s$`, `$1`***REMOVED***,
			***REMOVED***`(?i)([lre])ves$`, `$***REMOVED***1***REMOVED***f`***REMOVED***,
			***REMOVED***`(?i)([^fo])ves$`, `$***REMOVED***1***REMOVED***fe`***REMOVED***,
			***REMOVED***`(?i)(hive)s$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(drive)s$`, `$1`***REMOVED***,
			***REMOVED***`(?i)(^analy)ses$`, `$***REMOVED***1***REMOVED***sis`***REMOVED***,
			***REMOVED***`(?i)(analy|diagno|^ba|(p)arenthe|(p)rogno|(s)ynop|(t)he)ses$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***sis`***REMOVED***,
			***REMOVED***`(?i)([ti])a$`, `$***REMOVED***1***REMOVED***um`***REMOVED***,
			***REMOVED***`(?i)(p)eople$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***erson`***REMOVED***,
			***REMOVED***`(?i)(m)en$`, `$***REMOVED***1***REMOVED***an`***REMOVED***,
			***REMOVED***`(?i)(c)hildren$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***hild`***REMOVED***,
			***REMOVED***`(?i)(n)ews$`, `$***REMOVED***1***REMOVED***$***REMOVED***2***REMOVED***ews`***REMOVED***,
			***REMOVED***`eaus$`, `eau`***REMOVED***,
			***REMOVED***`^(.*us)$`, `$1`***REMOVED***,
			***REMOVED***`(?i)s$`, ``***REMOVED***,
		***REMOVED***,
		Irregular: []*irregularItem***REMOVED***
			***REMOVED***`foes`, `foe`***REMOVED***,
			***REMOVED***`waves`, `wave`***REMOVED***,
			***REMOVED***`curves`, `curve`***REMOVED***,
			***REMOVED***`atlases`, `atlas`***REMOVED***,
			***REMOVED***`beefs`, `beef`***REMOVED***,
			***REMOVED***`brothers`, `brother`***REMOVED***,
			***REMOVED***`cafes`, `cafe`***REMOVED***,
			***REMOVED***`children`, `child`***REMOVED***,
			***REMOVED***`cookies`, `cookie`***REMOVED***,
			***REMOVED***`corpuses`, `corpus`***REMOVED***,
			***REMOVED***`cows`, `cow`***REMOVED***,
			***REMOVED***`ganglions`, `ganglion`***REMOVED***,
			***REMOVED***`genies`, `genie`***REMOVED***,
			***REMOVED***`genera`, `genus`***REMOVED***,
			***REMOVED***`graffiti`, `graffito`***REMOVED***,
			***REMOVED***`hoofs`, `hoof`***REMOVED***,
			***REMOVED***`loaves`, `loaf`***REMOVED***,
			***REMOVED***`men`, `man`***REMOVED***,
			***REMOVED***`monies`, `money`***REMOVED***,
			***REMOVED***`mongooses`, `mongoose`***REMOVED***,
			***REMOVED***`moves`, `move`***REMOVED***,
			***REMOVED***`mythoi`, `mythos`***REMOVED***,
			***REMOVED***`niches`, `niche`***REMOVED***,
			***REMOVED***`numina`, `numen`***REMOVED***,
			***REMOVED***`occiputs`, `occiput`***REMOVED***,
			***REMOVED***`octopuses`, `octopus`***REMOVED***,
			***REMOVED***`opuses`, `opus`***REMOVED***,
			***REMOVED***`oxen`, `ox`***REMOVED***,
			***REMOVED***`penises`, `penis`***REMOVED***,
			***REMOVED***`people`, `person`***REMOVED***,
			***REMOVED***`sexes`, `sex`***REMOVED***,
			***REMOVED***`soliloquies`, `soliloquy`***REMOVED***,
			***REMOVED***`testes`, `testis`***REMOVED***,
			***REMOVED***`trilbys`, `trilby`***REMOVED***,
			***REMOVED***`turfs`, `turf`***REMOVED***,
			***REMOVED***`potatoes`, `potato`***REMOVED***,
			***REMOVED***`heroes`, `hero`***REMOVED***,
			***REMOVED***`teeth`, `tooth`***REMOVED***,
			***REMOVED***`geese`, `goose`***REMOVED***,
			***REMOVED***`feet`, `foot`***REMOVED***,
		***REMOVED***,
	***REMOVED***
	prepare(Singular)
***REMOVED***

// prepare rule, e.g., compile the pattern.
func prepare(r Rule) error ***REMOVED***
	var reString string

	switch r ***REMOVED***
	case Plural:
		// Merge global uninflected with singularsUninflected
		rules[r].Uninflected = merge(uninflected, uninflectedPlurals)
	case Singular:
		// Merge global uninflected with singularsUninflected
		rules[r].Uninflected = merge(uninflected, uninflectedSingulars)
	***REMOVED***

	// Set InflectorRule.compiledUninflected by joining InflectorRule.Uninflected into
	// a single string then compile it.
	reString = fmt.Sprintf(`(?i)(^(?:%s))$`, strings.Join(rules[r].Uninflected, `|`))
	rules[r].compiledUninflected = regexp.MustCompile(reString)

	// Prepare irregularMaps
	irregularMaps[r] = make(cache, len(rules[r].Irregular))

	// Set InflectorRule.compiledIrregular by joining the irregularItem.word of Inflector.Irregular
	// into a single string then compile it.
	vIrregulars := make([]string, len(rules[r].Irregular))
	for i, item := range rules[r].Irregular ***REMOVED***
		vIrregulars[i] = item.word
		irregularMaps[r][item.word] = item.replacement
	***REMOVED***
	reString = fmt.Sprintf(`(?i)(.*)\b((?:%s))$`, strings.Join(vIrregulars, `|`))
	rules[r].compiledIrregular = regexp.MustCompile(reString)

	// Compile all patterns in InflectorRule.Rules
	rules[r].compiledRules = make([]*compiledRule, len(rules[r].Rules))
	for i, item := range rules[r].Rules ***REMOVED***
		rules[r].compiledRules[i] = &compiledRule***REMOVED***item.replacement, regexp.MustCompile(item.pattern)***REMOVED***
	***REMOVED***

	// Prepare caches
	caches[r] = make(cache)

	return nil
***REMOVED***

// merge slice a and slice b
func merge(a []string, b []string) []string ***REMOVED***
	result := make([]string, len(a)+len(b))
	copy(result, a)
	copy(result[len(a):], b)

	return result
***REMOVED***

// Pluralize returns string s in plural form.
func Pluralize(s string) string ***REMOVED***
	return getInflected(Plural, s)
***REMOVED***

// Singularize returns string s in singular form.
func Singularize(s string) string ***REMOVED***
	return getInflected(Singular, s)
***REMOVED***

func getInflected(r Rule, s string) string ***REMOVED***
	mutex.Lock()
	defer mutex.Unlock()
	if v, ok := caches[r][s]; ok ***REMOVED***
		return v
	***REMOVED***

	// Check for irregular words
	if res := rules[r].compiledIrregular.FindStringSubmatch(s); len(res) >= 3 ***REMOVED***
		var buf bytes.Buffer

		buf.WriteString(res[1])
		buf.WriteString(s[0:1])
		buf.WriteString(irregularMaps[r][strings.ToLower(res[2])][1:])

		// Cache it then returns
		caches[r][s] = buf.String()
		return caches[r][s]
	***REMOVED***

	// Check for uninflected words
	if rules[r].compiledUninflected.MatchString(s) ***REMOVED***
		caches[r][s] = s
		return caches[r][s]
	***REMOVED***

	// Check each rule
	for _, re := range rules[r].compiledRules ***REMOVED***
		if re.MatchString(s) ***REMOVED***
			caches[r][s] = re.ReplaceAllString(s, re.replacement)
			return caches[r][s]
		***REMOVED***
	***REMOVED***

	// Returns unaltered
	caches[r][s] = s
	return caches[r][s]
***REMOVED***
