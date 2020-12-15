project_name=taskman
version=v1.0
current_dir=$(shell pwd)
remote_docker_image_registry=ccr.ccs.tencentyun.com/webankpartners/wecube-app
s3_server_url=http://10.10.10.1:9000
s3_access_key=access_key
s3_secret_key=secret_key

clean:
	rm -rf $(current_dir)/target
	rm -rf $(current_dir)/taskman-ui/dist

.PHONY:build
build:
	mkdir -p repository
	docker run --rm --name taskman-build -e SASS_BINARY_SITE=https://npm.taobao.org/mirrors/node-sass -v /data/wecube_repository:/usr/src/mymaven/repository -v $(current_dir)/build/maven_settings.xml:/usr/share/maven/ref/settings-docker.xml  -v $(current_dir):/usr/src/mymaven -w /usr/src/mymaven maven:3.3-jdk-8 mvn -U clean install -Dmaven.test.skip=true -s /usr/share/maven/ref/settings-docker.xml dependency:resolve

image:
	docker build -t $(project_name):$(version) .
	
.PHONY:package
package:
	rm -rf package
	mkdir -p package
	cd package && docker save -o image.tar $(project_name):$(version) 
	
	rm -rf taskman-ui/dist/*
	cd taskman-ui && cnpm install && npm run plugin
	cd taskman-ui && zip -r ui.zip dist/* && mv ui.zip ../package/ui.zip
	
	cp src/main/resources/database/init.sql package/init.sql
	
	git checkout -- .
	sh build/build_plugin_xml.sh $(version)
	cd package && cp ../register.xml .
	
	cd package && rm -f wecube-plugins-taskman-$(version).zip
	cd package && zip -r wecube-plugins-taskman-$(version).zip .
	
	docker run --name minio-client-taskman -v `pwd`/package:/package -itd --entrypoint=/bin/sh minio/mc
	docker exec minio-client-taskman mc config host add wecubeS3 $(s3_server_url) $(s3_access_key) $(s3_secret_key) wecubeS3
	docker exec minio-client-taskman mc cp /package/wecube-plugins-$(project_name)-$(version).zip wecubeS3/wecube-plugin-package-bucket
	docker rm -f minio-client-taskman
	
	docker rmi $(project_name):$(version)
