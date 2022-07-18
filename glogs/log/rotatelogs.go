// package rotatelogs is a port of File-RotateLogs from Perl
// (https://metacpan.org/release/File-RotateLogs), and it allows
// you to automatically rotate output files when you write to them
// according to the filename pattern that you can specify.
package log

import (
	"cloud-utlis/glogs/log/fileutil"
	"cloud-utlis/glogs/log/strftime"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

func (c clockFn) Now() time.Time {
	return c()
}

// New creates a new RotateLogs object. A log filename pattern
// must be passed. Optional `option` parameters may be passed
func New(p string, options ...Option) (*RotateLogs, error) {
	globPattern := p
	for _, re := range patternConversionRegexps {
		globPattern = re.ReplaceAllString(globPattern, "*")
	}

	pattern, err := strftime.New(p)
	if err != nil {
		return nil, errors.Wrap(err, `invalid strftime pattern`)
	}

	var clock Clock = Local
	rotationTime := 24 * time.Hour
	var rotationSize int64
	var rotationCount uint
	var linkName string
	var maxAge time.Duration
	var handler Handler
	var forceNewFile bool
	var noBuffer bool

	for _, o := range options {
		switch o.Name() {
		case optkeyClock:
			clock = o.Value().(Clock)
		case optkeyLinkName:
			linkName = o.Value().(string)
		case optkeyMaxAge:
			maxAge = o.Value().(time.Duration)
			if maxAge < 0 {
				maxAge = 0
			}
		case optkeyRotationTime:
			rotationTime = o.Value().(time.Duration)
			if rotationTime < 0 {
				rotationTime = 0
			}
		case optkeyRotationSize:
			rotationSize = o.Value().(int64)
			if rotationSize < 0 {
				rotationSize = 0
			}
		case optkeyRotationCount:
			rotationCount = o.Value().(uint)
		case optkeyHandler:
			handler = o.Value().(Handler)
		case optkeyForceNewFile:
			forceNewFile = true
		case optkeyNoBufferWrite:
			noBuffer = true
		}
	}

	if maxAge > 0 && rotationCount > 0 {
		//return nil, errors.New("options MaxAge and RotationCount cannot be both set")
		rotationCount = 0
	}

	if maxAge == 0 && rotationCount == 0 {
		// if both are 0, give maxAge a sane default
		maxAge = 7 * 24 * time.Hour
	}

	// 获取初始化位移
	generation := 0
	baseFn := fileutil.GenerateFn(pattern, clock, rotationTime)
	fileNameExt := filepath.Base(baseFn)
	fileNameRegexp := regexp.MustCompile(fileNameExt[:len(fileNameExt)-4] + `.(\d+).log`)
	filePath := filepath.Dir(p)
	files, _ := ioutil.ReadDir(filePath)
	for _, item := range files {
		regexRes := fileNameRegexp.FindStringSubmatch(item.Name())
		if len(regexRes) == 2 {
			num, _ := strconv.Atoi(regexRes[len(regexRes)-1])
			if generation < num {
				generation = num
			}
		}
	}

	rl := &RotateLogs{
		generation:    generation,
		clock:         clock,
		eventHandler:  handler,
		globPattern:   globPattern,
		linkName:      linkName,
		maxAge:        maxAge,
		pattern:       pattern,
		rotationTime:  rotationTime,
		rotationSize:  rotationSize,
		rotationCount: rotationCount,
		forceNewFile:  forceNewFile,
	}
	if noBuffer {
		rl.outFh = fileutil.GenerateDirectIO()

		// 无缓冲需要定时30s刷新sync
		ticker := time.NewTicker(30 * time.Second)
		go func() {
			for {
				select {
				case <-ticker.C:
					rl.tickerSync()
				}
			}
		}()

		// 无缓冲需要信号处理
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			switch <-quit {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP:
				ticker.Stop()
				_ = rl.Close()
				return
			}
		}()
	} else {
		rl.outFh = fileutil.GenerateBuffIO()
	}

	return rl, nil
}

