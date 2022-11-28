package agent

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
)

func Run() {
	fmt.Println("Running agent")

	systrayContent := func() {
		logoIcon, err := os.ReadFile("apiteam-logo.png")
		if err != nil {
			panic(err)
		}

		systray.SetIcon(logoIcon)
		systray.SetTitle("APITeam Agent")
		systray.SetTooltip("APITeam Agent")

		// Add non clickable menu item with name and icon
		mTitle := systray.AddMenuItem("About APITeam Agent", "About APITeam Agent")
		mTitle.SetIcon(logoIcon)

		systray.AddSeparator()
		mAbortAll := systray.AddMenuItem("Abort All", "Abort All")
		mQuit := systray.AddMenuItem("Quit", "Quit APITeam Agent")

		setJobCountFunc := func(count int) {
			if count == 0 {
				mAbortAll.Hide()
			} else {
				mAbortAll.SetTitle(fmt.Sprintf("Abort All (%d)", count))
				mAbortAll.Show()
			}
		}

		setJobCountFunc(0)

		go func() {
			<-mTitle.ClickedCh
			// Open the URL in the default browser
			err := browser.OpenURL("https://apiteam.cloud/agent")
			if err != nil {
				panic(err)
			}
		}()

		setupChildProcesses()
		serverStoppedCh := runAgentServer(mQuit.ClickedCh, mAbortAll.ClickedCh, setJobCountFunc)

		// Wait for the server to stop before exiting
		<-serverStoppedCh
		systray.Quit()
	}

	systray.Run(systrayContent, func() {})

}
