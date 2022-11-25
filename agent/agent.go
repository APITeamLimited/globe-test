package agent

import (
	"fmt"
	"os"

	"github.com/getlantern/systray"
	"github.com/pkg/browser"
)

func Run() ***REMOVED***
	fmt.Println("Running agent")

	systrayContent := func() ***REMOVED***
		logoIcon, err := os.ReadFile("agent/icon.png")
		if err != nil ***REMOVED***
			panic(err)
		***REMOVED***

		systray.SetIcon(logoIcon)
		systray.SetTitle("APITeam Agent")
		systray.SetTooltip("APITeam Agent")

		// Add non clickable menu item with name and icon
		mTitle := systray.AddMenuItem("About APITeam Agent", "About APITeam Agent")
		mTitle.SetIcon(logoIcon)

		systray.AddSeparator()
		mAbortAll := systray.AddMenuItem("Abort All", "Abort All")
		mQuit := systray.AddMenuItem("Quit", "Quit APITeam Agent")

		setJobCountFunc := func(count int) ***REMOVED***
			if count == 0 ***REMOVED***
				mAbortAll.Hide()
			***REMOVED*** else ***REMOVED***
				mAbortAll.SetTitle(fmt.Sprintf("Abort All (%d)", count))
				mAbortAll.Show()
			***REMOVED***
		***REMOVED***

		setJobCountFunc(0)

		go func() ***REMOVED***
			<-mTitle.ClickedCh
			// Open the URL in the default browser
			err := browser.OpenURL("https://apiteam.cloud/agent")
			if err != nil ***REMOVED***
				panic(err)
			***REMOVED***
		***REMOVED***()

		serverStoppedCh := runAgentServer(mQuit.ClickedCh, mAbortAll.ClickedCh, setJobCountFunc)

		// Wait for the server to stop before exiting
		<-serverStoppedCh
		systray.Quit()
	***REMOVED***

	systray.Run(systrayContent, func() ***REMOVED******REMOVED***)

***REMOVED***
