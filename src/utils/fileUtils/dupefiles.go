package fileUtils

import "fmt"

type Dupe struct {
	dupecount int
	dupeFile  string
	dupeRec   string
	dupeOrig  string
}

func (env *Env) chk4dupe(file string, exf get.EXIF) (dup *Dupe, err error) {
	dup = &Dupe{}

	// check if physical target file exists -- targetpath, dupePath -- dupeFile  string
	// check if DB file record exists -- filelist -- dupeRec   string
	//Check if this file entry already exists in db

	//fileIDnCount := dbl.FileIDnCount{}
	fileIDnCount, err := dbconn.GetFileID(fileList)
	if err != nil {
		log.Errorf("Error getting existing file id  -- err: %v\n", err)
		//err = errors.New("skip: error checking if file already exists") //<to-do> check if file also exists on filesystem tgt folder
		return nil, err
	}

	// check if filehash already exists in DB
	fileminmaxid := dbl.FileMinMaxID{}
	fileminmaxid, err = dbconn.GetMinMaxID(fmt.Sprint(fHash))
	dupcount := 0
	if err != nil {
		log.Errorf("Error getting min max id: %v\t error: %v\t dupoffset %v\n", fileminmaxid, err, dupcount)
	} else {
		// file exists, increment duplicate count accordingly
		dupcount = fileminmaxid.MaxDup + 1

		//change file path based on if file is dupe
		dupePath = env.Config.DupePath

		log.Warnf("File is duplicate - Fileminmaxid: %v\t  dupoffset: %v\n", fileminmaxid, dupcount)
	}

	log.Infof("min max ID: %v \n ", fileminmaxid)

	// get fileID of original file from DB -- What criterion? -- dupeOrig  string
	//	if multiple matches, get count of all matches -- dupecount int

	//dupecount int
	//dupeFile  string
	//dupeRec   string
	//dupeOrig  string

	return dup, err
}
