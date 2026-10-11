#!/usr/bin/env bash
# Symlink sibling repos into providers/ for local dev (providers/ is gitignored).
# Also soft-links content-pipelines-mcp Cursor skills into .cursor/skills/.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
providers="${root}/providers"
mkdir -p "${providers}"

link() {
	local name="$1"
	local target="$2"
	local link_path="${providers}/${name}"
	if [[ -L "${link_path}" ]] && [[ "$(readlink "${link_path}")" == "${target}" ]]; then
		echo "providers/${name} -> ${target} (ok)"
		return 0
	fi
	if [[ -L "${link_path}" ]]; then
		ln -sfn "${target}" "${link_path}"
		echo "providers/${name} -> ${target} (retargeted)"
		return 0
	fi
	if [[ -e "${link_path}" ]]; then
		echo "providers/${name} exists and is not the expected symlink; remove it and retry" >&2
		exit 1
	fi
	ln -sfn "${target}" "${link_path}"
	echo "providers/${name} -> ${target}"
}

link_skill() {
	local name="$1"
	local target="$2"
	local link_path="${root}/.cursor/skills/${name}"
	mkdir -p "${root}/.cursor/skills"
	if [[ -L "${link_path}" ]] && [[ "$(readlink "${link_path}")" == "${target}" ]]; then
		echo ".cursor/skills/${name} -> ${target} (ok)"
		return 0
	fi
	if [[ -L "${link_path}" ]]; then
		ln -sfn "${target}" "${link_path}"
		echo ".cursor/skills/${name} -> ${target} (retargeted)"
		return 0
	fi
	if [[ -e "${link_path}" ]]; then
		echo ".cursor/skills/${name} exists and is not a symlink; remove the overlay to adopt the MCP skill" >&2
		exit 1
	fi
	ln -sfn "${target}" "${link_path}"
	echo ".cursor/skills/${name} -> ${target}"
}

# From dss/blog/providers/: ../../../ is Xynova/, then ai/content-pipelines.
link content-pipelines ../../../ai/content-pipelines
link content-pipelines-mcp ../../../ai/content-pipelines-mcp

# From dss/blog/.cursor/skills/: ../../../../ is Xynova/, then ai/content-pipelines-mcp skills.
link_skill essay-command-center ../../../../ai/content-pipelines-mcp/.cursor/skills/essay-command-center
