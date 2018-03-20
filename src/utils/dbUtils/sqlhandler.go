package dbUtils

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
)

type SQLHandler struct {
	*sql.DB
}

func (handler *SQLHandler) GetAvailableFiles() ([]Filelist, error) {
	return handler.sendQueryFilelist("select * from filelist")
}

func (handler *SQLHandler) GetFileByName(filename string) (Filelist, error) {

	sqlStr := fmt.Sprintf("select * from filelist where filename = '%s'", filename)
	row := handler.QueryRow(sqlStr) //? for mysql or sqlite and it used to be $1 for pq
	f := Filelist{}
	err := row.Scan(&f.Fileid, &f.Filename, &f.Filesuffix, &f.Filelocation, &f.Filesize, &f.Filehash, &f.Filedate, &f.Rowaction, &f.Rowactiondatetime)
	return f, err
}

func (handler *SQLHandler) GetMinMaxID(filehash string) (FileMinMaxID, error) {

	sqlStr := fmt.Sprintf("SELECT MIN(fl.fileid) as minid, MAX(fl.fileid) as maxid , MIN(ei.createdate) AS mindt, MAX(ei.createdate) AS maxdt, MIN(fl.dupecount) as mindup, MAX(fl.dupecount) as maxdup FROM filelist fl, exifinfo ei WHERE   fl.fileid = ei.fileid and fl.filehash = '%s'", filehash)
	row := handler.QueryRow(sqlStr) //? for mysql or sqlite and it used to be $1 for pq
	fmm := FileMinMaxID{}
	err := row.Scan(&fmm.MinID, &fmm.MaxID, &fmm.MinDate, &fmm.MaxDate, &fmm.MinDup, &fmm.MaxDup)
	return fmm, err
}

func (handler *SQLHandler) GetFileID(f Filelist) (FileIDnCount, error) {
	sqlStr := fmt.Sprintf("SELECT fl.fileid as fileid, Count(fl.fileid) as filecount FROM filelist fl WHERE fl.Filename = '%v' AND  fl.Filesuffix = '%v' AND fl.Filelocation = '%v' AND fl.Filesize = '%v' AND fl.Filehash = '%v' AND fl.Filedate = '%v'", f.Filename, f.Filesuffix, f.Filelocation, f.Filesize, f.Filehash, f.Filedate)

	log.Println(sqlStr) //? for mysql or sqlite and it used to be $1 for pq

	row := handler.QueryRow(sqlStr)

	fic := FileIDnCount{}
	err := row.Scan(&fic.FileID, &fic.FileCount)
	return fic, err
}

func (handler *SQLHandler) GetFilesByType(fileType string) ([]Filelist, error) {
	sqlStr := fmt.Sprintf("select * from filelist where filesuffix = '%s'", fileType)
	return handler.sendQueryFilelist(sqlStr)
}

//Add
func (handler *SQLHandler) AddFile(f Filelist) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into filelist (filename,dupecount,filesuffix,filelocation,filesize,filehash,filedate,rowaction) values ('%s','%d','%s','%s',%v,'%s','%v','%s')", f.Filename, f.Dupecount, f.Filesuffix, f.Filelocation, f.Filesize, f.Filehash, f.Filedate, f.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
}

func (handler *SQLHandler) AddDupes(d DupeList) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into dupes (fileid,dupefileid,rowaction) values ('%v','%v','%v')", d.Fileid, d.Dupefileid, d.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err

}

func (handler *SQLHandler) AddEXIF(e EXIFInfoList) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into exifinfo (fileid, createdate, modifydate, accessdate, comments, tags, faces, facecoords, imageDescription, make, model, orientation, gpsaltitude, gpsaltituderef, gpslatitude, gpslatituderef, gpslongitude, gpslongituderef, gpsmapdatum, gpsprocessingmethod, gpsdatetime, gpsversionid, rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v' )", e.Fileid, e.Createdate, e.Modifydate, e.Accessdate, e.Comments, e.Tags, e.Faces, e.Facecoords, e.ImageDescription, e.Make, e.Model, e.Orientation, e.GPSAltitude, e.GPSAltitudeRef, e.GPSLatitude, e.GPSLatitudeRef, e.GPSLongitude, e.GPSLongitudeRef, e.GPSMapDatum, e.GPSProcessingMethod, e.GPSDateTime, e.GPSVersionID, e.Rowaction)

	/* mySqlDB, err := handler.Prepare("Insert into exifinfo (fileid, createdate, modifydate, accessdate, comments, tags, faces, facecoords, imageDescription, make, model, orientation, gpsaltitude, gpsaltituderef, gpslatitude, gpslatituderef, gpslongitude, gpslongituderef, gpsmapdatum, gpsprocessingmethod, gpsdatetime, gpsversionid, rowaction) values ( e.Fileid, e.Createdate, e.Modifydate, e.Accessdate, e.Comments, e.Tags, e.Faces, e.Facecoords, e.ImageDescription, e.Make, e.Model, e.Orientation, e.GPSAltitude, e.GPSAltitudeRef, e.GPSLatitude, e.GPSLatitudeRef, e.GPSLongitude, e.GPSLongitudeRef, e.GPSMapDatum, e.GPSProcessingMethod, e.GPSDateTime, e.GPSVersionID, e.Rowaction)")
	if err != nil {
		log.Printf("Errror preparing sqlStr for exifinfo")
		return 0, err
	} */
	log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	//res, err := mySqlDB.Exec(1)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}
	return rowid, err
}

