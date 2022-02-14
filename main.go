package main

import (
	"os"
	"time"
	"flag"
	"fmt"
	"syscall"
	"math/rand"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/ghetzel/go-stockutil/sliceutil"
)

func main() {
	var fd int
	var err error
	var diskSize int64

	blockDev := flag.String("blockDev", "", "Block device (REQUIRED)")
	readFreq := flag.String("readFreq", "5m", "Frequency to check for record update (default 5m)")

	flag.Parse()

	flag.VisitAll(func(f *flag.Flag) {
		reqVars := []string{"blockDev"}
		if sliceutil.ContainsString(reqVars, f.Name) && f.Value.String() == "" {
			fmt.Printf("Required variable missing: %s\n", f.Name)
			os.Exit(1)
		}
	})

	sleepDur, err := time.ParseDuration(*readFreq)

	if err != nil {
		fmt.Print(err.Error(), "\n")
		return
	}

	diskInfo, err := diskfs.OpenWithMode(*blockDev, diskfs.ReadOnly)

	if err != nil {
		fmt.Print(err.Error(), "\n")
		return
	}

	diskSize = diskInfo.Size

	diskInfo = nil

	for {

		fd, err = syscall.Open(*blockDev, syscall.O_RDONLY, 0777)

		if err != nil {
			fmt.Print(err.Error(), "\n")
			return
		}

		readByte := rand.Int63n(diskSize - 1)

		fmt.Printf("Reading byte %d from device %s\n", readByte, *blockDev)

		_, err = syscall.Seek(fd, readByte, 0)

		if err != nil {
			fmt.Print(err.Error(), "\n")
			return
		}

		buffer := make([]byte, 1)

		_, err = syscall.Read(fd, buffer)

		if err != nil {
			fmt.Print(err.Error(), "\n")
		}

		err = syscall.Close(fd)

		if err != nil {
			fmt.Print(err.Error(), "\n")
		}

		time.Sleep(sleepDur)
	}
}
