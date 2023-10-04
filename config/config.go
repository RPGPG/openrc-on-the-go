package config

import (
	"errors"
	"os"

	"github.com/gookit/ini/v2"
)

func LoadConfigCli(configPath string) (map[string]string, int, string, error) {
	var err = ini.LoadFiles(configPath)
	if err != nil {
		return nil, -999, "", err
	}
	return ini.StringMap("services"), ini.Int("mon_interval_seconds"),
		ini.String("mon_name"), nil
}

func GenerateExampleConfig(configPath string) error {
	content := []byte("; cli-checker params\n; seconds interval for refresh in --mon mode\nmon_interval_seconds = 5\n; set other than 'nouse' to display name on --mon mode\n; use //hostname// to display machine name\nmon_name = //hostname//\n\n; telegram-checker params\ntelegram_message_password = pass123\ntelegram_apikey = YOUR_TELE_KEY\n\n; ntfy.sh notifier params\nntfy_topic = your_topic\nntfy_interval_seconds = 5\n\n[services]\n; some example services to check by program\n; names on the left are not important, just unique\n; on the right are names used by openrc\nservice_net = networking\nconn = sshd\nwebsites = nginx\ncontainers = docker\n")
	err := os.WriteFile(configPath, content, 0600)
	return err
}

// telegram can use the same config file as CLI, as long as it provides necessary keys
func LoadConfigTelegram(configPath string) (map[string]string, string, string, error) {
	var err = ini.LoadFiles(configPath)
	if err != nil {
		return nil, "", "", err
	}
	return ini.StringMap("services"),
		ini.String("telegram_message_password"),
		ini.String("telegram_apikey"),
		nil
}

func LoadConfigNtfy(configPath string) (map[string]string, string, int, error) {
	var err = ini.LoadFiles(configPath)
	if err != nil {
		return nil, "", -999, err
	}
	return ini.StringMap("services"),
		ini.String("ntfy_topic"),
		ini.Int("ntfy_interval_seconds"),
		nil
}

func GetConfigPath(cmdArgs []string) (string, error) {
	if len(cmdArgs) > 0 {
		if len(cmdArgs[0]) > 5 {
			if cmdArgs[0][:6] == "--ini=" {
				return cmdArgs[0][6:], nil
			}
		}
	}
	return "", errors.New("provide config path with --ini=[[path]]")
}
