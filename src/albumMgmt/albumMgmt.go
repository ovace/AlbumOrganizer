package main

//package experiments

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
	//"log"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	th "./utils/thumbnails"
	dbl "github.com/ovace/albumMgmt/pkg/databaselayer"
	fc "github.com/ovace/pkg/fileCopy"
	gh "github.com/ovace/pkg/fileHash"
	fm "github.com/ovace/pkg/fileMove"
	fs "github.com/ovace/pkg/fileStats"
	get "github.com/ovace/pkg/goexiftool"
	dw "github.com/ovace/pkg/walk"
	log "github.com/sirupsen/logrus"
)

//Initializers
const (
	oUsage = `
		Usage: %s [-n][-cp] <src> <dest>
		-dr=[true|false]      	dryrun [default=true]
		-cp=[true|false]     	copy [default=false]
		-v=[true|false]		verbose [default=false]
		-src=<srcfile>   	source path [default=""]
		-dest=<dest location>  	destination path [default=""]
	`
)

const (
	DT_GPS         = "2006:01:02  15:04:05Z"
	DT_ANSIC       = "Mon Jan _2 15:04:05 2006"
	DT_UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	DT_RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	DT_RFC822      = "02 Jan 06 15:04 MST"
	DT_RFC822Z     = "02 Jan 06 15:04 -0700" // RFC822 with numeric zone
	DT_RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	DT_RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	DT_RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700" // RFC1123 with numeric zone
	DT_RFC3339     = "2006-01-02T15:04:05Z07:00"
	DT_RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	DT_Kitchen     = "3:04PM"
	// Handy time stamps.
	Stamp      = "Jan _2 15:04:05"
	StampMilli = "Jan _2 15:04:05.000"
	StampMicro = "Jan _2 15:04:05.000000"
	StampNano  = "Jan _2 15:04:05.000000000"
)

var (
	dryRun       bool
	verbose      bool
	justCopy     bool
	srcPath      string
	destPath     string
	dupePath     string
	dupcount     int
	srcFilerowid int64
	tgtrowid     int64
	hashMap      map[uint32][]string
	lock         sync.Mutex
)

// type

//Configuration file
type Config struct {
	OUsage       string `json:"oUsage"`
	DtFormWTz    string `json: "dtFormWTz"`
	DtFormWoTz   string `json: "dtFormWoTz"`
	FolderFormat string `json: "folderFormat"`
	Logfile      string `json:  "logfile"`
	argsConf
	dbConf
}
type argsConf struct {
	DryRun   bool   `json: "dryRun"`
	Verbose  bool   `json:"verbose"`
	JustCopy bool   `json:"justCopy"`
	SrcPath  string `json:"srcPath"`
	DestPath string `json:"destPath"`
	DupePath string `json:"dupePath"`
}

type dbConf struct {
	DB_USER    string `json: "DB_USER"`
	DB_PASS    string `json: "DB_PASS"`
	DB_HOST    string `json: "DB_HOST"`
	DB_SCHEMA  string `json: "DB_SCHEMA"`
	DB_OPTIONS string `json: "DB_OPTIONS"`
}

// end of config

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
type FileInfo struct {
	fileName string
	dirPath  string
	fileHash uint32
	dupCnt   int
	totDup   int
	atime    string
	mtime    string
	ctime    string
	dModel   string
}
type ImageInfo struct {
	fileID   int
	fileHash uint32
	comments string
	a2d      get.AppliedToDimensions
}
type FaceInfo struct {
	fileID     int
	fileHash   uint32
	regionList get.RegionList
}
type DestInfo struct {
	fileHash uint32
	dupCnt   int
	totDup   int
	fileName string
	dirPath  string
	atime    string
	mtime    string
	ctime    string
}
type Report struct {
	tFiles int
	tBytes int
	eFiles int
	eBytes int
	uFiles int
	uBytes int
}
type Env struct {
	//dbconn *sql.DB
	dbconn   dbl.MediaDBHandler
	fHash    uint32
	srcID    int64
	dupeID   int64
	tgtID    int64
	thnailID int64
	faceID   int64
	Report   Report
	Config   *Config
}
type Fps struct {
	j   int
	fps []string
}
type FileTime struct {
	stm     string
	ttm     time.Time
	sepoc   string
	tepoc   time.Time
	sorigin string
	torigin time.Time
	screate string
	smodify string
	saccess string
	tcreate time.Time
	tmodify time.Time
	taccess time.Time
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
}

