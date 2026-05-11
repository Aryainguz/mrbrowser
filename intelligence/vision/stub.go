// Package vision provides a stub for future OCR and visual element detection.
// Currently, the interface is defined but not implemented.
// Integration with Tesseract OCR is planned for Phase 5.
package vision

import "github.com/mrbrowser/mrbrowser/core/browser"

// Analyzer defines the interface for visual page analysis.
type Analyzer interface {
	// AnalyzeScreenshot processes a screenshot and returns detected elements.
	AnalyzeScreenshot(imgBytes []byte) ([]*browser.Element, error)

	// ExtractText runs OCR on a screenshot or region and returns the detected text.
	ExtractText(imgBytes []byte) (string, error)
}

// StubAnalyzer is a no-op implementation of Analyzer.
// Replace with TesseractAnalyzer in Phase 5.
type StubAnalyzer struct{}

// NewStubAnalyzer returns a stub analyzer.
func NewStubAnalyzer() *StubAnalyzer {
	return &StubAnalyzer{}
}

// AnalyzeScreenshot is not yet implemented.
func (s *StubAnalyzer) AnalyzeScreenshot(_ []byte) ([]*browser.Element, error) {
	return nil, nil
}

// ExtractText is not yet implemented.
func (s *StubAnalyzer) ExtractText(_ []byte) (string, error) {
	return "", nil
}
