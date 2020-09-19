# GOTROLLER

MPRIS controller to Polybar wrote in Go.

![screenshot](https://i.imgur.com/FFmjBnu.png)

![screenshot of the gui](https://i.imgur.com/LYGRtex.png)


## Features
- Play/Pause
- Volume control (scrollling)
- Playlist control
- Restart song (right click the previous song button)
- Player selector
- GUI with progress bar and thumbnail
- "Disable" mode

## Instalation
Clone this repository, inside the repository folder run `go build`, then copy 
the gotroller binary to your `/usr/bin/`.

Now define the module in your polybar config:
```
[module/gotroller]
type = custom/script
exec = gotroller
interval = 1
format-underline = #8be9fd
```
(you can customize the format-underline and the interval).

You can open the GUI by running `gotroller gui`.

