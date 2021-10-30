package main

import (
	"flag"
	"fmt"
	"himawari8Capturer/pkg/himawari8"
	"image"
	"io"
	"log"
	"os"
	"time"
)

const tileWidth = 550
const tileUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/%d/%02d/%02d/%02d%d000_%d_%d.png"
const shoreLineUrlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/coastline/%s_%d_%d.png"

type Signal struct {
	X     int
	Y     int
	Image image.Image
}

const (
	None uint = iota
	Red
	Green
	Yellow
)

func main() {
	var quality uint
	var timestamp int64
	var shorelinesColor uint
	flag.UintVar(&quality, "q", 2, "Image quality (1: 550x550, 2: 1100x1100, 4: 2200x2200, 8: 4400x4400, 16: 8800x8800)")
	flag.UintVar(&shorelinesColor, "l", 0, "Shorelines color (0: None, 1: Red, 2: Green, 3: Yellow)")
	flag.Int64Var(&timestamp, "t", 0, "Unix timestamp(ms)")
	flag.Parse()
	var q himawari8.Quality
	switch quality {
	case 1:
		q = himawari8.Low
	case 2:
		q = himawari8.HD
	case 4:
		q = himawari8.FHD
	case 8:
		q = himawari8.QHD
	case 16:
		q = himawari8.UHD
	default:
		log.Fatalf("unknown image quality %d", quality)
	}
	var c himawari8.ShorelineColor
	switch shorelinesColor {
	case 0:
		c = himawari8.Ignore
	case 1:
		c = himawari8.Red
	case 2:
		c = himawari8.Green
	case 3:
		c = himawari8.Yellow
	default:
		log.Fatalf("unknown shorelines color %d", shorelinesColor)
	}
	var t time.Time
	if timestamp != 0 {
		t = time.UnixMilli(int64(timestamp)).UTC()
	} else {
		t = time.Now().UTC().Add(-20 * time.Minute)
	}
	year, month, day := t.Date()
	hour, minute := t.Hour(), t.Minute()
	img, err := himawari8.GetImage(q, t, c)
	if err != nil {
		log.Fatalf("get image failed %s", err)
	}
	// save to os.
	filename := fmt.Sprintf("himawari8_%d%02d%02dT%02d%d000Z.png", year, month, day, hour, minute/10)
	log.Printf("saving image to %s", filename)
	save(img, filename)
	log.Printf("all done!")
}

func save(reader io.Reader, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, reader)
	return err
}
