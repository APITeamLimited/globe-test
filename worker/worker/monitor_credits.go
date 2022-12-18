package worker

import (
	"time"

	"github.com/APITeamLimited/globe-test/lib"
	"github.com/APITeamLimited/globe-test/worker/libWorker"
)

// Monitors credits in account and stops execution if there are no credits left
func monitorCredits(gs libWorker.BaseGlobalState, creditsManager *lib.CreditsManager) {
	// Every second check if we have enough credits to continue
	ticker := time.NewTicker(100 * time.Millisecond)

	if creditsManager == nil {
		return
	}

	go func() {
		for range ticker.C {
			creditsManager.ForceDeductCredits(gs.FuncModeInfo().Instance100MSUnitRate, true)
		}
	}()
}
