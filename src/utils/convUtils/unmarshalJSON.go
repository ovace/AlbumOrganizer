package convUtils

import (
	"encoding/json"
	"fmt"
	"os"
)

type E struct {
	EF []EXIF
}
type EXIF struct {
	SourceFile               string
	ExifToolVersion          float64
	FileName                 string
	Directory                string
	FileSize                 string
	FileModifyDate           string
	FileAccessDate           string
	FileCreateDate           string
	FilePermissions          string
	FileType                 string
	FileTypeExtension        string
	MIMEType                 string
	ExifByteOrder            string
	Make                     string
	Model                    string
	Orientation              string
	XResolution              float64
	YResolution              float64
	ResolutionUnit           string
	Software                 string
	ModifyDate               string
	YCbCrPositioning         string
	ExposureTime             string
	FNumber                  float64
	ExposureProgram          string
	ISO                      int64
	ExifVersion              string
	DateTimeOriginal         string
	CreateDate               string
	ComponentsConfiguration  string
	ShutterSpeedValue        string
	ApertureValue            float64
	BrightnessValue          float64
	ExposureCompensation     float64
	MaxApertureValue         float64
	MeteringMode             string
	LightSource              string
	Flash                    string
	FocalLength              string
	UserComment              string
	FlashpixVersion          string
	ColorSpace               string
	ExifImageWidth           int64
	ExifImageHeight          int64
	InteropIndex             string
	InteropVersion           string
	SensingMethod            string
	SceneType                string
	ExposureMode             string
	WhiteBalance             string
	SceneCaptureType         string
	ImageUniqueID            string
	Compression              string
	ThumbnailOffset          int64
	ThumbnailLength          int64
	XMPToolkit               string
	CurrentIPTCDigest        string
	EnvelopeRecordVersion    int64
	CodedCharacterSet        string
	ApplicationRecordVersion int64
	IPTCDigest               string
	ImageWidth               int64
	ImageHeight              int64
	EncodingProcess          string
	BitsPerSample            int64
	ColorComponents          int64
	YCbCrSubSampling         string
	Aperture                 float64
	ImageSize                string
	Megapixels               float64
	ShutterSpeed             string
	ThumbnailImage           string
	FocalLength35efl         string
	LightValue               float64
	RegionInfo               RegionInfo
}
type RegionInfo struct {
	AppliedToDimensions AppliedToDimensions
	RegionList          []RegionList
}

type AppliedToDimensions struct {
	H    int64
	Unit string
	W    int64
}

type RegionList struct {
	Area Area
	Name string
	Type string
}
type Area struct {
	H    float64
	Unit string
	W    float64
	X    float64
	Y    float64
}

type person struct {
	Name   string
	Traits struct {
		Age     int
		Gender  int
		Traits2 struct {
			EyeColor string
		}
	}
}