func main() {
	// capture the time of start of process
	t0 := time.Now()

	env := &Env{}
	config := env.Config

	config = env.loadConfigJson("config/albumMgmt.json")

	log.Infof("Configuration read from JSON file -- conifg:: %v\n", config)

	log.Infof("Env at start of main: %v\n", env)

	// Get the commandline parameters

	flag.BoolVar(&justCopy, "cp", config.JustCopy, "Copy, don't move. Default is Move") // Copy files instead of move. Default is move
	flag.BoolVar(&dryRun, "dr", config.DryRun, "Dryrun")                                //Dryrun: only simulate copy or move.
	flag.BoolVar(&verbose, "v", config.Verbose, "verbose summary statistics")           // VERBOSE causes detailed summary statistics to be logged.
	flag.StringVar(&srcPath, "src", config.SrcPath, "source file with fully qualified path")
	flag.StringVar(&destPath, "dest", config.DestPath, "Destination fully qualified path")

	// Process the commandline parameters
	lock.Lock()
	flag.Parse()
	// make sure src and dest is provided
	if srcPath == "" || destPath == "" {
		fmt.Printf(oUsage, os.Args[0])
		os.Exit(0)
	}

	env.Config.DryRun = dryRun
	env.Config.Verbose = verbose
	env.Config.JustCopy = justCopy
	env.Config.SrcPath = srcPath
	env.Config.DestPath = destPath

	lock.Unlock()

	// start logging to a file
	//logf, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logf, err := os.OpenFile(config.Logfile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err == nil {
		log.SetOutput(logf)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	defer logf.Close()

	log.Infoln("Log started")
	// logging setup complete

	// Initiate DB connection
	lock.Lock()
	dbconn, err := dbl.GetDatabaseHandler(dbl.MYSQL, config.DB_USER+":"+config.DB_PASS+"@"+config.DB_HOST+"/"+config.DB_SCHEMA+"?"+config.DB_OPTIONS)
	//dbconn, err := dbl.GetDatabaseHandler(dbl.MYSQL, "mediadb:mediaDB@tcp(192.168.10.190:3306)/mediadb?parseTime=true")
	//dbconn, err := dbl.GetDatabaseHandler(dbl.SQLITE, "mediaDB.db")
	//dbconn, err := dbl.GetDatabaseHandler(dbl.POSTGRESQL, "user=postgres dbname=dino sslmode=disable")
	//dbconn, err := dbl.GetDatabaseHandler(dbl.MONGODB, "mongodb://127.0.0.1")
	if err != nil {
		log.Fatalf("Errror getting database connection. err: ", err)
	}
	//send back to the receiver method Env

	env.dbconn = dbconn
	lock.Unlock()

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
	wg.Wait()

	lock.Lock()
	env.Report.tFiles = TFiles
	env.Report.tBytes = TBytes
	lock.Unlock()
	log.Infof("report %v\n", env)

	//process each file in the hashMap
	err = env.processfileList(hashMap)
	if err != nil {
		log.Errorf("error processing files. err: %v\n", err)
		os.Exit(0)
	}
	// get end time
	t1 := time.Now()
	// get execution end time
	𝛥t := float64(t1.Sub(t0)) / 1e9

	// print optional verbose summary report <<to-do>> need to get this report done right --separat into function
	if verbose {
		tFiles := env.Report.tFiles
		tBytes := env.Report.tBytes
		log.Infof("     total: %8d files (%7.2f%%), %13d bytes (%7.2f%%)\n",
			tFiles, 100.0, tBytes, 100.0)
		log.Infof("  examined: %8d files, %13d bytes in %.4f seconds\n",
			env.Report.eFiles, env.Report.eBytes, 𝛥t)
		log.Infof("  Unique Files: %8d files, in %.4f seconds\n",
			len(hashMap), 𝛥t)
		// print to terminal
		fmt.Printf("     total: %8d files (%7.2f%%), %13d bytes (%7.2f%%)\n",
			tFiles, 100.0, tBytes, 100.0)
		fmt.Printf("  examined: %8d files, %13d bytes in %.4f seconds\n",
			env.Report.eFiles, env.Report.eBytes, 𝛥t)
		fmt.Printf("  Unique Files: %8d files, in %.4f seconds\n",
			len(hashMap), 𝛥t)

	}
	// Add session info to the sessionInfo table in db
	sessionInfo := dbl.SessionInfo{
		//Sessionid		,									//int       `bson:"-"`				//int(11) NOT NULL AUTO_INCREMENT,
		Starttime:     t0,                //time.Time `bson:"start_time"`		//timestamp,
		Stoptime:      t1,                //time.Time `bson:"stop_time"`		//timestamp,
		TotalFiles:    env.Report.tFiles, //int64   	`bson:"total_files"`	//int,
		TotalBytes:    env.Report.tBytes, //int64   	`bson:"total_bytes"`	//int,
		ExaminedFiles: env.Report.eFiles, //int64   	`bson:"examined_files"`	//int,
		ExaminedBytes: env.Report.eBytes, //int64   	`bson:"examined_bytes"`	//int,
		UniqueFiles:   len(hashMap),      //int64   	`bson:"unique_files"`	//int,
		UniqueBytes:   0,                 //int64   	`bson:"unique_bytes"`	//int,
		RunTime:       𝛥t,                //float64   `bson:"run_time"`		//float,
		LogFile:       config.Logfile,    //string    `bson:"log_file"`		//varchar(254),
		Rowaction:     "A",               //string    `bson:"row_action"`		// char(2) NOT NULL,
		//Rowactiondatetime		, 							//time.Time `bson:"row_ts"`			//timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	}
	sirowid, err := dbconn.AddSessionInfo(sessionInfo)
	if err != nil {
		log.Errorf("Error adding file record to database: %v\t error: %v \n", sessionInfo, err)
	}
	log.Infoln("Session info Table Row ID: ", sirowid)
	// end of SessionInfo db insert
	log.Infoln("SessionInfo: %v\n firowid: %v\n", sessionInfo, sirowid)

} //Close main

func (env *Env) processfileList(hm map[uint32][]string) (err error) {

	var Err error
	wg := new(sync.WaitGroup)

	// process each file in the hash map
	for fHash, fps := range hm {

		log.Infof(" file Hash: %v\t dupes: %v\t g: %v\n", fHash, len(fps), fps)

		env.fHash = fHash

		log.Infof("Env in ProcessfileList: %v\n", env)

		// process each file with same hash

		wg.Add(1)
		chRetData := make(chan error)
		defer close(chRetData)

		go func(fps []string) {
			defer wg.Done()
			err = env.processEachFile(fps)
			if err != nil {
				log.Errorf("error processing files. err: %v\n", err)
			}
			chRetData <- err

		}(fps)

		Err = <-chRetData
		//Err = err

	} // close map iterator loop "fHash"
	log.Errorf("ProcessFileList error: %v\n", Err)
	wg.Wait()
	return Err
} // close processfileList

func (env *Env) processEachFile(fps []string) (err error) {

	var Err error

	for j := range fps {

		err := env.manageFileinfo(j, fps)
		if err != nil {
			log.Errorf("manageFileInfo error. err: %v\n", err)
		}

		Err = err

	} // close map iterator loop "j"
	log.Errorf("ProcessEachFile error: %v\n", Err)
	return Err
} // close processEachFile

func (env *Env) manageFileinfo(j int, fps []string) (err error) {
	dbconn := env.dbconn
	fHash := env.fHash

	//Get file EXIF information
	exfe, err := fs.StatExifStruct(fps[j])
	if err != nil {
		log.Errorf("Exif error: %v\n", err)
	}
	log.Infof("EXIF: %v\n", exfe)

	// find the date and time of the file from exif info
	ftm, err := env.time4file(exfe)
	if err != nil {
		log.Errorf("time4file: %v\n", err)
	}

	// Get the file stats
	fi, err := os.Stat(fps[j])
	if err != nil {
		log.Errorf("os.filestat error: %v\n", err)
	}
	fSize := fi.Size()
	//fModDate := ftm.tmodify	//fi.ModTime().Format(env.Config.DtFormWoTz)
	fName := filepath.Base(fps[j])
	ext := strings.ToLower(filepath.Ext(fps[j]))
	dir := filepath.Dir(fps[j])

	/* // Create the source file information object
	fileInfo := &FileInfo{
		fileName: fps[j],
		fileHash: fHash,
		dirPath:  dir,
		dupCnt:   j,
		totDup:   len(fps) - 1,
		atime:    ftm.saccess,
		mtime:    ftm.smodify,
		ctime:    ftm.screate,
		dModel:   exfe.Model,
	}

	log.Infof("fileInfo: %v\n", fileInfo) */
	// create source file data object
	fileList := dbl.Filelist{
		//	Fileid:            			,//int			`bson:"file_id"`
		Filename:     fName,             //string		`bson:"file_name"`
		Dupecount:    0,                 //int			`bson:"dupe_count"`
		Filesuffix:   ext,               //string		`bson:"file_suffix"`
		Filelocation: dir,               //string		`bson:"file_loc"`
		Filesize:     fSize,             //int64  		`bson:"file_siz"`
		Filehash:     fmt.Sprint(fHash), //string		`bson:"file_hash"`
		Filedate:     ftm.ttm,           //time.Time	`bson:"file_date"`
		Rowaction:    "A",               //string		`bson:"row_action"`
		//Rowactiondatetime: t0,         //time.Time	`bson:"row_ts"`
	}

	// Generate target file information
	tgtfileName := fmt.Sprint(fHash) + "_" + exfe.Model + "_" + strconv.Itoa(dupcount) + ext
	// Generate the output folder name
	tgtfolder := ftm.ttm.Format(env.Config.FolderFormat)
	log.Infof("fileName: %v\n", tgtfileName)
	tgtPath := filepath.Join(destPath, dupePath, tgtfolder) //, tgtfileName) //ToDo: Move this to a function

	// create dest file data object
	tgtInfo := dbl.TgtInfoList{
		//Targetid:  ,  								//int       `bson:"-"`           //","int(11)","NO","PRI","","auto_increment"
		Fileid:       int(env.srcID),    //int       `bson:"file_id"`     //","int(11)","NO","MUL","",""
		Filename:     tgtfileName,       //string    `bson:"file_name"`   //","varchar(254)","NO","","",""
		Filesuffix:   ext,               //string    `bson:"file_suffix"` //","varchar(5)","YES","","",""
		Filelocation: tgtPath,           //string    `bson:"file_loc"`    //","varchar(254)","YES","","",""
		Filehash:     fmt.Sprint(fHash), //string    `bson:"file_hash"`   //","varchar(254)","YES","","",""
		Filesize:     float32(fSize),    //float32   `bson:"file_siz"`    //","float","YES","","",""
		Filedate:     ftm.ttm,           //time.Time `bson:"file_date"`   //","datetime","YES","","",""
		Fileaction:   "T",               //T = To-Do	//string    `bson:"file_action"` //","char(2)","NO","","",""
		Validated:    0,                 //int       `bson:"validated"`   //","tinyint(1)","NO","","0",""
		Rowaction:    "A",               //string    `bson:"row_action"`  //","char(2)","NO","","",""
		//Rowactiondatetime:  ,  						//time.Time `bson:"row_ts"`      //","datetime","NO","","CURRENT_TIMESTAMP",""
	}
	// check if file already exists, ignore it
	dupePath = "" //initially set the dupe path to blank.

	//<<Placeholder call env.chk4dupe

	//process further only if file does not exist
	if !fileIDnCount.FileID.Valid {
		// collect metrics for reporting
		//var lock sync.Mutex // moved to global scope
		lock.Lock()
		EFiles := env.Report.eFiles
		EBytes := env.Report.eBytes

		EFiles++
		EBytes += int(fSize)

		env.Report.eFiles = EFiles
		env.Report.eBytes = EBytes
		lock.Unlock()

		// Add file info into database fileinfo table
		if !dryRun { // add file to db table only is
			//1> dryrun is false
			//2> DB record does not exist <<TO_DO>>
			//3> target file does not exist <<TO_DO>>

			srcFilerowid, err = dbconn.AddFile(fileList)
			if err != nil {
				log.Errorf("Error adding file record to database: %v\t error: %v \n", fileList, err)
			}
			log.Infoln("Source file Row ID: ", srcFilerowid)
			env.srcID = srcFilerowid

			log.Infof("Env after addFilelist: %v\n", env)
			// end of filelist db insert

		}

		// Get face data and coordinates, <<to-do>> update the Album database
		rI := exfe.RegionInfo

		imageInfo := dbl.ImageInfo{
			Fileid:    env.srcID,                   //int       `bson:"file_id"`		// int(11),
			Filehash:  fmt.Sprint(fHash),           //string    `bson:"file_hash"`		// varchar(254) DEFAULT NULL,
			DimH:      rI.AppliedToDimensions.H,    //int   	`bson:"image_dimH"`		// float DEFAULT NULL,
			DimW:      rI.AppliedToDimensions.W,    //int   	`bson:"image_dimW"`		// float DEFAULT NULL,
			DimUnit:   rI.AppliedToDimensions.Unit, //string    `bson:"image_dimunits"` // varchar(5) DEFAULT NULL,
			Rowaction: "A",                         //string    `bson:"row_action"`		// char(2) NOT NULL,

		}
		log.Infof("Image Info and Range Dimensions: %v\n", imageInfo)

		if !dryRun { // add info to db table only is dryrun is false
			//go func(imageInfo dbl.ImageInfo) {
			irowid, err := dbconn.AddImageInfo(imageInfo)
			if err != nil {
				log.Errorf("Error adding file record to database: %v\t error: %v \n", imageInfo, err)
			}
			log.Println("Face info Table Row ID: ", irowid)
			// end of addRegionInfo db insert
			log.Infof("faceInfo: %v\n firowid: %v\n", imageInfo, irowid)
			//}(imageInfo)
		}

		// Add faceinfo to database
		// Get face data and coordinates, <<to-do>> update the Album database
		rIL := rI.RegionList
		var fcI FaceInfo
		//if there is REgion information in EXIF, process the list
		if len(rIL) > 0 {
			for r := 0; r < len(rIL); r++ {

				// Create the faceinfo object
				faceInfo := dbl.RegionInfo{
					//Regionid:				, 				//	int       `bson:"-"`				//` int(11) NOT NULL AUTO_INCREMENT,
					Fileid:    env.srcID,         //	int       `bson:"file_id"`		// int(11),
					Filehash:  fmt.Sprint(fHash), //	string    `bson:"file_hash"`	// varchar(254) DEFAULT NULL,
					Name:      rIL[r].Name,       //	string    `bson:"region_name"`// varchar(254) NOT NULL,
					Typ:       rIL[r].Type,       //	string    `bson:"region_type"`// varchar(5) DEFAULT NULL,
					AreaH:     rIL[r].Area.H,     //	float32   `bson:"area_H"`		// float DEFAULT NULL,
					AreaW:     rIL[r].Area.W,     //	float32   `bson:"area_W"`		// float DEFAULT NULL,
					AreaX:     rIL[r].Area.X,     //	float32   `bson:"area_X"`		// float DEFAULT NULL,
					AreaY:     rIL[r].Area.Y,     // float32   `bson:"area_Y"`		//` float DEFAULT NULL,
					AreaUnit:  rIL[r].Area.Unit,  //	string    `bson:"area_units"` // varchar(5) DEFAULT NULL,
					Rowaction: "A",               //	string    `bson:"row_action"`	// char(2) NOT NULL,
				}
				if !dryRun { // add info to db table only is dryrun is false

					firowid, err := dbconn.AddRegionInfo(faceInfo)
					if err != nil {
						log.Errorf("Error adding file record to database: %v\t error: %v \n", faceInfo, err)
					}
					//log.Infoln("Face info Table Row ID: ", firowid)
					// end of addRegionInfo db insert
					log.Infof("faceInfo: %v\n firowid: %v\n", faceInfo, firowid)

					env.faceID = firowid
					log.Infof("Env after addRegionInfo: %v\n", env)
				}
			}
		} else {
			log.Infof("EXIF RegionInfo-RegionList missin")
		}
		log.Infof("Face Info: %v\n", fcI)

		if j > 0 { // If file is a duplicate
			// and if it doesnt already exist at target <<TO-DO>>

			// Add dupe file info into database dupes table
			dupelist := dbl.DupeList{
				//Dupeid:  ,  						  //int       `bson:"-"`           //,"int(11)","NO","PRI","","auto_increment"
				Fileid:     int(env.srcID),     //int       `bson:"file_id"`     //,"int(11)","NO","MUL","",""
				Dupefileid: fileminmaxid.MinID, //int       `bson:"dupefile_id"` //,"int(11)","NO","","",
				Rowaction:  "A",                //string    `bson:"row_action"`  //,"char(2)","NO","","",""
				//Rowactiondatetime:  ,  		  //time.Time `bson:"row_ts"`      //,"datetime","NO","","CURRENT_TIMESTAMP",""
			}
			if !dryRun { // add file to db table only is dryrun is false

				duperowid, err := dbconn.AddDupes(dupelist)
				if err != nil {
					log.Errorf("Error adding file record to database: %v\t error: %v \n", dupelist, err)
				}

				log.Println("Dupe file Row ID: ", duperowid)

				env.dupeID = duperowid
				log.Infof("Env after addDupes: %v\n", env)
			}
		}
		// get gps iformation from EXIF
		gps := exfe //.GPSInfo

		//parse GPS DateTime field
		gpsdtm, err := time.Parse(DT_GPS, gps.GPSDateTime)
		if err != nil {
			log.Errorf("timeconv error GPSDateTime: %v\t Form: %v\t Err: %v\n", gps.GPSDateTime, DT_GPS, err)
		}

		// Creat  exifinfo data object
		exiflist := dbl.EXIFInfoList{
			//Exifid:           ,						//int       `bson:"-"`           //","int(11)","NO","PRI","","auto_increment"
			//Rowactiondatetime: 						//time.Time `bson:"row_ts"`      //","datetime","NO","","CURRENT_TIMESTAMP",""
			Accessdate:          ftm.taccess,             //time.Time `bson:"access_time"` //","datetime","YES","","",""
			Comments:            exfe.UserComment,        //string    `bson:"comments"`    //","varchar(254)","YES","","",""
			Createdate:          ftm.tcreate,             //time.Time `bson:"create_time"` //","datetime","YES","","",""
			Facecoords:          "FaceCoords",            //string    `bson:"face_coords"` //","varchar(254)","YES","","",""
			Faces:               "Faces",                 //string    `bson:"faces"`       //","varchar(254)","YES","","",""
			Fileid:              int(env.srcID),          //int       `bson:"file_id"`     //","int(11)","NO","MUL","",""
			GPSAltitude:         gps.GPSAltitude,         //string   `bson:"gps_altitude"`          //` varchar(254) DEFAULT NULL,
			GPSAltitudeRef:      gps.GPSAltitudeRef,      //string   `bson:"gps_altitude_ref"`      //` varchar(254) DEFAULT NULL,
			GPSDateTime:         gpsdtm,                  //string 	 `bson:"gps_datetime"`          //` datetime,
			GPSLatitude:         gps.GPSLatitude,         //string   `bson:"gps_latitude"`          //` varchar(254) DEFAULT NULL,
			GPSLatitudeRef:      gps.GPSLatitudeRef,      //string    `bson:"gps_latitude_ref"`      //` varchar(254) DEFAULT NULL,
			GPSLongitude:        gps.GPSLongitude,        //int       `bson:"gps_longitude"`         //` varchar(254) DEFAULT NULL,
			GPSLongitudeRef:     gps.GPSLongitudeRef,     //string    `bson:"gps_longitude_ref"`     //` varchar(254) DEFAULT NULL,
			GPSMapDatum:         gps.GPSMapDatum,         //string    `bson:"gps_map_datum"`         //` varchar(254) DEFAULT NULL,
			GPSProcessingMethod: gps.GPSProcessingMethod, // string    `bson:"gps_processing_method"` //` varchar(254) DEFAULT NULL,
			GPSVersionID:        gps.GPSVersionID,        //string    `bson:"gps_version_id"`        //` varchar(254) DEFAULT NULL,
			ImageDescription:    exfe.ImageDescription,   //string    `bson:"image_description"`     //` varchar(254) DEFAULT NULL,
			Make:                exfe.Make,               //string    `bson:"make"`                  //` varchar(254) DEFAULT NULL,
			Model:               exfe.Model,              //string    `bson:"model"`                 //` varchar(254) DEFAULT NULL,
			Modifydate:          ftm.tmodify,             //time.Time `bson:"mod_time"`    //","datetime","YES","","",""
			Orientation:         exfe.Orientation,        //int       `bson:"orientation"`           //` int DEFAULT NULL,
			Rowaction:           "A",                     //string    `bson:"row_action"`  //","char(2)","NO","","",""
			Tags:                "Tags",                  //string    `bson:"tags"`        //","varchar(254)","YES","","",""
		}

		/* gpslist := dbl.EXIFInfoList.GPSInfo{

			GPSAltitude:         gps.GPSAltitude,         //float32   `bson:"gps_altitude"`          //` float DEFAULT NULL,
			GPSAltitudeRef:      gps.GPSAltitudeRef,      //float32   `bson:"gps_altitude_ref"`      //` float DEFAULT NULL,
			GPSLatitude:         gps.GPSLatitude,         //float32   `bson:"gps_latitude"`          //` float DEFAULT NULL,
			GPSLatitudeRef:      gps.GPSLatitudeRef,      //string    `bson:"gps_latitude_ref"`      //` varchar(254) DEFAULT NULL,
			GPSLongitude:        gps.GPSLongitude,        //int       `bson:"gps_longitude"`         //` int DEFAULT NULL,
			GPSLongitudeRef:     gps.GPSLongitudeRef,     //string    `bson:"gps_longitude_ref"`     //` varchar(254) DEFAULT NULL,
			GPSMapDatum:         gps.GPSMapDatum,         //string    `bson:"gps_map_datum"`         //` varchar(254) DEFAULT NULL,
			GPSProcessingMethod: gps.GPSProcessingMethod, // string   `bson:"gps_processing_method"` //` varchar(254) DEFAULT NULL,
			GPSDateTime:         gps.GPSDateTime,         //string	 `bson:"gps_datetime"`          //` datetime,
			GPSVersionID:        gps.GPSVersionID,        //string    `bson:"gps_version_id"`        //` varchar(254) DEFAULT NULL,

		}
		*/
		// Add EXIF Info into database exifinfo table
		if !dryRun { // add file to db table only is dryrun is false
			//go func(exiflist dbl.EXIFInfoList) {
			exifrowid, err := dbconn.AddEXIF(exiflist)
			if err != nil {
				log.Errorf("Error adding file record to database: %v\t error: %v \n", exiflist, err)
			}
			log.Println("EXIF Table Row ID: ", exifrowid)
			// end of exifInfolist db insert
			//}(exiflist)
		}
		// Add Target  info into database targetinfo table
		if !dryRun { // add file to db table only is dryrun is false

			tgtrowid, err = dbconn.AddTgtInfo(tgtInfo)
			if err != nil {
				log.Errorf("Error adding file record to database: %v\t error: %v \n", tgtInfo, err)
			}
			log.Infof("Targetinfo Table Row ID: %v\n", tgtrowid)
			// end of tgtInfolist db insert

			env.tgtID = tgtrowid
			log.Infof("ChRetData after addTgtInfo: %v\n", env.tgtID)
			log.Infof("Env after addTgtInfo: %v\n", env)
		}

		//Only generate thumbnails for original files, not dupes
		if (dupcount <= 0) && (!dryRun) {
			// check if its an Image file
			imageTypes := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp"}
			if env.stringInSlice(strings.ToLower(ext), imageTypes) {

				//generate thumbnails
				thSizes := make(map[string]uint)

				thSizes["2s"] = uint((exfe.ImageHeight / 100) * 12) //240 //240x180  //orig:2048x1536	12%
				thSizes["me"] = uint((exfe.ImageHeight / 100) * 39) //792x594		39%
				thSizes["th"] = uint((exfe.ImageHeight / 100) * 7)  //144x108		7%
				thSizes["sm"] = uint((exfe.ImageHeight / 100) * 28) //576x432		28%
				thSizes["sq"] = 120                                 //120x120 	6%,8%

				wg := new(sync.WaitGroup)
				for thKey, thVal := range thSizes {
					log.Infoln("Thumbnail Tag:", thKey, "Thumbnail Size:", thVal)

					thImgName := fmt.Sprint(fHash) + "_" + exfe.Model + "_" + strconv.Itoa(dupcount) + "_" + thKey + ext
					thPath := filepath.Join(destPath, "thumbnails", tgtfolder, tgtfileName)

					log.Infof("ThImgName: %v\t ThSize: %v\t ThPath: %v\n", thImgName, thVal, thPath)

					wg.Add(1)
					go func(imageFile string, thImgName string, thPath string, thVal uint) {
						defer wg.Done()
						err = th.env.generateThumbImg(imageFile, thImgName, thPath, thVal)
						if err != nil {
							log.Errorf("Error generating thumbnail: %v\t error: %v \n", thImgName, err)
						}
					}(fps[j], thImgName, thPath, thVal)

					// Insert thumbnail info into table thumbnails
					thlist := dbl.ThumbList{
						//Thid:  ,  						 //int       `bson:"-"`           //","int(11)","NO","PRI","","auto_increment"
						Fileid:     int(env.srcID),    //int       `bson:"file_id"`     //","int(11)","NO","MUL","",""
						Targetid:   int(env.tgtID),    //int       `bson:"tgt_id"`      //,"int(11)","NO","MUL","",""
						Filename:   thImgName,         //string    `bson:"file_name"`   //","varchar(254)","NO","","",""
						Filesuffix: ext,               //string    `bson:"file_suffix"` //","varchar(254)","YES","","",""
						Filesize:   exfe.ImageSize,    //string    `bson:"file_size"`   //","varchar(254)","YES","","",""
						Thsize:     fmt.Sprint(thVal), //string    `bson:"th_size"`     //","varchar(254)","NO","","",""
						Fileloc:    thPath,            //string    `bson:"file_loc"`    //","varchar(254)","YES","","",""
						Rowaction:  "A",               //string    `bson:"row_action"`  //","char(2)","NO","","",""
						//Rowactiondatetime:  ,  		 //time.Time `bson:"row_ts"`      //","datetime","NO","","CURRENT_TIMESTAMP",""
					}
					if !dryRun { // add file to db table only is dryrun is false
						wg.Add(1)
						go func(thlist dbl.ThumbList) {
							defer wg.Done()
							throwid, err := dbconn.AddThumb(thlist)
							if err != nil {
								log.Errorf("Error adding file record to database: %v\t error: %v \n", thlist, err)
							}
							log.Infoln("thumbnails Table Row ID: ", throwid)
							// end of thumbnails db insert
						}(thlist)
					}
				}
				wg.Wait()
			}
		}
		if !dryRun { // Copy or move only if dryrun is false
			// Send source and dest info to copy function
			srcFile := filepath.Join(dir, fileList.Filename)
			destFile := filepath.Join(tgtPath, tgtInfo.Filename)
			wg := new(sync.WaitGroup)
			log.Infof("for Copy or Move: Source file: %v\n Target file: %v\n", srcFile, destFile)
			wg.Add(1)
			go func(destFile string, srcFile string) {
				defer wg.Done()
				if justCopy {
					fileValid, bytesCopied, err := fc.CopyFile(destFile, srcFile)
					if err != nil {
						log.Errorln(err)
					}
					log.Infof("Is file valid: %v\t Bytes Copied: %v\t err: %v\n", fileValid, bytesCopied, err)
				} else { //if justCopy is false than move the file
					fileValid, err := fm.MoveFile(destFile, srcFile)
					if err != nil {
						log.Errorln(err)
					}
					log.Infof("Is file valid:: %v\t  err: %v\n", fileValid, err)
				}
			}(destFile, srcFile)
			//<-fileDone
			wg.Wait()
		}

		//ToDo: update file exif with original file name and folder name.
		//ToDo: set file owner and permissions
		//ToDo: generate the final statistics
	} // file exist check loop end
	log.Warnf("Existing file ID:%v\t Not Null:%v\t fileCount:%v\t err:%v\n", fileIDnCount.FileID.Int64, fileIDnCount.FileID.Valid, fileIDnCount.FileCount, err)
	return err

} // close manageFileinfo

// Helper Functions

func (env *Env) time4file(exf get.EXIF /*map[string]interface{}*/) (ftm *FileTime, err error) {
	ftm = &FileTime{}

	epoctime := "0001:01:01 00:00:00-00:00"
	etm, err := time.Parse(env.Config.DtFormWTz, epoctime)
	if err != nil {
		log.Errorf("timeconv error epoctime: %v\t Form: %v\t Err: %v\n", epoctime, env.Config.DtFormWoTz, err)
	}
	ftm.sepoc = epoctime
	ftm.tepoc = etm

	ectime := exf.CreateDate
	ectm, err := time.Parse(env.Config.DtFormWoTz, ectime)
	if err != nil {
		log.Errorf("timeconv error ctime: %v\t Form: %v\t Err: %v\n", ectime, env.Config.DtFormWoTz, err)
	}

	fctime := exf.FileCreateDate
	fctm, err := time.Parse(env.Config.DtFormWTz, fctime)
	if err != nil {
		log.Errorf("timeconv error fctime: %v\t Form: %v\t Err: %v\n", fctime, env.Config.DtFormWTz, err)
	}

	//Check which time EXIF or Filestat time exists
	// if both exist, use which ever is earlier
	ctm := ectm
	ctime := ectime
	if ctime == "" || ctm.IsZero() || ctm.After(fctm) {
		ctm = fctm
		ctime = fctime
	}

	ftm.screate = ctime
	ftm.tcreate = ctm

	emtime := exf.ModifyDate
	emtm, err := time.Parse(env.Config.DtFormWoTz, emtime)
	if err != nil {
		log.Errorf("timeconv error emtime: %v\t Form: %v\t Err: %v\n", emtime, env.Config.DtFormWoTz, err)
	}

	fmtime := exf.FileModifyDate
	fmtm, err := time.Parse(env.Config.DtFormWTz, fmtime)
	if err != nil {
		log.Errorf("timeconv error fmtime: %v\t Form: %v\t Err: %v\n", fmtime, env.Config.DtFormWTz, err)
	}

	//Check which time EXIF or Filestat time exists
	// if both exist, use which ever is later
	mtm := emtm
	mtime := emtime
	if mtime == "" || mtm.IsZero() || mtm.Before(fmtm) {
		mtm = fmtm
		mtime = fmtime
	}

	ftm.smodify = mtime
	ftm.tmodify = mtm

	eotime := exf.DateTimeOriginal
	eotm, err := time.Parse(env.Config.DtFormWoTz, eotime)
	if err != nil {
		log.Errorf("timeconv error otime: %v\t Form: %v\t Err: %v\n", eotime, env.Config.DtFormWoTz, err)
	}
	ftm.sorigin = eotime
	ftm.torigin = eotm

	fatime := exf.FileAccessDate
	fatm, err := time.Parse(env.Config.DtFormWTz, fatime)
	if err != nil {
		log.Errorf("timeconv error atime: %v\t Form: %v\t Err: %v\n", fatime, env.Config.DtFormWTz, err)
	}
	ftm.saccess = fatime
	ftm.taccess = fatm

	tval := eotm
	sval := eotime
	if sval == "" || tval.IsZero() {
		tval = ectm
		sval = ectime
		if sval == "" || tval.IsZero() {
			tval = emtm
			sval = emtime
			if sval == "" || tval.IsZero() {
				tval = fctm
				sval = fctime
				if sval == "" || tval.IsZero() {
					tval = fmtm
					sval = fmtime
					if sval == "" || tval.IsZero() {
						tval = fatm
						sval = fatime
						if sval == "" || tval.IsZero() {
							tval = etm
							sval = epoctime
						}
					}
				}
			}
		}
	}
	log.Infof("DateVal: %v\n", tval)
	ftm.ttm = tval
	return ftm, err
}

func (env *Env) decodeConfig(filename string) (image.Config, string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return image.Config{}, "", err
	}
	defer f.Close()
	return image.DecodeConfig(bufio.NewReader(f))
}

