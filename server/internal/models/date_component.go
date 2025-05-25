// Package models contains the data structures for the AllMiTools server
package models

// DayTool returns the tool info for the day component
func DayTool() ToolInfo {
	return ToolInfo{
		Name:        "day",
		Description: "Get the current day of the month",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters:  []ToolParameter{}, // No parameters needed
	}
}

// MonthTool returns the tool info for the month component
func MonthTool() ToolInfo {
	return ToolInfo{
		Name:        "month",
		Description: "Get the current month as a string",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters:  []ToolParameter{}, // No parameters needed
	}
}

// YearTool returns the tool info for the year component
func YearTool() ToolInfo {
	return ToolInfo{
		Name:        "year",
		Description: "Get the current year",
		Version:     "1.0.0",
		Author:      "AllMiTools Team",
		Parameters:  []ToolParameter{}, // No parameters needed
	}
}
