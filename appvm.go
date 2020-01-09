/**
 * @author Mikhail Klementev jollheef<AT>riseup.net
 * @license GNU GPLv3
 * @date July 2018
 * @brief appvm launcher
 */

package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/go-cmd/cmd"
	"github.com/jollheef/go-system"
	"github.com/olekukonko/tablewriter"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func list(l *libvirt.Libvirt) {
	domains, err := l.Domains()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Started VM:")
	for _, d := range domains {
		if strings.HasPrefix(d.Name, "appvm") {
			fmt.Println("\t", d.Name[6:])
		}
	}

	fmt.Println("\nAvailable VM:")
	files, err := ioutil.ReadDir(configDir + "/nix")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		switch f.Name() {
		case "base.nix":
			continue
		case "local.nix":
			continue

		}

		fmt.Println("\t", f.Name()[0:len(f.Name())-4])
	}
}

func copyFile(from, to string) (err error) {
	source, err := os.Open(from)
	if err != nil {
		return
	}
	defer source.Close()

	destination, err := os.Create(to)
	if err != nil {
		return
	}

	_, err = io.Copy(destination, source)
	if err != nil {
		destination.Close()
		return
	}

	return destination.Close()
}

func prepareTemplates(appvmPath string) (err error) {
	if _, err = os.Stat(appvmPath + "/nix/local.nix"); os.IsNotExist(err) {
		err = ioutil.WriteFile(configDir+"/nix/local.nix", local_nix_template, 0644)
		if err != nil {
			return
		}
	}

	return
}

func streamStdOutErr(command *cmd.Cmd) {
	for {
		select {
		case line := <-command.Stdout:
			fmt.Println(line)
		case line := <-command.Stderr:
			fmt.Fprintln(os.Stderr, line)
		}
	}
}

func generateVM(path, name string, verbose bool) (realpath, reginfo, qcow2 string, err error) {
	command := cmd.NewCmdOptions(cmd.Options{Buffered: false, Streaming: true},
		"nix-build", "<nixpkgs/nixos>", "-A", "config.system.build.vm",
		"-I", "nixos-config="+path+"/nix/"+name+".nix", "-I", path)

	if verbose {
		go streamStdOutErr(command)
	}

	status := <-command.Start()
	if status.Error != nil || status.Exit != 0 {
		log.Println(status.Error, status.Stdout, status.Stderr)
		if status.Error != nil {
			err = status.Error
		} else {
			s := fmt.Sprintf("ret code: %d, out: %v, err: %v",
				status.Exit, status.Stdout, status.Stderr)
			err = errors.New(s)
		}
		return
	}

	realpath, err = filepath.EvalSymlinks("result/system")
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadFile("result/bin/run-nixos-vm")
	if err != nil {
		return
	}

	match := regexp.MustCompile("regInfo=.*/registration").FindSubmatch(bytes)
	if len(match) != 1 {
		err = errors.New("should be one reginfo")
		return
	}

	reginfo = string(match[0])

	syscall.Unlink("result")

	qcow2 = os.Getenv("HOME") + "/appvm/.fake.qcow2"
	if _, err = os.Stat(qcow2); os.IsNotExist(err) {
		system.System("qemu-img", "create", "-f", "qcow2", qcow2, "512M")
		err = os.Chmod(qcow2, 0400) // qemu run with -snapshot, we only need it for create /dev/vda
		if err != nil {
			return
		}
	}

	return
}

func isRunning(l *libvirt.Libvirt, name string) bool {
	_, err := l.DomainLookupByName("appvm_" + name) // yep, there is no libvirt error handling
	// VM is destroyed when stop so NO VM means STOPPED
	return err == nil
}

func generateAppVM(l *libvirt.Libvirt,
	nixName, vmName, appvmPath, sharedDir string,
	verbose, online bool) (err error) {

	realpath, reginfo, qcow2, err := generateVM(appvmPath, nixName, verbose)
	if err != nil {
		return
	}

	xml := generateXML(vmName, online, realpath, reginfo, qcow2, sharedDir)
	_, err = l.DomainCreateXML(xml, libvirt.DomainStartValidate)
	return
}

