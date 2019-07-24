package runner

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"text/template"
	"time"

	"github.com/jhump/protoreflect/desc"
)

// call template data
type callTemplateData struct {
	WorkerID           string // unique worker ID
	RequestNumber      int64  // unique incremented request number for each request
	FullyQualifiedName string // fully-qualified name of the method call
	UUID               string // random UUID v4
	MethodName         string // shorter call method name
	ServiceName        string // the service name
	InputName          string // name of the input message type
	OutputName         string // name of the output message type
	IsClientStreaming  bool   // whether this call is client streaming
	IsServerStreaming  bool   // whether this call is server streaming
	Timestamp          string // timestamp of the call in RFC3339 format
	TimestampUnix      int64  // timestamp of the call as unix time
}

// generateUUID sourced from https://stackoverflow.com/a/25736155/210827.
func generateUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Printf("Rand error: %+v\n", err)
	}
	return fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// newCallTemplateData returns new call template data
func newCallTemplateData(mtd *desc.MethodDescriptor, workerID string, reqNum int64) *callTemplateData {
	now := time.Now()

	return &callTemplateData{
		WorkerID:           workerID,
		RequestNumber:      reqNum,
		FullyQualifiedName: mtd.GetFullyQualifiedName(),
		UUID:               generateUUID(),
		MethodName:         mtd.GetName(),
		ServiceName:        mtd.GetService().GetName(),
		InputName:          mtd.GetInputType().GetName(),
		OutputName:         mtd.GetOutputType().GetName(),
		IsClientStreaming:  mtd.IsClientStreaming(),
		IsServerStreaming:  mtd.IsServerStreaming(),
		Timestamp:          now.Format(time.RFC3339),
		TimestampUnix:      now.Unix(),
	}
}

func (td *callTemplateData) execute(data string) (*bytes.Buffer, error) {
	t := template.Must(template.New("call_template_data").Parse(data))
	var tpl bytes.Buffer
	err := t.Execute(&tpl, td)
	return &tpl, err
}

func (td *callTemplateData) executeData(data string) ([]byte, error) {
	if len(data) > 0 {
		input := []byte(data)
		tpl, err := td.execute(data)
		if err == nil {
			input = tpl.Bytes()
		}

		return input, nil
	}

	return []byte{}, nil
}

func (td *callTemplateData) executeMetadata(metadata string) (*map[string]string, error) {
	var mdMap map[string]string

	if len(metadata) > 0 {
		input := []byte(metadata)
		tpl, err := td.execute(metadata)
		if err == nil {
			input = tpl.Bytes()
		}

		err = json.Unmarshal(input, &mdMap)
		if err != nil {
			return nil, err
		}
	}

	return &mdMap, nil
}
