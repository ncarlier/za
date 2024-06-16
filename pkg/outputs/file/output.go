package file

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/ncarlier/za/pkg/conditional"
	"github.com/ncarlier/za/pkg/events"
	"github.com/ncarlier/za/pkg/outputs"
	"github.com/ncarlier/za/pkg/serializers"
)

// Output for file writing
type Output struct {
	Files []string `toml:"files"`

	writer     io.Writer
	closers    []io.Closer
	serializer serializers.Serializer
	condition  conditional.Expression
}

// SetSerializer set data serializer
func (o *Output) SetSerializer(serializer serializers.Serializer) {
	o.serializer = serializer
}

// SetCondition set condition expression
func (o *Output) SetCondition(condition conditional.Expression) {
	o.condition = condition
}

// Connect activate the output writer
func (o *Output) Connect() error {
	writers := []io.Writer{}

	if len(o.Files) == 0 {
		o.Files = []string{"stdout"}
	}

	for _, file := range o.Files {
		if file == "stdout" {
			writers = append(writers, os.Stdout)
		} else {
			fd, err := os.Create(file)
			if err != nil {
				return err
			}

			of := bufio.NewWriter(fd)

			writers = append(writers, of)
			o.closers = append(o.closers, fd)
		}
	}
	o.writer = io.MultiWriter(writers...)
	slog.Debug("using FILE output", "uri", o.Files)
	return nil
}

// Close the output writer
func (o *Output) Close() error {
	var err error
	for _, c := range o.closers {
		errClose := c.Close()
		if errClose != nil {
			err = errClose
		}
	}
	return err
}

// Description returns description
func (o *Output) Description() string {
	return "Send page view to file(s)"
}

// SendEvent send event to the Output
func (o *Output) SendEvent(event events.Event) error {
	if !o.condition.Match(event) {
		return nil
	}
	b, err := o.serializer.Serialize(event)
	if err != nil {
		return fmt.Errorf("unable to serialize page view: %v", err)
	}

	if _, err = o.writer.Write(b); err != nil {
		return fmt.Errorf("unable to write page view to file output: %v", err)
	}

	return nil
}

func init() {
	outputs.Add("file", func() outputs.Output {
		return &Output{}
	})
}
