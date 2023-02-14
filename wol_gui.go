package main

import (
        "log"
	"strconv"
	"net"
        "fyne.io/fyne/v2"
        "fyne.io/fyne/v2/app"
        "fyne.io/fyne/v2/container"
        "fyne.io/fyne/v2/layout"
        "fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/data/binding"
)

func main() {
        myApp := app.New()
        myWindow := myApp.NewWindow("wake on lan")
        myWindow.Resize(fyne.NewSize(400, 208))
        myWindow.CenterOnScreen()

        mac := widget.NewLabel("MAC address:")
        mac_v := widget.NewEntry()
        mac_v.SetPlaceHolder("default: 1A:2B:3C:4D:5E:6F")
        ip := widget.NewLabel("IP address:")
        ip_v := widget.NewEntry()
        ip_v.SetPlaceHolder("default: 255.255.255.255")
        port := widget.NewLabel("UDP port:")
        port_v := widget.NewEntry()
        port_v.SetPlaceHolder("default: 9")
	
	exel := widget.NewLabel("Wake status:")
	exestr := binding.NewString()
	exestr.Set("Not started.")
	exetext := widget.NewLabelWithData(exestr)

        grid := container.New(layout.NewFormLayout(), mac, mac_v, ip, ip_v, port, port_v, exel, exetext)

        window := container.NewVBox(grid, widget.NewButton("Start UP", func() {
		execf := wakeUp(mac_v.Text, ip_v.Text, port_v.Text)
		if execf {
			exestr.Set("Magic packet sent successfully.")
		} else {
			exestr.Set("Magic packet sending failed.")
		}
        }))

        myWindow.SetContent(window)
        myWindow.ShowAndRun()
        tidyUp()
}

func tidyUp() {
        log.Println("Exited")
}

func wakeUp(mac string, ip string, port string) bool {

        flag := 0

        if len(ip) == 0 {ip = "255.255.255.255"}
        ip_address := net.ParseIP(ip)
        if ip_address == nil {
                flag = flag + 1
        }

        if len(mac) == 0 {mac = "1A:2B:3C:4D:5E:6F"}
        mac_address, err:= net.ParseMAC(mac)
        if err != nil {
                flag = flag + 2
        }

        if len(port) == 0 {port = "9"}
        uport, err := strconv.Atoi(port)
        if err != nil || uport > 65535 || uport < 1 {
                flag = flag + 4
        }

        if flag != 0 {
                log.Println("----------------------------------------------------------")
                log.Println("ip, mac, port values is invalid!")
                log.Println("ip address: ", ip)
                log.Println("mac address: ", mac)
                log.Println("UDP port: ", port)
                log.Println("----------------------------------------------------------")
                return false
        }

	var macpak []byte = make([]byte, 6)
        for i,v := range mac_address {
                macpak[i] = v
        }

        var packheader []byte = []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}
        var magic_packet []byte
        for i := 0; i < 16; i++ {
                magic_packet = append(magic_packet, macpak...)
        }
        magic_packet = append(packheader, magic_packet...)

        socket, err := net.DialUDP(
                "udp",
                nil,
                &net.UDPAddr{
                        IP: ip_address,
                        Port: uport,
                },
        )
        if err != nil {
                log.Println(err)
                return false
        }

        defer socket.Close()

        _, err = socket.Write(magic_packet)
        if err != nil {
                log.Println(err)
                return false
        }
	log.Println("sender magic packet to : ",  mac_address, "!")
	return true
}

