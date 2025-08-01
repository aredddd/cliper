# Cliper

Cliper 是一款轻量级的 macOS 剪贴板历史记录工具，使用 Go 语言开发。它可以在 macOS 状态栏中运行，记录你的复制历史，通过点击状态栏图标即可查看历史记录并重新复制内容。

## 功能特点

- 在 macOS 状态栏显示，界面简洁美观
- 自动记录剪贴板历史
- 点击历史记录项目即可复制到剪贴板
- 显示复制时间信息
- 轻量级设计，占用系统资源少

## 安装方法

### 直接下载

1. 从 [Releases](https://github.com/lilithgames/cliper/releases) 页面下载最新版本
2. 解压后运行 Cliper 文件

### 从源码构建

```bash
git clone https://github.com/lilithgames/cliper.git
cd cliper
go build -o Cliper ./cmd/cliper
```

## 使用方法

1. 运行 Cliper 应用
2. 在状态栏中找到剪切板图标 (📎)
3. 点击图标查看剪切板历史
4. 点击任意历史记录项即可将其复制到剪贴板

## 开发

### 依赖

- Go 1.16+
- [github.com/caseymrm/menuet](https://github.com/caseymrm/menuet) - macOS 状态栏应用框架
- [github.com/atotto/clipboard](https://github.com/atotto/clipboard) - 剪贴板操作库

### 构建

```bash
go build -o Cliper ./cmd/cliper
```

## License

MIT
