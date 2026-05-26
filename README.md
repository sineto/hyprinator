<p align="center">
  <a href="https://gitlab.com/sn1o/devbox/hyprinator">
    <img src="https://pub-c1bae9fd4d3844cc89f5ad7dfb067717.r2.dev/assets/hyprinator-logo.png" alt="Hyprinator" />
  </a>
  <br />
  Hyprinator is a simple tool used to automatically setup<br /> windows in <a href="https://hypr.land">Hyprland</a> with YAML files.
  <br />
</p>

---
## Install
1. clone repository:
```bash
git clone https://gitlab.com/sn1o/devbox/hyprinator.git
cd hyprinator
```
2. build the binary:
```bash
go build -o hyprinator
```

## Usage
```bash
hyprinator setup [space-name] [flags]
```

### Commands
| Command     | Description                                                                     |
|-------------|---------------------------------------------------------------------------------|
| **`setup`** | automatically setup environment of windows and workspaces based on `.yaml` file |

### Options
| Option                | Description                                                                       |
|-----------------------|-----------------------------------------------------------------------------------|
| **`--config`**        | receive a `string` with the path to the custom config file (only support: `yaml`) |

### Setup configuration
| Attribute        | Type         | Description                                                                          |
|------------------|--------------|--------------------------------------------------------------------------------------|
| **`app`**        | string       | The app name that must be open/run - (eg.: `alacritty`, `brave`, `vscodium`, etc.)   |
| **`workspace`**  | number       | Hyprland workspace where the app should be opened                                    |
| **`monitor`**    | string       | Monitor name where the app should be moved to - (eg.: `DP-1`, `HDMI-1`, etc.)        |
| **`cwd`**        | string       | Current working directory path to set by terminal - (only for terminal applications) |
| **`exec`**       | array[string] | List of terminal commands to execute - (onl for terminal applications)               |
| **`focus`**      | boolean      | What application should be focused at the setups ending                              |
| **`delay`**      | number       | How much time the CLI need to wait until execute a new iteration                     |

#### Example
```yaml
# /home/john-doe/.hyprinator.yaml
---
main:
  - app: alacritty
    cwd: /home/john-doe/projects
    exec:
      - tmuxp load projects /home/john-doe/.tmuxp-sessions.yaml
    workspace: 1
    monitor: HDMI-1
    delay: 1000
      
  - app: brave
    workspace: 2
    monitor: DP-2
    delay: 1000
    focus: true  
```

To execute the previous `.hyprinator.yaml` configuration:
```bash
hyprinator setup main --config=$HOME/john-doe/.hyprinator
```
## Features
- Execute a list of applications via Hyprland APIs
- Can create many spaces with different setups
- Move windows to any workspaces defined by the configuration
- Windows focusing

## TODO
- [ ] Better user-level configuration
- [ ] Ensure a default `space` is executed if no space name is passed as an argument to `setup` command
- [ ] Ensure execute `cwd` and `exec` configuration attributes only if `app` is a terminal application
- [ ] Ensure run the application on a fallback monitor if the name passed does not exist or is incorrect
- [ ] Implement a better handling for `exec` list of commands to execute on terminal application
- [ ] Implement tests

### Future features
- [ ] Define windows layouts (tabbed, horizontal, vertical)
- [ ] Define window position layout
- [ ] Define window dimensions
- [ ] Support `json` and `toml` configuration file types
- [ ] Emit pop-up notifications with log level type (`info`, `warning`, `error`)
- [ ] Integration with Tmux
- [ ] Integration with Podman (maybe)

## SELF PROMOTION
Do you like this project? Come on:
- Star and follow the repository on [Gitlab](https://gitlab.com/sn1o/devbox/hyprinator).

## LICENSE
[MIT License](LICENSE)