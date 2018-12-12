package persistence

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/incubator/common/maybe"
	"io/ioutil"
	"math"
	"os"
	"strings"
	"time"
)

const (
	FROM_PERSISTENCE_MODE_RECOVER = iota
	FROM_PERSISTENCE_MODE_REBOOT
)

type Persistentable interface {
	SetStoreExpiration(time.Duration) maybe.MaybeError
	SetLoadExpiration(time.Duration) maybe.MaybeError
	GetStoreExpiration() time.Duration
	GetLoadExpiration() time.Duration
	Persistent() maybe.MaybeError
}

type CommomPersistentable struct {
	storeExpiration time.Duration
	loadExpiration  time.Duration
}

func (this *CommomPersistentable) SetStoreExpiration(e time.Duration) (err maybe.MaybeError) {
	if e <= 0 {
		err.Error(fmt.Errorf("illegal store expiration: %d", e))
		return
	}
	this.storeExpiration = e

	err.Error(nil)
	return
}

func (this *CommomPersistentable) SetLoadExpiration(e time.Duration) (err maybe.MaybeError) {
	if e <= 0 {
		err.Error(fmt.Errorf("illegal load expiration: %d", e))
		return
	}
	this.loadExpiration = e

	err.Error(nil)
	return
}

func (this CommomPersistentable) GetStoreExpiration() time.Duration {
	return this.storeExpiration
}

func (this CommomPersistentable) GetLoadExpiration() time.Duration {
	return this.loadExpiration
}

func ToPersistence(expire time.Duration, version int64, space string, layer int32, class string, content []byte) (err maybe.MaybeError) {
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
	if version < 0 {
		err.Error(fmt.Errorf("illegal version: %d", version))
		return
	}
	l := len(content)
	if l == 0 {
		err.Error(errors.New("empty content"))
		return
	}

	filePath := fmt.Sprintf("%s/%d/%s_%10d_%d", space, layer, class, version, time.Now().Unix())

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

	fileDir := fmt.Sprintf("%s/%d/", space, layer, class, version)
	filePrefix := fmt.Sprintf("%s_%10d", class, version)

	dir, e := ioutil.ReadDir(fileDir)
	if e != nil {
		err.Error(e)
		return
	}

	for _, f := range dir {
		if !f.IsDir() && strings.Contains(f.Name(), filePrefix) {
			if f == nil || f.ModTime().Before(f.ModTime()) {
				if time.Now().Sub(f.ModTime()) > expire {
					e := os.Remove(f.Name())
					if e != nil {
						err.Error(e)
						return
					}
				}
			}
			break
		}
	}

	return
}

func FromPersistence(mode int, expire time.Duration, version int64, space string, layer int32, class string) (ret maybe.MaybeBytes) {
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
	if version < 0 {
		ret.Error(fmt.Errorf("illegal version: %d", version))
		return
	}

	var filePrefix string
	switch mode {
	case FROM_PERSISTENCE_MODE_REBOOT:
		filePrefix = class
		expire = time.Duration(math.MaxInt64)
	case FROM_PERSISTENCE_MODE_RECOVER:
		filePrefix = fmt.Sprintf("%s_%10d", class, version)
	default:
		ret.Error(fmt.Errorf("unknow mode: %d", mode))
		return
	}

	fileDir := fmt.Sprintf("%s/%d/", space, layer, class, version)

	dir, err := ioutil.ReadDir(fileDir)
	if err != nil {
		ret.Error(err)
		return
	}

	var fi os.FileInfo
	for _, f := range dir {
		if !f.IsDir() && strings.Contains(f.Name(), filePrefix) {
			if fi == nil || fi.ModTime().Before(f.ModTime()) {
				fi = f
			}
			break
		}
	}

	if fi == nil {
		ret.Error(fmt.Errorf("file not found: %s/%s_<time>", fileDir, filePrefix))
		return
	}

	if time.Now().Sub(fi.ModTime()) > expire {
		ret.Error(fmt.Errorf("file expired: modTime=%d, expire=%d", fi.ModTime(), expire))
		return
	}

	file, e := os.OpenFile(fi.Name(), os.O_RDONLY, 0644)
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
