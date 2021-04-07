# GOTROLLER

MPRIS controller to Polybar wrote in Go.

![screenshot](https://i.imgur.com/YJ0pMbG.png)

![screenshot of the gui](https://i.imgur.com/HrEGG2E.png)

## Features

- Play/Pause
- Volume control (scrollling)
- Playlist control
- Restart song (right click the previous song button)
- Player selector
- GUI with progress bar and thumbnail
- "Disable" mode

## Instalation
First, make sure you have polybar and Font Awesome 5.

Then clone the repository and open it's folder:

> git clone https://github.com/Pauloo27/gotroller.git

> cd gotroller

If you don't want the GUI you can use dmenu to select the player to be 
displayed. To do that, run:
> make install-cli

If you want the GUI, run (the GUI takes sometime to compile):
> make install

_You can start the GUI by running gotroller-gui or clicking the "menu" icon
in the bar_

Now, add gotroller as a module in your `~/.config/polybar/config`:
```
[module/gotroller]
type = custom/script
exec = gotroller polybar-gui
tail = true
interval = 1
format-underline = #8be9fd
```

If you don't want to use the GUI, change `exec = gotroller polybar-gui` to 
`exec = gotroller polybar-dmenu`

Finally, restart polybar.

## License

<img src="https://i.imgur.com/AuQQfiB.png" alt="GPL Logo" height="100px" />

This project is licensed under [GNU General Public License v2.0](./LICENSE).

This program is free software; you can redistribute it and/or modify 
it under the terms of the GNU General Public License as published by 
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

