# -*- coding: utf-8 -*-
'''
Author: gaoyong gaoyong06@qq.com
Date: 2023-08-09 10:49:39
LastEditors: gaoyong gaoyong06@qq.com
LastEditTime: 2023-08-09 11:57:37
FilePath: \image_hub\article_tag_items.py
Description:
'''
# 遍历D:/work/wechat_download_data目录及其子目录下的所有html文件，解析html 得到css class 为article-tag__item的标签内的文本内容，写入一个txt文件
# 例如：会找到下面的html
# <span class="article-tag__item">#情侣头像</span>
# 然后会将#情侣头像写入txt文件
import os
from bs4 import BeautifulSoup

def parse_html_file(file_path):
    with open(file_path, 'r', encoding='utf-8') as file:
        html_content = file.read()
        soup = BeautifulSoup(html_content, 'html.parser')
        tags = soup.find_all(class_='article-tag__item')
        return [tag.text.strip() for tag in tags]

def traverse_directory(directory):
    html_files = []
    for root, _, files in os.walk(directory):
        for file in files:
            if file.endswith('.html'):
                html_files.append(os.path.join(root, file))
    return html_files

def write_to_file(file_path, content):
    with open(file_path, 'a', encoding='utf-8') as file:  # 以追加模式打开文件
        for tag in content:
            file.write(tag + '\n')

def main():
    directory = 'D:/work/wechat_download_data'  # 替换为你要遍历的目录路径
    output_file = 'article_tag_items_output.txt'  # 替换为输出文件路径及文件名

    html_files = traverse_directory(directory)
    total_files = len(html_files)
    processed_files = 0

    existing_tags = set()
    if os.path.exists(output_file):
        with open(output_file, 'r', encoding='utf-8') as file:
            existing_tags = {line.strip() for line in file}
    
    for file_path in html_files:
        processed_files += 1
        print(f'Parsing file: {file_path} ({processed_files}/{total_files})')

        result = parse_html_file(file_path)
        if result:
            print(f'Found {len(result)} tag(s)')
        else:
            print('No tag found')

        new_tags = [tag for tag in result if tag not in existing_tags]  # 获取新的标签
        write_to_file(output_file, new_tags)
        existing_tags.update(new_tags)  # 将新标签添加到已有标签集合中
    
    print(f'Parsing completed. Results written to file: {output_file}')

if __name__ == '__main__':
    main()