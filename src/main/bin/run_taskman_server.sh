JAVA_OPTS="-Xms128m -Xmx256m"

echo "taskman-server on...."

if find -name run_taskman_serve.pid | grep "run_taskman_server.pid";
then
  echo "taskman-server is running..."
  exit
fi 

CLASSPATH="$CLASSPATH":"./taskman-0.0.1-SNAPSHOT.jar"
echo $CLASSPATH
LIBPATH="./lib"
if [ -d "$LIBPATH" ]; then
  for i in "$LIBPATH"/*.jar; do
    CLASSPATH="$CLASSPATH":"$i"
  done
fi

echo "Using CLASSPATH:   $CLASSPATH"

nohup java $JAVA_OPTS \
    -classpath $CLASSPATH \
    com.webank.taskman.TaskManApplication >  output 2>&1 &

if [ ! -z "run_taskman_serve.pid" ]; then
  echo $!> run_taskman_serve.pid
fi