// Write satisfies the io.Writer interface. It writes to the
// appropriate file handle that is currently being used.
// If we have reached rotation time, the target file gets
// automatically rotated, and also purged if necessary.
func (rl *RotateLogs) Write(p []byte) (n int, err error) {
	// Guard against concurrent writes
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	err = rl.getWriterNolock(false, false)
	if err != nil {
		return 0, errors.Wrap(err, `failed to acquite target io.Writer`)
	}

	return rl.outFh.Write(p)
}

// must be locked during this operation
func (rl *RotateLogs) getWriterNolock(bailOnRotateFail, useGenerationalNames bool) error {
	generation := rl.generation
	previousFn := rl.curFn

	// This filename contains the name of the "NEW" filename
	// to log to, which may be newer than rl.currentFilename
	baseFn := fileutil.GenerateFn(rl.pattern, rl.clock, rl.rotationTime)
	filename := baseFn
	var forceNewFile bool

	fi, err := os.Stat(rl.curFn)
	sizeRotation := false
	if err == nil && rl.rotationSize > 0 && rl.rotationSize <= fi.Size() {
		forceNewFile = true
		sizeRotation = true
	}

	if baseFn != rl.curBaseFn {
		if rl.curBaseFn != "" {
			generation = 0
		}
		// even though this is the first write after calling New(),
		// check if a new file needs to be created
		if rl.forceNewFile {
			forceNewFile = true
		}
	} else {
		if !useGenerationalNames && !sizeRotation {
			// nothing to do
			return nil
		}
		forceNewFile = true
		generation++
	}
	if forceNewFile {
		// A new file has been requested. Instead of just using the
		// regular strftime pattern, we create a new file name using
		// generational names such as "foo.1", "foo.2", "foo.3", etc
		var name = filename
		for {
			if generation == 0 {
				name = filename
				rl.rotationFirst = true
			} else {
				if strings.HasSuffix(filename, ".log") {
					if !rl.rotationFirst {
						generation++
						rl.rotationFirst = true
					}
					tmpName := filename[:len(filename)-4]
					newFileName := fmt.Sprintf("%s.%d.log", tmpName, generation)
					rl.outFh.Close()
					err := os.Rename(filename, newFileName)
					if err != nil {
						fmt.Println(fmt.Sprintf("Failed to rename file %s to %s, the error is ", filename, newFileName), err)
					}
				} else {
					name = fmt.Sprintf("%s.%d", filename, generation)
				}
			}
			if _, err := os.Stat(name); err != nil {
				filename = name

				break
			}
			generation++
		}
	}

	err = rl.outFh.CreateFile(filename)
	if err != nil {
		return errors.Wrapf(err, `failed to create a new file %v`, filename)
	}

	if err = rl.rotateNolock(filename); err != nil {
		err = errors.Wrap(err, "failed to rotate")
		if bailOnRotateFail {
			// Failure to rotate is a problem, but it's really not a great
			// idea to stop your application just because you couldn't rename
			// your log.
			//
			// We only return this error when explicitly needed (as specified by bailOnRotateFail)
			//
			// However, we *NEED* to close `fh` here
			if rl.outFh != nil { // probably can't happen, but being paranoid
				_ = rl.outFh.Close()
			}

			return err
		}
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}

	rl.curBaseFn = baseFn
	rl.curFn = filename
	rl.generation = generation

	if h := rl.eventHandler; h != nil {
		go h.Handle(&FileRotatedEvent{
			prev:    previousFn,
			current: filename,
		})
	}

	return nil
}

// CurrentFileName returns the current file name that
// the RotateLogs object is writing to
func (rl *RotateLogs) CurrentFileName() string {
	rl.mutex.RLock()
	defer rl.mutex.RUnlock()

	return rl.curFn
}

var patternConversionRegexps = []*regexp.Regexp{
	regexp.MustCompile(`%[%+A-Za-z]`),
	regexp.MustCompile(`\*+`),
}

