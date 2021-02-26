# init

The `init` folder contains system init (systemd, upstart, sysv) and process manager/supervisor (runit, supervisord) configs for multiple platforms.

## Conventions

On unix systems agent configurations shall be kept under `/etc/seashell` and data under `/var/lib/seashell/`. These directories are assumed to exist, and the Seashell binary is assumed to be located at `/usr/bin/seashell`.

## Agent configuration

The following configuration files are provided as examples:

    * `demo/server.yml`
    * `demo/client.yml`

Place one of these under `/etc/seashell.d` depending on the host's role. You should use `server.yml` to configure a host as a server or `client.yml` to configure a host as a client.

## systemd

On systems using `systemd`, the basic unit file under `systemd/seashell.service` starts and stops the seashell agent. Place it under `/etc/systemd/system/seashell.service`.

You can control Seashell with `systemctl start|stop|restart seashell`.