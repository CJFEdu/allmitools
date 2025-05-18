// Package mocks provides mock data for testing the AllMiTools server
package mocks

// MockToolInfo represents a collection of mock tool information for testing
var MockToolInfo = map[string]interface{}{
	"random-number": map[string]interface{}{
		"id":          "random-number",
		"name":        "Random Number Generator",
		"description": "Returns a random integer within a specified range.",
		"docPath":     "/docs/random-number",
		"toolPath":    "/tools/random-number",
		"parameters": []map[string]interface{}{
			{
				"name":        "min",
				"description": "Minimum value (inclusive).",
				"optional":    true,
				"default":     "0",
			},
			{
				"name":        "max",
				"description": "Maximum value (inclusive).",
				"optional":    true,
				"default":     "100",
			},
		},
	},
	"date": map[string]interface{}{
		"id":          "date",
		"name":        "Current Date Formatter",
		"description": "Returns the current date, optionally offset and formatted.",
		"docPath":     "/docs/date",
		"toolPath":    "/tools/date",
		"parameters": []map[string]interface{}{
			{
				"name":        "format",
				"description": "Date format (e.g., 'YYYY-MM-DD', 'MM/DD/YYYY', 'RFC3339').",
				"optional":    true,
				"default":     "RFC3339",
			},
			{
				"name":        "offset",
				"description": "Number of days to offset from current date (e.g., 1, -7).",
				"optional":    true,
				"default":     "0",
			},
		},
	},
	"day": map[string]interface{}{
		"id":          "day",
		"name":        "Current Day",
		"description": "Returns the current day of the month.",
		"docPath":     "/docs/day",
		"toolPath":    "/tools/day",
		"parameters":  []map[string]interface{}{},
	},
	"month": map[string]interface{}{
		"id":          "month",
		"name":        "Current Month",
		"description": "Returns the current month.",
		"docPath":     "/docs/month",
		"toolPath":    "/tools/month",
		"parameters":  []map[string]interface{}{},
	},
	"year": map[string]interface{}{
		"id":          "year",
		"name":        "Current Year",
		"description": "Returns the current year.",
		"docPath":     "/docs/year",
		"toolPath":    "/tools/year",
		"parameters":  []map[string]interface{}{},
	},
}

// MockToolOutputs represents mock outputs for each tool for testing
var MockToolOutputs = map[string]interface{}{
	"random-number": "42",
	"date":          "2025-05-18T16:40:53-05:00",
	"day":           "18",
	"month":         "May",
	"year":          "2025",
}

// MockAllTools returns a list of all mock tools for testing
func MockAllTools() []map[string]interface{} {
	var tools []map[string]interface{}
	
	for _, tool := range MockToolInfo {
		tools = append(tools, tool.(map[string]interface{}))
	}
	
	return tools
}
