/*
Define an interface Notifier with two methods:

    SendAlert(message string) string

    IsAvailable() bool — returns true if the notifier can send messages (e.g., is configured).

Implement three notifiers:

    SlackNotifier

    EmailNotifier

    PagerDutyNotifier

Each notifier should hold a Configured field (bool), which changes behavior:

    If not configured, IsAvailable() should return false.

    SendAlert() should still print an error-like message (e.g., “Slack not configured”).

Write a function NotifyAll(notifiers []Notifier, message string) that:

    Iterates through each notifier and only sends alerts if IsAvailable() returns true.

Use pointer receivers so the Configured state can be changed or tracked properly.
*/

package main

import "fmt"

type Notifier interface {
	SendAlert(message string) string
	IsAvailable() bool
}

// ----------
type SlackNotifier struct {
	Configured bool
}

type EmailNotifier struct {
	Configured bool
}

type PagerDutyNotifier struct {
	Configured bool
}

// ----------
func (s *SlackNotifier) IsAvailable() bool {
	return s.Configured
}

func (s *SlackNotifier) SendAlert(message string) string {
	defaultMessage := "Slack is not configured"

	if s.Configured {
		return fmt.Sprintf("Slack Alert: %s", message)
	}
	return defaultMessage
}

func (e *EmailNotifier) IsAvailable() bool {
	return e.Configured
}

func (e *EmailNotifier) SendAlert(message string) string {
	defaultMessage := "Email Notifier is not configured"

	if e.Configured {
		return fmt.Sprintf("Email Alert: %s", message)

	}

	return defaultMessage
}

func (p *PagerDutyNotifier) IsAvailable() bool {
	return p.Configured
}

func (p *PagerDutyNotifier) SendAlert(message string) string {
	defaultMessage := "Pager Duty is not configured"

	if p.Configured {
		return fmt.Sprintf("Pager Duty Alert: %s", message)

	}

	return defaultMessage
}

// ----------
func NotifyAll(notifiers []Notifier, message string) {
	for _, notifier := range notifiers {
		if notifier.IsAvailable() {
			fmt.Println(notifier.SendAlert(message))
		} else {
			fmt.Println(notifier.SendAlert(message))
		}
	}
}

// ----------
func main() {
	slackNotifierObj := &SlackNotifier{
		Configured: true,
	}

	emailNotifierObj := &EmailNotifier{
		Configured: true,
	}

	pagerDutyObj := &PagerDutyNotifier{
		Configured: false,
	}

	notifiers := []Notifier{slackNotifierObj, emailNotifierObj, pagerDutyObj}

	NotifyAll(notifiers, "🚨 High CPU usage detected on server 🚨")
}
