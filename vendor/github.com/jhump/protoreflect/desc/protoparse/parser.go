package protoparse

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/internal"
	"github.com/jhump/protoreflect/desc/protoparse/ast"
)

//go:generate goyacc -o proto.y.go -p proto proto.y

func init() ***REMOVED***
	protoErrorVerbose = true

	// fix up the generated "token name" array so that error messages are nicer
	setTokenName(_STRING_LIT, "string literal")
	setTokenName(_INT_LIT, "int literal")
	setTokenName(_FLOAT_LIT, "float literal")
	setTokenName(_NAME, "identifier")
	setTokenName(_ERROR, "error")
	// for keywords, just show the keyword itself wrapped in quotes
	for str, i := range keywords ***REMOVED***
		setTokenName(i, fmt.Sprintf(`"%s"`, str))
	***REMOVED***
***REMOVED***

func setTokenName(token int, text string) ***REMOVED***
	// NB: this is based on logic in generated parse code that translates the
	// int returned from the lexer into an internal token number.
	var intern int
	if token < len(protoTok1) ***REMOVED***
		intern = protoTok1[token]
	***REMOVED*** else ***REMOVED***
		if token >= protoPrivate ***REMOVED***
			if token < protoPrivate+len(protoTok2) ***REMOVED***
				intern = protoTok2[token-protoPrivate]
			***REMOVED***
		***REMOVED***
		if intern == 0 ***REMOVED***
			for i := 0; i+1 < len(protoTok3); i += 2 ***REMOVED***
				if protoTok3[i] == token ***REMOVED***
					intern = protoTok3[i+1]
					break
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***

	if intern >= 1 && intern-1 < len(protoToknames) ***REMOVED***
		protoToknames[intern-1] = text
		return
	***REMOVED***

	panic(fmt.Sprintf("Unknown token value: %d", token))
***REMOVED***

// FileAccessor is an abstraction for opening proto source files. It takes the
// name of the file to open and returns either the input reader or an error.
type FileAccessor func(filename string) (io.ReadCloser, error)

// FileContentsFromMap returns a FileAccessor that uses the given map of file
// contents. This allows proto source files to be constructed in memory and
// easily supplied to a parser. The map keys are the paths to the proto source
// files, and the values are the actual proto source contents.
func FileContentsFromMap(files map[string]string) FileAccessor ***REMOVED***
	return func(filename string) (io.ReadCloser, error) ***REMOVED***
		contents, ok := files[filename]
		if !ok ***REMOVED***
			return nil, os.ErrNotExist
		***REMOVED***
		return ioutil.NopCloser(strings.NewReader(contents)), nil
	***REMOVED***
***REMOVED***

// Parser parses proto source into descriptors.
type Parser struct ***REMOVED***
	// The paths used to search for dependencies that are referenced in import
	// statements in proto source files. If no import paths are provided then
	// "." (current directory) is assumed to be the only import path.
	//
	// This setting is only used during ParseFiles operations. Since calls to
	// ParseFilesButDoNotLink do not link, there is no need to load and parse
	// dependencies.
	ImportPaths []string

	// If true, the supplied file names/paths need not necessarily match how the
	// files are referenced in import statements. The parser will attempt to
	// match import statements to supplied paths, "guessing" the import paths
	// for the files. Note that this inference is not perfect and link errors
	// could result. It works best when all proto files are organized such that
	// a single import path can be inferred (e.g. all files under a single tree
	// with import statements all being relative to the root of this tree).
	InferImportPaths bool

	// LookupImport is a function that accepts a filename and
	// returns a file descriptor, which will be consulted when resolving imports.
	// This allows a compiled Go proto in another Go module to be referenced
	// in the proto(s) being parsed.
	//
	// In the event of a filename collision, Accessor is consulted first,
	// then LookupImport is consulted, and finally the well-known protos
	// are used.
	//
	// For example, in order to automatically look up compiled Go protos that
	// have been imported and be able to use them as imports, set this to
	// desc.LoadFileDescriptor.
	LookupImport func(string) (*desc.FileDescriptor, error)

	// LookupImportProto has the same functionality as LookupImport, however it returns
	// a FileDescriptorProto instead of a FileDescriptor.
	//
	// It is an error to set both LookupImport and LookupImportProto.
	LookupImportProto func(string) (*dpb.FileDescriptorProto, error)

	// Used to create a reader for a given filename, when loading proto source
	// file contents. If unset, os.Open is used. If ImportPaths is also empty
	// then relative paths are will be relative to the process's current working
	// directory.
	Accessor FileAccessor

	// If true, the resulting file descriptors will retain source code info,
	// that maps elements to their location in the source files as well as
	// includes comments found during parsing (and attributed to elements of
	// the source file).
	IncludeSourceCodeInfo bool

	// If true, the results from ParseFilesButDoNotLink will be passed through
	// some additional validations. But only constraints that do not require
	// linking can be checked. These include proto2 vs. proto3 language features,
	// looking for incorrect usage of reserved names or tags, and ensuring that
	// fields have unique tags and that enum values have unique numbers (unless
	// the enum allows aliases).
	ValidateUnlinkedFiles bool

	// If true, the results from ParseFilesButDoNotLink will have options
	// interpreted. Any uninterpretable options (including any custom options or
	// options that refer to message and enum types, which can only be
	// interpreted after linking) will be left in uninterpreted_options. Also,
	// the "default" pseudo-option for fields can only be interpreted for scalar
	// fields, excluding enums. (Interpreting default values for enum fields
	// requires resolving enum names, which requires linking.)
	InterpretOptionsInUnlinkedFiles bool

	// A custom reporter of syntax and link errors. If not specified, the
	// default reporter just returns the reported error, which causes parsing
	// to abort after encountering a single error.
	//
	// The reporter is not invoked for system or I/O errors, only for syntax and
	// link errors.
	ErrorReporter ErrorReporter

	// A custom reporter of warnings. If not specified, warning messages are ignored.
	WarningReporter WarningReporter
