package main

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/Kagami/go-face"
)

// This example shows the basic usage of the package: create an
// recognizer, recognize faces, classify them using few known ones.
func main() {
	if len(os.Args) != 3 {
		log.Fatal("find-faces <model dir> <image dir>")
	}
	// Init the recognizer.
	rec, err := face.NewRecognizer(os.Args[1])
	if err != nil {
		log.Fatalf("Can't init face recognizer: %v", err)
	}
	// Free the resources when you're finished.
	defer rec.Close()

	files, err := ioutil.ReadDir(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		// Recognize faces on that image.
		n := filepath.Join(os.Args[2], file.Name())
		faces, err := rec.RecognizeFile(n)
		if err != nil {
			log.Printf("Can't recognize %s: %v\n", file.Name(), err)
			continue
		}
		if len(faces) == 0 {
			log.Printf("No face in %s", file.Name())
			continue
		}

		fi, err := os.Open(n)
		if err != nil {
			log.Fatal("Cannot open input file '", n, "':", err)
		}
		defer fi.Close()

		img, _, err := image.Decode(fi)
		if err != nil {
			log.Fatal("Cannot decode image at '", n, "':", err)
		}

		for i, face := range faces {
			log.Printf("%s's #%d faces: (%d, %d) - (%d-%d)\n", file.Name(), i + 1,
				face.Rectangle.Min.X, face.Rectangle.Min.Y, face.Rectangle.Max.X, face.Rectangle.Max.Y)

			fo, err := os.Create(fmt.Sprintf("%s_face_%d.jpg", path.Base(file.Name()), i))
			if err != nil {
				log.Fatalf("Cannot create output image: %v", err)
			}
			defer fo.Close()
			nImg := imaging.Crop(img, face.Rectangle)
			if err := jpeg.Encode(fo, nImg, nil); err != nil {
				log.Fatalf("can not encode image: %v", err)
			}
		}
	}
}
