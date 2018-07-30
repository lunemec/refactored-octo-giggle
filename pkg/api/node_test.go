package api_test

import (
	"encoding/json"
	"testing"

	"refactored-octo-giggle/pkg/api"

	"github.com/stretchr/testify/assert"
)

func testNode(t *testing.T) *api.Node {
	t.Helper()

	facet1 := api.Node{Name: "facet1"}
	facet2 := api.Node{Name: "facet2"}
	facet3 := api.Node{Name: "facet3", Parent: &facet1}
	facet4 := api.Node{Name: "facet4", Parent: &facet3}
	facet5 := api.Node{
		Name:   "facet5",
		Count:  50,
		Parent: &facet3,
	}
	facet6 := api.Node{
		Name:   "facet6",
		Count:  20,
		Parent: &facet4,
	}
	facet7 := api.Node{
		Name:   "facet7",
		Count:  30,
		Parent: &facet4,
	}
	root := api.Node{}
	root.Children = []*api.Node{
		&facet1,
		&facet2,
	}
	facet1.Parent = &root
	facet2.Parent = &root
	facet1.Children = []*api.Node{&facet3}
	facet3.Children = []*api.Node{&facet4, &facet5}
	facet4.Children = []*api.Node{&facet6, &facet7}
	return &root
}

func TestUnmarshalJSON(t *testing.T) {
	var node api.Node
	jsonData := `{
		"facet1": {
			"facet3": {
				"facet4": {
					"facet6": {
						"count": 20
					},
					"facet7": {
						"count": 30
					}
				},
				"facet5": {
					"count": 50
				}
			}
		}, 
		"facet2": {
			"count": 0
		}
	}`

	err := json.Unmarshal([]byte(jsonData), &node)
	assert.NoError(t, err, "unmarshal should not return error")

	// Sadly this can't be compared using assert.Equal and data from testNode()
	// because the children and parents are pointers, it compares pointer location
	// equality.
	// Also, since we parse it into map[string] which is not ordered, we can't
	// be certain that node.Children are in the same order, so we have to detect it.
	assert.Equal(t, "", node.Name, "node unmarshaled incorrectly")
	assert.True(t, node.Children[0].Name != "", "root node child name should not be empty")
	assert.Equal(t, 2, len(node.Children), "incorrect number of children for root node")
}

func TestToMap(t *testing.T) {
	expected := map[string]float64{
		"facet1": 10,
		"facet2": 10,
	}

	node := api.Node{
		Children: []*api.Node{
			&api.Node{
				Name: "facet1",
				Children: []*api.Node{
					&api.Node{
						Name:  "facet2",
						Count: 10,
					},
				},
			},
		},
	}

	assert.Equal(t, expected, node.ToMap(), "Node counts are incorrect")
}

func TestSumChildren(t *testing.T) {
	node := api.Node{
		Children: []*api.Node{
			&api.Node{
				Name: "facet1",
				Children: []*api.Node{
					&api.Node{
						Name:  "facet2",
						Count: 10,
					},
				},
			},
		},
	}
	assert.Equal(t, float64(10), node.SumChildren(), "children sum is incorrect")
}

// TestSumChildrenAll test that root's count is 100.
func TestSumChildrenAll(t *testing.T) {
	assert.Equal(t, float64(100), testNode(t).SumChildren(), "children sum is incorrect")
}
