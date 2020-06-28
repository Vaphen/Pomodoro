package main

import (
	"fmt"
	"time"

	"aaronroehl.info/pomodoro/icons"
	"github.com/getlantern/systray"
)

type tickCallback func(remainingTime time.Duration)

type pomodoro struct {
	ticker   *time.Ticker
	duration time.Duration
	started  time.Time
	callback tickCallback
	stopped  bool
}

func (p *pomodoro) start(callback tickCallback) {
	p.ticker = time.NewTicker(time.Second)
	p.callback = callback
	p.started = time.Now()
	p.resume()
}

func (p *pomodoro) stop() {
	p.stopped = true
}

func (p *pomodoro) resume() {
	p.stopped = false
	go func(p *pomodoro) {
		for {
			select {
			case a := <-p.ticker.C:
				if p.stopped {
					p.started = p.started.Add(time.Second)
					continue
				}
				timeLeft := p.started.Add(p.duration).Sub(a)
				if timeLeft < 0 {
					p.stop()
				}
				p.callback(timeLeft)
			}
		}
	}(p)
}

func newPomodoro(duration time.Duration) *pomodoro {
	return &pomodoro{
		ticker:   nil,
		duration: duration,
	}
}

func main() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetTitle("🍅25:00")
	startItem, _ := addMenuItem("(Re-)start", "Starts the timer", icons.Starticon)
	pauseItem, _ := addMenuItem("Pause", "Pauses the timer", icons.Pauseicon)
	resumeItem, _ := addMenuItem("Resume", "Resume the timer", icons.Resumeicon)
	quitItem, _ := addMenuItem("Quit", "Quit the app", icons.Quiticon)

	resumeItem.Disable()
	pauseItem.Disable()

	pom := newPomodoro(time.Minute*25 + time.Second)

	for {
		select {
		case <-startItem.ClickedCh:
			pom.start(tick)
			startItem.Disable()
			pauseItem.Enable()
		case <-pauseItem.ClickedCh:
			pom.stop()
			resumeItem.Enable()
			startItem.Enable()
			pauseItem.Disable()
		case <-resumeItem.ClickedCh:
			pom.resume()
			resumeItem.Disable()
			pauseItem.Enable()
		case <-quitItem.ClickedCh:
			systray.Quit()
		}
	}
}

func tick(d time.Duration) {
	systray.SetTitle(fmt.Sprintf("🍅%s", d.Truncate(time.Second)))
}

func addMenuItem(title, tooltip string, icon []byte) (menuItem *systray.MenuItem, err error) {
	menuItem = systray.AddMenuItem(title, tooltip)
	menuItem.SetIcon(icon)
	return
}

func onExit() {}
