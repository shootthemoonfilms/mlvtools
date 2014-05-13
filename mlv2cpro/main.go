package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
)

var (
	mlvDumpPath      = flag.String("mlvdump", "."+string(os.PathSeparator)+"mlv_dump", "Path to mlv_dump binary")
	raw2gpcfPath     = flag.String("raw2gpcf", "."+string(os.PathSeparator)+"raw2gpcf", "Path to raw2gpcf binary")
	outDir           = flag.String("outdir", ".", "Output directory")
	extension        = flag.String("extension", "mov", "File extension")
	threading        = flag.Bool("threading", false, "Use multi-threading")
	keepFiles        = flag.Bool("keepfiles", true, "Keep source files after transcoding")
	scalingParameter string
	wg               sync.WaitGroup
)

func main() {
	flag.Parse()

	args := flag.Args()

	if *threading {
		log.Print("Setting maximum parallelism")
		runtime.GOMAXPROCS(MaxParallelism())
	}

	// Determine if we're using local directory or list of provided ones.
	var dirs []string
	if len(args) < 1 {
		log.Print("Using current working directory to scan")
		dirs = make([]string, 1)
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		dirs[0] = cwd
		log.Print("Found cwd " + cwd)
	} else {
		dirs = args[:]
	}

	// Iterate through directories
	for idx := range dirs {
		scanDir(dirs[idx])
	}

	if *threading {
		log.Print("Waiting for threads to finish")
		wg.Wait()
		log.Print("Run completed")
	}
}

func scanDir(dirName string) {
	files, _ := ioutil.ReadDir(dirName)
	for _, f := range files {
		//log.Print(f.Name())
		fullPath := dirName + string(os.PathSeparator) + f.Name()
		if FileExists(fullPath) && (strings.HasSuffix(f.Name(), ".mlv") || strings.HasSuffix(f.Name(), ".MLV")) {
			log.Print("Processing " + f.Name())
			if *threading {
				wg.Add(1)
				go processFile(dirName, f.Name())
			} else {
				processFile(dirName, f.Name())
			}
		}
	}
}

func processFile(pathName, fileName string) {
	if *threading {
		defer wg.Done()
	}

	// Figure base filename without suffix
	baseFileName := ""
	if strings.HasSuffix(fileName, ".mlv") {
		baseFileName = strings.TrimSuffix(fileName, ".mlv")
	} else if strings.HasSuffix(fileName, ".MLV") {
		baseFileName = strings.TrimSuffix(fileName, ".MLV")
	}

	log.Print("Processing " + fileName + " in '" + pathName + "'")
	origPath := pathName + string(os.PathSeparator) + fileName
	rawPath := *outDir + string(os.PathSeparator) + baseFileName + ".RAW"
	outPath := *outDir + string(os.PathSeparator) + baseFileName + "." + *extension
	_ = os.MkdirAll(*outDir, 0755)

	//
	//	MLV_DUMP
	//	Convert MLV to RAW
	//

	rawArgs := []string{
		"-o", rawPath,
		"-r", // output "legacy" raw format
		origPath,
	}

	rawCommand := exec.Cmd{
		Path: *mlvDumpPath,
		Args: rawArgs,
	}
	rawCommand.Stdout = os.Stdout
	rawCommand.Stderr = os.Stderr
	if err := rawCommand.Start(); err != nil {
		log.Print(err.Error())
		return
	}
	if err := rawCommand.Wait(); err != nil {
		log.Print(err.Error())
		return
	}

	//
	//	RAW2GPCF
	//

	args := []string{
		rawPath,
		outPath,
	}

	args = append(args, outPath)
	command := exec.Cmd{
		Path: *raw2gpcfPath,
		Args: args,
	}
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		log.Print(err.Error())
		return
	}
	if err := command.Wait(); err != nil {
		log.Print(err.Error())
		return
	}

	if !*keepFiles {
		log.Print("Removing original MLV file")
		os.Remove(origPath)
		log.Print("Removing intermediate RAW file")
		os.Remove(rawPath)
	}

	log.Print("Successfully processed")
}

// FileExists reports whether the named file exists.
func FileExists(name string) bool {
	st, err := os.Stat(name)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	if st.IsDir() {
		return false
	}
	return true
}

func MaxParallelism() int {
	maxProcs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	if maxProcs < numCPU {
		return maxProcs
	}
	return numCPU
}
