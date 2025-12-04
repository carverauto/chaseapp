package push

import (
	"archive/zip"
	"bytes"
	"crypto/sha1" //nolint:gosec // Safari manifest requires SHA1
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"path"
	"time"

	"chaseapp.tv/api/internal/config"
)

var iconPNG = []byte{ // 1x1 transparent PNG
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
	0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
	0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
	0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
	0x54, 0x78, 0x9c, 0x63, 0x00, 0x01, 0x00, 0x00,
	0x05, 0x00, 0x01, 0x0d, 0x0a, 0x2d, 0xb4, 0x00,
	0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44, 0xae,
	0x42, 0x60, 0x82,
}

// BuildSafariPackage creates a Safari push package ZIP in memory.
func BuildSafariPackage(cfg config.PushConfig) ([]byte, error) {
	if cfg.SafariPushID == "" || cfg.SafariWebURL == "" {
		return nil, errors.New("safari push configuration missing")
	}

	files := map[string][]byte{}

	website := map[string]any{
		"websiteName":         "ChaseApp",
		"websitePushID":       cfg.SafariPushID,
		"allowedDomains":      []string{cfg.SafariWebURL},
		"urlFormatString":     cfg.SafariWebURL + "/%@",
		"authenticationToken": "auth-token-placeholder",
		"webServiceURL":       cfg.SafariWebURL,
		"date":                time.Now().UTC().Format(time.RFC3339),
	}
	websiteJSON, err := json.Marshal(website)
	if err != nil {
		return nil, fmt.Errorf("marshal website.json: %w", err)
	}
	files["website.json"] = websiteJSON

	// Minimal iconset
	iconNames := []string{
		"icon_16x16.png",
		"icon_16x16@2x.png",
		"icon_32x32.png",
		"icon_32x32@2x.png",
		"icon_128x128.png",
		"icon_128x128@2x.png",
	}
	for _, name := range iconNames {
		files[path.Join("icon.iconset", name)] = iconPNG
	}

	manifest, err := buildManifest(files)
	if err != nil {
		return nil, err
	}
	files["manifest.json"] = manifest

	// Signature placeholder (would be PKCS7 in production)
	files["signature"] = []byte{}

	buf := &bytes.Buffer{}
	zipWriter := zip.NewWriter(buf)
	for name, content := range files {
		w, err := zipWriter.Create(name)
		if err != nil {
			return nil, fmt.Errorf("zip create %s: %w", name, err)
		}
		if _, err := w.Write(content); err != nil {
			return nil, fmt.Errorf("write %s: %w", name, err)
		}
	}
	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("close zip: %w", err)
	}

	return buf.Bytes(), nil
}

func buildManifest(files map[string][]byte) ([]byte, error) {
	manifest := make(map[string]string)
	for name, content := range files {
		if name == "manifest.json" || name == "signature" {
			continue
		}
		h := sha1.Sum(content)
		manifest[name] = hex.EncodeToString(h[:])
	}
	return json.Marshal(manifest)
}
