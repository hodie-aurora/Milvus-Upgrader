# Milvus-Upgrader

A Python tool to upgrade Milvus clusters (v2.2.3+) on Linux.

## Features

- Supports Operator and Helm deployments (Operator first, Helm TBD)
- Backup, upgrade, and rollback capabilities
- Designed for offline environments

## Usage

```bash
python3 milvus_upgrader.py -i my-release -t v2.2.5 -m operator
```
