#!/bin/sh
set -eu

BASE_URL=${BASE_URL:-http://127.0.0.1:8080}

curl --fail --silent --show-error "$BASE_URL/health"
curl --fail --silent --show-error "$BASE_URL/feed" >/dev/null
curl --fail --silent --show-error "$BASE_URL/playlists" >/dev/null
curl --fail --silent --show-error "$BASE_URL/playlist/1" >/dev/null

audio_file=$(mktemp "${TMPDIR:-/tmp}/spotify-audio.XXXXXX")
trap 'rm -f "$audio_file"' EXIT HUP INT TERM
curl --fail --silent --show-error "$BASE_URL/songs/night-signals.wav" --output "$audio_file"
audio_header=$(dd if="$audio_file" bs=1 count=4 2>/dev/null)
test "$audio_header" = "RIFF"

printf '\nSpotify backend smoke check passed.\n'
