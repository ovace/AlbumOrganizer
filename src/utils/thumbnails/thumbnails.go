package thumbnails

import (
)

func (env *Env) generateThumbImg(imgFile string, thImgName string, thPath string, thSize uint) (err error) {

	var img image.Image
	_, format, err := env.decodeConfig(imgFile)

	if err != nil {
		log.Fatal(err)
	}
	// open "test.jpg"
	file, err := os.Open(imgFile)
	if err != nil {
		log.Fatal(err)
	}
	switch format {
	case "jpeg":
		img, err = jpeg.Decode(file)
	case "png":
		img, err = png.Decode(file)
	case "gif":
		img, err = gif.Decode(file)
	case "bmp":
		img, err = bmp.Decode(file)
	default:
		err = errors.New("Unsupported file type")
		log.Errorf("Unsupport image type %v\n", format)
	}
	//img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(thSize, 0, img, resize.Lanczos3)

	if _, err := os.Stat(thPath); os.IsNotExist(err) {
		os.MkdirAll(thPath, os.ModePerm)
	}

	out, err := os.Create(thPath + thImgName)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	switch format {
	case "jpeg":
		jpeg.Encode(out, m, nil)
	case "png":
		png.Encode(out, m)
	case "gif":
		gif.Encode(out, m, nil)
	case "bmp":
		bmp.Encode(out, m)
	default:
		err = errors.New("Unsupported file type")
	}
	// write new image to file
	//jpeg.Encode(out, m, nil)

	return err
}
