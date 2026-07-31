#!/usr/bin/env python3
from pathlib import Path

root = Path(__file__).resolve().parents[2]
path = root / "internal/win/gui_main_windows.go"
text = path.read_text(encoding="utf-8")
old = "\tidInfoMeta     = 118\n"
new = "\tidInfoMeta     = 118\n\tidHeader       = 119 // shared title style for options and macro windows\n"
if old not in text:
    raise SystemExit("expected idInfoMeta marker not found")
path.write_text(text.replace(old, new, 1), encoding="utf-8", newline="\n")
print("restored shared idHeader")
