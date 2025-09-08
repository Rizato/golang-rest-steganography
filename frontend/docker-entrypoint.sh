#!/bin/sh

# Replace the placeholder in config.js with the actual backend URL
if [ ! -z "$BACKEND_URL" ]; then
  sed -i "s|BACKEND_URL_PLACEHOLDER|$BACKEND_URL|g" /usr/share/nginx/html/config.js
fi

# Start nginx
exec "$@"