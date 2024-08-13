package userinput

func (u *UserInput) GetWindowWidth() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.windowWidth
}

func (u *UserInput) GetWindowHeight() int {
	u.lock.RLock()
	defer u.lock.RUnlock()
	return u.windowHeight
}
