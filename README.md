
赛题内容：多个注册中心，链接多个应用，题目要求优化分治逻辑，让分系统均匀承压

赛题公式化描述：注册中心实际加载的内存与理想的内存占比 乘以 注册中心的内存标准差 加 多个应用的标准差 尽可能小

赛题解释：让每个注册中心加载的服务数相近，每个中心链接的应用相近，尽量不要重复加载服务

算法描述：多目标优化问题，具体为非线性的整数规划组合优化问题，多维装箱问题模型，计算应用间相似度，对应用分组、合并

思路解法：1.仅量不要重复加载服务数据 2.注册中心链接的应用相近 3.动态部分建模

获得成果：并查集等算法的灵活应用，融合了深度学习的思路，不断迭代的去优化整个搜索过程，完成一个实际应用问题的解答


成就：阿里巴巴集团内部选手排行榜-季军

排行：https://tianchi.aliyun.com/forum/postDetail?spm=5176.12586969.0.0.21ec62b5YPzs4V&postId=117150

导师解析：https://tianchi.aliyun.com/forum/postDetail?postId=109732

我的代码：https://code.alibaba-inc.com/linfan.wty/tianchi/blob/master/pkg/cmd/playermock/manager.go

看到一位做法相近同学的思路：https://blog.csdn.net/qq_37969433/article/details/107384816


```sh
# 需要 go 1.13 及以上，请启用 go modules
go build -o main demo/pkg/cmd/playermock
```

请使用附件中的 Dockerfile 编译成 Docker 镜像运行，详细请参看[这里](https://code.aliyun.com/middleware-contest-2020/pilotset/blob/master/LOCAL_DEV.md)。如果需要在本机运行，
需要建立 /root/input, /root/output 以及 /root/data 目录，将测试数据解压至 /root/data，启动 demo 后再启动评分程序。
