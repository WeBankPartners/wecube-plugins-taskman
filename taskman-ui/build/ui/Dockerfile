FROM nginx
LABEL maintainer = "Webank CTB Team"
# RUN mkdir -p /root/wecube-docs
COPY dist /root/app/
COPY build/ui/nginx.tpl /etc/nginx/conf.d/nginx.tpl
COPY build/ui/start_taskman.sh /etc/nginx/start_taskman.sh
RUN chmod +x /etc/nginx/start_taskman.sh
# RUN rm -rf /root/wecube-docs/.git