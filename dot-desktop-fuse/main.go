package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type ddf struct {
	Path string
	fs.Inode
}

func (r *ddf) OnAdd(ctx context.Context) {
	files, err := ioutil.ReadDir(r.Path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		b, err := ioutil.ReadFile(filepath.Join(r.Path, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		if bytes.Contains(b, []byte("Actions=")) {
			b = bytes.ReplaceAll(b, []byte("Actions="), []byte("Actions=appvm;"))
		} else {
			b = bytes.ReplaceAll(b, []byte("Exec="), []byte("Actions=appvm;\nExec="))
		}

		raw := string(regexp.MustCompile("Exec=[a-zA-Z0-9]*").Find(b))

		var app string
		if len(raw) > 5 {
			app = string(raw)[5:]
		} else {
			log.Println("Can't find Exec entry for", file.Name())
			continue
		}

		b = append(b, []byte("\n[Desktop Action appvm]\n")...)
		b = append(b, []byte("Name=Open in appvm\n")...)
		b = append(b, []byte("Exec=appvm start "+app+"\n")...)

		ch := r.NewPersistentInode(
			ctx, &fs.MemRegularFile{
				Data: b,
				Attr: fuse.Attr{Mode: 0444},
			}, fs.StableAttr{})
		r.AddChild(file.Name(), ch, false)
	}
}

var _ = (fs.NodeOnAdder)((*ddf)(nil))

func setupSigintHandler(server *fuse.Server) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Unmount()
		os.Exit(1)
	}()
}

func main() {
	const from = "/var/run/current-system/sw/share/applications/"

	to := filepath.Join(os.Getenv("HOME"), "/.local/share/applications")
	os.MkdirAll(to, 0755)

	server, err := fs.Mount(to, &ddf{Path: from}, nil)
	if err != nil {
		log.Fatal(err)
	}
	setupSigintHandler(server)
	server.Wait()
}
