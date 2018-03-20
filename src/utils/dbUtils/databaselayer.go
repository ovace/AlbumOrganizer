package dbUtils


import (
	"database/sql"
	"errors"
	"time"
)

const (
	MYSQL uint8 = iota
	SQLITE
	POSTGRESQL
	MONGODB
)

type MediaDBHandler interface {
	GetAvailableFiles() ([]Filelist, error)
	GetFileByName(string) (Filelist, error)
	GetFilesByType(string) ([]Filelist, error)
	GetMinMaxID(filehash string) (FileMinMaxID, error)
	GetFileID(f Filelist) (FileIDnCount, error)
	//CrtFilelist(Filelist, string) error
	AddFile(Filelist) (rowid int64, err error)
	UpdFile(Filelist, string) (rowid int64, err error)
	//DelFile(Filelist, string) error
	//CrtDupes(DupeList) (rowid int64, err error)
	AddDupes(DupeList) (rowid int64, err error)
	//UpdDupes(DupeList) (rowid int64, err error)
	//DelDupes(DupeList) (rowid int64, err error)
	//CrtEXIF(EXIFInfoList) (rowid int64, err error)
	AddEXIF(EXIFInfoList) (rowid int64, err error)
	//UpdEXIF(EXIFInfoList) (rowid int64, err error)
	//DelEXIF(EXIFInfoList) (rowid int64, err error)
	//CrtTgtInfo(TgtInfoList) (rowid int64, err error)
	AddTgtInfo(TgtInfoList) (rowid int64, err error)
	//UpdTgtInfo(TgtInfoList) (rowid int64, err error)
	//DelTgtInfo(TgtInfoList) (rowid int64, err error)
	//CrtThumb(ThumbList) (rowid int64, err error)
	AddThumb(ThumbList) (rowid int64, err error)
	//UpdThumb(ThumbList) (rowid int64, err error)
	//DelThumb(ThumbList) (rowid int64, err error)
	//CrtImageInfo(i ImgInfo)(rowid int64, err error)
	AddImageInfo(i ImageInfo) (rowid int64, err error)
	//UpdateImageInfo(i ImageInfo)(rowid int64, err error)
	//DelImageInfo(i ImageInfo)(rowid int64, err error)
	//CrtRegionInfo(r RegionInfo)(rowid int64, err error)
	AddRegionInfo(r RegionInfo) (rowid int64, err error)
	//UpdateRegionInfo(r RegionInfo)(rowid int64, err error)
	//DelRegionInfo(r RegionInfo)(rowid int64, err error)
	AddSessionInfo(s SessionInfo) (rowid int64, err error)
}

type FileMinMaxID struct {
	MinID   int       `bson:"min-id"`
	MaxID   int       `bson:"max-id"`
	MinDate time.Time `bson:"min_time"`
	MaxDate time.Time `bson:"max_time"`
	MinDup  int       `bson:"min-dup"`
	MaxDup  int       `bson:"max-id"`
}

type FileIDnCount struct {
	FileID    sql.NullInt64 `bson:"file-id"`
	FileCount sql.NullInt64 `bson:"file-count"`
}

type Filelist struct {
	Fileid            int       `bson:"-"`
	Filename          string    `bson:"file_name"`
	Dupecount         int       `bson:"dupe_count-"`
	Filesuffix        string    `bson:"file_suffix"`
	Filelocation      string    `bson:"file_loc"`
	Filesize          int64     `bson:"file_size"`
	Filehash          string    `bson:"file_hash"`
	Filedate          time.Time `bson:"file_date"`
	Rowaction         string    `bson:"row_action"`
	Rowactiondatetime time.Time `bson:"row_ts"`
}

type DupeList struct {
	Dupeid            int       `bson:"-"`           //,"int(11)","NO","PRI","","auto_increment"
	Fileid            int       `bson:"file_id"`     //,"int(11)","NO","MUL","",""
	Dupefileid        int       `bson:"dupefile_id"` //,"int(11)","NO","","",""
	Rowaction         string    `bson:"row_action"`  //,"char(2)","NO","","",""
	Rowactiondatetime time.Time `bson:"row_ts"`      //,"datetime","NO","","CURRENT_TIMESTAMP",""
}

