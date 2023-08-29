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

////////////////////////////////////////////////////////////////////////////////////
/*               Image helpers                                                    */
////////////////////////////////////////////////////////////////////////////////////

// Gets minimum image dimension (img1, img2). Used for new image creation
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

// Limits pixel's color value by specified value to prevent overflow
func clamp(value, min, max uint32) uint32 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}

// Returns pixel3 (RGBA) = pixel1 (RGBA) + pixel2 (RGBA)
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

////////////////////////////////////////////////////////////////////////////////////
/*               File helpers                                                     */
////////////////////////////////////////////////////////////////////////////////////

// Opens file and return image from it
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

////////////////////////////////////////////////////////////////////////////////////
/*               Operations                                                       */
////////////////////////////////////////////////////////////////////////////////////

// Returns img3 = img1 + img2
func mergeImage(img1 image.Image, img2 image.Image) (image.Image, error) {
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

// Creates new img file and fills it with img3 = img1 + img
func mergeFile(f1 os.FileInfo, f2 os.FileInfo, dir_in_1 string, dir_in_2 string, dir_merged string) {
	img1, err1 := openImage(dir_in_1 + f1.Name())
	img2, err2 := openImage(dir_in_2 + f2.Name())
	if err1 == nil && err2 == nil {
		img_res, err := mergeImage(img1, img2)
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

// Returns cropped image from img with coords: [xx, xx+boundsReq.Max.X], [yy, yy+boundsReq.Max.X]
func cutImage(img *image.Image, xx int, yy int, boundsReq image.Rectangle) image.Image {
	boundsImg := (*img).Bounds()
	imgCrop := image.NewRGBA(image.Rect(boundsReq.Min.X, boundsReq.Min.Y, boundsReq.Max.X, boundsReq.Max.Y))
	for y := boundsImg.Min.Y; y < boundsReq.Max.Y; y++ {
		for x := boundsImg.Min.X; x < boundsReq.Max.X; x++ {
			if x+xx > boundsImg.Max.X || y+yy > boundsImg.Max.Y {
				return nil
			}
			imgCrop.Set(x, y, (*img).At(x+xx, y+yy))
		}
	}
	return imgCrop
}

// Create n cropped files from given file and given image dimension.
// Example: if original file has 1024*1024 resolution and 256*256 dimension given,
// then 4 files with 256*256 resolution will be created.
func cutFile(f os.FileInfo, boundsReq image.Rectangle, dir_in string, dir_result string) {
	img, err := openImage(dir_in + f.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	if img.Bounds().Max.X < boundsReq.Max.X || img.Bounds().Max.Y < boundsReq.Max.Y {
		fmt.Println("image to small")
		return
	}
	bounds := img.Bounds()

	yy := bounds.Min.Y
	xx := bounds.Min.X
	cycle := 1

	for {
		imgCrop := cutImage(&img, xx, yy, boundsReq)
		// Reach the end of file
		if imgCrop == nil {
			return
		}
		ext := filepath.Ext(strings.ToLower(f.Name()))
		f, err := os.Create(dir_result + f.Name()[:len(f.Name())-len(ext)] + "_" + strconv.Itoa(cycle) + filepath.Ext(f.Name()))
		if err == nil {
			if filepath.Ext(strings.ToLower(f.Name())) == ".jpeg" || filepath.Ext(strings.ToLower(f.Name())) == ".jpg" {
				jpeg.Encode(f, imgCrop, nil)
			} else if filepath.Ext(strings.ToLower(f.Name())) == ".png" {
				png.Encode(f, imgCrop)
			} else {
				fmt.Println("wrong type")
			}
			f.Close()
		} else {
			fmt.Println(err)
		}

		xx += boundsReq.Max.X
		if xx+boundsReq.Max.X > bounds.Max.X {
			xx = 0
			yy += boundsReq.Max.Y
		}
		cycle++
	}
}

// Returns image with increased brightness by factor (8 = 3EV).
func increaseBrightnessImage(img image.Image, factor float64) image.Image {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			newR := clampFloat(float64(r)*factor, 0, 65535)
			newG := clampFloat(float64(g)*factor, 0, 65535)
			newB := clampFloat(float64(b)*factor, 0, 65535)

			/*

				if newR > 20000 && newG > 20000 && newB > 20000 {
					newColor := color.RGBA64{R: uint16(newR), G: uint16(newG), B: uint16(newB), A: uint16(a)}
					newImg.Set(x, y, newColor)
				} else {
					newColor := color.RGBA64{R: uint16(0), G: uint16(0), B: uint16(0), A: uint16(a)}
					newImg.Set(x, y, newColor)
				}

			*/

			newColor := color.RGBA64{R: uint16(newR), G: uint16(newG), B: uint16(newB), A: uint16(a)}
			newImg.Set(x, y, newColor)
		}
	}

	return newImg
}

func increaseBrightnessFile(fInfo os.FileInfo, factor int, dir_in string, dir_result string) {
	img, err := openImage(dir_in + fInfo.Name())
	if err != nil {
		fmt.Println(err)
		return
	}
	imgBrigthness := increaseBrightnessImage(img, float64(factor))
	ext := filepath.Ext(strings.ToLower(fInfo.Name()))
	f, err := os.Create(dir_result + fInfo.Name()[:len(fInfo.Name())-len(ext)] + "_" + strconv.Itoa(factor) + filepath.Ext(fInfo.Name()))
	if err == nil {
		if filepath.Ext(strings.ToLower(f.Name())) == ".jpeg" || filepath.Ext(strings.ToLower(f.Name())) == ".jpg" {
			jpeg.Encode(f, imgBrigthness, nil)
		} else if filepath.Ext(strings.ToLower(f.Name())) == ".png" {
			png.Encode(f, imgBrigthness)
		} else {
			fmt.Println("wrong type")
		}
		f.Close()
	} else {
		fmt.Println(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////
/*               Processes, which should be called by frontend                    */
////////////////////////////////////////////////////////////////////////////////////

// Main helper, provoke all needed function
func processFiles(dirIn1, dirIn2, dirOut string, action func(os.FileInfo, os.FileInfo, string, string, string)) {
	fmt.Println(dirIn1, dirIn2)
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()

	fileList1, err := ioutil.ReadDir(dirIn1)
	if err != nil {
		panic(err)
	}

	var fileList2 []os.FileInfo
	if dirIn2 != "" {
		fileList2, err = ioutil.ReadDir(dirIn2)
		if err != nil {
			panic(err)
		}
	}

	work := make(chan int, len(fileList1))
	for i := range fileList1 {
		work <- i
	}
	close(work)

	wg := sync.WaitGroup{}
	for cpu := 0; cpu < runtime.NumCPU(); cpu++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range work {
				if i < len(fileList1) {
					if len(fileList2) > 0 {
						action(fileList1[i], fileList2[i%len(fileList2)], dirIn1, dirIn2, dirOut)
					} else {
						action(fileList1[i], nil, dirIn1, "", dirOut)
					}
				}
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	fmt.Printf("Time spent: %s\n", elapsed)
}

// Creates Dir3 and fills it with images = (All images from dir1) + (All images from dir2).
// If count(dir2 images) < count(dir1 images), then merge algorithm proccess cyclically,
// repeating merge dir2 images to remaining dir1 images.
func ProcessMerge(dirIn1, dirIn2, dirMerged string) {
	processFiles(dirIn1, dirIn2, dirMerged, mergeFile)
}

// Cut all files from dir, placing resilt in dir_result
func ProcessCut(dirIn, dirResult string, x, y int) {
	cut := func(fileInfo os.FileInfo, _ os.FileInfo, dirIn, _, dirResult string) {
		cutFile(fileInfo, image.Rectangle{image.Point{0, 0}, image.Point{x, y}}, dirIn, dirResult)
	}
	processFiles(dirIn, "", dirResult, cut)
}

func ProcessBrightness(dirIn, dirResult string, factor int) {
	factorFunc := func(fileInfo os.FileInfo, _ os.FileInfo, dirIn, _, dirResult string) {
		increaseBrightnessFile(fileInfo, factor, dirIn, dirResult)
	}
	processFiles(dirIn, "", dirResult, factorFunc)
}