func stupidProgressBar() {
	const length = 70
	for {
		time.Sleep(time.Second / 4)
		fmt.Printf("\r%s]\r[", strings.Repeat(" ", length))
		for i := 0; i <= length-2; i++ {
			time.Sleep(time.Second / 20)
			fmt.Printf("+")
		}
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func isAppvmConfigurationExists(appvmPath, name string) bool {
	return fileExists(appvmPath + "/nix/" + name + ".nix")
}

func start(l *libvirt.Libvirt, name string, verbose, online, stateless bool,
	args, open string) {

	appvmPath := configDir

	statelessName := fmt.Sprintf("tmp_%d_%s", rand.Int(), name)

	sharedDir := os.Getenv("HOME") + "/appvm/"
	if stateless {
		sharedDir += statelessName
	} else {
		sharedDir += name
	}

	os.MkdirAll(sharedDir, 0700)

	vmName := "appvm_"
	if stateless {
		vmName += statelessName
	} else {
		vmName += name
	}

	if open != "" {
		filename := sharedDir + "/" + filepath.Base(open)
		err := copyFile(open, filename)
		if err != nil {
			log.Println("Can't copy file")
			return
		}

		args += "/home/user/" + filepath.Base(open)
	}

	if args != "" {
		err := ioutil.WriteFile(sharedDir+"/"+".args", []byte(args), 0700)
		if err != nil {
			log.Println("Can't write args")
			return
		}
	}

	if !isAppvmConfigurationExists(appvmPath, name) {
		log.Println("No configuration exists for app, " +
			"trying to generate")
		err := generate(name, "", "", false)
		if err != nil {
			log.Println("Can't auto generate")
			return
		}
	}

	// Copy templates
	err := prepareTemplates(appvmPath)
	if err != nil {
		log.Fatal(err)
	}

	if !isRunning(l, vmName) {
		if !verbose {
			go stupidProgressBar()
		}

		err = generateAppVM(l, name, vmName, appvmPath, sharedDir,
			verbose, online)
		if err != nil {
			log.Fatal(err)
		}
	}

	cmd := exec.Command("virt-viewer", "-c", "qemu:///system", vmName)
	cmd.Start()
}

func stop(l *libvirt.Libvirt, name string) {
	dom, err := l.DomainLookupByName("appvm_" + name)
	if err != nil {
		if libvirt.IsNotFound(err) {
			log.Println("Appvm not found or already stopped")
			return
		} else {
			log.Fatal(err)
		}
	}
	err = l.DomainShutdown(dom)
	if err != nil {
		log.Fatal(err)
	}
}

func drop(name string) {
	appDataPath := fmt.Sprintf(os.Getenv("HOME") + "/appvm/" + name)
	os.RemoveAll(appDataPath)
}

func autoBalloon(l *libvirt.Libvirt, memoryMin, adjustPercent uint64) {
	domains, err := l.Domains()
	if err != nil {
		log.Fatal(err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Application VM", "Used memory", "Current memory", "Max memory", "New memory"})
	for _, d := range domains {
		if strings.HasPrefix(d.Name, "appvm_") {
			name := d.Name[6:]

			memoryUsedRaw, err := ioutil.ReadFile(os.Getenv("HOME") + "/appvm/" + name + "/.memory_used")
			if err != nil {
				log.Println(err)
				continue
			}
			memoryUsedMiB, err := strconv.Atoi(string(memoryUsedRaw[0 : len(memoryUsedRaw)-1]))
			if err != nil {
				log.Println(err)
				continue
			}
			memoryUsed := memoryUsedMiB * 1024

			_, memoryMax, memoryCurrent, _, _, err := l.DomainGetInfo(d)
			if err != nil {
				log.Println(err)
				continue
			}

			memoryNew := uint64(float64(memoryUsed) * (1 + float64(adjustPercent)/100))

			if memoryNew > memoryMax {
				memoryNew = memoryMax - 1
			}

			if memoryNew < memoryMin {
				memoryNew = memoryMin
			}

			err = l.DomainSetMemory(d, memoryNew)
			if err != nil {
				log.Println(err)
				continue
			}

			table.Append([]string{name,
				fmt.Sprintf("%d", memoryUsed),
				fmt.Sprintf("%d", memoryCurrent),
				fmt.Sprintf("%d", memoryMax),
				fmt.Sprintf("%d", memoryNew)})
		}
	}
	table.Render()
}

func search(name string) {
	command := exec.Command("nix", "search", name)
	bytes, err := command.Output()
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(bytes), "\n") {
		fmt.Println(line)
	}
	return
}

func sync() {
	err := exec.Command("nix-channel", "--update").Run()
	if err != nil {
		log.Fatalln(err)
	}

	err = exec.Command("nix", "search", "-u").Run()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Done")
}

func cleanupStatelessVMs(l *libvirt.Libvirt) {
	domains, err := l.Domains()
	if err != nil {
		log.Fatal(err)
	}

	dirs, err := ioutil.ReadDir(appvmHomesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range dirs {
		if !strings.HasPrefix(f.Name(), "tmp_") {
			continue
		}

		alive := false
		for _, d := range domains {
			if d.Name == "appvm_"+f.Name() {
				alive = true
			}
		}
		if !alive {
			os.RemoveAll(appvmHomesDir + f.Name())
		}
	}
}

var configDir = os.Getenv("HOME") + "/.config/appvm/"
var appvmHomesDir = os.Getenv("HOME") + "/appvm/"

func main() {
	rand.Seed(time.Now().UnixNano())

	os.Mkdir(os.Getenv("HOME")+"/appvm", 0700)

	os.MkdirAll(configDir+"/nix", 0700)

	err := writeBuiltinApps(configDir + "/nix")
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(configDir+"/nix/base.nix", baseNix(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	kingpin.Command("list", "List applications")
	autoballonCommand := kingpin.Command("autoballoon", "Automatically adjust/reduce app vm memory")
	minMemory := autoballonCommand.Flag("min-memory", "Set minimal memory (megabytes)").Default("1024").Uint64()
	adjustPercent := autoballonCommand.Flag("adj-memory", "Adjust memory amount (percents)").Default("20").Uint64()

	startCommand := kingpin.Command("start", "Start application")
	startName := startCommand.Arg("name", "Application name").Required().String()
	startQuiet := startCommand.Flag("quiet", "Less verbosity").Bool()
	startArgs := startCommand.Flag("args", "Command line arguments").String()
	startOpen := startCommand.Flag("open", "Pass file to application").String()
	startOffline := startCommand.Flag("offline", "Disconnect").Bool()
	startStateless := startCommand.Flag("stateless", "Do not use default state directory").Bool()

	stopName := kingpin.Command("stop", "Stop application").Arg("name", "Application name").Required().String()
	dropName := kingpin.Command("drop", "Remove application data").Arg("name", "Application name").Required().String()

	generateCommand := kingpin.Command("generate", "Generate appvm definition")
	generateName := generateCommand.Arg("name", "Nix package name").Required().String()
	generateBin := generateCommand.Arg("bin", "Binary").Default("").String()
	generateVMName := generateCommand.Flag("vm", "Use VM Name").Default("").String()
	generateBuildVM := generateCommand.Flag("build", "Build VM").Bool()

	searchCommand := kingpin.Command("search", "Search for application")
	searchName := searchCommand.Arg("name", "Application name").Required().String()

	kingpin.Command("sync", "Synchronize remote repos for applications")

	var l *libvirt.Libvirt
	if kingpin.Parse() != "generate" {
		c, err := net.DialTimeout(
			"unix",
			"/var/run/libvirt/libvirt-sock",
			time.Second,
		)
		if err != nil {
			log.Fatal(err)
		}

		l = libvirt.New(c)
		if err := l.Connect(); err != nil {
			log.Fatal(err)
		}
		defer l.Disconnect()

		cleanupStatelessVMs(l)
	}

	switch kingpin.Parse() {
	case "list":
		list(l)
	case "search":
		search(*searchName)
	case "generate":
		generate(*generateName, *generateBin, *generateVMName,
			*generateBuildVM)
	case "start":
		start(l, *startName,
			!*startQuiet, !*startOffline, *startStateless,
			*startArgs, *startOpen)
	case "stop":
		stop(l, *stopName)
	case "drop":
		drop(*dropName)
	case "autoballoon":
		autoBalloon(l, *minMemory*1024, *adjustPercent)
	case "sync":
		sync()
	}
}