type cleanupGuard struct {
	enable bool
	fn     func()
	mutex  sync.Mutex
}

func (g *cleanupGuard) Enable() {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	g.enable = true
}

func (g *cleanupGuard) Run() {
	g.fn()
}

// Rotate forcefully rotates the log files. If the generated file name
// clash because file already exists, a numeric suffix of the form
// ".1", ".2", ".3" and so forth are appended to the end of the log file
//
// Thie method can be used in conjunction with a signal handler so to
// emulate servers that generate new log files when they receive a
// SIGHUP
func (rl *RotateLogs) Rotate() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()
	err := rl.getWriterNolock(true, true)
	return err
}

func (rl *RotateLogs) rotateNolock(filename string) error {
	lockfn := filename + `_lock`
	fh, err := os.OpenFile(lockfn, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		// Can't lock, just return
		return err
	}

	var guard cleanupGuard
	guard.fn = func() {
		fh.Close()
		os.Remove(lockfn)
	}
	defer guard.Run()

	if rl.linkName != "" {
		tmpLinkName := filename + `_symlink`

		// Change how the link name is generated based on where the
		// target location is. if the location is directly underneath
		// symlink with a relative path
		linkDest := filename
		linkDir := filepath.Dir(rl.linkName)

		baseDir := filepath.Dir(filename)
		if strings.Contains(rl.linkName, baseDir) {
			tmp, err := filepath.Rel(linkDir, filename)
			if err != nil {
				return errors.Wrapf(err, `failed to evaluate relative path from %#v to %#v`, baseDir, rl.linkName)
			}

			linkDest = tmp
		}

		if err := os.Symlink(linkDest, tmpLinkName); err != nil {
			return errors.Wrap(err, `failed to create new symlink`)
		}

		// the directory where rl.linkName should be created must exist
		_, err := os.Stat(linkDir)
		if err != nil { // Assume err != nil means the directory doesn't exist
			if err := os.MkdirAll(linkDir, 0755); err != nil {
				return errors.Wrapf(err, `failed to create directory %s`, linkDir)
			}
		}

		if err := os.Rename(tmpLinkName, rl.linkName); err != nil {
			return errors.Wrap(err, `failed to rename new symlink`)
		}
	}

	if rl.maxAge <= 0 && rl.rotationCount <= 0 {
		return errors.New("panic: maxAge and rotationCount are both set")
	}

	matches, err := filepath.Glob(rl.globPattern)
	if err != nil {
		return err
	}

	cutoff := rl.clock.Now().Add(-1 * rl.maxAge)

	// the linter tells me to pre allocate this...
	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		// Ignore lock files
		if strings.HasSuffix(path, "_lock") || strings.HasSuffix(path, "_symlink") {
			continue
		}

		fi, err := os.Stat(path)
		if err != nil {
			continue
		}

		fl, err := os.Lstat(path)
		if err != nil {
			continue
		}

		if rl.maxAge > 0 && fi.ModTime().After(cutoff) {
			continue
		}

		if rl.rotationCount > 0 && fl.Mode()&os.ModeSymlink == os.ModeSymlink {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	if rl.rotationCount > 0 {
		// Only delete if we have more than rotationCount
		if rl.rotationCount >= uint(len(toUnlink)) {
			return nil
		}

		toUnlink = toUnlink[:len(toUnlink)-int(rl.rotationCount)]
	}

	if len(toUnlink) <= 0 {
		return nil
	}

	guard.Enable()
	go func() {
		// unlink files on a separate goroutine
		for _, path := range toUnlink {
			os.Remove(path)
		}
	}()

	return nil
}

// Close satisfies the io.Closer interface. You must
// call this method if you performed any writes to
// the object.
func (rl *RotateLogs) Close() error {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	if rl.outFh == nil {
		return nil
	}

	rl.outFh.Close()

	return nil
}

func (rl *RotateLogs) tickerSync() {
	rl.mutex.Lock()
	defer rl.mutex.Unlock()

	_ = rl.outFh.SyncFile()
}
