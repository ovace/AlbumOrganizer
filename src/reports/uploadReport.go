package reports

// This utility copies photos into a date-specific folder

import (
	"fmt"
)

const form = "2006:01:02 15:04:05"
const folderFormat = "2006/01"

type pathEntry struct {
	filename string
}

// pathWalker walks the filesystem, queueing pathEntry items onto the queue.
type pathWalker struct {
	myCounters
	queue chan pathEntry
}
type moveEntry struct {
	source   string
	dest     string
	filehash uint32

	//	dupcnt   uint32
}
type fileMover struct {
	sourcePath string
	destPath   string
	isCopy     bool
	dryRun     bool
	myCounters
	queue chan moveEntry
}
type myCounters struct {
	readDirCounter   int
	readFileCounter  int
	writeDirCounter  int
	writeFileCounter int
	dupCounter       int
	failCounter      int
}

func (c *myCounters) UploadReport() {

	totDirRead := c.readDirCounter
	totDirWrite := c.writeDirCounter
	totFilRead := c.readFileCounter
	totFilWrite := c.writeFileCounter

	fmt.Printf("Total dir scanned: %d\n Total dir created: %d\n Total files read: %d\n Total Files transfered: %d\n",
		totDirRead, totDirWrite, totFilRead, totFilWrite)

	return
}
