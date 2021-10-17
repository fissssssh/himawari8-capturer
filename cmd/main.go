package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const urlTemplate = "https://himawari8.nict.go.jp/img/D531106/%dd/550/%d/%02d/%02d/%02d%d000_%d_%d.png"
const tileWidth = 550

type Signal struct {
	X     int
	Y     int
	Image image.Image
}

func main() {
	var quality uint
	var timestamp int64
	flag.UintVar(&quality, "q", 2, "Image quality (1: 550x550, 2: 1100x1100, 4: 2200x2200, 8: 4400x4400, 16: 8800x8800)")
	flag.Int64Var(&timestamp, "t", 0, "Unix timestamp(ms)")
	flag.Parse()
	if quality != 1 && quality != 2 && quality != 4 && quality != 8 && quality != 16 {
		log.Printf("unknown image quality %d", quality)
		os.Exit(1)
	}
	control := make(chan struct{}, 5)
	var ch = make(chan Signal)
	var t time.Time
	if timestamp != 0 {
		t = time.UnixMilli(int64(timestamp)).UTC()
	} else {
		t = time.Now().UTC().Add(-20 * time.Minute)
	}
	year, month, day := t.Date()
	hour, minute := t.Hour(), t.Minute()
	level := int(quality)
	r := image.NewRGBA(image.Rect(0, 0, level*tileWidth, level*tileWidth))
	for x := 0; x < level; x++ {
		for y := 0; y < level; y++ {
			control <- struct{}{}
			go func(x int, y int) {
				url := fmt.Sprintf(urlTemplate, level, year, month, day, hour, minute/10, x, y)
				log.Printf("Get %d-%d tile image from %s...", x, y, url)
				data, err := download(url)
				if err != nil {
					log.Fatal(err)
				}
				img, _, err := image.Decode(bytes.NewReader(data))
				if err != nil {
					log.Fatal(err)
				}
				log.Printf("Get %d-%d tile image done!", x, y)
				<-control
				ch <- Signal{x, y, img}
			}(x, y)
		}
	}

	// compose images.
	for i := 0; i < level*level; i++ {
		signal := <-ch
		x, y, img := signal.X, signal.Y, signal.Image
		draw.Draw(r, img.Bounds().Add(image.Pt(x*tileWidth, y*tileWidth)), img, img.Bounds().Min, draw.Src)
	}
	log.Printf("Tile images were composed!")

	// save to os.
	filename := fmt.Sprintf("himawari8_%d%02d%02dT%02d%d000Z.png", year, month, day, hour, minute/10)
	log.Printf("Saving image to %s", filename)
	var buf bytes.Buffer
	png.Encode(&buf, r)
	save(&buf, filename)
	log.Printf("All done!")
}

func download(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func save(reader io.Reader, filename string) error {
	out, err := os.Create(filename)
	defer out.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(out, reader)
	return err
}
