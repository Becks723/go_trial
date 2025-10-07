## 任务截图
![](task.png)

## 发布说明
- 0.1.0 - 2025-10-08  
初版。
  - 最基本的用户名+评论爬取（控制台打印）
  - 分页爬取
  - 无需cookie

## 技术细节
- 不同于document页面爬取，b站评论爬的是json。
- 新版b站评论接口使用了wbi签名鉴权，需要计算`w_rid`和`wts`字段。参考了https://github.com/SocialSisterYi/bilibili-API-collect/blob/master/docs/misc/sign/wbi.md