package plugin

import (
	"errors"

	"github.com/pocketbase/pocketbase/core"
)

type Plugin interface {
	/**
	 * Preload is called before the app is setup.
	 * This is a good place to load configurations.
	 */
	OnPreload(app core.App) error

	/**
	 * Load is called after the app is setup.
	 * This is a good place to register commands
	 * and hooks.
	 */
	OnLoad(app core.App) error

	/**
	 * Serve is called after the app starts and
	 * has finished setting up the HTTP server.
	 * This is a good place to register routes.
	 */
	OnServe(event *core.ServeEvent) error

	/**
	 * Get plugin info
	 */
	Info() PluginInfo
}

type PluginInfo struct {
	Name        string
	Version     string
	Description string
}

var plugins = []Plugin{}

func Register(plugin Plugin) {
	plugins = append(plugins, plugin)
}

func TriggerPreload(app core.App) (err error) {
	for _, plugin := range plugins {
		if pluginErr := plugin.OnPreload(app); pluginErr != nil {
			err = errors.Join(err, pluginErr)
		}
	}
	return
}

func TriggerLoad(app core.App) (err error) {
	for _, plugin := range plugins {
		if pluginErr := plugin.OnLoad(app); pluginErr != nil {
			err = errors.Join(err, pluginErr)
		}
	}
	return
}

func TriggerServe(event *core.ServeEvent) (err error) {
	for _, plugin := range plugins {
		if pluginErr := plugin.OnServe(event); pluginErr != nil {
			err = errors.Join(err, pluginErr)
		}
	}
	return
}

func GetPlugins() []Plugin {
	return plugins
}
