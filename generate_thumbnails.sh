#!/bin/bash
# 将D:/work/wechat_download_data/images目录下的所有图片在D:/work/wechat_download_data/thumbnails生成相关的缩略图

# 原图目录和缩略图目录
original_dir="D:/work/wechat_download_data/images"
thumbnail_dir="D:/work/wechat_download_data/thumbnails"

# 遍历原图目录中的所有图片文件
find "$original_dir" -type f -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.gif" -iname "*.webp" | while read -r original_image; do
    # 构造缩略图路径
    thumbnail_image="${original_image/$original_dir/$thumbnail_dir}"

    # 创建缩略图目录（如果不存在）
    mkdir -p "$(dirname "$thumbnail_image")"

    # 使用 ImageMagick 生成缩略图
    convert "$original_image" -resize 200x200 "$thumbnail_image"

    # 输出生成的缩略图路径
    echo "Generated thumbnail: $thumbnail_image"

    # 输出进度信息
    processed_count=$(find "$thumbnail_dir" -type f | wc -l)
    total_count=$(find "$original_dir" -type f -iname "*.jpg" -o -iname "*.jpeg" -o -iname "*.png" -o -iname "*.gif" -iname "*.webp" | wc -l)
    echo "Progress: $processed_count / $total_count"
done