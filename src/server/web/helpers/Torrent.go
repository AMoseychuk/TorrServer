package helpers

import (
	"fmt"
	"sort"
	"time"

	"server/settings"
	"server/torr"
	"server/utils"

	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func Add(bts *torr.BTServer, magnet metainfo.Magnet, save bool) error {
	fmt.Println("Adding torrent", magnet.String())
	_, err := bts.AddTorrent(magnet, func(torr *torr.Torrent) {
		if torr, _ := settings.LoadTorrentDB(magnet.InfoHash.HexString()); torr != nil {
			return
		}
		torDb := new(settings.Torrent)
		torDb.Name = torr.Name()
		torDb.Hash = torr.Hash().HexString()
		torDb.Size = torr.Length()
		torDb.Magnet = magnet.String()
		torDb.Timestamp = time.Now().Unix()
		files := torr.Stats().FileStats
		sort.Slice(files, func(i, j int) bool {
			return files[i].Path < files[j].Path
		})
		for _, f := range files {
			ff := settings.File{
				Id:   f.Id,
				Name: f.Path,
				Size: f.Length,
			}
			torDb.Files = append(torDb.Files, ff)
		}

		if save {
			err := settings.SaveTorrentDB(torDb)
			if err != nil {
				fmt.Println("Error add torrent to db:", err)
			}
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func FindFileLink(fileLink string, torr *torrent.Torrent) *torrent.File {
	for _, f := range torr.Files() {
		if utils.CleanFName(f.Path()) == fileLink {
			return f
		}
	}
	return nil
}

func FindFile(fileInd int, tor *torr.Torrent) *torrent.File {
	files := tor.Files()
	if len(files) == 0 || fileInd < 0 || fileInd >= len(files) {
		return nil
	}
	return files[fileInd]
}
