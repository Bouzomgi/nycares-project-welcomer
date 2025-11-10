package main

import (
	"fmt"
	"time"

	"github.com/Bouzomgi/nycares-project-welcomer/internal/models"
	log "github.com/sirupsen/logrus"
)

func parseProjectDate(dateStr string) (time.Time, error) {
	projectDate, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		log.Errorf("invalid project date format: %v", err)
		return time.Time{}, err
	}
	return projectDate, nil
}

func CheckProjectDate(project models.Project) error {
	// Parse project date
	projectDate, err := parseProjectDate(project.Date)
	if err != nil {
		return fmt.Errorf("invalid project date for %s: %w", project.Name, err)
	}

	// Load EST timezone
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return fmt.Errorf("failed to load EST timezone: %w", err)
	}

	// Convert projectDate and now to EST
	nowEST := time.Now().In(loc)
	log.Infof("Current EST date: %s", nowEST.Format("2006-01-02"))

	// Compare only the date components
	today := time.Date(nowEST.Year(), nowEST.Month(), nowEST.Day(), 0, 0, 0, 0, loc)
	projectDay := time.Date(projectDate.Year(), projectDate.Month(), projectDate.Day(), 0, 0, 0, 0, loc)
	log.Infof("Project day EST: %s", projectDay.Format("2006-01-02"))

	// Subtract days
	daysUntil := int(projectDay.Sub(today).Hours() / 24)
	log.Infof("Days until project %s: %d", project.Name, daysUntil)

	switch {
	case daysUntil > 7:
		return fmt.Errorf("project too far")
	case daysUntil <= 0:
		return fmt.Errorf("project passed")
	}
	return nil
}
