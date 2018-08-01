package api

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/json-iterator/go"
	"github.com/kr/pretty"
	"github.com/pkg/errors"
)

// InputData is top level container for our Nodes.
type InputData struct {
	Data Node `json:"data"`
}

// Node represents input facet tree.
type Node struct {
	Name  string
	Count float64

	Parent   *Node
	Children []*Node
}

// BufferedChallengeHandler implements the same functionality as StreamingChallengeHandler
// but it buffers the entire body, and parses the json as a whole. This version
// is much more readable and extendable than the Streaming version.
func BufferedChallengeHandler(rw http.ResponseWriter, req *http.Request) error {
	var (
		out OutputJSON
	)

	rootNode, err := unmarshal(req.Body)
	defer closer(req.Body)

	if err != nil {
		return errors.Wrap(err, "unable to parse facets json")
	}
	out.Result = mapToSlice(rootNode.ToMap())

	enc := jsoniter.NewEncoder(rw)
	return enc.Encode(&out)
}

// unmarshal reads the input reader into a buffer and returns the Root Node
// containing the entire node tree.
func unmarshal(r io.Reader) (*Node, error) {
	var (
		input InputData
	)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.Wrap(err, "unable to read json body")
	}
	err = jsoniter.Unmarshal(b, &input)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse json")
	}

	return &input.Data, nil
}

// UnmarshalJSON implements json.Unmarshaler interface for our Node.
func (n *Node) UnmarshalJSON(b []byte) error {
	var v map[string]interface{}

	err := jsoniter.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	err = n.FromMap(v)
	if err != nil {
		return errors.Wrap(err, "error converting to node tree")
	}
	return nil
}

// FromMap builds the node tree from parsed json objects.
func (n *Node) FromMap(m map[string]interface{}) error {
	var (
		err error
		i   int
	)

	// Create slice of nodes of size len(input_map) to avoid reallocations.
	n.Children = make([]*Node, len(m))
	for k, v := range m {
		node := Node{
			Name:   k,
			Parent: n,
		}
		if inner, ok := v.(map[string]interface{}); ok {
			// If inner map contains key "count", it is not a node, but our
			// "attributes" map/struct.
			if v, ok := inner["count"]; ok {
				count, ok := v.(float64)
				if !ok {
					err = errors.Errorf("count value is invalid type: %+v %T", v, v)
					return err
				}
				node.Count = count
			} else {
				err := node.FromMap(inner)
				if err != nil {
					return err
				}
			}
		}
		n.Children[i] = &node
		i++
	}
	return err
}

// SumChildren returns the sum of all child counts.
func (n *Node) SumChildren() (sum float64) {
	// Last child returns its count.
	if len(n.Children) == 0 {
		return n.Count
	}
	for _, child := range n.Children {
		sum += child.SumChildren()
	}
	return
}

// ToMap goes over the Node tree and returns flattened map of node names and
// their count.
// The name might be better, since this might lead one to believe it is the
// exact opposite of FromMap and that they could be used FromMap(ToMap()),
// that is not the case.
func (n *Node) ToMap() (out map[string]float64) {
	out = make(map[string]float64)
	// No more children, this is either empty root or last child.
	// If it is root (n.Name == "" && n.Count == 0), we do nothing.
	if len(n.Children) == 0 && !n.IsRoot() {
		out[n.Name] = n.Count
		return
	}
	// Iterate over children and add their counts to my map.
	for _, child := range n.Children {
		childMap := child.ToMap()
		for k, v := range childMap {
			out[k] = v
		}
	}
	// Add the sum of my children to my map under this node name, but not for root.
	if !n.IsRoot() {
		out[n.Name] = n.SumChildren()
	}
	return
}

// IsRoot returns true if node name is empty.
func (n *Node) IsRoot() bool {
	return n.Name == ""
}

// nolint: vet
// String implements stringer interface.
func (n *Node) String() string {
	return pretty.Sprint(n)
}
