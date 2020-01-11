## Google Cloud Functions for [Eve-App](https://github.com/dariusbakunas/eve-app) project

Deploy commands:

```shell script
$ gcloud functions deploy --runtime=go111 --trigger-topic=eve-esi-cron Esi --env-vars-file .env.yaml
$ gcloud functions deploy --runtime=go111 --trigger-topic=esi-character-wallet-transactions ProcessCharacterWalletTransactions --env-vars-file .env.yaml
$ gcloud functions deploy --runtime=go111 --trigger-topic=esi-character-journal-entries ProcessCharacterJournalEntries --env-vars-file .env.yaml
$ gcloud functions deploy --runtime=go111 --trigger-topic=esi-character-skills ProcessCharacterSkills --env-vars-file .env.yaml
$ gcloud functions deploy --runtime=go111 --trigger-topic=esi-character-skill-queue ProcessCharacterSkillQueue --env-vars-file .env.yaml
$ gcloud functions deploy --runtime=go111 --trigger-topic=esi-character-market-orders ProcessCharacterMarketOrders --env-vars-file .env.yaml
```