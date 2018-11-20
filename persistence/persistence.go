package persistence

import (
	"../common/maybe"
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

func ToPersistence(space string, layer int32, class string, id int64, content []byte) (err maybe.MaybeError) {
	if space == "" {
		err.Error(errors.New("empty space"))
		return
	}
	if layer < 0 {
		err.Error(fmt.Errorf("illegal layer id: %d", layer))
		return
	}
	if class == "" {
		err.Error(errors.New("empty class"))
		return
	}
	if id < 0 {
		err.Error(fmt.Errorf("illegal id: %d", id))
		return
	}
	l := len(content)
	if l == 0 {
		err.Error(errors.New("empty content"))
		return
	}

	filePath := fmt.Sprintf("%s/%d/%s_%d", space, layer, class, id)

	file, e := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if e != nil {
		err.Error(e)
		return
	}
	writter := bufio.NewWriterSize(file, l)

	_, e = writter.Write(content)
	if e != nil {
		err.Error(e)
		return
	}

	return
}

func FromPersistence(space string, layer int32, class string, id int64) (ret maybe.MaybeBytes) {
	if space == "" {
		ret.Error(errors.New("empty space"))
		return
	}
	if layer < 0 {
		ret.Error(fmt.Errorf("illegal layer id: %d", layer))
		return
	}
	if class == "" {
		ret.Error(errors.New("empty class"))
		return
	}
	if id < 0 {
		ret.Error(fmt.Errorf("illegal id: %d", id))
		return
	}

	filePath := fmt.Sprintf("%s/%d/%s_%d", space, layer, class, id)

	file, e := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if e != nil {
		ret.Error(e)
		return
	}
	reader := bufio.NewReader(file)

	content, e := ioutil.ReadAll(reader)
	if e != nil {
		ret.Error(e)
		return
	}

	ret.Value(content)
	return
}
