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

	"github.com/disintegration/imaging"
	sfInput "github.com/t8nlab/surface/input"
)

/**
 * Native Image Resizing (Supports Local, Remote, and Base64 Output)
 */
func ImageResize(input map[string]any) (any, error) {
	srcPath, err := sfInput.GetString(input, "src")
	if err != nil { return nil, err }

	distPath, _ := sfInput.GetString(input, "dist") // Optional

	width, _ := sfInput.GetInt(input, "width")
	height, _ := sfInput.GetInt(input, "height")
	quality, _ := sfInput.GetInt(input, "quality")
	if quality == 0 { quality = 85 }

	var src image.Image
	if strings.HasPrefix(srcPath, "http") {
		src, err = downloadAndOpen(srcPath)
	} else {
		src, err = imaging.Open(srcPath)
	}
	if err != nil { return nil, err }

	dst := imaging.Resize(src, width, height, imaging.Lanczos)

	// Decision: Save to file or return Base64?
	if distPath != "" {
		err = saveImage(dst, distPath, quality)
		if err != nil { return nil, err }
		return map[string]any{"status": "ok", "path": distPath}, nil
	}

	// Buffer output (Return as Base64)
	b64, err := encodeToBase64(dst, quality)
	if err != nil { return nil, err }
	return map[string]any{"status": "ok", "base64": b64}, nil
}

/**
 * Native Smart Cropping (Supports Local, Remote, and Base64 Output)
 */
func ImageCrop(input map[string]any) (any, error) {
	srcPath, err := sfInput.GetString(input, "src")
	if err != nil { return nil, err }

	distPath, _ := sfInput.GetString(input, "dist") // Optional

	width, _ := sfInput.GetInt(input, "width")
	height, _ := sfInput.GetInt(input, "height")

	var src image.Image
	var errImg error
	if strings.HasPrefix(srcPath, "http") {
		src, errImg = downloadAndOpen(srcPath)
	} else {
		src, errImg = imaging.Open(srcPath)
	}
	if errImg != nil { return nil, errImg }

	dst := imaging.Fill(src, width, height, imaging.Center, imaging.Lanczos)

	if distPath != "" {
		err = saveImage(dst, distPath, 85)
		if err != nil { return nil, err }
		return map[string]any{"status": "ok", "path": distPath}, nil
	}

	b64, err := encodeToBase64(dst, 85)
	if err != nil { return nil, err }
	return map[string]any{"status": "ok", "base64": b64}, nil
}

func encodeToBase64(img image.Image, quality int) (string, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: quality})
	if err != nil { return "", err }

	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func downloadAndOpen(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil { return nil, err }
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: %s", resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil { return nil, err }

	return imaging.Decode(bytes.NewReader(data))
}

func saveImage(img image.Image, path string, quality int) error {
	ext := strings.ToLower(filepath.Ext(path))
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, 0755)
	}

	f, err := os.Create(path)
	if err != nil { return err }
	defer f.Close()

	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Encode(f, img, &jpeg.Options{Quality: quality})
	case ".png":
		return png.Encode(f, img)
	case ".gif":
		return gif.Encode(f, img, nil)
	default:
		return imaging.Save(img, path)
	}
}