type EXIFInfoList struct {
	Exifid int `bson:"-"`       //","int(11)","NO","PRI","","auto_increment"
	Fileid int `bson:"file_id"` //","int(11)","NO","MUL","",""
	//*GPSInfo                   // Embedded
	Accessdate          time.Time `bson:"access_time"`           //","datetime","YES","","",""
	Comments            string    `bson:"comments"`              //","varchar(254)","YES","","",""
	Createdate          time.Time `bson:"create_time"`           //","datetime","YES","","",""
	Facecoords          string    `bson:"face_coords"`           //","varchar(254)","YES","","",""
	Faces               string    `bson:"faces"`                 //","varchar(254)","YES","","",""
	GPSAltitude         string    `bson:"gps_altitude"`          //` float DEFAULT NULL,
	GPSAltitudeRef      string    `bson:"gps_altitude_ref"`      //` float DEFAULT NULL,
	GPSDateTime         time.Time `bson:"gps_datetime"`          //` datetime,
	GPSLatitude         string    `bson:"gps_latitude"`          //` float DEFAULT NULL,
	GPSLatitudeRef      string    `bson:"gps_latitude_ref"`      //` varchar(254) DEFAULT NULL,
	GPSLongitude        string    `bson:"gps_longitude"`         //` string DEFAULT NULL,
	GPSLongitudeRef     string    `bson:"gps_longitude_ref"`     //` varchar(254) DEFAULT NULL,
	GPSMapDatum         string    `bson:"gps_map_datum"`         //` varchar(254) DEFAULT NULL,
	GPSProcessingMethod string    `bson:"gps_processing_method"` //` varchar(254) DEFAULT NULL,
	GPSVersionID        string    `bson:"gps_version_id"`        //` varchar(254) DEFAULT NULL,
	ImageDescription    string    `bson:"image_description"`     //` varchar(254) DEFAULT NULL,
	Make                string    `bson:"make"`                  //` varchar(254) DEFAULT NULL,
	Model               string    `bson:"model"`                 //` varchar(254) DEFAULT NULL,
	Modifydate          time.Time `bson:"mod_time"`              //","datetime","YES","","",""
	Orientation         string    `bson:"orientation"`           //` varchar(254) DEFAULT NULL,
	Rowaction           string    `bson:"row_action"`            //","char(2)","NO","","",""
	Rowactiondatetime   time.Time `bson:"row_ts"`                //","datetime","NO","","CURRENT_TIMESTAMP",""
	Tags                string    `bson:"tags"`                  //","varchar(254)","YES","","",""
}

type GPSInfo struct {
	GPSAltitude         float32   `bson:"gps_altitude"`          //` float DEFAULT NULL,
	GPSAltitudeRef      float32   `bson:"gps_altitude_ref"`      //` float DEFAULT NULL,
	GPSLatitude         float32   `bson:"gps_latitude"`          //` float DEFAULT NULL,
	GPSLatitudeRef      string    `bson:"gps_latitude_ref"`      //` varchar(254) DEFAULT NULL,
	GPSLongitude        int       `bson:"gps_longitude"`         //` int DEFAULT NULL,
	GPSLongitudeRef     string    `bson:"gps_longitude_ref"`     //` varchar(254) DEFAULT NULL,
	GPSMapDatum         string    `bson:"gps_map_datum"`         //` varchar(254) DEFAULT NULL,
	GPSProcessingMethod string    `bson:"gps_processing_method"` //` varchar(254) DEFAULT NULL,
	GPSDateTime         time.Time `bson:"gps_datetime"`          //` datetime,
	GPSVersionID        string    `bson:"gps_version_id"`        //` varchar(254) DEFAULT NULL,
}

