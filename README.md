Deploy commands:

```shell script
$ gcloud functions deploy --runtime=go111 --trigger-topic=eve-esi-cron Esi --set-env-vars GOOGLE_CLOUD_PROJECT={PROJECT_ID},CLOUD_SQL_CONNECTION_NAME={SQL_CONNECTION_NAME}
```