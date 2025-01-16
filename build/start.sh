#!/bin/bash

sed -i "s~{{TASKMAN_MYSQL_HOST}}~$TASKMAN_MYSQL_HOST~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MYSQL_PORT}}~$TASKMAN_MYSQL_PORT~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MYSQL_SCHEMA}}~$TASKMAN_MYSQL_SCHEMA~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MYSQL_USER}}~$TASKMAN_MYSQL_USER~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MYSQL_PWD}}~$TASKMAN_MYSQL_PWD~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_LOG_LEVEL}}~$TASKMAN_LOG_LEVEL~g" /app/taskman/conf/default.json
sed -i "s~{{GATEWAY_URL}}~$GATEWAY_URL~g" /app/taskman/conf/default.json
sed -i "s~{{JWT_SIGNING_KEY}}~$JWT_SIGNING_KEY~g" /app/taskman/conf/default.json
sed -i "s~{{SUB_SYSTEM_CODE}}~$SUB_SYSTEM_CODE~g" /app/taskman/conf/default.json
sed -i "s~{{SUB_SYSTEM_KEY}}~$SUB_SYSTEM_KEY~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MAIL_SENDER_NAME}}~$TASKMAN_MAIL_SENDER_NAME~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MAIL_AUTH_SENDER}}~$TASKMAN_MAIL_AUTH_SENDER~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MAIL_AUTH_SERVER}}~$TASKMAN_MAIL_AUTH_SERVER~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MAIL_AUTH_PWD}}~$TASKMAN_MAIL_AUTH_PWD~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MAIL_AUTH_SSL}}~$TASKMAN_MAIL_AUTH_SSL~g" /app/taskman/conf/default.json
sed -i "s~{{MINIO_ADDRESS}}~$MINIO_ADDRESS~g" /app/taskman/conf/default.json
sed -i "s~{{MINIO_ACCESS_KEY}}~$MINIO_ACCESS_KEY~g" /app/taskman/conf/default.json
sed -i "s~{{MINIO_SECRET_KEY}}~$MINIO_SECRET_KEY~g" /app/taskman/conf/default.json
sed -i "s~{{ENCRYPT_SEED}}~$ENCRYPT_SEED~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_WEB_URL}}~$TASKMAN_WEB_URL~g" /app/taskman/conf/default.json
sed -i "s~{{TASKMAN_MENU_API_PERMISSION_ENABLE}}~$TASKMAN_MENU_API_PERMISSION_ENABLE~g" /app/taskman/conf/default.json

./taskman-server