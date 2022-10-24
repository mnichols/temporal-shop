package testhelper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/google/go-cmp/cmp/cmpopts"
)

var SortSlices = cmpopts.SortSlices(func(x, y interface{}) bool {
	return fmt.Sprint(x) < fmt.Sprint(y)
})

func MustReadCloser(obj interface{}) io.ReadCloser {
	data, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return ioutil.NopCloser(bytes.NewReader(data))
}
