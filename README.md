# chatgpt - azure openai

### 怎么用

```bash
cp .env.example .env

vim .env
# PASSWORD=认证信息，在打开页面后输入才可使用
# AZURE_OPENAI_KEY=
# ENDPOINT=https://xxx.openai.azure.com/openai/deployments/xxx/chat/completions?api-version=2023-05-15

# 配置好这三个启动就可以用啦
go build
./chatgpt
```
