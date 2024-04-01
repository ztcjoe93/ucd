-- Pull in the wezterm API
local wezterm = require 'wezterm'

-- This table will hold the configuration.
local config = {}

-- In newer versions of wezterm, use the config_builder which will
-- help provide clearer error messages
if wezterm.config_builder then
  config = wezterm.config_builder()
end

-- This is where you actually apply your config choices

config.front_end = "OpenGL"
config.enable_wayland = true
config.hide_tab_bar_if_only_one_tab = true

config.font = wezterm.font_with_fallback {
  'Monaco',
  'Symbols Nerd Font'
}
config.font_size = 10.0
config.color_scheme = 'Omni (Gogh)'
config.window_background_opacity = 0.5

-- and finally, return the configuration to wezterm
return config
