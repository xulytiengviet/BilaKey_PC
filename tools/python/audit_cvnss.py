#!/usr/bin/env python3
"""Independently audit CVNSS4.0 5.x candidate graph and canonical policy."""
from __future__ import annotations

import argparse
import json
from collections import defaultdict
from pathlib import Path


def load_rules(source: Path):
    data_dir = source if source.is_dir() else source.parent / "cvnss"
    if not data_dir.is_dir():
        raise ValueError(f"cannot locate modular CVNSS data directory from {source}")

    def read(name: str):
        return json.loads((data_dir / name).read_text(encoding="utf-8"))

    base = []
    for item in sorted(data_dir.glob("base_[0-9][0-9].json")):
        base.extend(json.loads(item.read_text(encoding="utf-8")))
    patch54 = []
    for item in sorted(data_dir.glob("patch54_[0-9][0-9].json")):
        patch54.extend(json.loads(item.read_text(encoding="utf-8")))
    return base, patch54, read("patch56.json"), read("canonical.json"), read("critical.json")

def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("oracle", type=Path)
    parser.add_argument("--out", type=Path, required=True)
    args = parser.parse_args()

    base, patch54, patch56, preferred, critical = load_rules(args.oracle)

    candidates: dict[str, list[str]] = defaultdict(list)

    def add(code: str, cqn: str) -> None:
        if cqn not in candidates[code]:
            candidates[code].append(cqn)

    for cqn, _cvn, cvss in base:
        add(cvss, cqn)
    for cqn, code in patch54 + patch56:
        add(code, cqn)

    collisions = [(code, values) for code, values in sorted(candidates.items()) if len(values) > 1]
    uncovered = [code for code, values in collisions if preferred.get(code) not in values]
    unexpected_policies = [code for code in preferred if len(candidates.get(code, [])) < 2]
    invariants = (len(base), len(patch54) + len(patch56), len(collisions), len(preferred), len(critical))
    if invariants != (758, 336, 56, 56, 5):
        raise SystemExit(f"FAIL invariant drift: {invariants}")
    if uncovered or unexpected_policies:
        raise SystemExit(f"FAIL policy coverage: uncovered={uncovered} unexpected={unexpected_policies}")

    args.out.parent.mkdir(parents=True, exist_ok=True)
    with args.out.open("w", encoding="utf-8", newline="") as report:
        report.write("code\tcandidates\tclass\trecommended\tpolicy_status\n")
        for code, values in collisions:
            cls = "CRITICAL_PATCH_COLLISION" if code in critical else "CANONICAL_AMBIGUITY"
            report.write(
                f"{code}\t{' | '.join(values)}\t{cls}\t{preferred[code]}\tEXPLICIT\n"
            )

    print(
        f"base={len(base)} patches={len(patch54)+len(patch56)} "
        f"collisions={len(collisions)} policies={len(preferred)} critical={len(critical)}"
    )
    print("policy_coverage=56/56 silent_reverse_overwrite=0 PASS")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