func main() {
	/* 	jsonString1 := `{  "SourceFile": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/20130713_103256_3264x2448_SPH_L710_001.jpg",
	   "ExifToolVersion": 10.61,
	   "FileName": "20130713_103256_3264x2448_SPH_L710_001.jpg",
	   "Directory": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org",
	   "FileSize": "3.7 MB",
	   "FileModifyDate": "2017:05:31 21:53:26-05:00",
	   "FileAccessDate": "2017:05:31 21:53:26-05:00",
	   "FileCreateDate": "2017:05:31 21:53:26-05:00",
	   "FilePermissions": "rw-rw-rw-",
	   "FileType": "JPEG",
	   "FileTypeExtension": "jpg",
	   "MIMEType": "image/jpeg",
	   "ExifByteOrder": "Little-endian (Intel, II)",
	   "Make": "SAMSUNG",
	   "Model": "SPH-L710",
	   "Orientation": "Horizontal (normal)",
	   "XResolution": 72,
	   "YResolution": 72,
	   "ResolutionUnit": "inches",
	   "Software": "L710VPBMD4",
	   "ModifyDate": "2017:02:20 18:31:00",
	   "YCbCrPositioning": "Centered",
	   "ExposureTime": "1/165",
	   "FNumber": 2.6,
	   "ExposureProgram": "Aperture-priority AE",
	   "ISO": 80,
	   "ExifVersion": "0220",
	   "DateTimeOriginal": "2013:07:13 10:32:56",
	   "CreateDate": "2013:07:13 10:32:56",
	   "ComponentsConfiguration": "Y, Cb, Cr, -",
	   "ShutterSpeedValue": "1/165",
	   "ApertureValue": 2.6,
	   "BrightnessValue": 5.55859375,
	   "ExposureCompensation": 0,
	   "MaxApertureValue": 2.6,
	   "MeteringMode": "Center-weighted average",
	   "LightSource": "Unknown",
	   "Flash": "No Flash",
	   "FocalLength": "3.7 mm",
	   "UserComment": "",
	   "FlashpixVersion": "0100",
	   "ColorSpace": "sRGB",
	   "ExifImageWidth": 3264,
	   "ExifImageHeight": 2448,
	   "InteropIndex": "R98 - DCF basic file (sRGB)",
	   "InteropVersion": "0100",
	   "SensingMethod": "One-chip color area",
	   "SceneType": "Directly photographed",
	   "ExposureMode": "Auto",
	   "WhiteBalance": "Auto",
	   "SceneCaptureType": "Standard",
	   "ImageUniqueID": "GDFI02",
	   "Compression": "JPEG (old-style)",
	   "ThumbnailOffset": 4962,
	   "ThumbnailLength": 6768,
	   "XMPToolkit": "XMP Core 5.1.2",
	   "RegionInfo": {
	      "AppliedToDimensions": {
	         "H": 2448,
	         "Unit": "pixel",
	         "W": 3264
	      },
	      "RegionList": [
	         {
	            "Area": {
	               "H": 0.0408497,
	               "Unit": "normalized",
	               "W": 0.0248162,
	               "X": 0.619945,
	               "Y": 0.495507
	            },
	            "Name": "Shaheen Ahmed",
	            "Type": "Face"
	         },
	         {
	            "Area": {
	               "H": 0.0433007,
	               "Unit": "normalized",
	               "W": 0.0266544,
	               "X": 0.679075,
	               "Y": 0.461601
	            },
	            "Name": "Iftakar Ahmed (uncle)",
	            "Type": "Face"
	         },
	         {
	            "Area": {
	               "H": 0.0555556,
	               "Unit": "normalized",
	               "W": 0.0343137,
	               "X": 0.41973,
	               "Y": 0.461193
	            },
	            "Name": "Kulsum Asghar Anis",
	            "Type": "Face"
	         }
	      ]
	   },
	   "CurrentIPTCDigest": "b443520a10119da99c2550175e6d0efb",
	   "EnvelopeRecordVersion": 4,
	   "CodedCharacterSet": "UTF8",
	   "ApplicationRecordVersion": 4,
	   "IPTCDigest": "b443520a10119da99c2550175e6d0efb",
	   "ImageWidth": 3264,
	   "ImageHeight": 2448,
	   "EncodingProcess": "Baseline DCT, Huffman coding",
	   "BitsPerSample": 8,
	   "ColorComponents": 3,
	   "YCbCrSubSampling": "YCbCr4:2:2 (2 1)",
	   "Aperture": 2.6,
	   "ImageSize": "3264x2448",
	   "Megapixels": 8.0,
	   "ShutterSpeed": "1/165",
	   "ThumbnailImage": "(Binary data 6768 bytes, use -b option to extract)",
	   "FocalLength35efl": "3.7 mm",
	   "LightValue": 10.4
	}`
	*/
	/* 	jsonString2 := `{
	  "SourceFile": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/some_description_with_subs/20131215_204143_2448x3264_SPH_L710_000.jpg",
	  "ExifToolVersion": 10.61,
	  "FileName": "20131215_204143_2448x3264_SPH_L710_000.jpg",
	  "Directory": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/some_description_with_subs",
	  "FileSize": "2.8 MB",
	  "FileModifyDate": "2017:05:31 21:53:35-05:00",
	  "FileAccessDate": "2017:05:31 21:53:35-05:00",
	  "FileCreateDate": "2017:05:31 21:53:35-05:00",
	  "FilePermissions": "rw-rw-rw-",
	  "FileType": "JPEG",
	  "FileTypeExtension": "jpg",
	  "MIMEType": "image/jpeg",
	  "JFIFVersion": 1.01,
	  "ExifByteOrder": "Little-endian (Intel, II)",
	  "Make": "SAMSUNG",
	  "Model": "SPH-L710",
	  "Orientation": "Horizontal (normal)",
	  "XResolution": 72,
	  "YResolution": 72,
	  "ResolutionUnit": "inches",
	  "Software": "L710VPBMD4",
	  "ModifyDate": "2017:02:23 10:13:50",
	  "YCbCrPositioning": "Centered",
	  "ExposureTime": "1/15",
	  "FNumber": 2.6,
	  "ExposureProgram": "Aperture-priority AE",
	  "ISO": 250,
	  "ExifVersion": "0220",
	  "DateTimeOriginal": "2013:12:15 20:41:43",
	  "CreateDate": "2013:12:15 20:41:43",
	  "ComponentsConfiguration": "Y, Cb, Cr, -",
	  "ShutterSpeedValue": "1/15",
	  "ApertureValue": 2.6,
	  "BrightnessValue": 0.3671875,
	  "ExposureCompensation": 0,
	  "MaxApertureValue": 2.6,
	  "MeteringMode": "Center-weighted average",
	  "LightSource": "Unknown",
	  "Flash": "Fired",
	  "FocalLength": "3.7 mm",
	  "FlashpixVersion": "0100",
	  "ColorSpace": "sRGB",
	  "ExifImageWidth": 2448,
	  "ExifImageHeight": 3264,
	  "InteropIndex": "R98 - DCF basic file (sRGB)",
	  "InteropVersion": "0100",
	  "SensingMethod": "One-chip color area",
	  "SceneType": "Directly photographed",
	  "ExposureMode": "Auto",
	  "WhiteBalance": "Auto",
	  "SceneCaptureType": "Standard",
	  "ImageUniqueID": "GDFI02",
	  "Compression": "JPEG (old-style)",
	  "ThumbnailOffset": 872,
	  "ThumbnailLength": 6502,
	  "XMPToolkit": "XMP Core 5.1.2",
	  "RegionInfo": {
	    "AppliedToDimensions": {
	      "H": 3264,
	      "Unit": "pixel",
	      "W": 2448
	    },
	    "RegionList": [{
	      "Area": {
	        "H": 0.231311,
	        "Unit": "normalized",
	        "W": 0.257353,
	        "X": 0.752042,
	        "Y": 0.876991
	      },
	      "Name": "Nishat Mamnoon",
	      "Type": "Face"
	    }]
	  },
	  "CurrentIPTCDigest": "b443520a10119da99c2550175e6d0efb",
	  "EnvelopeRecordVersion": 4,
	  "CodedCharacterSet": "UTF8",
	  "ApplicationRecordVersion": 4,
	  "IPTCDigest": "b443520a10119da99c2550175e6d0efb",
	  "ImageWidth": 2448,
	  "ImageHeight": 3264,
	  "EncodingProcess": "Baseline DCT, Huffman coding",
	  "BitsPerSample": 8,
	  "ColorComponents": 3,
	  "YCbCrSubSampling": "YCbCr4:4:0 (1 2)",
	  "Aperture": 2.6,
	  "ImageSize": "2448x3264",
	  "Megapixels": 8.0,
	  "ShutterSpeed": "1/15",
	  "ThumbnailImage": "(Binary data 6502 bytes, use -b option to extract)",
	  "FocalLength35efl": "3.7 mm",
	  "LightValue": 5.3
	}` */

	/* jsonString3 := `{
	  "SourceFile": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/some_description_with_subs/sub1/20140524_142818_1944x2592_Xoom_000.jpg",
	  "ExifToolVersion": 10.61,
	  "FileName": "20140524_142818_1944x2592_Xoom_000.jpg",
	  "Directory": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/some_description_with_subs/sub1",
	  "FileSize": "3.8 MB",
	  "FileModifyDate": "2017:05:31 21:53:36-05:00",
	  "FileAccessDate": "2017:05:31 21:53:36-05:00",
	  "FileCreateDate": "2017:05:31 21:53:36-05:00",
	  "FilePermissions": "rw-rw-rw-",
	  "FileType": "JPEG",
	  "FileTypeExtension": "jpg",
	  "MIMEType": "image/jpeg",
	  "JFIFVersion": 1.01,
	  "ExifByteOrder": "Big-endian (Motorola, MM)",
	  "ImageDescription": "",
	  "Make": "Motorola Mobility",
	  "Model": "Xoom",
	  "Orientation": "Horizontal (normal)",
	  "XResolution": 72,
	  "YResolution": 72,
	  "ResolutionUnit": "inches",
	  "Software": "",
	  "ModifyDate": "2014:05:24 14:28:18",
	  "Artist": "",
	  "YCbCrPositioning": "Centered",
	  "Copyright": "",
	  "ExposureTime": "1/270",
	  "FNumber": 2.8,
	  "ExposureProgram": "Not Defined",
	  "ISO": 100,
	  "ExifVersion": "0220",
	  "DateTimeOriginal": "2014:05:24 14:28:18",
	  "CreateDate": "2014:05:24 14:28:18",
	  "ComponentsConfiguration": "Y, Cb, Cr, -",
	  "CompressedBitsPerPixel": 4,
	  "ExposureCompensation": 0,
	  "MaxApertureValue": 1.7,
	  "SubjectDistance": "0 m",
	  "MeteringMode": "Other",
	  "LightSource": "Unknown",
	  "Flash": "No Flash",
	  "FocalLength": "4.4 mm",
	  "FlashpixVersion": "0100",
	  "ColorSpace": "sRGB",
	  "ExifImageWidth": 1944,
	  "ExifImageHeight": 2592,
	  "InteropIndex": "Unknown ()",
	  "InteropVersion": "0110",
	  "FileSource": "Digital Camera",
	  "CustomRendered": "Normal",
	  "ExposureMode": "Auto",
	  "WhiteBalance": "Auto",
	  "DigitalZoomRatio": 1,
	  "SceneCaptureType": "Standard",
	  "UserComment": "",
	  "MakerNoteUnknownText": "(Binary data 2048 bytes, use -b option to extract)",
	  "GPSVersionID": "2.2.0.0",
	  "GPSTimeStamp": "19:28:18",
	  "GPSDateStamp": "2014:05:24",
	  "Compression": "JPEG (old-style)",
	  "ThumbnailOffset": 3063,
	  "ThumbnailLength": 38289,
	  "ImageWidth": 1944,
	  "ImageHeight": 2592,
	  "EncodingProcess": "Baseline DCT, Huffman coding",
	  "BitsPerSample": 8,
	  "ColorComponents": 3,
	  "YCbCrSubSampling": "YCbCr4:2:0 (2 2)",
	  "Aperture": 2.8,
	  "GPSDateTime": "2014:05:24 19:28:18Z",
	  "ImageSize": "1944x2592",
	  "Megapixels": 5.0,
	  "ShutterSpeed": "1/270",
	  "ThumbnailImage": "(Binary data 38289 bytes, use -b option to extract)",
	  "FocalLength35efl": "4.4 mm",
	  "LightValue": 11.0
	}` */

	jsonString4 := `[{
  "SourceFile": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org/20121027_193625_1936x2592_iPad_001.jpg",
  "ExifToolVersion": 10.61,
  "FileName": "20121027_193625_1936x2592_iPad_001.jpg",
  "Directory": "c:/Users/mamnoon/My_Workspace/goProjects/src/github.com/ovace/albumMgmt/Pictures/org",
  "FileSize": "1565 kB",
  "FileModifyDate": "2017:05:31 21:53:25-05:00",
  "FileAccessDate": "2017:05:31 21:53:25-05:00",
  "FileCreateDate": "2017:05:31 21:53:25-05:00",
  "FilePermissions": "rw-rw-rw-",
  "FileType": "JPEG",
  "FileTypeExtension": "jpg",
  "MIMEType": "image/jpeg",
  "JFIFVersion": 1.01,
  "ExifByteOrder": "Big-endian (Motorola, MM)",
  "Make": "Apple",
  "Model": "iPad",
  "Orientation": "Horizontal (normal)",
  "XResolution": 72,
  "YResolution": 72,
  "ResolutionUnit": "inches",
  "Software": "5.1.1",
  "ModifyDate": "2017:02:22 15:57:55",
  "YCbCrPositioning": "Centered",
  "ExposureTime": "1/15",
  "FNumber": 2.4,
  "ExposureProgram": "Program AE",
  "ISO": 1000,
  "ExifVersion": "0221",
  "DateTimeOriginal": "2012:10:27 19:36:25",
  "CreateDate": "2012:10:27 19:36:25",
  "ComponentsConfiguration": "Y, Cb, Cr, -",
  "ShutterSpeedValue": "1/15",
  "ApertureValue": 2.4,
  "BrightnessValue": -3.011695906,
  "MeteringMode": "Multi-segment",
  "Flash": "No flash function",
  "FocalLength": "4.3 mm",
  "SubjectArea": "1295 967 699 696",
  "FlashpixVersion": "0100",
  "ColorSpace": "sRGB",
  "ExifImageWidth": 1936,
  "ExifImageHeight": 2592,
  "SensingMethod": "One-chip color area",
  "ExposureMode": "Auto",
  "WhiteBalance": "Auto",
  "FocalLengthIn35mmFormat": "35 mm",
  "SceneCaptureType": "Standard",
  "Sharpness": "Normal",
  "Compression": "JPEG (old-style)",
  "ThumbnailOffset": 714,
  "ThumbnailLength": 9294,
  "XMPToolkit": "XMP Core 5.1.2",
  "RegionInfo": {
    "AppliedToDimensions": {
      "H": 2592,
      "Unit": "pixel",
      "W": 1936
    },
    "RegionList": [{
      "Area": {
        "H": 0.254244,
        "Unit": "normalized",
        "W": 0.284091,
        "X": 0.425103,
        "Y": 0.43152
      },
      "Name": "Osaid Mamnoon",
      "Type": "Face"
    }]
  },
  "CurrentIPTCDigest": "b443520a10119da99c2550175e6d0efb",
  "EnvelopeRecordVersion": 4,
  "CodedCharacterSet": "UTF8",
  "ApplicationRecordVersion": 4,
  "IPTCDigest": "b443520a10119da99c2550175e6d0efb",
  "ImageWidth": 1936,
  "ImageHeight": 2592,
  "EncodingProcess": "Baseline DCT, Huffman coding",
  "BitsPerSample": 8,
  "ColorComponents": 3,
  "YCbCrSubSampling": "YCbCr4:2:0 (2 2)",
  "Aperture": 2.4,
  "ImageSize": "1936x2592",
  "Megapixels": 5.0,
  "ScaleFactor35efl": 8.2,
  "ShutterSpeed": "1/15",
  "ThumbnailImage": "(Binary data 9294 bytes, use -b option to extract)",
  "CircleOfConfusion": "0.004 mm",
  "FOV": "54.4 deg",
  "FocalLength35efl": "4.3 mm (35 mm equivalent: 35.0 mm)",
  "HyperfocalDistance": "2.08 m",
  "LightValue": 3.1
}]`

	var decoded []EXIF

	if err := json.Unmarshal([]byte(jsonString4), &decoded); err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	fmt.Println(decoded[0].SourceFile)
	fmt.Println(decoded[0].RegionInfo.RegionList[len(decoded[0].RegionInfo.RegionList)-1].Name)
	fmt.Println(decoded[0].RegionInfo.RegionList[len(decoded[0].RegionInfo.RegionList)-1].Area.X)
	fmt.Println(len(decoded[0].RegionInfo.RegionList))
	fmt.Println(cap(decoded[0].RegionInfo.RegionList))
}
