package main

import (
	"testing"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
)

// helper to generate date strings
func formatDate(daysFromNow int) string {
	return time.Now().AddDate(0, 0, daysFromNow).Format("2006-01-02")
}

func TestCheckProjectDate_TodayOrPast(t *testing.T) {
	project := models.Project{Name: "TodayProject", Date: formatDate(0)}
	err := CheckProjectDate(project)
	if err == nil || !contains(err.Error(), "project passed") {
		t.Errorf("expected 'project passed' error, got %v", err)
	}
}

func TestCheckProjectDate_NearFuture(t *testing.T) {
	project := models.Project{Name: "FutureProject", Date: formatDate(3)}
	err := CheckProjectDate(project)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCheckProjectDate_FarFuture(t *testing.T) {
	project := models.Project{Name: "FarProject", Date: formatDate(10)}
	err := CheckProjectDate(project)
	if err == nil || !contains(err.Error(), "project too far") {
		t.Errorf("expected 'project too far' error, got %v", err)
	}
}

func TestCheckProjectDate_InvalidFormat(t *testing.T) {
	project := models.Project{Name: "InvalidProject", Date: "2025-13-01"}
	err := CheckProjectDate(project)
	if err == nil || !contains(err.Error(), "invalid project date") {
		t.Errorf("expected 'invalid project date' error, got %v", err)
	}
}

// simple substring check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || (len(s) > len(substr) && string(s[:len(substr)]) == substr) || string(s[len(s)-len(substr):]) == substr || string(s) == substr)
}
