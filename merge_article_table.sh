#!/bin/bash

# 在image_hub数据库中，有很多数据表，一个是：tbl_article， 一个是tbl_article_后缀, 其中后缀由数字和英文字母组成
# 读出image_hub数据库中的所有tbl_article_后缀 数据表，然后全部写入tbl_article数据表中
# 即：tbl_article 是所有 tbl_article_后缀 的数据合并表
# tbl_article和tbl_article_后缀 数据表的字段都是相同的，当出现写入tbl_article时，出现主键相同的重复记录时，则用所以新的所有字段值覆盖旧的所有字段值

# MySQL connection information
host="192.168.1.5"
port="3306"
user="root"
passwd="root"
database="image_hub"

# Get a list of all tbl_article_ suffix tables
suffix_tables=$(mysql -h $host -P $port -u $user -p$passwd $database -e "SHOW TABLES LIKE 'tbl_article_%';" | tail -n +2)

# Count the total number of suffix tables
total_tables=$(echo "$suffix_tables" | wc -l)

# Initialize a counter for the progress
counter=0

# Loop through each suffix table and insert its records into tbl_article
for table in "$suffix_tables"; do
    counter=$((counter+1))
    echo "Processing table $counter of $total_tables: $table"
    
    mysql -h $host -P $port -u $user -p$passwd $database -e "INSERT INTO tbl_article (sn, mid, idx, biz, author, wechat_id, title, tags, sections, local_path, publish_time, created_at, updated_at, deleted_at) SELECT sn, mid, idx, biz, author, wechat_id, title, tags, sections, local_path, publish_time, created_at, updated_at, deleted_at FROM "$table" ON DUPLICATE KEY UPDATE tbl_article.sn = VALUES(sn), tbl_article.mid = VALUES(mid), tbl_article.idx = VALUES(idx), tbl_article.biz = VALUES(biz), tbl_article.author = VALUES(author), tbl_article.wechat_id = VALUES(wechat_id), tbl_article.title = VALUES(title), tbl_article.tags = VALUES(tags), tbl_article.sections = VALUES(sections), tbl_article.local_path = VALUES(local_path), tbl_article.publish_time = VALUES(publish_time), tbl_article.created_at = VALUES(created_at), tbl_article.updated_at = VALUES(updated_at), tbl_article.deleted_at = VALUES(deleted_at);"
    
    echo "Table $counter of $total_tables: $table processed"
done