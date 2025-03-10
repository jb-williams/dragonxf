# Dragonxf
* This is a simple terminal d20 based dragon slayer game. Its all about the least amount of deaths your hero can have while trying to kill the dragon and advance.
* It's expected for the hero to die many times and you can keep going killing the dragon as many times as you want.
* You can also reset the save file at any point.

+-------------------------------------------------+
        Chadwick Strongpants, Dragon Slayer
+-------------------------------------------------+
        How many Deaths to slay the Dragon!
      Advancement at 100xp, 400xp, 700xp, etc...
+-------------------------------------------------+

Usage: dragonxf [flags]
Flags:
  -l, --list    List current stats
  -s, --stealth Attempt to stealth
  -a, --attack  Attack the dragon
  -r, --reset   Reset the game
  -h, --help    Show this help message

* make:
* `echo` shows all variables that the Makefile will use
* `makedir` will make the bin and pkg dir in the current directory to build the
	* `build` will build in the current repo
	* `install` will build in the current repo, then move it to the local bin path, then clean up the current repo.
