JAVA_OPTS="-Xms128m -Xmx256m"

echo "taskman-server on...."

if find -name run_taskman_serve.pid | grep "run_taskman_server.pid";
then
  echo "taskman-server is running..."
  kill -9 `cat run_taskman_serve.pid`
  rm -rf run_taskman_serve.pid
  exit
fi
nohup java -jar taskman-0.0.1-SNAPSHOT.jar >  output 2>&1 &
if [ ! -z "run_taskman_serve.pid" ]; then
  echo $!> run_taskman_serve.pid
fi
