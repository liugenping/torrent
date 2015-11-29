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
