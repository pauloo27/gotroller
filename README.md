# GOTROLLER

MPRIS controller, with GUI and Polybar integrations.

_wrote in go btw._

![screenshot](https://i.imgur.com/YJ0pMbG.png)

_polybar module_

![screenshot of the gui](https://i.imgur.com/uRI5Gos.png)

_waybar module + gui_

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

## Config

You can set the song title and artist max length, if the length is greater than
the defined one, it will be limited and "..." will be appended. The default
value is 30 (song title) and 20 (artist).

To set it, you can set a system env or just create the file
`~/.config/gotroller.env` with the following content:
```bash
GOTROLLER_MAX_ARTIST_SIZE=20
GOTROLLER_MAX_TITLE_SIZE=30
GOTROLLER_GUI_MAX_ARTIST_SIZE=20
GOTROLLER_GUI_MAX_TITLE_SIZE=30
```

_If the value is 0 or negative, the length will not be limited._


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

