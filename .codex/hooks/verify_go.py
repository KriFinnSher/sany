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


# state_path derives the shared per-session state location.
def state_path(event: dict[str, object], root: Path, suffix: str) -> Path:
    key = f"{event['session_id']}:{root}"
    name = hashlib.sha256(key.encode()).hexdigest()
    return Path(tempfile.gettempdir()) / "codex-go-validation" / f"{name}{suffix}"


# run_checks runs both repository validations and returns their failures.
def run_checks(root: Path) -> list[str]:
    failures = []
    for command in (["make", "lint"], ["make", "test"]):
        result = subprocess.run(command, cwd=root, capture_output=True, text=True)
        if result.returncode != 0:
            output = (result.stdout + result.stderr).strip()
            failures.append(f"$ {' '.join(command)}\n{output[-6000:]}")
    return failures


# main validates Go changes and continues the agent when validation fails.
def main() -> None:
    event = json.load(sys.stdin)
    root = repository_root(str(event["cwd"]))
    snapshot = state_path(event, root, f"-{event['turn_id']}.json")
    pending = state_path(event, root, ".pending")

    changed = False
    if snapshot.exists():
        before = json.loads(snapshot.read_text())["files"]
        snapshot.unlink()
        changed = before != go_files(root)

    if not changed and not pending.exists():
        print(json.dumps({"continue": True}))
        return

    pending.touch()
    failures = run_checks(root)
    if failures:
        print(json.dumps({
            "decision": "block",
            "reason": "Go validation failed. Fix the issue and let the hook rerun.\n\n" + "\n\n".join(failures),
        }))
        return

    pending.unlink()
    print(json.dumps({"continue": True, "systemMessage": "Go lint and tests passed."}))


if __name__ == "__main__":
    main()