***REMOVED***

// ParseFiles parses the named files into descriptors. The returned slice has
// the same number of entries as the give filenames, in the same order. So the
// first returned descriptor corresponds to the first given name, and so on.
//
// All dependencies for all specified files (including transitive dependencies)
// must be accessible via the parser's Accessor or a link error will occur. The
// exception to this rule is that files can import standard Google-provided
// files -- e.g. google/protobuf/*.proto -- without needing to supply sources
// for these files. Like protoc, this parser has a built-in version of these
// files it can use if they aren't explicitly supplied.
//
// If the Parser has no ErrorReporter set and a syntax or link error occurs,
// parsing will abort with the first such error encountered. If there is an
// ErrorReporter configured and it returns non-nil, parsing will abort with the
// error it returns. If syntax or link errors are encountered but the configured
// ErrorReporter always returns nil, the parse fails with ErrInvalidSource.
func (p Parser) ParseFiles(filenames ...string) ([]*desc.FileDescriptor, error) ***REMOVED***
	accessor := p.Accessor
	if accessor == nil ***REMOVED***
		accessor = func(name string) (io.ReadCloser, error) ***REMOVED***
			return os.Open(name)
		***REMOVED***
	***REMOVED***
	paths := p.ImportPaths
	if len(paths) > 0 ***REMOVED***
		acc := accessor
		accessor = func(name string) (io.ReadCloser, error) ***REMOVED***
			var ret error
			for _, path := range paths ***REMOVED***
				f, err := acc(filepath.Join(path, name))
				if err != nil ***REMOVED***
					if ret == nil ***REMOVED***
						ret = err
					***REMOVED***
					continue
				***REMOVED***
				return f, nil
			***REMOVED***
			return nil, ret
		***REMOVED***
	***REMOVED***
	lookupImport, err := p.getLookupImport()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	protos := map[string]*parseResult***REMOVED******REMOVED***
	results := &parseResults***REMOVED***
		resultsByFilename:      protos,
		recursive:              true,
		validate:               true,
		createDescriptorProtos: true,
	***REMOVED***
	errs := newErrorHandler(p.ErrorReporter, p.WarningReporter)
	parseProtoFiles(accessor, filenames, errs, results, lookupImport)
	if err := errs.getError(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if p.InferImportPaths ***REMOVED***
		// TODO: if this re-writes one of the names in filenames, lookups below will break
		protos = fixupFilenames(protos)
	***REMOVED***
	linkedProtos, err := newLinker(results, errs).linkFiles()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if p.IncludeSourceCodeInfo ***REMOVED***
		for name, fd := range linkedProtos ***REMOVED***
			pr := protos[name]
			fd.AsFileDescriptorProto().SourceCodeInfo = pr.generateSourceCodeInfo()
			internal.RecomputeSourceInfo(fd)
		***REMOVED***
	***REMOVED***
	fds := make([]*desc.FileDescriptor, len(filenames))
	for i, name := range filenames ***REMOVED***
		fd := linkedProtos[name]
		fds[i] = fd
	***REMOVED***
	return fds, nil
***REMOVED***

// ParseFilesButDoNotLink parses the named files into descriptor protos. The
// results are just protos, not fully-linked descriptors. It is possible that
// descriptors are invalid and still be returned in parsed form without error
// due to the fact that the linking step is skipped (and thus many validation
// steps omitted).
//
// There are a few side effects to not linking the descriptors:
//   1. No options will be interpreted. Options can refer to extensions or have
//      message and enum types. Without linking, these extension and type
//      references are not resolved, so the options may not be interpretable.
//      So all options will appear in UninterpretedOption fields of the various
//      descriptor options messages.
//   2. Type references will not be resolved. This means that the actual type
//      names in the descriptors may be unqualified and even relative to the
//      scope in which the type reference appears. This goes for fields that
//      have message and enum types. It also applies to methods and their
//      references to request and response message types.
//   3. Type references are not known. For non-scalar fields, until the type
//      name is resolved (during linking), it is not known whether the type
//      refers to a message or an enum. So all fields with such type references
//      will not have their Type set, only the TypeName.
//
// This method will still validate the syntax of parsed files. If the parser's
// ValidateUnlinkedFiles field is true, additional checks, beyond syntax will
// also be performed.
//
// If the Parser has no ErrorReporter set and a syntax error occurs, parsing
// will abort with the first such error encountered. If there is an
// ErrorReporter configured and it returns non-nil, parsing will abort with the
// error it returns. If syntax errors are encountered but the configured
// ErrorReporter always returns nil, the parse fails with ErrInvalidSource.
func (p Parser) ParseFilesButDoNotLink(filenames ...string) ([]*dpb.FileDescriptorProto, error) ***REMOVED***
	accessor := p.Accessor
	if accessor == nil ***REMOVED***
		accessor = func(name string) (io.ReadCloser, error) ***REMOVED***
			return os.Open(name)
		***REMOVED***
	***REMOVED***
	lookupImport, err := p.getLookupImport()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	protos := map[string]*parseResult***REMOVED******REMOVED***
	errs := newErrorHandler(p.ErrorReporter, p.WarningReporter)
	results := &parseResults***REMOVED***
		resultsByFilename:      protos,
		validate:               p.ValidateUnlinkedFiles,
		createDescriptorProtos: true,
	***REMOVED***
	parseProtoFiles(accessor, filenames, errs, results, lookupImport)
	if err := errs.getError(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	if p.InferImportPaths ***REMOVED***
		// TODO: if this re-writes one of the names in filenames, lookups below will break
		protos = fixupFilenames(protos)
	***REMOVED***
	fds := make([]*dpb.FileDescriptorProto, len(filenames))
	for i, name := range filenames ***REMOVED***
		pr := protos[name]
		fd := pr.fd
		if p.InterpretOptionsInUnlinkedFiles ***REMOVED***
			// parsing options will be best effort
			pr.lenient = true
			// we don't want the real error reporter see any errors
			pr.errs.errReporter = func(err ErrorWithPos) error ***REMOVED***
				return err
			***REMOVED***
			_ = interpretFileOptions(pr, poorFileDescriptorish***REMOVED***FileDescriptorProto: fd***REMOVED***)
		***REMOVED***
		if p.IncludeSourceCodeInfo ***REMOVED***
			fd.SourceCodeInfo = pr.generateSourceCodeInfo()
		***REMOVED***
		fds[i] = fd
	***REMOVED***
	return fds, nil
***REMOVED***

// ParseToAST parses the named files into ASTs, or Abstract Syntax Trees. This
// is for consumers of proto files that don't care about compiling the files to
// descriptors, but care deeply about a non-lossy structured representation of
// the source (since descriptors are lossy). This includes formatting tools and
// possibly linters, too.
//
// If the requested filenames include standard imports (such as
// "google/protobuf/empty.proto") and no source is provided, the corresponding
// AST in the returned slice will be nil. These standard imports are only
// available for use as descriptors; no source is available unless it is
// provided by the configured Accessor.
//
// If the Parser has no ErrorReporter set and a syntax error occurs, parsing
// will abort with the first such error encountered. If there is an
// ErrorReporter configured and it returns non-nil, parsing will abort with the
// error it returns. If syntax errors are encountered but the configured
// ErrorReporter always returns nil, the parse fails with ErrInvalidSource.
func (p Parser) ParseToAST(filenames ...string) ([]*ast.FileNode, error) ***REMOVED***
	accessor := p.Accessor
	if accessor == nil ***REMOVED***
		accessor = func(name string) (io.ReadCloser, error) ***REMOVED***
			return os.Open(name)
		***REMOVED***
	***REMOVED***
	lookupImport, err := p.getLookupImport()
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	protos := map[string]*parseResult***REMOVED******REMOVED***
	errs := newErrorHandler(p.ErrorReporter, p.WarningReporter)
	parseProtoFiles(accessor, filenames, errs, &parseResults***REMOVED***resultsByFilename: protos***REMOVED***, lookupImport)
	if err := errs.getError(); err != nil ***REMOVED***
		return nil, err
	***REMOVED***
	ret := make([]*ast.FileNode, 0, len(filenames))
	for _, name := range filenames ***REMOVED***
		ret = append(ret, protos[name].root)
	***REMOVED***
	return ret, nil
***REMOVED***

func (p Parser) getLookupImport() (func(string) (*dpb.FileDescriptorProto, error), error) ***REMOVED***
	if p.LookupImport != nil && p.LookupImportProto != nil ***REMOVED***
		return nil, ErrLookupImportAndProtoSet
	***REMOVED***
	if p.LookupImportProto != nil ***REMOVED***
		return p.LookupImportProto, nil
	***REMOVED***
	if p.LookupImport != nil ***REMOVED***
		return func(path string) (*dpb.FileDescriptorProto, error) ***REMOVED***
			value, err := p.LookupImport(path)
			if value != nil ***REMOVED***
				return value.AsFileDescriptorProto(), err
			***REMOVED***
			return nil, err
		***REMOVED***, nil
	***REMOVED***
	return nil, nil
***REMOVED***

func fixupFilenames(protos map[string]*parseResult) map[string]*parseResult ***REMOVED***
	// In the event that the given filenames (keys in the supplied map) do not
	// match the actual paths used in 'import' statements in the files, we try
	// to revise names in the protos so that they will match and be linkable.
	revisedProtos := map[string]*parseResult***REMOVED******REMOVED***

	protoPaths := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	// TODO: this is O(n^2) but could likely be O(n) with a clever data structure (prefix tree that is indexed backwards?)
	importCandidates := map[string]map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	candidatesAvailable := map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
	for name := range protos ***REMOVED***
		candidatesAvailable[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
		for _, f := range protos ***REMOVED***
			for _, imp := range f.fd.Dependency ***REMOVED***
				if strings.HasSuffix(name, imp) ***REMOVED***
					candidates := importCandidates[imp]
					if candidates == nil ***REMOVED***
						candidates = map[string]struct***REMOVED******REMOVED******REMOVED******REMOVED***
						importCandidates[imp] = candidates
					***REMOVED***
					candidates[name] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
	for imp, candidates := range importCandidates ***REMOVED***
		// if we found multiple possible candidates, use the one that is an exact match
		// if it exists, and otherwise, guess that it's the shortest path (fewest elements)
		var best string
		for c := range candidates ***REMOVED***
			if _, ok := candidatesAvailable[c]; !ok ***REMOVED***
				// already used this candidate and re-written its filename accordingly
				continue
			***REMOVED***
			if c == imp ***REMOVED***
				// exact match!
				best = c
				break
			***REMOVED***
			if best == "" ***REMOVED***
				best = c
			***REMOVED*** else ***REMOVED***
				// HACK: we can't actually tell which files is supposed to match
				// this import, so arbitrarily pick the "shorter" one (fewest
				// path elements) or, on a tie, the lexically earlier one
				minLen := strings.Count(best, string(filepath.Separator))
				cLen := strings.Count(c, string(filepath.Separator))
				if cLen < minLen || (cLen == minLen && c < best) ***REMOVED***
					best = c
				***REMOVED***
			***REMOVED***
		***REMOVED***
		if best != "" ***REMOVED***
			prefix := best[:len(best)-len(imp)]
			if len(prefix) > 0 ***REMOVED***
				protoPaths[prefix] = struct***REMOVED******REMOVED******REMOVED******REMOVED***
			***REMOVED***
			f := protos[best]
			f.fd.Name = proto.String(imp)
			revisedProtos[imp] = f
			delete(candidatesAvailable, best)
		***REMOVED***
	***REMOVED***

	if len(candidatesAvailable) == 0 ***REMOVED***
		return revisedProtos
	***REMOVED***

	if len(protoPaths) == 0 ***REMOVED***
		for c := range candidatesAvailable ***REMOVED***
			revisedProtos[c] = protos[c]
		***REMOVED***
		return revisedProtos
	***REMOVED***

	// Any remaining candidates are entry-points (not imported by others), so
	// the best bet to "fixing" their file name is to see if they're in one of
	// the proto paths we found, and if so strip that prefix.
	protoPathStrs := make([]string, len(protoPaths))
	i := 0
	for p := range protoPaths ***REMOVED***
		protoPathStrs[i] = p
		i++
	***REMOVED***
	sort.Strings(protoPathStrs)
	// we look at paths in reverse order, so we'll use a longer proto path if
	// there is more than one match
	for c := range candidatesAvailable ***REMOVED***
		var imp string
		for i := len(protoPathStrs) - 1; i >= 0; i-- ***REMOVED***
			p := protoPathStrs[i]
			if strings.HasPrefix(c, p) ***REMOVED***
				imp = c[len(p):]
				break
			***REMOVED***
		***REMOVED***
		if imp != "" ***REMOVED***
			f := protos[c]
			f.fd.Name = proto.String(imp)
			revisedProtos[imp] = f
		***REMOVED*** else ***REMOVED***
			revisedProtos[c] = protos[c]
		***REMOVED***
	***REMOVED***

	return revisedProtos
***REMOVED***

func parseProtoFiles(acc FileAccessor, filenames []string, errs *errorHandler, parsed *parseResults, lookupImport func(string) (*dpb.FileDescriptorProto, error)) ***REMOVED***
	for _, name := range filenames ***REMOVED***
		parseProtoFile(acc, name, nil, errs, parsed, lookupImport)
		if errs.err != nil ***REMOVED***
			return
		***REMOVED***
	***REMOVED***
***REMOVED***

func parseProtoFile(acc FileAccessor, filename string, importLoc *SourcePos, errs *errorHandler, results *parseResults, lookupImport func(string) (*dpb.FileDescriptorProto, error)) ***REMOVED***
	if results.has(filename) ***REMOVED***
		return
	***REMOVED***
	if lookupImport == nil ***REMOVED***
		lookupImport = func(string) (*dpb.FileDescriptorProto, error) ***REMOVED***
			return nil, errors.New("no import lookup function")
		***REMOVED***
	***REMOVED***
	in, err := acc(filename)
	var result *parseResult
	if err == nil ***REMOVED***
		// try to parse the bytes accessed
		func() ***REMOVED***
			defer func() ***REMOVED***
				// if we've already parsed contents, an error
				// closing need not fail this operation
				_ = in.Close()
			***REMOVED***()
			result = parseProto(filename, in, errs, results.validate, results.createDescriptorProtos)
		***REMOVED***()
	***REMOVED*** else if d, lookupErr := lookupImport(filename); lookupErr == nil ***REMOVED***
		// This is a user-provided descriptor, which is acting similarly to a
		// well-known import.
		result = &parseResult***REMOVED***fd: proto.Clone(d).(*dpb.FileDescriptorProto)***REMOVED***
	***REMOVED*** else if d, ok := standardImports[filename]; ok ***REMOVED***
		// it's a well-known import
		// (we clone it to make sure we're not sharing state with other
		//  parsers, which could result in unsafe races if multiple
		//  parsers are trying to access it concurrently)
		result = &parseResult***REMOVED***fd: proto.Clone(d).(*dpb.FileDescriptorProto)***REMOVED***
	***REMOVED*** else ***REMOVED***
		if !strings.Contains(err.Error(), filename) ***REMOVED***
			// an error message that doesn't indicate the file is awful!
			// this cannot be %w as this is not compatible with go <= 1.13
			err = errorWithFilename***REMOVED***
				underlying: err,
				filename:   filename,
			***REMOVED***
		***REMOVED***
		// The top-level loop in parseProtoFiles calls this with nil for the top-level files
		// importLoc is only for imports, otherwise we do not want to return a ErrorWithSourcePos
		// ErrorWithSourcePos should always have a non-nil SourcePos
		if importLoc != nil ***REMOVED***
			// associate the error with the import line
			err = ErrorWithSourcePos***REMOVED***
				Pos:        importLoc,
				Underlying: err,
			***REMOVED***
		***REMOVED***
		_ = errs.handleError(err)
		return
	***REMOVED***

	results.add(filename, result)

	if errs.err != nil ***REMOVED***
		return // abort
	***REMOVED***

	if results.recursive ***REMOVED***
		fd := result.fd
		decl := result.getFileNode(fd)
		fnode, ok := decl.(*ast.FileNode)
		if !ok ***REMOVED***
			// no AST for this file? use imports in descriptor
			for _, dep := range fd.Dependency ***REMOVED***
				parseProtoFile(acc, dep, decl.Start(), errs, results, lookupImport)
				if errs.getError() != nil ***REMOVED***
					return // abort
				***REMOVED***
			***REMOVED***
			return
		***REMOVED***
		// we have an AST; use it so we can report import location in errors
		for _, decl := range fnode.Decls ***REMOVED***
			if dep, ok := decl.(*ast.ImportNode); ok ***REMOVED***
				parseProtoFile(acc, dep.Name.AsString(), dep.Name.Start(), errs, results, lookupImport)
				if errs.getError() != nil ***REMOVED***
					return // abort
				***REMOVED***
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***

type parseResults struct ***REMOVED***
	resultsByFilename map[string]*parseResult
	filenames         []string

	recursive, validate, createDescriptorProtos bool
***REMOVED***

func (r *parseResults) has(filename string) bool ***REMOVED***
	_, ok := r.resultsByFilename[filename]
	return ok
***REMOVED***

func (r *parseResults) add(filename string, result *parseResult) ***REMOVED***
	r.resultsByFilename[filename] = result
	r.filenames = append(r.filenames, filename)
***REMOVED***

type parseResult struct ***REMOVED***
	// handles any errors encountered during parsing, construction of file descriptor,
	// or validation
	errs *errorHandler

	// the root of the AST
	root *ast.FileNode
	// the parsed file descriptor
	fd *dpb.FileDescriptorProto

	// if set to true, enables lenient interpretation of options, where
	// unrecognized options will be left uninterpreted instead of resulting in a
	// link error
	lenient bool

	// a map of elements in the descriptor to nodes in the AST
	// (for extracting position information when validating the descriptor)
	nodes map[proto.Message]ast.Node

	// a map of uninterpreted option AST nodes to their relative path
	// in the resulting options message
	interpretedOptions map[*ast.OptionNode][]int32
***REMOVED***

func (r *parseResult) getFileNode(f *dpb.FileDescriptorProto) ast.FileDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(f.GetName())
	***REMOVED***
	return r.nodes[f].(ast.FileDeclNode)
***REMOVED***

func (r *parseResult) getOptionNode(o *dpb.UninterpretedOption) ast.OptionDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[o].(ast.OptionDeclNode)
***REMOVED***

func (r *parseResult) getOptionNamePartNode(o *dpb.UninterpretedOption_NamePart) ast.Node ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[o]
***REMOVED***

func (r *parseResult) getFieldNode(f *dpb.FieldDescriptorProto) ast.FieldDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[f].(ast.FieldDeclNode)
***REMOVED***

func (r *parseResult) getExtensionRangeNode(e *dpb.DescriptorProto_ExtensionRange) ast.RangeDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[e].(ast.RangeDeclNode)
***REMOVED***

func (r *parseResult) getMessageReservedRangeNode(rr *dpb.DescriptorProto_ReservedRange) ast.RangeDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[rr].(ast.RangeDeclNode)
***REMOVED***

func (r *parseResult) getEnumNode(e *dpb.EnumDescriptorProto) ast.Node ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[e]
***REMOVED***

func (r *parseResult) getEnumValueNode(e *dpb.EnumValueDescriptorProto) ast.EnumValueDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[e].(ast.EnumValueDeclNode)
***REMOVED***

func (r *parseResult) getEnumReservedRangeNode(rr *dpb.EnumDescriptorProto_EnumReservedRange) ast.RangeDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[rr].(ast.RangeDeclNode)
***REMOVED***

func (r *parseResult) getMethodNode(m *dpb.MethodDescriptorProto) ast.RPCDeclNode ***REMOVED***
	if r.nodes == nil ***REMOVED***
		return ast.NewNoSourceNode(r.fd.GetName())
	***REMOVED***
	return r.nodes[m].(ast.RPCDeclNode)
***REMOVED***

func (r *parseResult) putFileNode(f *dpb.FileDescriptorProto, n *ast.FileNode) ***REMOVED***
	r.nodes[f] = n
***REMOVED***

func (r *parseResult) putOptionNode(o *dpb.UninterpretedOption, n *ast.OptionNode) ***REMOVED***
	r.nodes[o] = n
***REMOVED***

func (r *parseResult) putOptionNamePartNode(o *dpb.UninterpretedOption_NamePart, n *ast.FieldReferenceNode) ***REMOVED***
	r.nodes[o] = n
***REMOVED***

func (r *parseResult) putMessageNode(m *dpb.DescriptorProto, n ast.MessageDeclNode) ***REMOVED***
	r.nodes[m] = n
***REMOVED***

func (r *parseResult) putFieldNode(f *dpb.FieldDescriptorProto, n ast.FieldDeclNode) ***REMOVED***
	r.nodes[f] = n
***REMOVED***

func (r *parseResult) putOneOfNode(o *dpb.OneofDescriptorProto, n *ast.OneOfNode) ***REMOVED***
	r.nodes[o] = n
***REMOVED***

func (r *parseResult) putExtensionRangeNode(e *dpb.DescriptorProto_ExtensionRange, n *ast.RangeNode) ***REMOVED***
	r.nodes[e] = n
***REMOVED***

func (r *parseResult) putMessageReservedRangeNode(rr *dpb.DescriptorProto_ReservedRange, n *ast.RangeNode) ***REMOVED***
	r.nodes[rr] = n
***REMOVED***

func (r *parseResult) putEnumNode(e *dpb.EnumDescriptorProto, n *ast.EnumNode) ***REMOVED***
	r.nodes[e] = n
***REMOVED***

func (r *parseResult) putEnumValueNode(e *dpb.EnumValueDescriptorProto, n *ast.EnumValueNode) ***REMOVED***
	r.nodes[e] = n
***REMOVED***

func (r *parseResult) putEnumReservedRangeNode(rr *dpb.EnumDescriptorProto_EnumReservedRange, n *ast.RangeNode) ***REMOVED***
	r.nodes[rr] = n
***REMOVED***

func (r *parseResult) putServiceNode(s *dpb.ServiceDescriptorProto, n *ast.ServiceNode) ***REMOVED***
	r.nodes[s] = n
***REMOVED***

func (r *parseResult) putMethodNode(m *dpb.MethodDescriptorProto, n *ast.RPCNode) ***REMOVED***
	r.nodes[m] = n
***REMOVED***

func parseProto(filename string, r io.Reader, errs *errorHandler, validate, createProtos bool) *parseResult ***REMOVED***
	beforeErrs := errs.errsReported
	lx := newLexer(r, filename, errs)
	protoParse(lx)
	if lx.res == nil || len(lx.res.Children()) == 0 ***REMOVED***
		// nil AST means there was an error that prevented any parsing
		// or the file was empty; synthesize empty non-nil AST
		lx.res = ast.NewEmptyFileNode(filename)
	***REMOVED***
	if lx.eof != nil ***REMOVED***
		lx.res.FinalComments = lx.eof.LeadingComments()
		lx.res.FinalWhitespace = lx.eof.LeadingWhitespace()
	***REMOVED***
	res := createParseResult(filename, lx.res, errs, createProtos)
	if validate && errs.err == nil ***REMOVED***
		validateBasic(res, errs.errsReported > beforeErrs)
	***REMOVED***

	return res
***REMOVED***

func createParseResult(filename string, file *ast.FileNode, errs *errorHandler, createProtos bool) *parseResult ***REMOVED***
	res := &parseResult***REMOVED***
		errs:               errs,
		root:               file,
		nodes:              map[proto.Message]ast.Node***REMOVED******REMOVED***,
		interpretedOptions: map[*ast.OptionNode][]int32***REMOVED******REMOVED***,
	***REMOVED***
	if createProtos ***REMOVED***
		res.createFileDescriptor(filename, file)
	***REMOVED***
	return res
***REMOVED***

func checkTag(pos *SourcePos, v uint64, maxTag int32) error ***REMOVED***
	if v < 1 ***REMOVED***
		return errorWithPos(pos, "tag number %d must be greater than zero", v)
	***REMOVED*** else if v > uint64(maxTag) ***REMOVED***
		return errorWithPos(pos, "tag number %d is higher than max allowed tag number (%d)", v, maxTag)
	***REMOVED*** else if v >= internal.SpecialReservedStart && v <= internal.SpecialReservedEnd ***REMOVED***
		return errorWithPos(pos, "tag number %d is in disallowed reserved range %d-%d", v, internal.SpecialReservedStart, internal.SpecialReservedEnd)
	***REMOVED***
	return nil
***REMOVED***

func checkExtensionsInFile(fd *desc.FileDescriptor, res *parseResult) error ***REMOVED***
	for _, fld := range fd.GetExtensions() ***REMOVED***
		if err := checkExtension(fld, res); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, md := range fd.GetMessageTypes() ***REMOVED***
		if err := checkExtensionsInMessage(md, res); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func checkExtensionsInMessage(md *desc.MessageDescriptor, res *parseResult) error ***REMOVED***
	for _, fld := range md.GetNestedExtensions() ***REMOVED***
		if err := checkExtension(fld, res); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	for _, nmd := range md.GetNestedMessageTypes() ***REMOVED***
		if err := checkExtensionsInMessage(nmd, res); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

func checkExtension(fld *desc.FieldDescriptor, res *parseResult) error ***REMOVED***
	// NB: It's a little gross that we don't enforce these in validateBasic().
	// But requires some minimal linking to resolve the extendee, so we can
	// interrogate its descriptor.
	if fld.GetOwner().GetMessageOptions().GetMessageSetWireFormat() ***REMOVED***
		// Message set wire format requires that all extensions be messages
		// themselves (no scalar extensions)
		if fld.GetType() != dpb.FieldDescriptorProto_TYPE_MESSAGE ***REMOVED***
			pos := res.getFieldNode(fld.AsFieldDescriptorProto()).FieldType().Start()
			return errorWithPos(pos, "messages with message-set wire format cannot contain scalar extensions, only messages")
		***REMOVED***
	***REMOVED*** else ***REMOVED***
		// In validateBasic() we just made sure these were within bounds for any message. But
		// now that things are linked, we can check if the extendee is messageset wire format
		// and, if not, enforce tighter limit.
		if fld.GetNumber() > internal.MaxNormalTag ***REMOVED***
			pos := res.getFieldNode(fld.AsFieldDescriptorProto()).FieldTag().Start()
			return errorWithPos(pos, "tag number %d is higher than max allowed tag number (%d)", fld.GetNumber(), internal.MaxNormalTag)
		***REMOVED***
	***REMOVED***

	return nil
***REMOVED***

func aggToString(agg []*ast.MessageFieldNode, buf *bytes.Buffer) ***REMOVED***
	buf.WriteString("***REMOVED***")
	for _, a := range agg ***REMOVED***
		buf.WriteString(" ")
		buf.WriteString(a.Name.Value())
		if v, ok := a.Val.(*ast.MessageLiteralNode); ok ***REMOVED***
			aggToString(v.Elements, buf)
		***REMOVED*** else ***REMOVED***
			buf.WriteString(": ")
			elementToString(a.Val.Value(), buf)
		***REMOVED***
	***REMOVED***
	buf.WriteString(" ***REMOVED***")
***REMOVED***

func elementToString(v interface***REMOVED******REMOVED***, buf *bytes.Buffer) ***REMOVED***
	switch v := v.(type) ***REMOVED***
	case bool, int64, uint64, ast.Identifier:
		_, _ = fmt.Fprintf(buf, "%v", v)
	case float64:
		if math.IsInf(v, 1) ***REMOVED***
			buf.WriteString(": inf")
		***REMOVED*** else if math.IsInf(v, -1) ***REMOVED***
			buf.WriteString(": -inf")
		***REMOVED*** else if math.IsNaN(v) ***REMOVED***
			buf.WriteString(": nan")
		***REMOVED*** else ***REMOVED***
			_, _ = fmt.Fprintf(buf, ": %v", v)
		***REMOVED***
	case string:
		buf.WriteRune('"')
		writeEscapedBytes(buf, []byte(v))
		buf.WriteRune('"')
	case []ast.ValueNode:
		buf.WriteString(": [")
		first := true
		for _, e := range v ***REMOVED***
			if first ***REMOVED***
				first = false
			***REMOVED*** else ***REMOVED***
				buf.WriteString(", ")
			***REMOVED***
			elementToString(e.Value(), buf)
		***REMOVED***
		buf.WriteString("]")
	case []*ast.MessageFieldNode:
		aggToString(v, buf)
	***REMOVED***
***REMOVED***

func writeEscapedBytes(buf *bytes.Buffer, b []byte) ***REMOVED***
	for _, c := range b ***REMOVED***
		switch c ***REMOVED***
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		case '"':
			buf.WriteString("\\\"")
		case '\'':
			buf.WriteString("\\'")
		case '\\':
			buf.WriteString("\\\\")
		default:
			if c >= 0x20 && c <= 0x7f && c != '"' && c != '\\' ***REMOVED***
				// simple printable characters
				buf.WriteByte(c)
			***REMOVED*** else ***REMOVED***
				// use octal escape for all other values
				buf.WriteRune('\\')
				buf.WriteByte('0' + ((c >> 6) & 0x7))
				buf.WriteByte('0' + ((c >> 3) & 0x7))
				buf.WriteByte('0' + (c & 0x7))
			***REMOVED***
		***REMOVED***
	***REMOVED***
***REMOVED***
