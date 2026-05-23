package helpers

import (
	"fmt"
	"image"
	"image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
)

// ImageSize defines standard image sizes
type ImageSize struct {
	Name   string
	Width  int
	Suffix string
}

// Standard image sizes
var StandardSizes = []ImageSize{
	{Name: "thumbnail", Width: 150, Suffix: "_thumb"},
	{Name: "medium", Width: 500, Suffix: "_medium"},
	{Name: "large", Width: 1000, Suffix: "_large"},
}

// ImageResult contains paths to processed images
type ImageResult struct {
	Original  string `json:"original"`
	Thumbnail string `json:"thumbnail,omitempty"`
	Medium    string `json:"medium,omitempty"`
	Large     string `json:"large,omitempty"`
}

// ProcessImageOptions configures image processing
type ProcessImageOptions struct {
	File      *multipart.FileHeader
	BasePath  string // Base storage path (e.g., "./public/img")
	Folder    string // Subfolder (e.g., "product")
	Resize    string // "all", "thumbnail", "medium", "large", or empty for original only
	Quality   int    // WebP quality (1-100, default 80)
	URLPrefix string // URL prefix for response (e.g., "/image")
}

// ProcessImage handles image upload, resize, and webp conversion
func ProcessImage(opts ProcessImageOptions) (*ImageResult, error) {
	// Set defaults
	if opts.BasePath == "" {
		opts.BasePath = "./public/img"
	}
	if opts.Quality == 0 {
		opts.Quality = 80
	}
	if opts.URLPrefix == "" {
		opts.URLPrefix = "/image"
	}

	// Determine full storage path
	storagePath := opts.BasePath
	if opts.Folder != "" {
		storagePath = filepath.Join(opts.BasePath, opts.Folder)
	}

	// Create directory if not exists
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(opts.File.Filename))
	allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
	if !allowedExts[ext] {
		return nil, fmt.Errorf("unsupported file format: %s (allowed: jpg, jpeg, png, gif)", ext)
	}

	// Open uploaded file
	src, err := opts.File.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Generate filename from original name (sanitized) + timestamp for uniqueness
	baseFilename := generateFilename(opts.File.Filename)

	// URL prefix with folder
	urlPrefix := opts.URLPrefix
	if opts.Folder != "" {
		urlPrefix = opts.URLPrefix + "/" + opts.Folder
	}

	result := &ImageResult{}

	// Handle GIF separately (preserve animation)
	if ext == ".gif" {
		return processGIF(src, storagePath, baseFilename, urlPrefix)
	}

	// Decode image
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Save original as webp
	originalPath := filepath.Join(storagePath, baseFilename+".webp")
	if err := saveWebP(img, originalPath, opts.Quality); err != nil {
		return nil, fmt.Errorf("failed to save original: %w", err)
	}
	result.Original = urlPrefix + "/" + baseFilename + ".webp"

	// Process resize based on options
	if opts.Resize == "" || opts.Resize == "original" {
		return result, nil
	}

	for _, size := range StandardSizes {
		// Skip if not requested
		if opts.Resize != "all" && opts.Resize != size.Name {
			continue
		}

		// Resize image
		resized := imaging.Resize(img, size.Width, 0, imaging.Lanczos)

		// Save resized as webp
		resizedPath := filepath.Join(storagePath, baseFilename+size.Suffix+".webp")
		if err := saveWebP(resized, resizedPath, opts.Quality); err != nil {
			return nil, fmt.Errorf("failed to save %s: %w", size.Name, err)
		}

		// Set result path
		resizedURL := urlPrefix + "/" + baseFilename + size.Suffix + ".webp"
		switch size.Name {
		case "thumbnail":
			result.Thumbnail = resizedURL
		case "medium":
			result.Medium = resizedURL
		case "large":
			result.Large = resizedURL
		}
	}

	return result, nil
}

// processGIF handles GIF files (preserves as GIF, no webp conversion for animation)
func processGIF(src io.Reader, storagePath, baseFilename, urlPrefix string) (*ImageResult, error) {
	// Decode GIF
	g, err := gif.DecodeAll(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode GIF: %w", err)
	}

	// Save as GIF (preserve animation)
	gifPath := filepath.Join(storagePath, baseFilename+".gif")
	outFile, err := os.Create(gifPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create GIF file: %w", err)
	}
	defer outFile.Close()

	if err := gif.EncodeAll(outFile, g); err != nil {
		return nil, fmt.Errorf("failed to encode GIF: %w", err)
	}

	return &ImageResult{
		Original: urlPrefix + "/" + baseFilename + ".gif",
	}, nil
}

// saveWebP saves an image as WebP format
func saveWebP(img image.Image, path string, quality int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return webp.Encode(file, img, &webp.Options{Quality: float32(quality)})
}

// DeleteImage removes image files from storage
func DeleteImage(basePath, folder, filename string) error {
	storagePath := basePath
	if folder != "" {
		storagePath = filepath.Join(basePath, folder)
	}

	// Remove all variants
	suffixes := []string{"", "_thumb", "_medium", "_large"}
	exts := []string{".webp", ".gif"}

	for _, suffix := range suffixes {
		for _, ext := range exts {
			path := filepath.Join(storagePath, filename+suffix+ext)
			os.Remove(path) // Ignore errors (file may not exist)
		}
	}

	return nil
}

// generateFilename creates filename from original name + timestamp
// Example: "My Photo.jpg" -> "my_photo_1706936400"
func generateFilename(originalName string) string {
	// Remove extension
	name := strings.TrimSuffix(originalName, filepath.Ext(originalName))

	// Sanitize: lowercase, replace spaces/special chars with underscore
	reg := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	name = reg.ReplaceAllString(strings.ToLower(name), "_")
	name = strings.Trim(name, "_")

	return fmt.Sprintf("%s", name)
}
