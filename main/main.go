package main

import (
	"fmt"
	"log"
	"os"

	"github.com/liugenping/torrent"
)

func main() {
	torrentCount := len(os.Args) - 1
	if torrentCount == 0 {
		fmt.Println("usage: gotorrent file.torrent[ file2.torrent[ ...]]")
		return
	}

	// Open torrent file
	torrent := torrent.Torrent{}
	torrent.TimeForExit = 5
	if err := torrent.Open(os.Args[1]); err != nil {
		log.Print(err)
		log.Print(os.Args[0])
		return

	}

	fmt.Println("Name:", torrent.Name)
	//fmt.Println("Announce URL:", torrent.trackers.Front().Value.(*Tracker).announceURL)
	fmt.Println("Comment:", torrent.Comment)
	fmt.Printf("Total size: %.2f MB\n", float64(torrent.TotalSize/1024/1024))
	//fmt.Printf("Downloaded: %.2f MB (%.2f%%)\n", float64(torrent.GetDownloadedSize()/1024/1024), float64(torrent.CompletedPieces)*100/float64(len(torrent.Pieces)))

	torrent.Download()
}



package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"time"

	"github.com/liugenping/torrent"
)

var t torrent.Torrent

type Download int

func (download *Download) GetDownloadRate(args struct{}, reply *([]string))error {
	rate := fmt.Sprintf("%.2f%c", float64(t.CompletedPieces)*100/float64(len(t.Pieces)), '%')
	*reply = append(*reply, rate)
	return nil
}

func (download *Download) GetDownloadTotalSize(args struct{}, reply *([]string))error {
	size := fmt.Sprintf("%.2f%c KB", float64(t.TotalSize/1024))
	*reply = append(*reply, size)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: torrent file.torrent")
		return
	}

	// Open torrent file
	if len(os.Args) == 2 {
		t.TimeForExit = 5 //default is 5 minutes exited after finish downlaod
	} else {
		i, err := strconv.Atoi(os.Args[2])
		if err != nil || i <= 0 {
			log.Print("parameter 2 should be a unsigned integer")
			return
		}
		t.TimeForExit = uint32(i)
	}
	if err := t.Open(os.Args[1]); err != nil {
		log.Print(err)
		return
	}

	fmt.Printf("Download Name:%s", t.Name)
	fmt.Printf("Total size: %.2f MB\n", float64(t.TotalSize/1024/1024))

	//start rpc service for client get download rate
	download := new(Download)
	rpc.Register(download)
	rpc.HandleHTTP()

	listener, err := net.Listen("tcp", ":9988")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	go http.Serve(listener, nil)

	//every 5s to echo download process
	go func() {
		for t.CompletedPieces != len(t.Pieces) {
			fmt.Printf("\r[%s] Downloaded: %.2f%c", t.Name, float64(t.CompletedPieces)*100/float64(len(t.Pieces)), '%')
			time.Sleep(time.Second * time.Duration(5))
		}
		fmt.Printf("\r[%s] Downloaded: %.2f%c", t.Name, float64(t.CompletedPieces)*100/float64(len(t.Pieces)), '%')
	}()

	t.Download()
	listener.Close()
}
