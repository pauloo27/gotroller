package polybar

func WithDmenu() {
	loadMaxSizes()
	startMainLoop("gotroller dmenu-select")
}
