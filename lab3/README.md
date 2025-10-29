# Memo简易备忘录

## 构建 & 运行

```bash
git clone https://github.com/Becks723/go_trial
cd go_trial/lab3/memo

# 构建
docker compose build

# 后台运行
docker compose up -d
```



## 项目依赖

- 语言：Golang1.25.1
- web框架：gin
  - Hertz实现位于分支`lab3/feature/hertz`中

- 数据库：MySQL8.0.42
- 数据库交互：gorm
- 文档：swaggo
- 部署：docker compose



## 项目结构

```shell
lab3/memo/
├── cmd         # main文件
├── config      # 配置
├── controller  # 控制层，包含所有api
├── docs        # swagger生成的api文档
├── dto         # dto对象
├── middleware  # 中间件
├── service     # 业务逻辑层
├── pkg         # 存放一些通用的包
│ ├── ctl        # 控制层的工具函数（如工厂生成response）
│ ├── e          # 错误代码
│ ├── util       # 工具函数
├── repository  # 数据库访问层（dao）
│ ├── model      # orm对象
```



## Bonus

1. 项目使用swagger自动生成api文档。

2. 三层架构设计:

   1. Controller层：接口层。封装各个api，接收req/query并返回resp，建立起一个基本的api实现框架，起到调度员的作用。见`controller`目录。
   2. Service层：核心业务逻辑。见`service`目录。
   3. DataAccess层：交互数据库，使Service层不必直接接触数据库对象。见`repository`目录。

   4. 在三层架构之间还有两种model：

      1. dto对象：本项目中dto对象负责连接接口层和业务逻辑层，起到解耦的作用。封装了req/query数据和resp数据。见`dto`目录。
      2. orm对象：负责映射到数据库。主要由dao使用，业务逻辑层中也有用到。见`repository/model`目录。

   5. 应对代码变动：

      本项目由于先做了gin实现，后期需要迁移为hertz实现。迁移变动位于commit [aa432276](https://github.com/Becks723/go_trial/commit/aa432276f69029a9c35325931ae99c9516c57f26)中。观察这个commit不难发现，实际上真正变化的部分仅仅有 Controller层、dto对象以及更上层的客户代码；而三层架构的下两层（Service和DA层）无需任何改动，达到了“隔离变化”的功能，符合三层架构的设计初衷。

3. 考虑数据库交互安全性：

   以占位符的形式调用gorm api，而不是直接将数据的值写进sql语句。如：

   ```go
   // 推荐
   db.Where("id = ? && name = ?", memoId, memoName)
   
   // 不推荐（伪代码）
   db.Where("id = "+memoId+" && name = "+memoName)
   ```

   这样能有效防止sql注入。

4. to be continued...

5. to be continued...
