// Package accessibility extracts the accessibility tree from a live browser page via CDP.
package accessibility

import (
	"context"
	"fmt"

	"github.com/chromedp/cdproto/accessibility"
	"github.com/chromedp/chromedp"
	"github.com/mrbrowser/mrbrowser/core/browser"
	"github.com/mrbrowser/mrbrowser/telemetry"
)

// AXNode represents a single node in the accessibility tree.
type AXNode struct {
	NodeID   string            `json:"node_id"`
	Role     string            `json:"role"`
	Name     string            `json:"name"`
	Value    string            `json:"value,omitempty"`
	Children []*AXNode         `json:"children,omitempty"`
	Props    map[string]string `json:"props,omitempty"`
}

// Tree holds the full accessibility tree of a page.
type Tree struct {
	Root  *AXNode   `json:"root"`
	Nodes []*AXNode `json:"nodes"`
}

// Extractor extracts the accessibility tree from a page.
type Extractor struct {
	page *browser.Page
	log  *telemetry.Logger
}

// NewExtractor creates a new accessibility tree extractor.
func NewExtractor(page *browser.Page) *Extractor {
	return &Extractor{
		page: page,
		log:  telemetry.New("a11y"),
	}
}

// Extract fetches the full accessibility tree via CDP.
func (e *Extractor) Extract() (*Tree, error) {
	e.log.Debug("Extracting accessibility tree")

	var axNodes []*accessibility.Node
	if err := chromedp.Run(e.page.Ctx(),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			axNodes, err = accessibility.GetFullAXTree().Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("get accessibility tree: %w", err)
	}

	tree := buildTree(axNodes)
	e.log.Debug("Accessibility tree extracted", telemetry.F("nodes", len(axNodes)))
	return tree, nil
}

// buildTree converts flat CDP AX nodes into a hierarchical tree.
func buildTree(axNodes []*accessibility.Node) *Tree {
	nodeMap := make(map[accessibility.NodeID]*AXNode, len(axNodes))

	// First pass: create all nodes.
	for _, n := range axNodes {
		ax := &AXNode{
			NodeID: string(n.NodeId),
			Props:  make(map[string]string),
		}

		if n.Role != nil {
			ax.Role = string(n.Role.Value)
		}

		if n.Name != nil {
			ax.Name = fmt.Sprintf("%v", n.Name.Value)
		}

		if n.Value != nil {
			ax.Value = fmt.Sprintf("%v", n.Value.Value)
		}

		nodeMap[n.NodeId] = ax
	}

	// Second pass: link parent → children.
	var roots []*AXNode
	childSet := make(map[accessibility.NodeID]bool)

	for _, n := range axNodes {
		parent, ok := nodeMap[n.NodeId]
		if !ok {
			continue
		}
		for _, childID := range n.ChildIds {
			if child, ok := nodeMap[childID]; ok {
				parent.Children = append(parent.Children, child)
				childSet[childID] = true
			}
		}
	}

	// Find root nodes (no parent).
	for _, n := range axNodes {
		if !childSet[n.NodeID] {
			if node, ok := nodeMap[n.NodeID]; ok {
				roots = append(roots, node)
			}
		}
	}

	var root *AXNode
	if len(roots) == 1 {
		root = roots[0]
	} else if len(roots) > 1 {
		root = &AXNode{Role: "document", Children: roots}
	}

	// Flatten for easy access
	all := make([]*AXNode, 0, len(nodeMap))
	for _, n := range nodeMap {
		all = append(all, n)
	}

	return &Tree{Root: root, Nodes: all}
}

// FindByRole returns all nodes in the tree with the given role.
func (t *Tree) FindByRole(role string) []*AXNode {
	var result []*AXNode
	walkTree(t.Root, func(n *AXNode) {
		if n.Role == role {
			result = append(result, n)
		}
	})
	return result
}

// FindByName returns all nodes whose name contains the given string (case-insensitive).
func (t *Tree) FindByName(name string) []*AXNode {
	nameLower := toLower(name)
	var result []*AXNode
	walkTree(t.Root, func(n *AXNode) {
		if contains(toLower(n.Name), nameLower) {
			result = append(result, n)
		}
	})
	return result
}

func walkTree(node *AXNode, fn func(*AXNode)) {
	if node == nil {
		return
	}
	fn(node)
	for _, child := range node.Children {
		walkTree(child, fn)
	}
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		result[i] = c
	}
	return string(result)
}

func contains(s, sub string) bool {
	if len(sub) == 0 {
		return true
	}
	if len(s) < len(sub) {
		return false
	}
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
