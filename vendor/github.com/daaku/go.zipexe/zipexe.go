// Package zipexe attempts to open an executable binary file as a zip file.
package zipexe

import (
	"archive/zip"
	"debug/elf"
	"debug/macho"
	"debug/pe"
	"errors"
	"io"
	"os"
)

// Opens a zip file by path.
func Open(path string) (*zip.Reader, error) ***REMOVED***
	_, rd, err := OpenCloser(path)
	return rd, err
***REMOVED***

// OpenCloser is like Open but returns an additional Closer to avoid leaking open files.
func OpenCloser(path string) (io.Closer, *zip.Reader, error) ***REMOVED***
	file, err := os.Open(path)
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	finfo, err := file.Stat()
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	zr, err := NewReader(file, finfo.Size())
	if err != nil ***REMOVED***
		return nil, nil, err
	***REMOVED***
	return file, zr, nil
***REMOVED***

// Open a zip file, specially handling various binaries that may have been
// augmented with zip data.
func NewReader(rda io.ReaderAt, size int64) (*zip.Reader, error) ***REMOVED***
	handlers := []func(io.ReaderAt, int64) (*zip.Reader, error)***REMOVED***
		zip.NewReader,
		zipExeReaderMacho,
		zipExeReaderElf,
		zipExeReaderPe,
	***REMOVED***

	for _, handler := range handlers ***REMOVED***
		zfile, err := handler(rda, size)
		if err == nil ***REMOVED***
			return zfile, nil
		***REMOVED***
	***REMOVED***
	return nil, errors.New("Couldn't Open As Executable")
***REMOVED***

// zipExeReaderMacho treats the file as a Mach-O binary
// (Mac OS X / Darwin executable) and attempts to find a zip archive.
func zipExeReaderMacho(rda io.ReaderAt, size int64) (*zip.Reader, error) ***REMOVED***
	file, err := macho.NewFile(rda)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var max int64
	for _, load := range file.Loads ***REMOVED***
		seg, ok := load.(*macho.Segment)
		if ok ***REMOVED***
			// Check if the segment contains a zip file
			if zfile, err := zip.NewReader(seg, int64(seg.Filesz)); err == nil ***REMOVED***
				return zfile, nil
			***REMOVED***

			// Otherwise move end of file pointer
			end := int64(seg.Offset + seg.Filesz)
			if end > max ***REMOVED***
				max = end
			***REMOVED***
		***REMOVED***
	***REMOVED***

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(rda, max, size-max)
	return zip.NewReader(section, section.Size())
***REMOVED***

// zipExeReaderPe treats the file as a Portable Exectuable binary
// (Windows executable) and attempts to find a zip archive.
func zipExeReaderPe(rda io.ReaderAt, size int64) (*zip.Reader, error) ***REMOVED***
	file, err := pe.NewFile(rda)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var max int64
	for _, sec := range file.Sections ***REMOVED***
		// Check if this section has a zip file
		if zfile, err := zip.NewReader(sec, int64(sec.Size)); err == nil ***REMOVED***
			return zfile, nil
		***REMOVED***

		// Otherwise move end of file pointer
		end := int64(sec.Offset + sec.Size)
		if end > max ***REMOVED***
			max = end
		***REMOVED***
	***REMOVED***

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(rda, max, size-max)
	return zip.NewReader(section, section.Size())
***REMOVED***

// zipExeReaderElf treats the file as a ELF binary
// (linux/BSD/etc... executable) and attempts to find a zip archive.
func zipExeReaderElf(rda io.ReaderAt, size int64) (*zip.Reader, error) ***REMOVED***
	file, err := elf.NewFile(rda)
	if err != nil ***REMOVED***
		return nil, err
	***REMOVED***

	var max int64
	for _, sect := range file.Sections ***REMOVED***
		if sect.Type == elf.SHT_NOBITS ***REMOVED***
			continue
		***REMOVED***

		// Check if this section has a zip file
		if zfile, err := zip.NewReader(sect, int64(sect.Size)); err == nil ***REMOVED***
			return zfile, nil
		***REMOVED***

		// Otherwise move end of file pointer
		end := int64(sect.Offset + sect.Size)
		if end > max ***REMOVED***
			max = end
		***REMOVED***
	***REMOVED***

	// No zip file within binary, try appended to end
	section := io.NewSectionReader(rda, max, size-max)
	return zip.NewReader(section, section.Size())
***REMOVED***
