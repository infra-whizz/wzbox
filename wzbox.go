package wzbox

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
)

// WzBox object
type WzBox struct {
	packageName    string
	structName     string
	outputFilename string
	compressed     bool
	files          map[string][]byte
}

// NewWzBox creates an instance of the WzBox
func NewWzBox() *WzBox {
	wzb := new(WzBox)
	wzb.files = make(map[string][]byte)
	return wzb
}

// SetCompression on or off
func (wzb *WzBox) SetCompression(c bool) *WzBox {
	wzb.compressed = c
	return wzb
}

// AddFile to the box processor
func (wzb *WzBox) AddFile(path string) *WzBox {
	wzb.files[path] = nil
	return wzb
}

// SetOutputFilename of the generated go code
func (wzb *WzBox) SetOutputFilename(name string) *WzBox {
	wzb.outputFilename = name
	return wzb
}

// SetStructName inside the go code to be reused later
func (wzb *WzBox) SetStructName(name string) *WzBox {
	wzb.structName = name
	return wzb
}

// SetPackageName inside the go code to be reused later
func (wzb *WzBox) SetPackageName(name string) *WzBox {
	wzb.packageName = name
	return wzb
}

// Generate code
func (wzb *WzBox) Generate() (string, error) {
	for fname := range wzb.files {
		data, err := ioutil.ReadFile(fname)
		if err != nil {
			return "", err
		}

		// Compress?
		if wzb.compressed {
			var buff bytes.Buffer
			gz := gzip.NewWriter(&buff)
			if _, err := gz.Write(data); err != nil {
				return "", err
			}

			if err := gz.Flush(); err != nil {
				return "", err
			}

			if err := gz.Close(); err != nil {
				return "", err
			}

			data = buff.Bytes()
		}

		wzb.files[fname] = data
	}

	return wzb.createSources(), nil
}

func (wzb *WzBox) createSources() string {
	var out bytes.Buffer

	out.WriteString(fmt.Sprintf("package %s\n\n", wzb.packageName))
	out.WriteString("import (\n\t\"bytes\"\n\t\"compress/gzip\"\n\t\"io/ioutil\"\n)\n\n")
	out.WriteString(fmt.Sprintf("// %s data container\ntype %s struct {\n\tcompressed bool\n\tdata       map[string][]byte\n}\n",
		wzb.structName, wzb.structName))
	out.WriteString(fmt.Sprintf("\n// New%s creates a new instance of %s\nfunc New%s() *%s {\n\tcnt := new(%s)\n",
		wzb.structName, wzb.structName, wzb.structName, wzb.structName, wzb.structName))

	if wzb.compressed {
		out.WriteString("\n\t// NOTE: This data is ZIP-compressed\n")
		out.WriteString("\tcnt.compressed = true\n")
	}

	out.WriteString("\tcnt.data = map[string][]byte{\n")
	for fname, fdata := range wzb.files {
		var arr bytes.Buffer
		l := len(fdata)
		w := 0
		for idx, c := range fdata {
			arr.WriteString(fmt.Sprintf("%d", int(c)))
			if idx < l {
				arr.WriteString(",")
			}
			w++
			if w > 10 {
				w = 0
				arr.WriteString("\n\t\t\t")
			} else if idx+1 < l {
				arr.WriteString(" ") // space after previous comma
			}
		}
		out.WriteString(fmt.Sprintf("\t\t\"%s\": {\n\t\t\t%s\n\t\t},", fname, arr.String()))
	}

	out.WriteString("\n\t}\n\treturn cnt\n}")
	out.WriteString("\n\n// Get file content\n")
	out.WriteString(fmt.Sprintf("func (cnt *%s) Get(name string) []byte {\n\tcontent := cnt.data[name]\n", wzb.structName))
	out.WriteString("\tif cnt.compressed {\n\t\tgz, _ := gzip.NewReader(bytes.NewReader(content))\n")
	out.WriteString("\t\tdefer gz.Close()\n\t\tcontent, _ = ioutil.ReadAll(gz)\n\t}\n\n\treturn content\n}\n")

	return out.String()
}
