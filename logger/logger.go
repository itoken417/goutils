package logger

import (
	"github.com/mitchellh/panicwrap"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var fh os.File
var MODE bool
var WRITE_CALLER bool

var errlog *log.Logger
var infolog *log.Logger

func Init(is_debug, write_caller bool) bool {
	MODE = is_debug
	WRITE_CALLER = write_caller
	errlog = log.New(os.Stdout, "[err]", log.LstdFlags)
	infolog = log.New(os.Stdout, "[log]", log.LstdFlags)

	if MODE {
		log.Println("DEBUG Mode")
	} else {
		_, sorce, _, _ := runtime.Caller(1)
		path := filepath.Dir(sorce)
		base := filepath.Base(path)
		file := "./logs/" + base + ".log"
		fh, err := os.OpenFile(file,
			os.O_CREATE|os.O_TRUNC|os.O_WRONLY,
			0666)
		if err != nil {
			panic("cannot open app log" + err.Error())
		}
		log.SetOutput(io.MultiWriter(fh))
		infolog.SetOutput(io.MultiWriter(fh))
		errlog.SetOutput(io.MultiWriter(fh))
		log.Println("Release Mode")

		exitStatus, _ := panicwrap.BasicWrap(panicHandler)
		if exitStatus >= 0 {
			defer Destory()
			os.Exit(exitStatus)
		}
	}
	return true
}
func panicHandler(output string) {
	// output contains the full output (including stack traces) of the
	// panic. Put it in a file or something.
	log.Printf("The child panicked:\n\n%s\n", output)
	defer Destory()
	os.Exit(1)
}
func Destory() {
	//	log.Println("log.Destory")
	fh.Close()
}

func print_caller() {
	if WRITE_CALLER {
		_, file, line, _ := runtime.Caller(2)
		log.Printf("%s:%d", file, line)
	}
}

func Log(a ...interface{}) {
	print_caller()
	for idx, str := range a {
		if idx > 0 {
			infolog.Println(">>>>", str)
		} else {
			infolog.Println(":", str)
		}
	}
}

func Dump(a ...interface{}) {
	print_caller()
	for idx, str := range a {
		if idx > 0 {
			infolog.Printf(">>>> %#v", str)
		} else {
			infolog.Printf(": %#v", str)
		}
	}
}

func ErrLog(a ...interface{}) {
	print_caller()
	for idx, str := range a {
		if idx > 0 {
			errlog.Panicln(">>>>", str)
		} else {
			errlog.Panicln(":", str)
		}
	}
}

// err 捕捉
func ErrCheck(err error) {
	if err != nil {
		print_caller()
		panic(err)
	}
}
