先把 [main.go (line 25)](C:/ProgramData/FromD/ProgramFile/typora/golang/client-go-eg/main.go:25) 设为：
controllerChoice := 2
启动后手动修改 Deployment：
kubectl scale deployment webapp-demo -n test1 --replicas=3
WebApp 中声明的是 replicas: 1，controller 应该把 Deployment 恢复为 1：
kubectl get deployment webapp-demo -n test1
再删除 Deployment：
kubectl delete deployment webapp-demo -n test1
controller 应该重新创建它。这部分代码已经存在，因为 controller 同时监听了 Deployment 事件。