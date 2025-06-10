package start

import (
	"os"
	"os/signal"
	"path"
	"runtime"
	"syscall"

	"github.com/nrfta/go-tiger/helpers"

	"github.com/nrfta/go-log"

	"github.com/air-verse/air/runner"
	"github.com/spf13/cobra"
)

var (
	cfgPath   string
	debugMode bool
	noUlimit  bool
)

func setUlimit() error {
	if runtime.GOOS == "windows" {
		return nil
	}

	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}
	rLimit.Max = 2048

	rLimit.Cur = rLimit.Max
	return syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit)
}

var StartCmd = &cobra.Command{
	Use:     "start",
	Aliases: []string{"s", "serve"},
	Short:   "Start serving a Go app (Air)",
	Run: func(_ *cobra.Command, _ []string) {
		if debugMode {
			log.Info("[debug] mode")
		}

		if !noUlimit {
			if debugMode {
				log.Info("[debug] set ulimit")
			}
			setUlimit()
		}

		if cfgPath == "" {
			defaultAirPath := path.Join(helpers.FindRootPath(), ".air.toml")

			if _, err := os.Stat(defaultAirPath); err == nil {
				cfgPath = defaultAirPath
			} else if os.IsNotExist(err) {
				appName := helpers.LoadConfig().Meta.ServiceName

				file, err := os.CreateTemp(os.TempDir(), appName+".*.toml")
				if err != nil {
					log.Fatal(err)
				}
				defer os.Remove(file.Name())

				file.WriteString(defaultAirConfig(appName))
				if err := file.Close(); err != nil {
					log.Fatal(err)
				}

				cfgPath = file.Name()
			}
		}

		if debugMode {
			log.Info("[debug] Using Config Path:", cfgPath)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		var err error
		r, err := runner.NewEngine(cfgPath, map[string]runner.TomlInfo{}, debugMode)
		if err != nil {
			log.Fatal(err)
			return
		}
		go func() {
			<-sigs
			r.Stop()
		}()

		defer func() {
			if e := recover(); e != nil {
				log.Fatalf("PANIC: %+v", e)
			}
		}()

		r.Run()
	},
}

func init() {
	StartCmd.Flags().BoolVarP(&debugMode, "debug", "d", false, "debug mode")
	StartCmd.Flags().StringVarP(&cfgPath, "config", "c", "", "config path")
	StartCmd.Flags().BoolVar(&noUlimit, "no-ulimit", false, "do not set ulimit")
}
