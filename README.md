# 1. wzbox

Although Go 1.16 allows you to embed static files inside the binary,
this is a backward-compat version to work with any older Go versions
until major Linux distributions will move to Go 1.16 version.

# 2. Usage

To use it straight-forward:

- Compile CLI command itself
- Specify your files as a list
- Specify struct name to access that later
- Specify the output file
- Compressed (zip)?

It will grab whatever you specified, compress (if you want) and
generate a byte array and generate a Go source with a class, which
is just ready to use as is.

# 3. Example

To embed into your binary "test.html" do the following:

	cd cmd
	go build
	./wzbox -f /path/to/test.html -s MyStruct --compress > foo.go

Then use this (edit package name to your project):

	mystruct := NewMyStruct()
	s := mystruct.Get("test.html")

It will automatically uncompress you your content.
Wasn't that hard, right? :-)
