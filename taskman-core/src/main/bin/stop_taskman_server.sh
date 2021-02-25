echo "Killing: `cat run_taskman_serve.pid`"
kill -9 `cat run_taskman_serve.pid`
rm -rf run_taskman_serve.pid
