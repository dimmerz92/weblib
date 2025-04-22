package weblib

import (
	"fmt"
	"mime"
	"mime/multipart"
	"path/filepath"
	"strings"
)

type File struct {
	Header    multipart.FileHeader
	Path      string
	Directory string
	Filename  string
}

// ExtractFullPath extracts the full relative path of an individual file submitted by a multipart form containing a
// file input with the directory attributes enabled.
//
// E.g.,
// <input type="file" webkitdirectory directory/>
//
// The output is sanitised to prevent path traversal and invalid characters.
// Note: Some browsers may sanitise the filename to include only the base filename (e.g., "file.txt").
// The Directory field may be empty or "." in such cases. Test with target browsers to confirm behaviour.
func ExtractFullPath(fileheader *multipart.FileHeader) (*File, error) {
	// parse Content-Disposition header
	_, params, err := mime.ParseMediaType(fileheader.Header.Get("Content-Disposition"))
	if err != nil {
		return nil, fmt.Errorf("invalid Content-Disposition header: %w", err)
	}

	// extract filename parameter
	filename, ok := params["filename"]
	if !ok {
		return nil, fmt.Errorf("filename not found in Content-Disposition header")
	}

	// trim quotes from filename
	separator := string(filepath.Separator)
	trimmed := TrimQuotes(filename)
	if trimmed == "" || trimmed == "." || trimmed == separator {
		return nil, fmt.Errorf("filename is empty after trimming quotes")
	}

	// validate filename is not empty
	if strings.HasSuffix(trimmed, separator+".") || strings.HasSuffix(trimmed, separator) {
		return nil, fmt.Errorf("invalid filename: empty or invalid base filename")
	}

	// sanitise the path and prevent path traversal attack
	cleanPath := filepath.Clean(trimmed)
	if strings.Contains(cleanPath, "..") || filepath.IsAbs(cleanPath) {
		return nil, fmt.Errorf("invalid filename: contains path traversal or absolute path")
	}

	// split the path into directory and filename components
	dir := filepath.Dir(cleanPath)
	fname := filepath.Base(cleanPath)

	return &File{
		Header:    *fileheader,
		Path:      cleanPath,
		Directory: dir,
		Filename:  fname,
	}, nil
}
