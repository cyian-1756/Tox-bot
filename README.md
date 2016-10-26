# go tox bot

A simple tox bot written in go

## Commands

![password] | Authenticates you to the bot

!check auth | Checks if you are authenticated

!unauth | unauthenticates you

!exit | exits

!shell | runs command

!open_tray | Opens the CD tray

!close_tray | Closes the CD tray

!screenshot | Takes screenshot

!os_check | Returns OS information (lsb_release -a)

!check_python_install | Checks what if any versions of python are installed

!check_go_install | Checks if go is installed

!check_mono_install | Checks if mono and mcs are installed

!detect_de | Detects the desktop environment on the end system 

!get_running_dir | Gets the directory the bot is running from

## Compiling

Install github.com/kitech/go-toxcore and github.com/vova616/screenshot

`go get github.com/kitech/go-toxcore`

`go get github.com/vova616/screenshot`

Compile with ./build.sh

## Usage

Run with `./bot`

The bots ID will be printed to the terminal

Add the bot by sending a tox friend request to the bot with the password in the message field (The default password is password)
