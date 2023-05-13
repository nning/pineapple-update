# pineapple-update

Automatically update pineapple EA

## Usage

Without any configuration, the latest version of pineapple EA will be downloaded
into the same directory as the updater. A symlink will be created from
`yuzu-ea.AppImage` to the latest downloaded version.

    ./pineapple-update

## Configuration

If you want to override this behaviour, copy `pineapple-update.example.yml` to
`.pineapple-update.yml` or `pineapple-update.yml`. Available options are:

- `targetFolder`: The folder where the AppImage files will be downloaded to.
- `symlink`: Controls whether a symlink will be created. Set to `false` to disable.
- `symlinkName`: The name of the symlink. (Default is `yuzu-ea.AppImage`)

## Automatic Updates

Copy contents of the `systemd` folder to `~/.config/systemd/user/`, then run:

    systemctl --user enable pineapple-update.timer
    systemctl --user start pineapple-update.timer