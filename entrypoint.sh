#!/bin/sh
set -e

# Create .netrc from env
if [[ ! -z "${NETRC_CONTENTS}" ]]; then
    echo "$NETRC_CONTENTS" | base64 -d -i > /root/.netrc
fi

exec "$@"