func (handler *SQLHandler) AddTgtInfo(t TgtInfoList) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into targetinfo (fileid,filename,filesuffix,filelocation,filehash,filesize,filedate,fileaction,validated,rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v' )", t.Fileid, t.Filename, t.Filesuffix, t.Filelocation, t.Filehash, t.Filesize, t.Filedate, t.Fileaction, t.Validated, t.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
	//return nil //err
}

func (handler *SQLHandler) AddThumb(th ThumbList) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into thumbnails (fileid, targetid, filename, filesuffix, filesize, thsize, fileloc, rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", th.Fileid, th.Targetid, th.Filename, th.Filesuffix, th.Filesize, th.Thsize, th.Fileloc, th.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
	//return nil // err
}

func (handler *SQLHandler) AddRegionInfo(r RegionInfo) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Insert into regionInfo (fileid, filehash, name, type, areaH, areaW, areaX, areaY, areaUnit, rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", r.Fileid, r.Filehash, r.Name, r.Typ, r.AreaH, r.AreaW, r.AreaX, r.AreaY, r.AreaUnit, r.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
	//return nil // err
}

func (handler *SQLHandler) AddImageInfo(i ImageInfo) (rowid int64, err error) {

	sqlStr := fmt.Sprintf("Insert into imageInfo (fileid, filehash, dimH, dimW, dimUnit, rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v')", i.Fileid, i.Filehash, i.DimH, i.DimW, i.DimUnit, i.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
}

func (handler *SQLHandler) AddSessionInfo(s SessionInfo) (rowid int64, err error) {

	sqlStr := fmt.Sprintf("Insert into sessionInfo (starttime, stoptime, totalFiles, totalBytes, examinedFiles, examinedBytes, uniqueFiles, uniqueBytes, runTime, logFile, rowaction) values ('%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v', '%v')", s.Starttime, s.Stoptime, s.TotalFiles, s.TotalBytes, s.ExaminedFiles, s.ExaminedBytes, s.UniqueFiles, s.UniqueBytes, s.RunTime, s.LogFile, s.Rowaction)
	//log.Println(sqlStr)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
}

//// end Add

func (handler *SQLHandler) UpdFile(f Filelist, fname string) (rowid int64, err error) {
	sqlStr := fmt.Sprintf("Update filelist set filesuffix = '%s' ,filename = '%s',filelocation = '%s',filedate = '%v' where filename = '%s'", f.Filesuffix, f.Filename, f.Filelocation, f.Filedate, fname)

	sqlStr, err = prepareSqlStr(sqlStr)
	if err != nil {
		log.Printf("Error: ", err)
		return 0, err
	}

	res, err := handler.Exec(sqlStr)
	if err != nil {
		return 0, err
	}
	rowid, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return rowid, err
}

func (handler *SQLHandler) sendQueryFilelist(q string) ([]Filelist, error) {
	Filelists := []Filelist{}
	rows, err := handler.Query(q)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		f := Filelist{}
		err := rows.Scan(&f.Fileid, &f.Filename, &f.Filesuffix, &f.Filelocation, &f.Filesize, &f.Filehash, &f.Filedate, &f.Rowaction, &f.Rowactiondatetime)
		if err != nil {
			log.Println(err)
			continue
		}
		Filelists = append(Filelists, f)
	}
	return Filelists, rows.Err()
}

/*
	&f.Fileid,&f.Filename,&f.Filesuffix,&f.Filelocation,&f.Filesize,&f.Filehash,&f.Filedate,&f.Rowaction,&f.Rowactiondatetime
*/

func prepareSqlStr(sqlStr string) (pSqlStr string, err error) {
	err = nil

	log.Printf("Input sqlStr: %v\n", sqlStr)

	strLen := len(sqlStr)
	log.Printf("Length of sqlStr: %v\n", strLen)

	//escape all " except first and last
	re := regexp.MustCompile(`(^[<^\"])(.*)([>^\"]$)`)

	str := re.ReplaceAllString(sqlStr, "$2")

	//log.Printf("intermediate sqlStr: %v\n", str)

	//strLen = len(str)
	//log.Printf("Length of str: %v\n", strLen)

	// escape all single slash
	re = regexp.MustCompile(`(\b|\:)(\\)([^\\])`)
	str = re.ReplaceAllString(str, "$1\\$2$3")

	// Create replacer with pairs as arguments.
	r := strings.NewReplacer("\"", "\\\"")

	// Replace all pairs.
	str = r.Replace(str)

	// REplace all intermediate ' with \'
	re = regexp.MustCompile(`([a-zA-Z0-9])(')(\s*)[^,][^\)]`)
	str = re.ReplaceAllString(str, "$1\\$2$3")

	log.Printf("output sqlStr: %v\n", str)

	strLen = len(str)
	log.Printf("Length of str: %v\n", strLen)

	return str, err

}
