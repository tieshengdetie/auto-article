#!/usr/bin/env sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(CDPATH= cd -- "$script_dir/.." && pwd)
save_script="$repo_root/skills/auto-media-writer/scripts/save_skill_article.py"

if [ ! -f "$save_script" ]; then
    echo "Cannot find auto-media-writer save script: $save_script" >&2
    exit 2
fi

python_bin="${PYTHON:-python3}"
exec "$python_bin" "$save_script" "$@"
