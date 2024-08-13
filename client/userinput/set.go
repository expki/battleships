package userinput

func (u *UserInput) setWindowWidth(width int) {
	if width == u.GetWindowWidth() {
		return
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	u.windowWidth = width
}

func (u *UserInput) setWindowHeight(height int) {
	if height == u.GetWindowHeight() {
		return
	}
	u.lock.Lock()
	defer u.lock.Unlock()
	u.windowHeight = height
}
