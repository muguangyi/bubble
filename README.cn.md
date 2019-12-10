# Bubble

[![Build Status](https://travis-ci.com/muguangyi/bubble.svg?branch=master)](https://travis-ci.com/muguangyi/bubble) [![codecov](https://codecov.io/gh/muguangyi/bubble/branch/master/graph/badge.svg)](https://codecov.io/gh/muguangyi/bubble)

> 非对称持续集成框架。

## 背景

如果让任何CI机器都能够执行任何类型的任务，则有一个预设前提：所有的机器必须具备同样的系统环境，像安装的软件，等等。但是对于小的团队或者公司很难为不同的CI需求维护很多机器。以`Unity`作为例子，`Unity`每年都会发布很多个版本，而不同的团队可能使用不同的版本。那么如果要建立一个CI的集群就要求每一台CI机器都要安装所有可能用到的版本，这会比较困难。另外一个选择就是一台机器只处理固定有限的任务，但是缺点就是物理资源不能被充分使用，并且缺少灵活性。

Bubble希望建立一个更为灵活的解决方案。Bubble能将一个任务分解为多个步骤，并在不同的机器上执行。同时，每一个物理机可以有不同的环境，而Bubble可以搜集每一个机器能执行什么样的工作，从而决定哪里执行目标步骤。

![bubble](doc/bubble.png)

### 功能

* 任务可分解。
* 动态分配任务步骤。
* 工作节点具有自述性。
* 内建很多命令：shell，zip，ftp，unity，email notification，等等。
* 支持变量和方法。
* 任务追踪，例如执行时间。
* 工作节点追踪。

## 安装

* 下载目标二进制[版本](https://github.com/muguangyi/bubble/releases).
* 解压到本地目录。

## 快速开始

Bubble有两种节点类型：`Master`和`Worker`.

|概念|描述|
|--:|:--|
|`Master`|主节点可以解析和分解任务，并拆分为步骤，以及分发到工作节点。|
|`Worker`|工作节点可以处理各种命令执行，以及向master同步结果。|

### 前置条件

* 首先本地启动[Redis](https://redis.io)。

### ① 运行master

* 进入master目录。
* 运行master。
  
  **Windows**:

  ```shell
  > bubble-master.exe
  ```

  **MacOS**/**Linux**:

  ```shell
  > ./bubble-master
  ```

### ② 运行worker

* 进入worker目录。
* 运行worker。
  
  **Windows**:

  ```shell
  > bubble-worker.exe
  ```

  **MacOS**/**Linux**:

  ```shell
  > ./bubble-worker
  ```

### ③ 访问Master主页

浏览器中访问localhost，如果显示以下结果，那么恭喜，Bubble已经可以工作了。

![result](doc/result.png)

### ④ 第一个任务

* 输入一个任务名，例如Test。
* 点击`CREATE`。

![first-job](doc/first-job.png)

### ⑤ 执行Job

* 点击`Setting`可以查看Bubble脚本。
  
  ```yml
  # .bubble.yml
  -
   action: shell
   script:
   - echo Hello Bubble!
  ```

* 点击`Trigger`可以触发执行当前任务。

## 文档

详细信息请查看[Wiki](https://github.com/muguangyi/bubble/wiki).

## 维护者

[@MuGuangyi](https://github.com/muguangyi)

## 许可证

[MIT](LICENSE) © MuGuangyi
