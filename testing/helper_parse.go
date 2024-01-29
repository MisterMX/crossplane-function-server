package testing

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"unicode"

	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func mustUnstructuredFromYAML(rawYAML []byte) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	if err := yaml.Unmarshal(rawYAML, u); err != nil {
		panic(err.Error())
	}
	return u
}

func mustUnstructuredFromJSON(rawJSON []byte) *unstructured.Unstructured {
	u := &unstructured.Unstructured{}
	if err := json.Unmarshal(rawJSON, u); err != nil {
		panic(err.Error())
	}
	return u
}

func unmarshalObjectsYAML(rawYAML []byte) ([]*unstructured.Unstructured, error) {
	buf := bytes.NewBuffer(rawYAML)
	return unmarshalObjectsReader(buf)
}

func unmarshalObjectsReader(in io.Reader) ([]*unstructured.Unstructured, error) {
	objects := []*unstructured.Unstructured{}
	reader := yaml.NewYAMLReader(bufio.NewReader(in))
	for {
		data, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return objects, err
		}
		if len(data) == 0 {
			continue
		}
		if isWhiteSpace(data) {
			continue
		}
		o := &unstructured.Unstructured{}
		if err := yaml.Unmarshal(data, o); err != nil {
			return nil, err
		}
		objects = append(objects, o)
	}
	return objects, nil
}

// isWhiteSpace determines whether the passed in bytes are all unicode white
// space.
func isWhiteSpace(bytes []byte) bool {
	empty := true
	for _, b := range bytes {
		if !unicode.IsSpace(rune(b)) {
			empty = false
			break
		}
	}
	return empty
}
