# pomodoro

This is a minior mutation of [pomodoro](https://github.com/progrium/macdriver/tree/main/macos/_examples/pomodoro), the modification includes:
1. change the work time to 40 minutes
2. change the break time to 10 minutes
3. made `localhost:58080/next` for trigger `nextClick` action, this would allow me to have some hotkey to cycle through timers.
	e.g: if you do `curl http://localhost:58080/next` the timer will jump to next one(from ready to work, etc).
4. add [topframe](https://github.com/progrium/topframe) to pop up the time up message

