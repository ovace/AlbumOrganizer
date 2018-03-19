package utils

import (
	"fmt"
	"os"
	"sync"

	dw "github.com/ovace/pkg/walk"
)

type pathEntry struct {
	filename string
}

// pathWalker walks the filesystem, queueing pathEntry items onto the queue.
type PathWalker struct {
	// MyCounters
	queue chan pathEntry
}
type HashPair struct {
	hash uint32
	path string
	// fDetail map[string]*FileInfo
}

func (env *Env) pathWalker() {
	// Part One: Build a map of file paths indexed by hash of each file
	//var lock sync.Mutex //Moved to global scope
	var TFiles, TBytes int // total files and bytes  //  moved to global env
	hashChan := make(chan HashPair, 1024)
	//hashDone := make(chan bool) // changed this to using WaitGroups
	wg := new(sync.WaitGroup)
	hashMap := make(map[uint32][]string)

	//defer close(hashDone)
	wg.Add(1)

	// Append each read file, indexed by the file hash value
	go func() {

		defer wg.Done()
		for sp := range hashChan {

			hashMap[sp.hash] = append(hashMap[sp.hash], sp.path)
		}
		//hashDone <- true
	}()

	fileVisitor := func(path string, info os.FileInfo, err error) error {
		if err != nil || info.Mode()&os.ModeType != 0 {
			return nil // skip special files
		}
		if size := info.Size(); size > 0 { // skip empty files
			lock.Lock()
			TFiles++
			TBytes += int(size)
			lock.Unlock()
			// fmt.Printf("path: %v\n",path)
			thisHash, _ := gh.FileHash(path)

			hashChan <- HashPair{thisHash, path}
		}
		return nil
	}

	fmt.Printf("info passed to Walk %v\t %v\n", srcPath, fileVisitor)
	dw.Walk(srcPath, fileVisitor)

	close(hashChan)
	//<-hashDone

}
