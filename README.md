# BlockChainSafeTesting
A tool for BlockChian's safe testing  
Usage:  
使用方法：  

1.Download codes  
git clone https://github.com/HanDaXia/BlockChainSafeTesting.git  

1.下载代码  
git clone https://github.com/HanDaXia/BlockChainSafeTesting.git  

2.Compile server nodes  
To facilitate developers'use, we provide a one-click compilation function that generates all docker image files by executing make instructions in the project directory.  Before compiling, make sure your golang verison is not smaller than v1.11.

2.编译服务节点  
为了方便开发者使用，我们提供了一键编译的功能，只需要在工程目录执行make指令就可以生成所有的docker镜像文件。编译前请确保go版本在v1.11及以上。

3.Start testing network  
Type following command in your project directory   
docker-compose –f docker/docekr-compose-sys.yaml up  

3.启动测试网络  
在工程目录输入以下指令:  
docker-compose –f docker/docekr-compose-sys.yaml up  

4.Open PC client for testing（fabricclient/main.exe)  
At present, we only provide Windows version of the client, the user opens the client, imports the corresponding files to be tested and submits them to the server, then the corresponding test results can be obtained.  

4.打开PC客户端进行测试（fabricclient/main.exe)  
目前我们仅提供windows版本的客户端，用户打开客户端，导入对应的待测文件并提交到服务器，即可获得对应的检测结果  

