---
id: glossary
title: Glossary
doc_type: glossary
status: accepted
normative: true
version_scope: v0_required
last_reviewed: 2026-05-06
---


# Glossary

## SpecOps

The repeatable workflow and toolchain for producing specification repositories from ideation, decisions, and agent-assisted drafting.

## SpecOps CLI

The command-line executable, tentatively `specops`, that installs scaffolding, manages runs, coordinates agent/relay backends, validates artifacts, applies accepted patches, and runs evals.

## Target spec repository

A repository containing the produced specification for one product/system. It includes canonical docs and a local `.specops/` operating scaffold.

## `.specops/`

The repo-local control plane. It contains process documentation, command templates, schemas, eval rubrics, and run scratch space. It is not the canonical spec itself.

## Canonical docs

Accepted source-of-truth files under `docs/`, including `CANON.md`, ADRs, subsystem specs, cross-cutting specs, interfaces, versions, and implementation phases.

## Run

A stateful SpecOps operation that ingests input and produces draft artifacts, such as a refinery note, ambiguity register, decision queue, spec delta, patch plan, and eval report.

## Spec delta

A typed intermediate artifact describing proposed additions, decisions, ambiguities, affected docs, and patch intent.

## Decision gate

The explicit point where a human accepts, rejects, defers, or amends proposed decisions before the CLI patches canonical docs.

## Refinery note

A human-readable provenance document that summarizes raw ideation, ambiguities, options, recommendations, affected docs, and patch plans.

## Delegated skill

A repo external to `mise-en-place` that exposes a stable installer contract so `mise-en-place` can orchestrate its installation.

## Dual target

Claude Code and Codex CLI payloads generated from the same skill/tool intent.
