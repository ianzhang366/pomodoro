package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
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
	VERSION          = "v0.0.1"
)

var version bool

func init() {
	flag.BoolVar(&version, "v", false, "print the binary version")
}

func main() {
	flag.Parse()

	if version {
		fmt.Println(VERSION)
		return
	}
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
			played := true
			for {
				select {
				case <-time.After(1 * time.Second):
					if state%2 == 0 && !played {
						go func() {
							if err := runTopframeBin(); err != nil {
								fmt.Println("failed to run topframe, err: ", err)
							}
						}()

						played = true
					}

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
						played = false
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

func runTopframeBin() error {
	return exec.Command("/usr/local/bin/topframe").Run()
}
