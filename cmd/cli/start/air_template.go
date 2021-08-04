package start

import "fmt"

func defaultAirConfig(appName string) string {
	template := `
[build]
  bin = "./bin/%s"
  cmd = "go build -o ./bin/%s ."
  exclude_dir = ["bin", "tests", "templates", "scripts", "db/migrations", "pkg/schemas"]
  kill_delay = 500
  send_interrupt = true
`

	return fmt.Sprintf(template, appName, appName)
}
