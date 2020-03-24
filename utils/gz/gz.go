package gz

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
)

func Marshal(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, err := gw.Write(data)
	gw.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Unmarshal(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	_, err := buf.Write(data)
	if err != nil {
		return nil, err
	}
	gr, err := gzip.NewReader(&buf)
	defer gr.Close()
	if err != nil {
		return nil, err
	}
	undatas, err := ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return undatas, nil
}
