# treeforge

[![Test](https://github.com/Qooh0/treeforge/actions/workflows/test.yml/badge.svg)](https://github.com/Qooh0/treeforge/actions/workflows/test.yml)
[![Release](https://github.com/Qooh0/treeforge/actions/workflows/release.yml/badge.svg)](https://github.com/Qooh0/treeforge/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Qooh0/treeforge)](https://goreportcard.com/report/github.com/Qooh0/treeforge)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![Latest Release](https://img.shields.io/github/v/release/Qooh0/treeforge)](https://github.com/Qooh0/treeforge/releases/latest)

**treeforge** は、フォルダツリーを実際のディレクトリとファイルに変換する、  
シンプルで安全なGoのCLIツールです。

ChatGPT（または任意のマークダウン文書）からツリー図をコピーして、コマンドを実行するだけ。  
もう `mkdir` や `touch` を手動で繰り返す必要はありません — 設計から実装へ、一瞬で移行できます。

---

## 🚀 特徴

- 📋 **コピペ対応** — ChatGPTのASCIIツリー出力をそのまま使える
- 🧩 **スマートなパース** — コメント（`# ...`）や装飾（`├─`、`│` など）を自動で無視
- 🪶 **デフォルトで安全** — `--apply` を指定しない限り*ドライラン*モードで動作
- 🔧 **設定可能** — 親ディレクトリやルートフォルダ名を選択可能
- 💻 **クロスプラットフォーム** — macOS、Linux、Windows で動作

---

## 📦 インストール
```bash
go install github.com/qooh0/treeforge@latest
```

`$HOME/go/bin`（または `$GOBIN`）がPATHに含まれていることを確認してください。

---

## 🧭 使い方

### 1️⃣ フォルダツリーをコピー

ChatGPTからツリー構造を取得するか、手動で作成します：
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

### 2️⃣ treeforge を実行

**ファイルから読み込む：**
```bash
# ツリーをファイルに保存
cat > tree.txt
# (ツリーを貼り付けて Ctrl+D)

# ドライラン（作成内容をプレビュー）
treeforge -i tree.txt

# 実際に構造を作成
treeforge -i tree.txt --apply
```

**クリップボードから（macOS/Linux）：**
```bash
# ツリーをクリップボードにコピーしてから：
pbpaste | treeforge --apply           # macOS
xclip -o | treeforge --apply          # Linux (X11)
wl-paste | treeforge --apply          # Linux (Wayland)
```

**標準入力から：**
```bash
# コマンドからパイプ
cat tree.txt | treeforge --apply

# または直接貼り付け
treeforge --apply
# (ツリーを貼り付けて Ctrl+D)
```

**カスタムオプション付き：**
```bash
# 親ディレクトリを指定
treeforge -i tree.txt --apply --parent ~/projects

# ルートフォルダ名を上書き
treeforge -i tree.txt --apply --root-name myapp

# 詳細出力
treeforge -i tree.txt --apply -v

# 既存ファイルを強制上書き
treeforge -i tree.txt --apply --force
```

---

## ⚙️ オプション

| オプション             | 説明                                         |
| ------------------ | ------------------------------------------ |
| `-i FILE`          | ファイルからツリーを読み込む（デフォルト: 標準入力）              |
| `--parent DIR`     | 親ディレクトリを指定（デフォルト: カレントディレクトリ）            |
| `--root-name NAME` | ルートフォルダ名を上書き（デフォルト: 1行目から取得）             |
| `--apply`          | 実際にファイル/ディレクトリを作成（デフォルト: ドライラン）          |
| `--force`          | 既存ファイルを上書き（ディレクトリは保持）                     |
| `-v`               | 詳細ログを出力                                    |

---

## 🧩 安全設計

- **デフォルトでドライラン** — `--apply` を指定するまで何も作成されない
- **コメント対応** — 行から `# コメント` を自動的に削除
- **装飾に寛容** — `├─`、`│`、`└─`、`|--`、タブ、スペースに対応
- **既存ファイルを保護** — 既存のファイルはスキップ（`--force` 指定時を除く）
- **冪等性** — 何度実行しても安全

---

## 💡 動機

ChatGPTやAIツールが「フォルダ構造」を出力したとき、  
それを `mkdir` や `touch` で手動再現するのは面倒です。

**treeforge** はその摩擦を取り除きます — プロジェクトレイアウト、サンプルコード、  
教材の例を瞬時に実体化できる、ミニマルなGoツールです。

**こんな用途に最適：**
- 🏗️ CLIやサービスのスケルトンを素早くプロトタイピング
- 📚 ドキュメントを実際のプロジェクトに変換
- 🎓 プロジェクト構造をインタラクティブに教える
- 🤖 AIの提案からプロジェクトの雛形を自動生成

---

## 🛠️ 開発
```bash
# クローンしてビルド
git clone https://github.com/qooh0/treeforge.git
cd treeforge
go build -o treeforge

# ローカルで実行
./treeforge -i sample.txt --apply
```

---

## 📄 ライセンス

Apache License 2.0  
© 2025 Qooh0 / Qadiff LLC


