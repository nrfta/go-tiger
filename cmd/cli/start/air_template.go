package start

import "fmt"

func defaultAirConfig(appName string) string {
	template := `
[build]
  bin = "./bin/%s"
  cmd = "go build -o ./bin/%s ."
  exclude_dir = ["bin", "tests", "templates", "scripts", "db", "build", "tmp"]
	exclude_regex = ["_test.go"]
  kill_delay = 500
  send_interrupt = true
	stop_on_error = true
[misc]
  clean_on_exit = true
`

	return fmt.Sprintf(template, appName, appName)
}
