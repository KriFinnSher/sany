#!/usr/bin/env sh
set -eu

# Run the server first: go run ./cmd/server
# Optional: BASE_URL=http://localhost:4040 ./api.sh
BASE_URL="${BASE_URL:-http://localhost:4040}"
tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

sample="$tmp_dir/sample.txt"
downloaded="$tmp_dir/downloaded.txt"
printf 'hello from api test\n' > "$sample"

assert_status() {
	actual="$1"
	expected="$2"
	label="$3"
	if [ "$actual" != "$expected" ]; then
		echo "$label: expected HTTP $expected, got $actual" >&2
		exit 1
	fi
}

response="$tmp_dir/upload.json"
status="$(curl -sS -o "$response" -w '%{http_code}' -F "file=@$sample;type=text/plain" "$BASE_URL/api/v1/files")"
assert_status "$status" 201 "upload"

link="$(sed -n 's/.*"link":"\([^"]*\)".*/\1/p' "$response")"
if [ -z "$link" ]; then
	echo 'upload: response did not contain a public link' >&2
	exit 1
fi

status="$(curl -sS -o "$downloaded" -w '%{http_code}' "$BASE_URL$link")"
assert_status "$status" 200 "download"
cmp -s "$sample" "$downloaded" || { echo 'download: content differs' >&2; exit 1; }

status="$(curl -sS -o /dev/null -w '%{http_code}' -X DELETE "$BASE_URL/api/v1/files")"
assert_status "$status" 405 "method validation"

status="$(curl -sS -o /dev/null -w '%{http_code}' -F 'wrong=value' "$BASE_URL/api/v1/files")"
assert_status "$status" 400 "form validation"

large="$tmp_dir/large.bin"
dd if=/dev/zero of="$large" bs=1048576 count=51 2>/dev/null
status="$(curl -sS -o /dev/null -w '%{http_code}' -F "file=@$large" "$BASE_URL/api/v1/files")"
assert_status "$status" 413 "file-size validation"

status="$(curl -sS -o /dev/null -w '%{http_code}' "$BASE_URL/api/v1/files?id=not-found")"
assert_status "$status" 404 "missing public file"

echo "API checks passed: $link"
