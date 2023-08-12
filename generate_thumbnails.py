# -*- coding: utf-8 -*-
'''
Author: gaoyong gaoyong06@qq.com
Date: 2023-08-12 21:39:38
LastEditors: gaoyong gaoyong06@qq.com
LastEditTime: 2023-08-12 21:57:00
FilePath: \image_hub\generate_thumbnails.py
Description: 生成缩略图
'''
import os
import concurrent.futures
from PIL import Image
import logging

# 配置日志
logging.basicConfig(level=logging.INFO, format='[%(levelname)s] %(asctime)s - %(message)s')

# 定义图片目录和缩略图目录
image_dir = 'D:/work/wechat_download_data/images'
thumbnail_dir = 'D:/work/wechat_download_data/thumbnails'

# 定义生成缩略图的大小
thumbnail_size = (200, 200)

# 递归遍历目录下的所有图片文件
def process_directory(directory):
    for filename in os.listdir(directory):
        filepath = os.path.join(directory, filename)
        if os.path.isdir(filepath):
            process_directory(filepath)
        else:
            # 生成缩略图并保存到对应目录
            generate_thumbnail(filepath)

# 生成缩略图并保存
def generate_thumbnail(image_path):
    try:
        image = Image.open(image_path)
        image.thumbnail(thumbnail_size)
        thumbnail_path = get_thumbnail_path(image_path)
        os.makedirs(os.path.dirname(thumbnail_path), exist_ok=True)
        image.save(thumbnail_path)
        logging.info(f'Generated thumbnail: {thumbnail_path}')
    except Exception as e:
        logging.error(f'Failed to generate thumbnail: {image_path} - {e}')

# 根据图片路径获取对应的缩略图路径
def get_thumbnail_path(image_path):
    relative_path = os.path.relpath(image_path, image_dir)
    return os.path.join(thumbnail_dir, relative_path)

# 处理单个图片文件
def process_image_file(image_path):
    generate_thumbnail(image_path)

# 控制线程数量处理图片
def process_images_with_thread_pool():
    with concurrent.futures.ThreadPoolExecutor(max_workers=4) as executor:  # 控制线程数量为4
        for root, dirs, files in os.walk(image_dir):
            for file in files:
                image_path = os.path.join(root, file)
                executor.submit(process_image_file, image_path)

# 主函数
if __name__ == '__main__':
    # 创建缩略图目录
    os.makedirs(thumbnail_dir, exist_ok=True)

    logging.info('Starting image processing...')
    process_images_with_thread_pool()
    logging.info('Image processing completed.')