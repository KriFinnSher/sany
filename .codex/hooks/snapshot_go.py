#!/usr/bin/env python3
import hashlib
import json
import subprocess
import sys
import tempfile
from pathlib import Path


# repository_root resolves the Git root from the active Codex directory.
def repository_root(cwd: str) -> Path:
    result = subprocess.run(
        ["git", "rev-parse", "--show-toplevel"],
        check=True,
        cwd=cwd,
        capture_output=True,
        text=True,
    )
    return Path(result.stdout.strip())


# go_files maps tracked and untracked Go files to their content hashes.
def go_files(root: Path) -> dict[str, str | None]:
    result = subprocess.run(
        ["git", "ls-files", "-co", "--exclude-standard", "--", "*.go"],
        check=True,
        cwd=root,
        capture_output=True,
        text=True,
    )
    return {name: file_hash(root / name) for name in result.stdout.splitlines()}


# file_hash returns None when a tracked file was deleted.
def file_hash(path: Path) -> str | None:
    try:
        return hashlib.sha256(path.read_bytes()).hexdigest()
    except FileNotFoundError:
        return None


# state_path isolates one Codex turn's snapshot in the system temp directory.
def state_path(event: dict[str, object], root: Path) -> Path:
    key = f"{event['session_id']}:{root}"
    name = hashlib.sha256(key.encode()).hexdigest()
    return Path(tempfile.gettempdir()) / "codex-go-validation" / f"{name}-{event['turn_id']}.json"


# main records the Go workspace state at the beginning of a turn.
def main() -> None:
    event = json.load(sys.stdin)
    root = repository_root(str(event["cwd"]))
    path = state_path(event, root)
    path.parent.mkdir(parents=True, exist_ok=True)

    # Keep snapshots outside the repository so hooks never change its working tree.
    path.write_text(json.dumps({"files": go_files(root)}))


if __name__ == "__main__":
    main()
