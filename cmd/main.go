package main

import (
	"flag"
	"fmt"
	"himawari8Capturer/core"
	"io"
	"log"
	"os"
	"time"
)

func main() {
	var quality uint
	var timestamp int64
	var shorelinesColor uint
	flag.UintVar(&quality, "q", 2, "Image quality (1: 550x550, 2: 1100x1100, 3: 2200x2200, 4: 4400x4400, 5: 8800x8800, 6: 11000*11000)")
	flag.UintVar(&shorelinesColor, "l", 0, "Shorelines color (0: None, 1: Red, 2: Green, 3: Yellow)")
	flag.Int64Var(&timestamp, "t", 0, "Unix timestamp(ms)")
	flag.Parse()
	var q core.Quality
	switch quality {
	case 1:
		q = core.Low
	case 2:
		q = core.HD
	case 3:
		q = core.FHD
	case 4:
		q = core.QHD
	case 5:
		q = core.UHD
	case 6:
		q = core.UHDPlus
	default:
		log.Fatalf("unknown image quality %d", quality)
	}
	log.Printf("Image Quality: %d", q)
	var c core.Shoreline
	switch shorelinesColor {
	case 0:
		c = core.Ignore
	case 1:
		c = core.Red
	case 2:
		c = core.Green
	case 3:
		c = core.Yellow
	default:
		log.Fatalf("unknown shorelines color %d", shorelinesColor)
	}
	if c != core.Ignore {
		log.Printf("Image Shorelines: %s", c)
	}
	var err error
	var t time.Time
	if timestamp != 0 {
		t = time.UnixMilli(int64(timestamp)).UTC()
		log.Printf("Getting image at %v", t)
	} else {
		t, err = core.GetLatestEarthTime()
		if err != nil {
			log.Fatalf("get latest date failed: %s", err)
		}
		log.Printf("Getting latest image at %v", t)
	}
	year, month, day := t.Date()
	hour, minute := t.Hour(), t.Minute()
	img, err := core.GetEarthWithShorelines(q, t, c)
	if err != nil {
		log.Fatalf("get image failed %s", err)
	}
	// save to os.
	filename := fmt.Sprintf("himawari8_%d%02d%02dT%02d%d000Z.png", year, month, day, hour, minute/10)
	fmt.Printf("saving image to %s\n", filename)
	out, err := os.Create(filename)
	if err != nil {
		log.Fatalf("save image failed: %s", err)
	}
	defer out.Close()
	_, err = io.Copy(out, img)
	if err != nil {
		log.Fatalf("save image failed: %s", err)
	} else {
		fmt.Printf("all done!\n")
	}
}
