package cli

import (
	"github.com/jllovet/decision-tree-cli/internal/model"
	"github.com/jllovet/decision-tree-cli/internal/tree"
)

// treeTemplate describes a named decision tree template.
type treeTemplate struct {
	Name        string
	Description string
	Build       func() *model.Tree
}

var templates = []treeTemplate{
	{
		Name:        "auth-flow",
		Description: "Authentication flow",
		Build: func() *model.Tree {
			t := model.NewTree("auth-flow")
			tree.AddNode(t, model.StartEnd, "Start")         // n1
			tree.AddNode(t, model.Decision, "Authenticated?") // n2
			tree.AddNode(t, model.Action, "Grant access")     // n3
			tree.AddNode(t, model.StartEnd, "End")            // n4
			tree.AddNode(t, model.Action, "Show login form")  // n5
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "yes")
			tree.ConnectNodes(t, "n3", "n4", "")
			tree.ConnectNodes(t, "n2", "n5", "no")
			return t
		},
	},
	{
		Name:        "approval",
		Description: "Approval workflow",
		Build: func() *model.Tree {
			t := model.NewTree("approval")
			tree.AddNode(t, model.StartEnd, "Start")          // n1
			tree.AddNode(t, model.Action, "Submit request")    // n2
			tree.AddNode(t, model.Decision, "Approved?")       // n3
			tree.AddNode(t, model.Action, "Process request")   // n4
			tree.AddNode(t, model.StartEnd, "End")             // n5
			tree.AddNode(t, model.Action, "Revise request")    // n6
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "")
			tree.ConnectNodes(t, "n3", "n4", "yes")
			tree.ConnectNodes(t, "n4", "n5", "")
			tree.ConnectNodes(t, "n3", "n6", "no")
			tree.ConnectNodes(t, "n6", "n2", "")
			return t
		},
	},
	{
		Name:        "troubleshooting",
		Description: "Troubleshooting guide",
		Build: func() *model.Tree {
			t := model.NewTree("troubleshooting")
			tree.AddNode(t, model.StartEnd, "Start")           // n1
			tree.AddNode(t, model.Decision, "Is it plugged in?") // n2
			tree.AddNode(t, model.Action, "Plug it in")         // n3
			tree.AddNode(t, model.StartEnd, "Done")             // n4
			tree.AddNode(t, model.Action, "Check settings")     // n5
			tree.AddNode(t, model.Decision, "Resolved?")        // n6
			tree.AddNode(t, model.Action, "Escalate")           // n7
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "no")
			tree.ConnectNodes(t, "n3", "n4", "")
			tree.ConnectNodes(t, "n2", "n5", "yes")
			tree.ConnectNodes(t, "n5", "n6", "")
			tree.ConnectNodes(t, "n6", "n4", "yes")
			tree.ConnectNodes(t, "n6", "n7", "no")
			return t
		},
	},
}

// findTemplate returns the template with the given name, or nil if not found.
func findTemplate(name string) *treeTemplate {
	for i := range templates {
		if templates[i].Name == name {
			return &templates[i]
		}
	}
	return nil
}
