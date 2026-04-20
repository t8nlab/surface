package sfImage

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	sfInput "github.com/t8nlab/surface/input"
)

/**
 * Single Image Pipeline Processing
 * Chains multiple operations without re-encoding
 */
func ImageProcess(input map[string]any) (any, error) {
	srcPath, err := sfInput.GetString(input, "src")
	if err != nil { return nil, err }

	outPath, _ := sfInput.GetString(input, "out")
	format, _ := sfInput.GetString(input, "format")
	quality, _ := sfInput.GetInt(input, "quality")
	if quality == 0 { quality = 85 }

	// Get steps
	steps, _ := input["steps"].([]any)

	// 1. OPEN
	var img image.Image
	if strings.HasPrefix(srcPath, "http") {
		img, err = downloadAndOpen(srcPath)
	} else if strings.HasPrefix(srcPath, "data:") {
		img, err = openBase64(srcPath)
	} else {
		img, err = imaging.Open(srcPath)
	}
	if err != nil { return nil, err }

	// 2. PROCESS STEPS
	// Maintain as NRGBA for best quality during manipulation
	workingImg := imaging.Clone(img)

	for _, s := range steps {
		step, ok := s.(map[string]any)
		if !ok { continue }

		action, _ := sfInput.GetString(step, "action")
		switch action {
		case "resize":
			w, _ := sfInput.GetInt(step, "width")
			h, _ := sfInput.GetInt(step, "height")
			workingImg = imaging.Resize(workingImg, w, h, imaging.Lanczos)
		case "crop":
			w, _ := sfInput.GetInt(step, "width")
			h, _ := sfInput.GetInt(step, "height")
			workingImg = imaging.Fill(workingImg, w, h, imaging.Center, imaging.Lanczos)
		case "grayscale":
			workingImg = imaging.Grayscale(workingImg)
		case "blur":
			sigma, _ := step["sigma"].(float64)
			if sigma == 0 { sigma = 1.0 }
			workingImg = imaging.Blur(workingImg, sigma)
		}
	}

	// 3. OUTPUT
	if outPath != "" {
		err = saveImage(workingImg, outPath, format, quality)
		if err != nil { return nil, err }
		return map[string]any{"status": "ok", "path": outPath}, nil
	}

	b64, err := encodeToFormat(workingImg, format, quality)
	if err != nil { return nil, err }
	return map[string]any{"status": "ok", "base64": b64}, nil
}

/**
 * Concurrent Batch Processing
 * Processes multiple images using parallel workers
 */
func ImageBatch(input map[string]any) (any, error) {
	items, ok := input["items"].([]any)
	if !ok { return nil, fmt.Errorf("items must be an array") }

	concurrency, _ := sfInput.GetInt(input, "concurrency")
	if concurrency <= 0 { concurrency = 4 }

	results := make([]any, len(items))
	var wg sync.WaitGroup
	queue := make(chan struct{idx int; data map[string]any}, len(items))

	for w := 0; w < concurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := range queue {
				res, err := ImageProcess(item.data)
				if err != nil {
					results[item.idx] = map[string]any{"error": err.Error(), "index": item.idx}
				} else {
					results[item.idx] = res
				}
			}
		}()
	}

	for i, it := range items {
		data, _ := it.(map[string]any)
		queue <- struct{idx int; data map[string]any}{i, data}
	}
	close(queue)
	wg.Wait()

	return results, nil
}

// LEGACY SUPPORT (Delegates to Pipeline)
func ImageResize(input map[string]any) (any, error) {
	w, _ := sfInput.GetInt(input, "width")
	h, _ := sfInput.GetInt(input, "height")
	input["steps"] = []any{
		map[string]any{"action": "resize", "width": w, "height": h},
	}
	return ImageProcess(input)
}

func ImageCrop(input map[string]any) (any, error) {
	w, _ := sfInput.GetInt(input, "width")
	h, _ := sfInput.GetInt(input, "height")
	input["steps"] = []any{
		map[string]any{"action": "crop", "width": w, "height": h},
	}
	return ImageProcess(input)
}

// HELPERS

func downloadAndOpen(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	return imaging.Decode(bytes.NewReader(data))
}

func openBase64(dataStr string) (image.Image, error) {
	parts := strings.Split(dataStr, ",")
	if len(parts) != 2 { return nil, fmt.Errorf("invalid base64 format") }
	
	unbased, err := base64.StdEncoding.DecodeString(parts[1])
	if err != nil { return nil, err }

	return imaging.Decode(bytes.NewReader(unbased))
}

func encodeToFormat(img image.Image, format string, quality int) (string, error) {
	var buf bytes.Buffer
	var mime string
	var err error

	f := strings.ToLower(format)
	if f == "" { f = "jpg" }

	switch f {
	case "webp":
		err = webp.Encode(&buf, img, &webp.Options{Quality: float32(quality)})
		mime = "image/webp"
	case "png":
		err = png.Encode(&buf, img)
		mime = "image/png"
	default:
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
		mime = "image/jpeg"
	}

	if err != nil { return "", err }
	return fmt.Sprintf("data:%s;base64,%s", mime, base64.StdEncoding.EncodeToString(buf.Bytes())), nil
}

func saveImage(img image.Image, path, format string, quality int) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) { os.MkdirAll(dir, 0755) }

	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()

	ext := strings.ToLower(filepath.Ext(path))
	if format != "" { ext = "." + strings.ToLower(format) }

	switch ext {
	case ".webp":
		return webp.Encode(f, img, &webp.Options{Quality: float32(quality)})
	case ".png":
		return png.Encode(f, img)
	case ".gif":
		return gif.Encode(f, img, nil)
	default:
		return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	}
}
