package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/andreykaipov/goobs"
	"github.com/andreykaipov/goobs/api/requests/sources"

	_ "image/jpeg"
	_ "image/png"

	pigo "github.com/esimov/pigo/core"
)

type OBS struct {
	Client     *goobs.Client
	Classifier *pigo.Pigo
}

func NewOBS(host, password string) (*OBS, error) {
	var err error
	obs := &OBS{}

	obs.Client, err = goobs.New(
		host,
		goobs.WithPassword(password), // optional
		goobs.WithDebug(os.Getenv("OBS_DEBUG") != ""), // optional
	)
	if err != nil {
		return nil, err
	}

	facefinderCascade, err := ioutil.ReadFile("./cascade/facefinder")
	if err != nil {
		return nil, fmt.Errorf("Failed reading cascade file: %w", err)
	}

	obs.Classifier, err = pigo.NewPigo().Unpack(facefinderCascade)
	if err != nil {
		return nil, fmt.Errorf("Error reading the cascade file: %w", err)
	}

	return obs, err
}

func (o *OBS) takeScreenshot(source string) (io.Reader, error) {
	screenshot, err := o.Client.Sources.TakeSourceScreenshot(&sources.TakeSourceScreenshotParams{
		CompressionQuality: -1,
		EmbedPictureFormat: "png",
		FileFormat:         "png",
		SourceName:         source,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed taking screenshot: %w", err)
	}

	// remove the `data:image/png;base64,...` prefix
	screenshotData := screenshot.Img[strings.IndexByte(screenshot.Img, ',')+1:]

	return base64.NewDecoder(base64.StdEncoding, strings.NewReader(screenshotData)), nil
}

// DetectFace detects a face from an OBS source name. Adapated from the
// minimal example from Pigo: https://github.com/esimov/pigo#api.
//
// Currently we don't bother checking the threshold score for a face. If there's
// even a slight hint of a face, we return true.
func (o *OBS) DetectFace(source string) (bool, error) {
	encodedImage, err := o.takeScreenshot(source)
	if err != nil {
		return false, fmt.Errorf("Failed detecting face: %w", err)
	}

	img, err := pigo.DecodeImage(encodedImage)
	if err != nil {
		return false, fmt.Errorf("Cannot decode image: %w", err)
	}

	pixels := pigo.RgbToGrayscale(img)
	cols, rows := img.Bounds().Max.X, img.Bounds().Max.Y

	params := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     1000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	angle := 0.0 // cascade rotation angle
	iou := 0.0   // intersection over union
	detections := o.Classifier.RunCascade(params, angle)
	detections = o.Classifier.ClusterDetections(detections, iou)

	return len(detections) > 0, nil
}
