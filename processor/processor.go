package processor

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Global Version
var Version = "0.1.0"

// Verbose enables verbose logging output
var Verbose = false

// Debug enables debug logging output
var Debug = false

// Trace enables trace logging output which is extremely verbose
var Trace = false

// Recursive to walk directories
var Recursive = false

// If set to true disables the use of memory maps
var NoMmap = false

// Do not print out results as they are processed
var NoStream = false

// Should the application print all hashes it knows about
var Hashes = false

// List of hashes that we want to process
var Hash = []string{}

// Format sets the output format of the formatter
var Format = ""

// FileOutput sets the file that output should be written to
var FileOutput = ""

// DirFilePaths is not set via flags but by arguments following the flags for file or directory to process
var DirFilePaths = []string{}

// FileListQueueSize is the queue of files found and ready to be processed
var FileListQueueSize = 1000

// Number of bytes in a size to enable memory maps or streaming
var StreamSize int64 = 1000000

var s_md5 = "md5"
var s_sha1 = "sha1"
var s_sha256 = "sha256"
var s_sha512 = "sha512"
var s_blake2b256 = "blake2b256"
var s_blake2b512 = "blake2b512"
var s_sha3224 = "sha3224"
var s_sha3256 = "sha3256"
var s_sha3384 = "sha3384"
var s_sha3512 = "sha3512"

// Process is the main entry point of the command line it sets everything up and starts running
func Process() {
	// Clean up any invalid arguments before setting everything up
	if len(DirFilePaths) == 0 {
		DirFilePaths = append(DirFilePaths, ".")
	}

	// Clean up hashes by setting all to lower
	h := []string{}
	for _, x := range Hash {
		h = append(h, strings.ToLower(x))
	}
	Hash = h

	// Check if the paths or files added exist and exit if not
	for _, f := range DirFilePaths {
		fpath := filepath.Clean(f)

		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			printError("file or directory does not exist: " + fpath)
			os.Exit(1)
		}
	}

	fileListQueue := make(chan string, FileListQueueSize)    // Files ready to be read from disk
	fileSummaryQueue := make(chan Result, FileListQueueSize) // Results ready to be printed

	// Spawn routine to start finding files on disk
	go func() {
		for _, f := range DirFilePaths {
			walkDirectory(f, fileListQueue)
		}
		close(fileListQueue)
	}()

	go fileProcessorWorker(fileListQueue, fileSummaryQueue)
	result := fileSummarize(fileSummaryQueue)

	if FileOutput == "" {
		fmt.Print(result)
	} else {
		_ = ioutil.WriteFile(FileOutput, []byte(result), 0600)
		fmt.Println("results written to " + FileOutput)
	}
}

func hasHash(hash string) bool {
	for _, x := range Hash {
		if x == "all" {
			return true
		}

		if x == hash {
			return true
		}
	}

	return false
}
