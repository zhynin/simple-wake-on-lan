package main
 
import (
        "flag"
        "log"
        "net"
        "os"
)

var (
        mac string
        ip string
        port int
)

// -ldflags=""：表示将引号里面的参数传给编译器
// -s：去掉符号信息
// -w：去掉DWARF调试信息
// -H: 以windows gui形式打包，不带dos窗口
// go build -ldflags="-s -w -H windowsgui" main.go

func main() {
        flag.StringVar(&mac, "mac", "1A:2B:3C:4D:5E:6F", "mac address.")
        flag.StringVar(&ip, "ip", "255.255.255.255", "ip address.")
        flag.IntVar(&port, "port", 9, "udp port 1-65535.")
        flag.Parse()

        ip_address := net.ParseIP(ip)
        if ip_address == nil {
                log.Println("address " + ip + " invalid IP address!")
                os.Exit(0)
        }

        mac_address, err:= net.ParseMAC(mac)
        if err != nil {
                log.Println(err)
                os.Exit(0)
        }

        if port > 65535 || port < 1 {
                log.Println("port", port, "invalid udp port!")
                os.Exit(0)
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
                        Port: port,
                },
        )
        if err != nil {
                log.Println(err)
                return
        }

        defer socket.Close()

        _, err = socket.Write(magic_packet)
        if err != nil {
                log.Println(err)
                return
        }
        log.Println("sender ok!")
}
