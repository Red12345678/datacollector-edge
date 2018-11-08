// Copyright 2018 StreamSets Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package textrecord

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/spf13/cast"
	"github.com/streamsets/datacollector-edge/api"
	"github.com/streamsets/datacollector-edge/api/dataformats"
	"github.com/streamsets/datacollector-edge/api/fieldtype"
	"github.com/streamsets/datacollector-edge/api/linkedhashmap"
	"github.com/streamsets/datacollector-edge/container/recordio"
	"io"
)

const (
	DefaultTextField = "text"
)

type TextWriterFactoryImpl struct {
	// TODO: Add needed configs
}

func (t *TextWriterFactoryImpl) CreateWriter(
	context api.StageContext,
	writer io.Writer,
) (dataformats.RecordWriter, error) {
	var recordWriter dataformats.RecordWriter
	recordWriter = newRecordWriter(context, writer)
	return recordWriter, nil
}

type TextWriterImpl struct {
	context api.StageContext
	writer  *bufio.Writer
}

func (textWriter *TextWriterImpl) WriteRecord(r api.Record) error {
	recordValue, _ := r.Get()
	textFieldValue, err := textWriter.getTextFieldPathValue(recordValue)
	if err != nil {
		return err
	}
	fmt.Fprintln(textWriter.writer, textFieldValue)
	return nil
}

func (textWriter *TextWriterImpl) getTextFieldPathValue(field *api.Field) (string, error) {
	var textFieldValue string
	if field.Value == nil {
		return textFieldValue, nil
	}
	var err error = nil
	switch field.Type {
	case fieldtype.MAP:
		fieldValue := field.Value.(map[string]*api.Field)
		textField := fieldValue[DefaultTextField]
		if textField.Type != fieldtype.STRING {
			err = errors.New("Invalid Field Type for Text Field path - " + textField.Type)
			return textFieldValue, err
		}
		textFieldValue = cast.ToString(textField.Value)
		return textFieldValue, err
	case fieldtype.LIST_MAP:
		listMapValue := field.Value.(*linkedhashmap.Map)
		textValue, found := listMapValue.Get(DefaultTextField)
		if !found {
			return textFieldValue, fmt.Errorf("invalid field path - %s", DefaultTextField)
		}
		textField := textValue.(*api.Field)
		if textField.Type != fieldtype.STRING {
			err = errors.New("Invalid Field Type for Text Field path - " + textField.Type)
			return textFieldValue, err
		}
		textFieldValue = cast.ToString(textField.Value)
		return textFieldValue, err
	default:
		err = errors.New("unsupported Field Type")
	}
	return textFieldValue, err
}

func (textWriter *TextWriterImpl) Flush() error {
	return recordio.Flush(textWriter.writer)
}

func (textWriter *TextWriterImpl) Close() error {
	return recordio.Close(textWriter.writer)
}

func newRecordWriter(context api.StageContext, writer io.Writer) *TextWriterImpl {
	return &TextWriterImpl{
		context: context,
		writer:  bufio.NewWriter(writer),
	}
}
