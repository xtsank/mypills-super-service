#!/bin/sh
set -eu

password_file=/run/secrets/admin_password
auth_file=/run/admin.htpasswd

umask 077
if [ -s "$password_file" ]; then
    case "$ADMIN_USER" in
        ''|*[!a-zA-Z0-9_.-]*)
            printf 'ADMIN_USER must contain only letters, digits, dots, underscores or hyphens.\n' >&2
            exit 1
            ;;
    esac
    hash=$(openssl passwd -6 -stdin < "$password_file")
    printf '%s:%s\n' "$ADMIN_USER" "$hash" > "$auth_file"
else
    : > "$auth_file"
fi

chown nginx:nginx "$auth_file"
chmod 0400 "$auth_file"
