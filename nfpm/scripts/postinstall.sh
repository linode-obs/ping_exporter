#!/bin/sh

# Source: https://nfpm.goreleaser.com/tips/#example-multi-platform-post-install-script

# Step 1, decide if we should use systemd or init/upstart
use_systemctl="True"
systemd_version=0
if ! command -V systemctl >/dev/null 2>&1; then
  use_systemctl="False"
else
    systemd_version=$(systemctl --version | head -1 | awk '{ print $2}')
fi

cleanup() {
    # This is where you remove files that were not needed on this platform / system
    if [ "${use_systemctl}" = "False" ]; then
        rm -f /etc/systemd/system/prometheus-ping-exporter.service
    else
        rm -f /etc/chkconfig/prometheus-ping-exporter
        rm -f /etc/init.d/prometheus-ping-exporter
    fi
}

cleanInstall() {
    # Step 3 (clean install), enable the service in the proper way for this platform
    if [ "${use_systemctl}" = "False" ]; then
        if command -V chkconfig >/dev/null 2>&1; then
          chkconfig --add prometheus-ping-exporter
        fi

        service prometheus-ping-exporter restart ||:
    else
        # rhel/centos7 cannot use ExecStartPre=+ to specify the pre start should be run as root
        # even if you want your service to run as non root.
        if [ "${systemd_version}" -lt 231 ]; then
            printf "\033[31m systemd version %s is less then 231, fixing the service file \033[0m\n" "${systemd_version}"
            sed -i "s/=+/=/g" /etc/systemd/system/prometheus-ping-exporter.service
        fi
        systemctl daemon-reload ||:
        systemctl unmask prometheus-ping-exporter ||:
        systemctl preset prometheus-ping-exporter ||:
        systemctl enable prometheus-ping-exporter ||:
        systemctl restart prometheus-ping-exporter ||:
    fi
}

# Step 2, check if this is a clean install
action="$1"
if  [ "$1" = "configure" ] && [ -z "$2" ]; then
  # Alpine linux does not pass args, and deb passes $1=configure
  action="install"
fi

case "$action" in
  "1" | "install")
    cleanInstall
    ;;
  *)
    cleanInstall
    ;;
esac

# Step 4, clean up unused files, yes you get a warning when you remove the package, but that is ok.
cleanup
