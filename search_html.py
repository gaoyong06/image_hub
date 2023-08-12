# -*- coding: utf-8 -*-
'''
Author: gaoyong gaoyong06@qq.com
Date: 2023-08-09 10:49:39
LastEditors: gaoyong gaoyong06@qq.com
LastEditTime: 2023-08-12 14:22:32
FilePath: \image_hub\article_tag_items.py
Description:
'''
# 遍历D:/work/wechat_download_data目录及其子目录下的所有html文件，查找文本关键字所在的html文件路径
# 在命令行中运行脚本时，需要传递两个参数：目录路径和要搜索的文本
# 使用方法：python search_html.py /path/to/directory Text to search
# 示例：python search_html.py D:/work/wechat_download_data "丨人间值得"
# 类似：grep  -r -l "丨人间值得" "D:/work/wechat_download_data/html/*.html"
import os
import sys
import concurrent.futures

def search_html_files(directory, text):
    html_files = []
    count = 0  # 统计已处理的文件数量
    total_files = count_total_files(directory)  # 获取总文件数量
    print("Searching for text in HTML files...")
    
    # 遍历目录下的所有文件和子目录
    with concurrent.futures.ThreadPoolExecutor() as executor:
        for root, dirs, files in os.walk(directory):
            for file in files:
                if file.endswith(".html"):
                    count += 1
                    filepath = os.path.join(root, file)
                    # 在文件中搜索文本
                    future = executor.submit(search_text_in_file, filepath, text)
                    future.add_done_callback(log_search_result(filepath))  # 添加回调函数记录搜索结果
                    if future.result():
                        html_files.append(filepath)
                    print_progress(count, total_files)  # 打印处理进度
    
    print("Total files processed:", count)
    return html_files

def count_total_files(directory):
    count = 0  # 统计文件数量
    
    for root, dirs, files in os.walk(directory):
        for file in files:
            if file.endswith(".html"):
                count += 1
    
    return count

def search_text_in_file(filepath, text):
    with open(filepath, 'r', encoding='utf-8') as file:
        contents = file.read()
        if text in contents:
            return True
    return False

def log_search_result(filepath):
    def callback(future):
        result = future.result()
        if result:
            print("Text found in:", filepath)
        else:
            print("Text not found in:", filepath)
    
    return callback

def print_progress(count, total):
    progress = count / total * 100
    print("Progress: {:.2f}% ({} / {})".format(progress, count, total))

def main():
    # 检查命令行参数
    if len(sys.argv) < 3:
        print("Usage: python search_html.py [directory] [text]")
        return
    
    directory = sys.argv[1]
    search_text = sys.argv[2]
    
    result = search_html_files(directory, search_text)
    if result:
        print("HTML files containing the text:")
        for filepath in result:
            print(filepath)
    else:
        print("No HTML files containing the text were found.")

if __name__ == "__main__":
    main()