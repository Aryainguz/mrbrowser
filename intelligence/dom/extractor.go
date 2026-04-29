package dom

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// Extractor extracts the full element tree from a live browser page.
type Extractor struct {
	page *browser.Page
	log  *telemetry.Logger
}

// NewExtractor creates an Extractor for the given page.
func NewExtractor(page *browser.Page) *Extractor {
	return &Extractor{
		page: page,
		log:  telemetry.New("dom"),
	}
}

// extractionScript is the JS injected into the page to collect element metadata in one round-trip.
const extractionScript = `
(function() {
	var results = [];
	var seen = new WeakSet();

	function getSection(el) {
		var ancestor = el.parentElement;
		while (ancestor) {
			var tag = ancestor.tagName ? ancestor.tagName.toLowerCase() : '';
			var role = ancestor.getAttribute ? (ancestor.getAttribute('role') || '') : '';
			if (tag === 'header' || role === 'banner') return 'header';
			if (tag === 'nav'    || role === 'navigation') return 'nav';
			if (tag === 'footer' || role === 'contentinfo') return 'footer';
			if (tag === 'main'   || role === 'main') return 'main';
			if (tag === 'aside'  || role === 'complementary') return 'aside';
			ancestor = ancestor.parentElement;
		}
		return 'body';
	}

	function getNearbyText(el) {
		var nearby = [];
		var parent = el.parentElement;
		if (parent) {
			Array.from(parent.children).forEach(function(sibling) {
				if (sibling !== el) {
					var t = (sibling.innerText || sibling.textContent || '').trim();
					if (t && t.length <= 80) nearby.push(t.substring(0, 60));
				}
			});
		}
		return nearby.slice(0, 5);
	}

	function getUniqueSelector(el) {
		if (el.id) return '#' + el.id.replace(/[^a-zA-Z0-9_-]/g, '\\\\$&');
		var path = [];
		var current = el;
		var depth = 0;
		while (current && current.tagName && depth < 6) {
			var tag = current.tagName.toLowerCase();
			var parent = current.parentElement;
			if (parent) {
				var siblings = Array.from(parent.children).filter(function(c) { return c.tagName === current.tagName; });
				if (siblings.length > 1) {
					tag += ':nth-of-type(' + (siblings.indexOf(current) + 1) + ')';
				}
			}
			path.unshift(tag);
			current = current.parentElement;
			depth++;
		}
		return path.join(' > ');
	}

	function isVisible(el) {
		var rect = el.getBoundingClientRect();
		if (rect.width === 0 && rect.height === 0) return false;
		try {
			var style = window.getComputedStyle(el);
			if (style.display === 'none' || style.visibility === 'hidden' || parseFloat(style.opacity) === 0) return false;
		} catch(e) {}
		return true;
	}

	var selectors = 'a, button, input, select, textarea, h1, h2, h3, h4, h5, h6, [role="button"], [role="link"], [role="textbox"], [role="combobox"], [role="checkbox"], [tabindex]';
	var elements;
	try {
		elements = document.querySelectorAll(selectors);
	} catch(e) { return '[]'; }

	var idx = 0;
	elements.forEach(function(el) {
		if (seen.has(el)) return;
		seen.add(el);

		var tag = el.tagName.toLowerCase();
		var inputType = el.getAttribute ? (el.getAttribute('type') || '') : '';
		var rect = el.getBoundingClientRect();
		var visible = isVisible(el);

		var ariaLabel = el.getAttribute ? (el.getAttribute('aria-label') || '') : '';
		var rawText = (el.innerText || el.textContent || '').trim();
		var text = (ariaLabel || rawText || el.getAttribute('placeholder') || '').substring(0, 200);

		var role = el.getAttribute ? (el.getAttribute('role') || '') : '';

		var isClickable = (tag === 'button' || tag === 'a' ||
			inputType === 'submit' || inputType === 'button' || inputType === 'reset' ||
			role === 'button' || role === 'link' ||
			el.onclick !== null || el.hasAttribute('tabindex'));

		results.push({
			node_id:     idx++,
			tag:         tag,
			input_type:  inputType,
			text:        text,
			role:        role,
			aria_label:  ariaLabel,
			visible:     visible,
			clickable:   isClickable,
			position:    {x: Math.round(rect.left), y: Math.round(rect.top), width: Math.round(rect.width), height: Math.round(rect.height)},
			attributes: {
				id:              el.id || '',
				class:           typeof el.className === 'string' ? el.className : '',
				name:            el.getAttribute ? (el.getAttribute('name') || '') : '',
				type:            inputType,
				href:            el.getAttribute ? (el.getAttribute('href') || '') : '',
				placeholder:     el.getAttribute ? (el.getAttribute('placeholder') || '') : '',
				'aria-label':    ariaLabel,
				disabled:        el.disabled ? 'true' : ''
			},
			placeholder:  el.getAttribute ? (el.getAttribute('placeholder') || '') : '',
			value:        el.value || '',
			selector:     getUniqueSelector(el),
			section:      getSection(el),
			nearby_text:  getNearbyText(el)
		});
	});

	return JSON.stringify(results);
})()
`

