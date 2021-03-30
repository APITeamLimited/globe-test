package ast

import "io"

// Print prints the given AST node to the given output. This operation
// basically walks the AST and, for each TerminalNode, prints the node's
// leading comments, leading whitespace, the node's raw text, and then
// any trailing comments. If the given node is a *FileNode, it will then
// also print the file's FinalComments and FinalWhitespace.
func Print(w io.Writer, node Node) error ***REMOVED***
	sw, ok := w.(stringWriter)
	if !ok ***REMOVED***
		sw = &strWriter***REMOVED***w***REMOVED***
	***REMOVED***
	var err error
	Walk(node, func(n Node) (bool, VisitFunc) ***REMOVED***
		if err != nil ***REMOVED***
			return false, nil
		***REMOVED***
		token, ok := n.(TerminalNode)
		if !ok ***REMOVED***
			return true, nil
		***REMOVED***

		err = printComments(sw, token.LeadingComments())
		if err != nil ***REMOVED***
			return false, nil
		***REMOVED***

		_, err = sw.WriteString(token.LeadingWhitespace())
		if err != nil ***REMOVED***
			return false, nil
		***REMOVED***

		_, err = sw.WriteString(token.RawText())
		if err != nil ***REMOVED***
			return false, nil
		***REMOVED***

		err = printComments(sw, token.TrailingComments())
		return false, nil
	***REMOVED***)
	if err != nil ***REMOVED***
		return err
	***REMOVED***

	if file, ok := node.(*FileNode); ok ***REMOVED***
		err = printComments(sw, file.FinalComments)
		if err != nil ***REMOVED***
			return err
		***REMOVED***
		_, err = sw.WriteString(file.FinalWhitespace)
		return err
	***REMOVED***

	return nil
***REMOVED***

func printComments(sw stringWriter, comments []Comment) error ***REMOVED***
	for _, comment := range comments ***REMOVED***
		if _, err := sw.WriteString(comment.LeadingWhitespace); err != nil ***REMOVED***
			return err
		***REMOVED***
		if _, err := sw.WriteString(comment.Text); err != nil ***REMOVED***
			return err
		***REMOVED***
	***REMOVED***
	return nil
***REMOVED***

// many io.Writer impls also provide a string-based method
type stringWriter interface ***REMOVED***
	WriteString(s string) (n int, err error)
***REMOVED***

// adapter, in case the given writer does NOT provide a string-based method
type strWriter struct ***REMOVED***
	io.Writer
***REMOVED***

func (s *strWriter) WriteString(str string) (int, error) ***REMOVED***
	if str == "" ***REMOVED***
		return 0, nil
	***REMOVED***
	return s.Write([]byte(str))
***REMOVED***