func (env *Env) SetField(obj interface{}, name string, value interface{}) error {

	structValue := reflect.ValueOf(obj).Elem()
	fieldVal := structValue.FieldByName(name)

	if !fieldVal.IsValid() {
		return fmt.Errorf("No such field: %s in obj", name)
	}

	if !fieldVal.CanSet() {
		return fmt.Errorf("Cannot set %s field value", name)
	}

	val := reflect.ValueOf(value)

	if fieldVal.Type() != val.Type() {

		if m, ok := value.(map[string]interface{}); ok {

			// if field value is struct
			if fieldVal.Kind() == reflect.Struct {
				return env.FillStruct(m, fieldVal.Addr().Interface())
			}

			// if field value is a pointer to struct
			if fieldVal.Kind() == reflect.Ptr && fieldVal.Type().Elem().Kind() == reflect.Struct {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				// fmt.Printf("recursive: %v %v\n", m,fieldVal.Interface())
				return env.FillStruct(m, fieldVal.Interface())
			}
		}

		return fmt.Errorf("Provided value type didn't match obj field type")
	}

	fieldVal.Set(val)
	return nil
}

func (env *Env) FillStruct(m map[string]interface{}, s interface{}) error {
	for k, v := range m {
		err := env.SetField(s, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *Env) saveConfigJson(c Config, filename string) error {
	bytes, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func (env *Env) loadConfigJson(file string) *Config {
	var config Config
	configFile, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}
	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("Config Parse Error: ", err)
	}

	env.Config = &config

	//jsonParser := json.NewDecoder(configFile)
	//jsonParser.Decode(&config)
	return &config
}

func (env *Env) saveConfigYaml(c Config, filename string) error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, bytes, 0644)
}

func (env *Env) loadConfigYaml(filename string) *Config { //(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	var config Config
	err = yaml.Unmarshal(bytes, &config)
	if err != nil {
		log.Fatal("Config File Missing. ", err)
	}

	return &config
}
