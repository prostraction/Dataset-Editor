package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// var mu sync.Mutex
var dir_first string
var dir_second string
var dir_merged string

func getMinBounds(bounds1 image.Rectangle, bounds2 image.Rectangle) image.Rectangle {
	var rect image.Rectangle
	if bounds1.Min.X < bounds2.Min.X {
		rect.Min.X = bounds1.Min.X
	} else {
		rect.Min.X = bounds2.Min.X
	}
	if bounds1.Max.X < bounds2.Max.X {
		rect.Max.X = bounds1.Max.X
	} else {
		rect.Max.X = bounds2.Max.X
	}
	if bounds1.Min.Y < bounds2.Min.Y {
		rect.Min.Y = bounds1.Min.Y
	} else {
		rect.Min.Y = bounds2.Min.Y
	}
	if bounds1.Max.Y < bounds2.Max.Y {
		rect.Max.Y = bounds1.Max.Y
	} else {
		rect.Max.Y = bounds2.Max.Y
	}
	return rect
}

func clamp(value, min, max uint32) uint32 {
	if value < min {
		return min
	} else if value > max {
		return max
	} else {
		return value
	}
}

func plusColors(c1 color.Color, c2 color.Color) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	r3 := clamp(r1+r2, 0, 65535)
	g3 := clamp(g1+g2, 0, 65535)
	b3 := clamp(b1+b2, 0, 65535)
	a3 := clamp(a1+a2, 0, 65535)
	rgba := color.RGBA64{uint16(r3), uint16(g3), uint16(b3), uint16(a3)}
	return color.Color(rgba)
}

func merge(img1 image.Image, img2 image.Image) (image.Image, error) {
	bounds := getMinBounds(img1.Bounds(), img2.Bounds())
	img := image.NewRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	if img != nil {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				img.Set(x, y, plusColors(img1.At(x, y), img2.At(x, y)))
			}
		}
		return img, nil
	}
	return nil, errors.New("image.NewRGMA: alloc failed")
}

func openImage(path string) (image.Image, error) {
	var img image.Image
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	if filepath.Ext(strings.ToLower(file.Name())) == ".jpeg" || filepath.Ext(strings.ToLower(file.Name())) == ".jpg" {
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, err
		}
	} else if filepath.Ext(strings.ToLower(file.Name())) == ".png" {
		img, err = png.Decode(file)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("wrong type of file")
	}
	return img, nil
}

func process(f1 os.FileInfo, f2 os.FileInfo) {
	img1, err1 := openImage(dir_first + f1.Name())
	img2, err2 := openImage(dir_second + f2.Name())
	if err1 == nil && err2 == nil {
		img_res, err := merge(img1, img2)
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.Create(dir_merged + f1.Name())
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		if filepath.Ext(strings.ToLower(f.Name())) == ".jpeg" || filepath.Ext(strings.ToLower(f.Name())) == ".jpg" {
			jpeg.Encode(f, img_res, nil)
		} else if filepath.Ext(strings.ToLower(f.Name())) == ".png" {
			png.Encode(f, img_res)
		} else {
			fmt.Println("wrong type")
		}
	} else {
		if err1 != nil {
			fmt.Println(err1.Error())
		}
		if err2 != nil {
			fmt.Println(err2.Error())
		}
	}
}

func cut(f os.FileInfo, boundsReq image.Rectangle) {
	img, err := openImage(dir_first + f.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	if img.Bounds().Max.X < boundsReq.Max.X || img.Bounds().Max.Y < boundsReq.Max.Y {
		fmt.Println("image to small")
		return
	}
	bounds := img.Bounds()
	imgCroppped := image.NewRGBA(image.Rect(boundsReq.Min.X, boundsReq.Min.Y, boundsReq.Max.X, boundsReq.Max.Y))
	yy := bounds.Min.Y
	xx := bounds.Min.X
	cycle := 1

	for {
		for y := bounds.Min.Y; y < boundsReq.Max.Y; y++ {
			for x := bounds.Min.X; x < boundsReq.Max.X; x++ {
				if x+xx > bounds.Max.X || y+yy > bounds.Max.Y {
					return
				}
				imgCroppped.Set(x, y, img.At(x+xx, y+yy))
			}
		}

		ext := filepath.Ext(strings.ToLower(f.Name()))
		f, err := os.Create(dir_merged + f.Name()[:len(f.Name())-len(ext)] + "_" + strconv.Itoa(cycle) + filepath.Ext(f.Name()))
		if err == nil {
			if filepath.Ext(strings.ToLower(f.Name())) == ".jpeg" || filepath.Ext(strings.ToLower(f.Name())) == ".jpg" {
				jpeg.Encode(f, imgCroppped, nil)
			} else if filepath.Ext(strings.ToLower(f.Name())) == ".png" {
				png.Encode(f, imgCroppped)
			} else {
				fmt.Println("wrong type")
			}
			f.Close()
		} else {
			fmt.Println(err)
		}

		xx += boundsReq.Max.X
		if xx+boundsReq.Max.X >= bounds.Max.X {
			xx = 0
			yy += boundsReq.Max.Y
		}
		cycle++
	}

}

func proc() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()

	dir_first = "./dataset_original/" //"E:\\dataset_cropped/"
	dir_second = "./noise_patterns/"  //"E:\\noise_patterns/"
	dir_merged = "./dataset/noise/"

	//entries, err := os.ReadDir(dir_first)

	f_list1, err := ioutil.ReadDir(dir_first)
	if err != nil {
		panic(err)
	}

	/*f_list2, err := ioutil.ReadDir(dir_second)
	if err != nil {
		panic(err)
	}*/

	stack := make([]int, len(f_list1))
	for i := 0; i < len(f_list1); i++ {
		stack[i] = i
	}
	work := make(chan int)
	//results := make(chan int)

	wg := sync.WaitGroup{}
	for cpu := 0; cpu < runtime.NumCPU(); cpu++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range work {
				if i < len(f_list1) {
					//process(f_list1[i], f_list2[i%len(f_list2)])
					cut(f_list1[i], image.Rectangle{image.Point{0, 0}, image.Point{1024, 1024}})
				}
			}
		}()
	}

	// send the work to the workers
	// this happens in a goroutine in order
	// to not block the main function, once
	// all 5 workers are busy

	go func() {
		for _, s := range stack {
			// could read the file from disk
			// here and pass a pointer to the file
			work <- s
		}
		// close the work channel after
		// all the work has been send
		close(work)

		// wait for the workers to finish
		// then close the results channel
		//wg.Wait()
		//close(results)
	}()
	wg.Wait()
	// collect the results
	// the iteration stops if the results
	// channel is closed and the last value
	// has been received

	//for result := range results {
	// could write the file to disk
	//	fmt.Println(result)
	//}
	elapsed := time.Since(start)
	fmt.Printf("Time spent: %s\n", elapsed)
}