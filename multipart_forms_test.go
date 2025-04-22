package weblib

import (
	"mime/multipart"
	"net/textproto"
	"reflect"
	"strings"
	"testing"
)

func TestExtractFullPath(t *testing.T) {
	tests := []struct {
		name        string
		header      *multipart.FileHeader
		expected    *File
		errContains string
	}{
		{
			name: "valid file with directory",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="folder/subfolder/file.txt"`},
				},
			},
			expected: &File{
				Header: multipart.FileHeader{
					Header: textproto.MIMEHeader{
						"Content-Disposition": []string{`form-data; name="file"; filename="folder/subfolder/file.txt"`},
					},
				},
				Path:      "folder/subfolder/file.txt",
				Directory: "folder/subfolder",
				Filename:  "file.txt",
			},
			errContains: "",
		},
		{
			name: "valid file without directory",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
				},
			},
			expected: &File{
				Header: multipart.FileHeader{
					Header: textproto.MIMEHeader{
						"Content-Disposition": []string{`form-data; name="file"; filename="file.txt"`},
					},
				},
				Path:      "file.txt",
				Directory: ".",
				Filename:  "file.txt",
			},
			errContains: "",
		},
		{
			name: "filename with quotes",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="\"file.txt\""`},
				},
			},
			expected: &File{
				Header: multipart.FileHeader{
					Header: textproto.MIMEHeader{
						"Content-Disposition": []string{`form-data; name="file"; filename="\"file.txt\""`},
					},
				},
				Path:      "file.txt",
				Directory: ".",
				Filename:  "file.txt",
			},
			errContains: "",
		},
		{
			name: "invalid Content-Disposition",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{``},
				},
			},
			expected:    nil,
			errContains: "invalid Content-Disposition header",
		},
		{
			name: "missing filename parameter",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"`},
				},
			},
			expected:    nil,
			errContains: "filename not found in Content-Disposition header",
		},
		{
			name: "empty filename after trimming",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="\"\""`},
				},
			},
			expected:    nil,
			errContains: "filename is empty after trimming quotes",
		},
		{
			name: "path traversal attempt",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="../file.txt"`},
				},
			},
			expected:    nil,
			errContains: "invalid filename: contains path traversal or absolute path",
		},
		{
			name: "absolute path",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="/folder/file.txt"`},
				},
			},
			expected:    nil,
			errContains: "invalid filename: contains path traversal or absolute path",
		},
		{
			name: "invalid base filename (dot)",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="folder/./."`},
				},
			},
			expected:    nil,
			errContains: "invalid filename: empty or invalid base filename",
		},
		{
			name: "invalid base filename (slash)",
			header: &multipart.FileHeader{
				Header: textproto.MIMEHeader{
					"Content-Disposition": []string{`form-data; name="file"; filename="folder/"`},
				},
			},
			expected:    nil,
			errContains: "invalid filename: empty or invalid base filename",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExtractFullPath(tt.header)

			if tt.errContains != "" {
				if err == nil {
					t.Errorf("expected error containing %q, got none", tt.errContains)
					return
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				} else if result != nil {
					t.Errorf("expected nil result, got %+v", result)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
