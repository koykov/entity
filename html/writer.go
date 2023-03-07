package html

import "io"

type Writer interface {
	io.Writer
	io.StringWriter
	io.ByteWriter
}
