package main

import (
	"log"
	"os"
	"runtime"
	"runtime/pprof"
)

func startCPUProfile(path string) func() {
	if path == "" {
		return func() {}
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("could not create cpu profile: %v", err)
	}
	pprof.StartCPUProfile(f)
	return func() {
		pprof.StopCPUProfile()
		f.Close()
	}
}

func writeMemProfile(path string) {
	if path == "" {
		return
	}
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("could not create memory profile: %v", err)
	}
	defer f.Close()
	runtime.GC()
	pprof.WriteHeapProfile(f)
}
