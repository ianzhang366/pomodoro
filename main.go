package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

const (
	defaultWorkTime  = 2400
	defaultBreakTime = 600
	defaultPort      = 58080
)

func main() {
	runtime.LockOSThread()

	app := cocoa.NSApp_WithDidLaunch(func(n objc.Object) {
		obj := cocoa.NSStatusBar_System().StatusItemWithLength(cocoa.NSVariableStatusItemLength)
		obj.Retain()
		obj.Button().SetTitle("▶️ Ready")

		nextClicked := make(chan bool)

		go func() {
			http.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
				nextClicked <- true
			})

			log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", defaultPort), nil))

		}()

		go func() {
			state := -1
			timer := defaultWorkTime
			countdown := false
			for {
				select {
				case <-time.After(1 * time.Second):
					if timer > 0 && countdown {
						timer = timer - 1
					}
					if timer <= 0 && state%2 == 1 {
						state = (state + 1) % 4
					}
				case <-nextClicked:
					state = (state + 1) % 4
					timer = map[int]int{
						0: defaultWorkTime,
						1: defaultWorkTime,
						2: 0,
						3: defaultBreakTime, // break timer
					}[state]
					if state%2 == 1 {
						countdown = true
					} else {
						countdown = false
					}
				}
				labels := map[int]string{
					0: "▶️ Ready %02d:%02d",
					1: "✴️ Working %02d:%02d",
					2: "✅ Finished %02d:%02d",
					3: "⏸️ Break %02d:%02d",
				}
				// updates to the ui should happen on the main thread to avoid strange bugs
				core.Dispatch(func() {
					obj.Button().SetTitle(fmt.Sprintf(labels[state], timer/60, timer%60))

				})
			}
		}()
		nextClicked <- true

		itemNext := cocoa.NSMenuItem_New()
		itemNext.SetTitle("Next")
		itemNext.SetAction(objc.Sel("nextClicked:"))
		cocoa.DefaultDelegateClass.AddMethod("nextClicked:", func(_ objc.Object) {
			nextClicked <- true
		})

		itemQuit := cocoa.NSMenuItem_New()
		itemQuit.SetTitle("Quit")
		itemQuit.SetAction(objc.Sel("terminate:"))

		menu := cocoa.NSMenu_New()

		menu.AddItem(itemNext)
		menu.AddItem(itemQuit)
		obj.SetMenu(menu)

	})
	app.Run()
}
