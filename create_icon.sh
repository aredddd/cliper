#!/bin/bash

# 创建临时目录
mkdir -p iconbuild/AppIcon.iconset

# 生成简单的文本图标（带有"CL"字样）
for size in 16 32 64 128 256 512 1024; do
  # 生成普通分辨率图标
  echo "Creating icon size ${size}x${size}"
  sips -s format png -z $size $size Cliper.app/Contents/MacOS/Cliper --out iconbuild/AppIcon.iconset/icon_${size}x${size}.png
done

# 使用iconutil创建.icns文件
iconutil -c icns iconbuild/AppIcon.iconset -o Cliper.app/Contents/Resources/AppIcon.icns

# 清理临时文件
rm -rf iconbuild