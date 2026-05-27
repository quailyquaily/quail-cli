#!/usr/bin/env bash
set -euo pipefail

repo="quailyquaily/quail-cli"
binary_name="quail-cli"

fail() {
  printf 'quail-cli install failed: %s\n' "$1" >&2
  exit 1
}

case "$(uname -s)" in
  Linux) os="Linux" ;;
  Darwin) os="Darwin" ;;
  CYGWIN*|MINGW*|MSYS*) os="Windows" ;;
  *) fail "unsupported OS: $(uname -s)" ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="x86_64" ;;
  arm64|aarch64) arch="arm64" ;;
  i386|i686) arch="i386" ;;
  *) fail "unsupported architecture: $(uname -m)" ;;
esac

if [ -n "${INSTALL_DIR:-}" ]; then
  install_dir="$INSTALL_DIR"
elif [ -n "${PREFIX:-}" ]; then
  install_dir="$PREFIX/bin"
else
  install_dir="$HOME/.local/bin"
fi

if [ "$os" = "Windows" ]; then
  archive_ext="zip"
  binary_file="${binary_name}.exe"
else
  archive_ext="tar.gz"
  binary_file="$binary_name"
fi

asset="${binary_name}_${os}_${arch}.${archive_ext}"
url="https://github.com/${repo}/releases/latest/download/${asset}"
tmp_dir="$(mktemp -d)"
archive_path="${tmp_dir}/${asset}"

cleanup() {
  rm -rf "$tmp_dir"
}
trap cleanup EXIT

printf 'Downloading %s\n' "$url"
if command -v curl >/dev/null 2>&1; then
  curl -fsSL --retry 3 -o "$archive_path" "$url"
elif command -v wget >/dev/null 2>&1; then
  wget -q -O "$archive_path" "$url"
else
  fail "curl or wget is required"
fi

if [ "$archive_ext" = "zip" ]; then
  command -v unzip >/dev/null 2>&1 || fail "unzip is required for Windows archives"
  unzip -q "$archive_path" -d "$tmp_dir"
else
  tar -xzf "$archive_path" -C "$tmp_dir"
fi

binary_path="$(find "$tmp_dir" -type f -name "$binary_file" | head -n 1)"
[ -n "$binary_path" ] || fail "binary not found in release archive"

mkdir -p "$install_dir"
cp "$binary_path" "${install_dir}/${binary_file}"
chmod 0755 "${install_dir}/${binary_file}"

printf 'Installed %s to %s\n' "$binary_name" "${install_dir}/${binary_file}"
"${install_dir}/${binary_file}" version || true

case ":$PATH:" in
  *":$install_dir:"*) ;;
  *)
    printf '\n%s is not in PATH.\n' "$install_dir"
    printf 'Add it with:\n'
    printf '  export PATH="%s:$PATH"\n' "$install_dir"
    ;;
esac

printf '\nNext step:\n'
printf '  "%s" login\n' "${install_dir}/${binary_file}"
