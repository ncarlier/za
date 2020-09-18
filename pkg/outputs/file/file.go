package file

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ncarlier/trackr/pkg/events"
	"github.com/ncarlier/trackr/pkg/logger"
	"github.com/ncarlier/trackr/pkg/outputs"
	"github.com/ncarlier/trackr/pkg/serializers"
)

// File output
type File struct {
	Files []string `toml:"files"`

	writer     io.Writer
	closers    []io.Closer
	serializer serializers.Serializer
}

var sampleConfig = `
  ## Files to write to, "stdout" is a specially handled file.
  files = ["stdout", "/tmp/access.log"]
  ## Data format to output.
  data_format = "json"
`

// SetSerializer set data serializer
func (f *File) SetSerializer(serializer serializers.Serializer) {
	f.serializer = serializer
}

// Connect activate the output writer
func (f *File) Connect() error {
	writers := []io.Writer{}

	if len(f.Files) == 0 {
		f.Files = []string{"stdout"}
	}

	for _, file := range f.Files {
		if file == "stdout" {
			writers = append(writers, os.Stdout)
		} else {
			fd, err := os.Create(file)
			if err != nil {
				return err
			}

			of := bufio.NewWriter(fd)

			writers = append(writers, of)
			f.closers = append(f.closers, fd)
		}
	}
	f.writer = io.MultiWriter(writers...)
	logger.Debug.Printf("using FILE output: %s\n", f.Files)
	return nil
}

// Close the output writer
func (f *File) Close() error {
	var err error
	for _, c := range f.closers {
		errClose := c.Close()
		if errClose != nil {
			err = errClose
		}
	}
	return err
}

// SampleConfig returns sample configuration
func (f *File) SampleConfig() string {
	return sampleConfig
}

// Description returns description
func (f *File) Description() string {
	return "Send page view to file(s)"
}

// SendEvent send event to the Output
func (f *File) SendEvent(event events.Event) error {
	b, err := f.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("unable to serialize page view: %v", err)
	}

	if _, err = f.writer.Write(b); err != nil {
		return fmt.Errorf("unable to write page view to file output: %v", err)
	}

	return nil
}

func init() {
	outputs.Add("file", func() outputs.Output {
		return &File{}
	})
}