type TgtInfoList struct {
	Targetid          int       `bson:"-"`           //","int(11)","NO","PRI","","auto_increment"
	Fileid            int       `bson:"file_id"`     //","int(11)","NO","MUL","",""
	Filename          string    `bson:"file_name"`   //","varchar(254)","NO","","",""
	Filesuffix        string    `bson:"file_suffix"` //","varchar(5)","YES","","",""
	Filelocation      string    `bson:"file_loc"`    //","varchar(254)","YES","","",""
	Filehash          string    `bson:"file_hash"`   //","varchar(254)","YES","","",""
	Filesize          float32   `bson:"file_siz"`    //","float","YES","","",""
	Filedate          time.Time `bson:"file_date"`   //","datetime","YES","","",""
	Fileaction        string    `bson:"file_action"` //","char(2)","NO","","",""
	Validated         int       `bson:"validated"`   //","tinyint(1)","NO","","0",""
	Rowaction         string    `bson:"row_action"`  //","char(2)","NO","","",""
	Rowactiondatetime time.Time `bson:"row_ts"`      //","datetime","NO","","CURRENT_TIMESTAMP",""
}
type ThumbList struct {
	Thid              int       `bson:"-"`           //","int(11)","NO","PRI","","auto_increment"
	Fileid            int       `bson:"file_id"`     //","int(11)","NO","MUL","",""
	Targetid          int       `bson:"tgt_id"`      //,"int(11)","NO","MUL","",""
	Filename          string    `bson:"file_name"`   //","varchar(254)","NO","","",""
	Filesuffix        string    `bson:"file_suffix"` //","varchar(254)","YES","","",""
	Filesize          string    `bson:"file_size"`   //","varchar(254)","YES","","",""
	Thsize            string    `bson:"th_size"`     //","varchar(254)","NO","","",""
	Fileloc           string    `bson:"file_loc"`    //","varchar(254)","YES","","",""
	Rowaction         string    `bson:"row_action"`  //","char(2)","NO","","",""
	Rowactiondatetime time.Time `bson:"row_ts"`      //","datetime","NO","","CURRENT_TIMESTAMP",""
}

type RegionInfo struct {
	Regionid          int       `bson:"-"`           //` int(11) NOT NULL AUTO_INCREMENT,
	Fileid            int64     `bson:"file_id"`     // int(11),
	Filehash          string    `bson:"file_hash"`   // varchar(254) DEFAULT NULL,
	Name              string    `bson:"region_name"` // varchar(254) NOT NULL,
	Typ               string    `bson:"region_type"` // varchar(5) DEFAULT NULL,
	AreaH             float64   `bson:"area_H"`      // float DEFAULT NULL,
	AreaW             float64   `bson:"area_W"`      // float DEFAULT NULL,
	AreaX             float64   `bson:"area_X"`      // float DEFAULT NULL,
	AreaY             float64   `bson:"area_Y"`      //` float DEFAULT NULL,
	AreaUnit          string    `bson:"area_units"`  // varchar(5) DEFAULT NULL,
	Rowaction         string    `bson:"row_action"`  // char(2) NOT NULL,
	Rowactiondatetime time.Time `bson:"row_ts"`      // timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
}

type ImageInfo struct {
	Imgid             int       `bson:"-"`              // int(11) NOT NULL AUTO_INCREMENT,
	Fileid            int64     `bson:"file_id"`        // int(11),
	Filehash          string    `bson:"file_hash"`      // varchar(254) DEFAULT NULL,
	DimH              int64     `bson:"image_dimH"`     // float DEFAULT NULL,
	DimW              int64     `bson:"image_dimW"`     // float DEFAULT NULL,
	DimUnit           string    `bson:"image_dimunits"` // varchar(5) DEFAULT NULL,
	Rowaction         string    `bson:"row_action"`     // char(2) NOT NULL,
	Rowactiondatetime time.Time `bson:"row_ts"`         // timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
}

type SessionInfo struct {
	Sessionid         int       `bson:"-"`              //int(11) NOT NULL AUTO_INCREMENT,
	Starttime         time.Time `bson:"start_time"`     //timestamp,
	Stoptime          time.Time `bson:"stop_time"`      //timestamp,
	TotalFiles        int       `bson:"total_files"`    //int,
	TotalBytes        int       `bson:"total_bytes"`    //int,
	ExaminedFiles     int       `bson:"examined_files"` //int,
	ExaminedBytes     int       `bson:"examined_bytes"` //int,
	UniqueFiles       int       `bson:"unique_files"`   //int,
	UniqueBytes       int       `bson:"unique_bytes"`   //int,
	RunTime           float64   `bson:"run_time"`       //float,
	LogFile           string    `bson:"log_file"`       //varchar(254),
	Rowaction         string    `bson:"row_action"`     //char(2) NOT NULL,
	Rowactiondatetime time.Time `bson:"row_ts"`         //timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
}

var DBTypeNotSupported = errors.New("The Database type provided is not supported...")

//factory function
func GetDatabaseHandler(dbtype uint8, connection string) (MediaDBHandler, error) {

	switch dbtype {
	case MYSQL:
		return NewMySQLHandler(connection)
	/* case MONGODB:
	return NewMongodbHandler(connection) */
	case SQLITE:
		return NewSQLiteHandler(connection)
		/* case POSTGRESQL:
		return NewPQHandler(connection) */
	}
	return nil, DBTypeNotSupported
}
