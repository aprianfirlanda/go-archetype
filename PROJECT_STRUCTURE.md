# 📁 Project Structure & File Inspection Guide

This document explains how to:

1. View folder structure
2. Dump file contents with filenames as headers

---

## 🌳 View Folder Structure

Use `tree` to visualize directory structure.

### Install `tree` (if not installed)
```bash
brew install tree
```

### Basic usage
```bash
tree
```

### Common options
```bash
# Show only directories
tree -d

# Limit depth
tree -L 2

# Ignore folders
tree -I "node_modules|.git|dist"
```

---

## 📄 View File Contents with Filename Headers

This command prints file contents and shows each filename as a header.

### ✅ Filtered (Recommended)
Only show `.go` and `.env` files, excluding unnecessary folders:

```bash
find . \
  \( -path ./node_modules -o -path ./.git \) -prune -o \
  -type f \( -name "*.go" -o -name "*.env" \) \
  -exec awk 'FNR==1{print "\n===== " FILENAME " ====="} {print}' {} +
```

---

### 📦 All Files (Be Careful ⚠️)

This will print **all files** (can be very large and messy):

```bash
find . \
  \( -path ./node_modules -o -path ./.git \) -prune -o \
  -type f \
  -exec awk 'FNR==1{print "\n===== " FILENAME " ====="} {print}' {} +
```

---

## 🚫 Excluding More Directories

You can extend exclusions:

```bash
find . \
  \( -path ./node_modules -o -path ./.git -o -path ./dist -o -path ./build \) -prune -o \
  -type f \
  -exec awk 'FNR==1{print "\n===== " FILENAME " ====="} {print}' {} +
```

---

## 💡 Tips

- Avoid dumping all files in large projects
- Filter by file type when possible
- Useful for:
    - Debugging
    - Code reviews
    - Sharing project snapshots
    - Feeding into LLMs

---

## 🚀 Optional: Save Output to File

```bash
find . \
  \( -path ./node_modules -o -path ./.git \) -prune -o \
  -type f \( -name "*.go" -o -name "*.env" \) \
  -exec awk 'FNR==1{print "\n===== " FILENAME " ====="} {print}' {} + \
  > output.txt
```

---

## ✅ Summary

| Task               | Command                  |
|--------------------|--------------------------|
| Folder structure   | `tree`                   |
| Filtered file dump | `find + awk`             |
| Full dump          | `find + awk (no filter)` |

---
