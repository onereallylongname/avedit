package projection

import "encoding/json"

// CloneResult holds the result of cloning a subtree.
type CloneResult struct {
	Root  *Node
	Nodes []*Node
}

// CloneSubtree deep-clones a node and all its descendants, assigning new IDs.
func CloneSubtree(proj *Projection, sourceID string) (*CloneResult, error) {
	source := proj.Get(sourceID)
	if source == nil {
		return nil, errNodeNotFound(sourceID)
	}

	var nodes []*Node
	root := cloneRecursive(proj, source, "", &nodes)
	return &CloneResult{Root: root, Nodes: nodes}, nil
}

func cloneRecursive(proj *Projection, node *Node, parentID string, collected *[]*Node) *Node {
	cloned := &Node{
		ID:       nextNodeID(),
		Kind:     node.Kind,
		ParentID: parentID,
		Children: make([]string, 0, len(node.Children)),
		Path:     nil, // paths are rebuilt after insertion
		Attrs:    cloneAttributes(node.Attrs),
	}
	*collected = append(*collected, cloned)

	for _, childID := range node.Children {
		child := proj.Get(childID)
		if child == nil {
			continue
		}
		clonedChild := cloneRecursive(proj, child, cloned.ID, collected)
		cloned.Children = append(cloned.Children, clonedChild.ID)
	}

	return cloned
}

// cloneAttributes deep-copies node attributes.
func cloneAttributes(attrs Attributes) Attributes {
	return Attributes{
		Native: deepClone(attrs.Native),
		Custom: cloneCustomMap(attrs.Custom),
	}
}

func cloneCustomMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	result := make(map[string]any, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

// deepClone performs a deep copy via JSON round-trip (safe for Avro values).
func deepClone(v any) any {
	if v == nil {
		return nil
	}
	data, err := json.Marshal(v)
	if err != nil {
		return v
	}
	var cloned any
	if err := json.Unmarshal(data, &cloned); err != nil {
		return v
	}
	return cloned
}
