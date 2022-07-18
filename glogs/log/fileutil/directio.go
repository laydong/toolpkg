package fileutil

import (
	"github.com/pkg/errors"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

const (
	LINUX_DIRECT          = 0x4000
	WINDOWS_NO_BUFF       = 0x20000000
	WINDOWS_WRITE_THROUGH = 0x20000000
)

type DirectIO struct {
	bufLen   int
	bufMax   int
	fileName string // 文件名
	counter  int64  // 计数器
	baseCap  []byte
	buf      []byte
	fh       *os.File
}

func GenerateDirectIO() IOWriter {
	directIO := new(DirectIO)
	directIO.buf = make([]byte, 4096, 4096)
	directIO.baseCap = Get4096ByteNbsp()
	directIO.bufMax = 4096
	return directIO
}

func (d *DirectIO) CreateFile(filename string) error {
	var fh *os.File
	var err error

	// 保留文件名在内存中
	d.fileName = filename

	// make sure the dir is existed, eg:
	// ./foo/bar/baz/hello.log must make sure ./foo/bar/baz is existed
	dirname := filepath.Dir(filename)
	if err = os.MkdirAll(dirname, 0755); err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dirname)
	}

	if runtime.GOOS == "linux" {
		fh, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY|LINUX_DIRECT, 0644)
		if err != nil {
			log.Printf("打开文件失败 err: %s", err.Error())
			return err
		}
	} else if runtime.GOOS == "windows" {
		fh, err = os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY|WINDOWS_NO_BUFF|WINDOWS_WRITE_THROUGH, 0644)
		if err != nil {
			log.Printf("打开文件失败 err: %s", err.Error())
			return err
		}
	} else {
		panic("not support this platform")
	}
	d.fh = fh
	return nil
}

// 需要补齐内存写入，防止下次打开的时候内存未对齐写入失败
// 30秒内没有日志写入的话需要刷新日志到文件，用空格补齐
func (d *DirectIO) SyncFile() error {
	if d.counter != 0 {
		copy(d.buf[d.bufLen:], d.baseCap[d.bufLen:])
		_, err := d.fh.Write(d.buf)
		if err != nil {
			log.Printf("sync file err: %s", err.Error())
			return err
		}
	}
	d.bufLen = 0
	d.counter = 0
	return nil
}

func (d *DirectIO) Write(data []byte) (int, error) {
	d.counter++
	needWrite := len(data)
	readIndex := 0
	for needWrite > 0 {
		canWrite := d.bufMax - d.bufLen
		if canWrite > needWrite {
			canWrite = needWrite
		}
		if canWrite == 0 {
			n, err := d.fh.Write(d.buf)
			if err != nil {
				return n, err
			}
			d.bufLen = 0
		} else {
			copy(d.buf[d.bufLen:], data[readIndex:readIndex+canWrite])
			readIndex += canWrite
			d.bufLen += canWrite
			needWrite -= canWrite
		}
	}
	return len(data), nil
}

func (d *DirectIO) Close() error {
	if d.fh == nil {
		return nil
	}
	if d.bufLen != 0 {
		name := d.fh.Name()
		err := d.fh.Close()
		if err != nil {
			return err
		}
		fh, _ := os.OpenFile(name, os.O_APPEND|os.O_WRONLY, 0666)
		_, _ = fh.Write(d.buf[:d.bufLen])
		_ = fh.Close()
	}

	d.buf = make([]byte, 4096, 4096)
	d.bufLen = 0
	d.fh = nil
	return nil
}
