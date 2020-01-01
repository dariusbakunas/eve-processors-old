Deploy commands:

```shell script
$ gcloud functions deploy --runtime=go111 --trigger-topic=eve-esi-cron Esi --env-vars-file .env.yaml
```