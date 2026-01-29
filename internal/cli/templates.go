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
			tree.ConnectNodes(t, "n5", "n2", "")
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
			tree.AddNode(t, model.StartEnd, "Escalated")        // n8
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "no")
			tree.ConnectNodes(t, "n3", "n4", "")
			tree.ConnectNodes(t, "n2", "n5", "yes")
			tree.ConnectNodes(t, "n5", "n6", "")
			tree.ConnectNodes(t, "n6", "n4", "yes")
			tree.ConnectNodes(t, "n6", "n7", "no")
			tree.ConnectNodes(t, "n7", "n8", "")
			return t
		},
	},
	{
		Name:        "bug-triage",
		Description: "Bug triage and incident response",
		Build: func() *model.Tree {
			t := model.NewTree("bug-triage")
			tree.AddNode(t, model.StartEnd, "Bug reported")                // n1
			tree.AddNode(t, model.Decision, "Reproducible?")               // n2
			tree.AddNode(t, model.IO, "Request reproduction steps")        // n3
			tree.AddNode(t, model.Decision, "Severity?")                   // n4
			tree.AddNode(t, model.Action, "Page on-call engineer")         // n5
			tree.AddNode(t, model.Decision, "Known issue?")                // n6
			tree.AddNode(t, model.Action, "Link to existing ticket")       // n7
			tree.AddNode(t, model.Action, "Create investigation ticket")   // n8
			tree.AddNode(t, model.Decision, "Root cause found?")           // n9
			tree.AddNode(t, model.Action, "Write fix")                     // n10
			tree.AddNode(t, model.Action, "Escalate to senior engineer")   // n11
			tree.AddNode(t, model.Decision, "Tests pass?")                 // n12
			tree.AddNode(t, model.Action, "Deploy fix")                    // n13
			tree.AddNode(t, model.Action, "Revise fix")                    // n14
			tree.AddNode(t, model.IO, "Notify reporter")                   // n15
			tree.AddNode(t, model.StartEnd, "Resolved")                    // n16
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "no")
			tree.ConnectNodes(t, "n3", "n2", "")
			tree.ConnectNodes(t, "n2", "n4", "yes")
			tree.ConnectNodes(t, "n4", "n5", "critical")
			tree.ConnectNodes(t, "n4", "n6", "normal")
			tree.ConnectNodes(t, "n5", "n6", "")
			tree.ConnectNodes(t, "n6", "n7", "yes")
			tree.ConnectNodes(t, "n7", "n15", "")
			tree.ConnectNodes(t, "n6", "n8", "no")
			tree.ConnectNodes(t, "n8", "n9", "")
			tree.ConnectNodes(t, "n9", "n10", "yes")
			tree.ConnectNodes(t, "n9", "n11", "no")
			tree.ConnectNodes(t, "n11", "n9", "")
			tree.ConnectNodes(t, "n10", "n12", "")
			tree.ConnectNodes(t, "n12", "n13", "yes")
			tree.ConnectNodes(t, "n12", "n14", "no")
			tree.ConnectNodes(t, "n14", "n10", "")
			tree.ConnectNodes(t, "n13", "n15", "")
			tree.ConnectNodes(t, "n15", "n16", "")
			return t
		},
	},
	{
		Name:        "hiring",
		Description: "Hiring pipeline",
		Build: func() *model.Tree {
			t := model.NewTree("hiring")
			tree.AddNode(t, model.StartEnd, "Application received")         // n1
			tree.AddNode(t, model.Decision, "Meets minimum qualifications?") // n2
			tree.AddNode(t, model.IO, "Send rejection email")               // n3
			tree.AddNode(t, model.Action, "Schedule phone screen")          // n4
			tree.AddNode(t, model.Decision, "Phone screen pass?")           // n5
			tree.AddNode(t, model.Action, "Schedule technical interview")   // n6
			tree.AddNode(t, model.Decision, "Technical pass?")              // n7
			tree.AddNode(t, model.Action, "Schedule team interview")        // n8
			tree.AddNode(t, model.Decision, "Team approval?")               // n9
			tree.AddNode(t, model.Action, "Prepare offer")                  // n10
			tree.AddNode(t, model.IO, "Send offer letter")                  // n11
			tree.AddNode(t, model.Decision, "Offer accepted?")              // n12
			tree.AddNode(t, model.Action, "Begin onboarding")              // n13
			tree.AddNode(t, model.StartEnd, "Hired")                       // n14
			tree.AddNode(t, model.Action, "Negotiate terms")               // n15
			tree.AddNode(t, model.StartEnd, "Candidate declined")          // n16
			tree.AddNode(t, model.StartEnd, "Not hired")                   // n17
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "no")
			tree.ConnectNodes(t, "n2", "n4", "yes")
			tree.ConnectNodes(t, "n4", "n5", "")
			tree.ConnectNodes(t, "n5", "n3", "no")
			tree.ConnectNodes(t, "n5", "n6", "yes")
			tree.ConnectNodes(t, "n6", "n7", "")
			tree.ConnectNodes(t, "n7", "n3", "no")
			tree.ConnectNodes(t, "n7", "n8", "yes")
			tree.ConnectNodes(t, "n8", "n9", "")
			tree.ConnectNodes(t, "n9", "n3", "no")
			tree.ConnectNodes(t, "n9", "n10", "yes")
			tree.ConnectNodes(t, "n10", "n11", "")
			tree.ConnectNodes(t, "n11", "n12", "")
			tree.ConnectNodes(t, "n12", "n13", "yes")
			tree.ConnectNodes(t, "n12", "n16", "no")
			tree.ConnectNodes(t, "n12", "n15", "counter")
			tree.ConnectNodes(t, "n15", "n11", "")
			tree.ConnectNodes(t, "n13", "n14", "")
			tree.ConnectNodes(t, "n3", "n17", "")
			return t
		},
	},
	{
		Name:        "medical-triage",
		Description: "Emergency room triage assessment",
		Build: func() *model.Tree {
			t := model.NewTree("medical-triage")
			tree.AddNode(t, model.StartEnd, "Patient arrives")                // n1
			tree.AddNode(t, model.Decision, "Conscious?")                     // n2
			tree.AddNode(t, model.Action, "Call code team")                   // n3
			tree.AddNode(t, model.Action, "Move to resuscitation bay")        // n4
			tree.AddNode(t, model.Decision, "Stabilized?")                    // n5
			tree.AddNode(t, model.Action, "Transfer to ICU")                  // n6
			tree.AddNode(t, model.StartEnd, "Critical care")                  // n7
			tree.AddNode(t, model.IO, "Record vitals")                        // n8
			tree.AddNode(t, model.Decision, "Vitals stable?")                 // n9
			tree.AddNode(t, model.Decision, "Pain level > 7?")                // n10
			tree.AddNode(t, model.Action, "Administer pain management")       // n11
			tree.AddNode(t, model.Decision, "Trauma or chest pain?")          // n12
			tree.AddNode(t, model.Action, "Fast-track to specialist")         // n13
			tree.AddNode(t, model.IO, "Notify attending physician")           // n14
			tree.AddNode(t, model.StartEnd, "Assessed")                       // n15
			tree.AddNode(t, model.IO, "Collect medical history")              // n16
			tree.AddNode(t, model.Decision, "Requires imaging?")              // n17
			tree.AddNode(t, model.Action, "Order imaging")                    // n18
			tree.AddNode(t, model.Action, "Assign to waiting room")           // n19
			tree.AddNode(t, model.IO, "Notify attending physician")           // n20
			tree.AddNode(t, model.StartEnd, "Assessed")                       // n21
			tree.SetRoot(t, "n1")
			// Unconscious path
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "no")
			tree.ConnectNodes(t, "n3", "n4", "")
			tree.ConnectNodes(t, "n4", "n5", "")
			tree.ConnectNodes(t, "n5", "n8", "yes")
			tree.ConnectNodes(t, "n5", "n6", "no")
			tree.ConnectNodes(t, "n6", "n7", "")
			// Conscious path
			tree.ConnectNodes(t, "n2", "n8", "yes")
			tree.ConnectNodes(t, "n8", "n9", "")
			tree.ConnectNodes(t, "n9", "n4", "no")
			tree.ConnectNodes(t, "n9", "n10", "yes")
			// Pain assessment
			tree.ConnectNodes(t, "n10", "n11", "yes")
			tree.ConnectNodes(t, "n11", "n12", "")
			tree.ConnectNodes(t, "n10", "n12", "no")
			// Trauma check
			tree.ConnectNodes(t, "n12", "n13", "yes")
			tree.ConnectNodes(t, "n13", "n14", "")
			tree.ConnectNodes(t, "n14", "n15", "")
			// Non-trauma path
			tree.ConnectNodes(t, "n12", "n16", "no")
			tree.ConnectNodes(t, "n16", "n17", "")
			tree.ConnectNodes(t, "n17", "n18", "yes")
			tree.ConnectNodes(t, "n18", "n20", "")
			tree.ConnectNodes(t, "n17", "n19", "no")
			tree.ConnectNodes(t, "n19", "n20", "")
			tree.ConnectNodes(t, "n20", "n21", "")
			return t
		},
	},
	{
		Name:        "loan-application",
		Description: "Loan approval decision process",
		Build: func() *model.Tree {
			t := model.NewTree("loan-application")
			tree.AddNode(t, model.StartEnd, "Application submitted")        // n1
			tree.AddNode(t, model.IO, "Pull credit report")                // n2
			tree.AddNode(t, model.Decision, "Credit score >= 650?")        // n3
			tree.AddNode(t, model.IO, "Send denial letter")                // n4
			tree.AddNode(t, model.Decision, "Debt-to-income < 40%?")      // n5
			tree.AddNode(t, model.Decision, "Employment verified?")        // n6
			tree.AddNode(t, model.IO, "Request employment docs")           // n7
			tree.AddNode(t, model.Decision, "Collateral required?")        // n8
			tree.AddNode(t, model.Action, "Order appraisal")              // n9
			tree.AddNode(t, model.Decision, "Appraisal sufficient?")      // n10
			tree.AddNode(t, model.Action, "Calculate loan terms")         // n11
			tree.AddNode(t, model.Decision, "Underwriter approved?")      // n12
			tree.AddNode(t, model.Action, "Generate loan documents")      // n13
			tree.AddNode(t, model.IO, "Send documents for signing")       // n14
			tree.AddNode(t, model.Action, "Disburse funds")               // n15
			tree.AddNode(t, model.StartEnd, "Loan closed")                // n16
			tree.AddNode(t, model.Action, "Flag for manual review")       // n17
			tree.AddNode(t, model.StartEnd, "Application denied")          // n18
			tree.SetRoot(t, "n1")
			tree.ConnectNodes(t, "n1", "n2", "")
			tree.ConnectNodes(t, "n2", "n3", "")
			tree.ConnectNodes(t, "n3", "n4", "no")
			tree.ConnectNodes(t, "n3", "n5", "yes")
			tree.ConnectNodes(t, "n5", "n4", "no")
			tree.ConnectNodes(t, "n5", "n6", "yes")
			tree.ConnectNodes(t, "n6", "n7", "no")
			tree.ConnectNodes(t, "n7", "n6", "")
			tree.ConnectNodes(t, "n6", "n8", "yes")
			tree.ConnectNodes(t, "n8", "n9", "yes")
			tree.ConnectNodes(t, "n9", "n10", "")
			tree.ConnectNodes(t, "n10", "n4", "no")
			tree.ConnectNodes(t, "n10", "n11", "yes")
			tree.ConnectNodes(t, "n8", "n11", "no")
			tree.ConnectNodes(t, "n11", "n12", "")
			tree.ConnectNodes(t, "n12", "n17", "no")
			tree.ConnectNodes(t, "n17", "n12", "")
			tree.ConnectNodes(t, "n12", "n13", "yes")
			tree.ConnectNodes(t, "n13", "n14", "")
			tree.ConnectNodes(t, "n14", "n15", "")
			tree.ConnectNodes(t, "n15", "n16", "")
			tree.ConnectNodes(t, "n4", "n18", "")
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