// jsElement mirrors the JSON shape returned by extractionScript.
type jsElement struct {
	NodeID      int               `json:"node_id"`
	Tag         string            `json:"tag"`
	InputType   string            `json:"input_type"`
	Text        string            `json:"text"`
	Role        string            `json:"role"`
	AriaLabel   string            `json:"aria_label"`
	Visible     bool              `json:"visible"`
	Clickable   bool              `json:"clickable"`
	Position    browser.Rect      `json:"position"`
	Attributes  map[string]string `json:"attributes"`
	Placeholder string            `json:"placeholder"`
	Value       string            `json:"value"`
	Selector    string            `json:"selector"`
	Section     string            `json:"section"`
	NearbyText  []string          `json:"nearby_text"`
}

// Extract walks the live DOM via JS injection and returns a slice of PageElements.
// It runs a single JS round-trip to collect all needed metadata efficiently.
func (e *Extractor) Extract() ([]*PageElement, error) {
	e.log.Debug("Extracting DOM elements")

	var raw string
	if err := chromedp.Run(e.page.Ctx(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var res interface{}
			if err := chromedp.Evaluate(extractionScript, &res).Do(ctx); err != nil {
				return err
			}
			if s, ok := res.(string); ok {
				raw = s
			}
			return nil
		}),
	); err != nil {
		return nil, fmt.Errorf("dom extraction script: %w", err)
	}

	elements, err := parseJSElements(raw)
	if err != nil {
		return nil, fmt.Errorf("parse dom elements: %w", err)
	}

	e.log.Debug("DOM extracted", telemetry.F("count", len(elements)))
	return elements, nil
}

// parseJSElements parses the JSON string returned by extractionScript.
func parseJSElements(raw string) ([]*PageElement, error) {
	if strings.TrimSpace(raw) == "" || raw == "[]" {
		return nil, nil
	}

	var jsEls []jsElement
	if err := json.Unmarshal([]byte(raw), &jsEls); err != nil {
		return nil, fmt.Errorf("json unmarshal: %w", err)
	}

	result := make([]*PageElement, 0, len(jsEls))
	for _, je := range jsEls {
		role := ARIARole(je.Tag, je.Role)
		semType := SemanticType(je.Tag, je.InputType)

		el := &PageElement{
			NodeID:      je.NodeID,
			Tag:         je.Tag,
			Type:        semType,
			Text:        strings.TrimSpace(je.Text),
			Role:        role,
			Visible:     je.Visible,
			Clickable:   je.Clickable,
			Position:    je.Position,
			Attributes:  je.Attributes,
			Placeholder: je.Placeholder,
			Value:       je.Value,
			Selector:    je.Selector,
			Section:     je.Section,
			NearbyText:  je.NearbyText,
		}
		result = append(result, el)
	}

	return result, nil
}

// GetNodes returns the raw CDP nodes (used internally by the accessibility tree).
func (e *Extractor) GetNodes() ([]*cdp.Node, error) {
	var nodes []*cdp.Node
	if err := chromedp.Run(e.page.Ctx(),
		chromedp.Nodes(`*`, &nodes, chromedp.ByQueryAll),
	); err != nil {
		return nil, err
	}
	return nodes, nil
}
