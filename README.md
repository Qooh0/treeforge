# treeforge

**treeforge** is a small Go CLI tool that turns a folder tree  
into real directories and files — safely and instantly.

Just copy the tree diagram from ChatGPT (or any markdown document) and run the command.  
No more manual `mkdir` or `touch` loops — go straight from design to implementation.

---

## 🚀 Features

- 📋 **Copy & Paste Friendly** — works directly with ChatGPT's ASCII tree output  
- 🧩 **Smart Parsing** — ignores comments (`# ...`) and decorations (`├─`, `│`, etc.)  
- 🪶 **Safe by Default** — runs in *dry-run mode* unless `--apply` is specified  
- 🔧 **Configurable** — choose parent directory and root folder name  
- 💻 **Cross-platform** — works on macOS, Linux, and Windows  

---

## 📦 Installation
```bash
go install github.com/qooh0/treeforge@latest
```

Make sure `$HOME/go/bin` (or `$GOBIN`) is in your PATH.

---

## 🧭 Usage

### 1️⃣ Copy a folder tree

Get a tree structure from ChatGPT or create one manually:
```text
myapp/
├─ src/
│  ├─ handlers/
│  │  ├─ user.go
│  │  └─ auth.go
│  ├─ models/
│  │  └─ user.go
│  ├─ middleware/
│  │  └─ logger.go
│  └─ main.go
├─ tests/
│  ├─ user_test.go
│  └─ auth_test.go
├─ config/
│  └─ config.yaml
├─ .env
├─ .gitignore
├─ go.mod
├─ Dockerfile
└─ README.md
```

### 2️⃣ Run treeforge

**From file:**
```bash
# Save tree to a file
cat > tree.txt
# (paste tree, then Ctrl+D)

# Dry-run (preview what would be created)
treeforge -i tree.txt

# Actually create the structure
treeforge -i tree.txt --apply
```

**From clipboard (macOS/Linux):**
```bash
# Copy tree to clipboard, then:
pbpaste | treeforge --apply           # macOS
xclip -o | treeforge --apply          # Linux (X11)
wl-paste | treeforge --apply          # Linux (Wayland)
```

**From stdin:**
```bash
# Pipe from command
cat tree.txt | treeforge --apply

# Or paste directly
treeforge --apply
# (paste tree, then Ctrl+D)
```

**With custom options:**
```bash
# Specify parent directory
treeforge -i tree.txt --apply --parent ~/projects

# Override root folder name
treeforge -i tree.txt --apply --root-name myapp

# Verbose output
treeforge -i tree.txt --apply -v

# Force overwrite existing files
treeforge -i tree.txt --apply --force
```

---

## ⚙️ Options

| Option             | Description                                          |
| ------------------ | ---------------------------------------------------- |
| `-i FILE`          | Read tree from file (default: stdin)                 |
| `--parent DIR`     | Parent directory (default: current directory)        |
| `--root-name NAME` | Override root folder name (default: from first line) |
| `--apply`          | Actually create files/directories (default: dry-run) |
| `--force`          | Overwrite existing files (directories are preserved) |
| `-v`               | Verbose logging                                      |

---

## 🧩 Safety Design

- **Dry-run by default** — nothing is created until `--apply` is specified
- **Comment-aware** — automatically strips `# comments` from lines
- **Decoration-tolerant** — handles `├─`, `│`, `└─`, `|--`, tabs, and spaces
- **Existing file protection** — skips files that already exist (unless `--force`)
- **Idempotent** — safe to re-run multiple times

---

## 💡 Motivation

When ChatGPT or AI tools output a "folder structure,"  
manually recreating it with `mkdir` and `touch` is tedious.

**treeforge** eliminates that friction — a minimal Go tool that  
lets you instantly materialize your project layout, sample code, or teaching examples.

**Perfect for:**
- 🏗️ Rapidly prototyping CLI or service skeletons
- 📚 Turning documentation into real projects  
- 🎓 Teaching project structures interactively
- 🤖 Automating project scaffolding from AI suggestions

---

## 🛠️ Development
```bash
# Clone and build
git clone https://github.com/qooh0/treeforge.git
cd treeforge
go build -o treeforge

# Run locally
./treeforge -i sample.txt --apply
```

---

## 📄 License

Apache License 2.0  
© 2025 Qooh0 / Qadiff LLC
