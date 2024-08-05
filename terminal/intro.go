package terminal

import (
	"math/rand"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
)

var titleFrames = []string{`
       ______   __  __    _   __  ______    ___     ____ 
      / ____/  / / / /   / | / / /_  __/   /   |   / __ \
     / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
    / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/ 
    \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|  
`, `
       ______   __  __    _   __  ______    ___     ____ 
      / ____/  / / / /   / | / / /_  __/   /   |   / __ \
      / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
    / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/ 
    \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|  
`, `
       ______   __  __    _   __  ______    ___     ____
      / ____/  / / / /   / | / / /_  __/   /   |   / __ \
     / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
   / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/
    \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|
`, `
       ______   __  __    _   __  ______    ___     ____
      / ____/  / / / /   / | / / /_  __/   /   |   / __ \
     / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
    / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/
   \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|
`, `
        ______   __  __    _   __  ______    ___     ____
      / ____/  / / / /   / | / / /_  __/   /   |   / __ \
     / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
    / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/
    \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|
`, `
       ______   __  __    _   __  ______    ___     ____
     / ____/  / / / /   / | / / /_  __/   /   |   / __ \
     / / __   / / / /   /  |/ /   / /     / /| |  / /_/ /
    / /_/ /  / /_/ /   / /|  /   / /     / ___ | / _, _/
    \____/   \____/   /_/ |_/   /_/     /_/  |_|/_/ |_|
`,
}

const defaultFrame = 0

type IntroModel struct {
	frames       []string
	currentFrame int
	timer        timer.Model
	quitting     bool
	keymap       KeyMap
}

func NewIntroModel() IntroModel {
	return IntroModel{
		frames:       titleFrames,
		timer:        timer.New(5 * time.Second),
		currentFrame: defaultFrame,
		quitting:     false,
		keymap:       DefaultKeyMap(),
	}
}

func (i *IntroModel) update(frame int, interval time.Duration) {
	i.currentFrame = frame
	i.timer.Interval = interval
}

func (i IntroModel) quit() (tea.Model, tea.Cmd) {
	i.quitting = true
	i.currentFrame = defaultFrame
	return i, tea.Quit
}

func (i IntroModel) Init() tea.Cmd {
	i.timer.Interval = 100 * time.Millisecond
	return i.timer.Init()
}

func (i IntroModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		if i.currentFrame != defaultFrame {
			i.update(defaultFrame, 200*time.Millisecond)
		} else {
			i.update(rand.Intn(len(titleFrames)), time.Duration(rand.Intn(100)*int(time.Millisecond)))
		}
		i.timer, cmd = i.timer.Update(msg)
		return i, cmd

	case timer.TimeoutMsg:
		return i.quit()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, i.keymap.Quit):
			return i.quit()
		}
	}
	return i, nil
}

func (i IntroModel) View() string {
	return i.frames[i.currentFrame]
}
