package main

/*
#include <unistd.h>
#include <stdint.h>
#include <stdlib.h>
#include <stdio.h>
#include <errno.h>
#include <string.h>
#include <syslog.h>
#include <sys/ioctl.h>
#include <sys/poll.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <linux/random.h>


int random_add_entropy(int fhandle, char *buf, int size) {
	struct {
		int ent_count;
		int size;
		unsigned char data[size];
	} entropy;
	entropy.ent_count = size * 8;
	entropy.size = size;
	memcpy((void*)entropy.data, (void*)buf, (size_t)size);
	if (ioctl(fhandle, RNDADDENTROPY, &entropy) != 0) {
		return(1);
	}
	return(0);
};
*/
import "C"

import (
	"log/syslog"
	"os"
	"time"
	"unsafe"
)

func checkError(err error, logger *syslog.Writer) {
	if err != nil {
		os.Exit(1)
	}
}

var inFileName = "/dev/TrueRNG"
var outFileName = "/dev/random"

func main() {
	logger, err := syslog.New(syslog.LOG_ERR, "TrueRNGd")
	checkError(err, logger)
	defer logger.Close()

	logger.Info("Starting up TrueRNGd")

	// open random
	var buffer [4096]byte
	inFile, err := os.Open(inFileName)
	checkError(err, logger)
	defer inFile.Close()

	outFile, err := os.Open(outFileName)
	checkError(err, logger)
	defer outFile.Close()

	for {
		bytesRead, err := inFile.Read(buffer[:])
		checkError(err, logger)

		result := C.random_add_entropy(C.int(outFile.Fd()), (*C.char)(unsafe.Pointer(&buffer[0])), C.int(bytesRead))
		if result != 0 {
			logger.Crit("Error adding entropy!\n")
			os.Exit(1)
		}
		time.Sleep(500 * time.Millisecond)
	}
}
