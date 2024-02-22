package output

import (
	"fmt"
	"os"
	"strings"
)

var (
	printDebug = false
)

func SetDebug(debug bool) {
	printDebug = debug
}

func Infoln(s any) {
	fmt.Println(s)
}

func Infof(s string, args ...any) {
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Printf(s, args...)
}

func Errorln(s any) {
	fmt.Fprintln(os.Stderr, s)
}

func Errorf(s string, args ...any) {
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Fprintf(os.Stderr, s, args...)
}

func Fatalln(s any) {
	Errorln(s)
	os.Exit(1)
}

func Fatalf(s string, args ...any) {
	Errorf(s, args...)
	os.Exit(1)
}

func Debugln(s any) {
	if printDebug {
		fmt.Println(s)
	}
}

func Debugf(s string, args ...any) {
	if !printDebug {
		return
	}
	// Avoid printing incomplete lines
	if !strings.HasSuffix(s, "\n") {
		s = s + "\n"
	}
	fmt.Printf(s, args...)
}
